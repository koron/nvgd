package aws

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/internal/ltsv"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

const (
	S3Token = "s3token"
)

var rxLastComponent = regexp.MustCompile(`[^/]+/?$`)

var s3config = S3Config{
	Default: S3BucketConfig{},
	Buckets: map[string]S3BucketConfig{},
}

var s3ObjHandler = &S3ObjHandler{
	Config: &s3config,
}

var s3ListHandler = &S3ListHandler{
	Config: &s3config,
}

func init() {
	protocol.MustRegister("s3obj", s3ObjHandler)
	protocol.MustRegister("s3list", s3ListHandler)
	config.RegisterProtocol("s3", &s3config)
}

// S3ObjHandler is AWS S3 object protocol handler
type S3ObjHandler struct {
	Config *S3Config
}

var _ protocol.Rangeable = (*S3ObjHandler)(nil)

func (ph *S3ObjHandler) newS3(u *url.URL) (svc *s3.S3, bucket, key string, err error) {
	bucket = u.Host
	key = u.Path
	conf := ph.Config.bucketConfig(bucket).awsConfig()
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, "", "", err
	}
	return s3.New(sess), bucket, key, nil
}

// Open opens a S3 URL.
func (ph *S3ObjHandler) Open(u *url.URL) (*resource.Resource, error) {
	svc, bucket, key, err := ph.newS3(u)
	if err != nil {
		return nil, err
	}
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return resource.New(out.Body), nil
}

func (ph *S3ObjHandler) Size(u *url.URL) (int, error) {
	svc, bucket, key, err := ph.newS3(u)
	if err != nil {
		return 0, err
	}
	out, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, err
	}
	if out.ContentLength == nil {
		return 0, fmt.Errorf("failed to get size of the object: bucket=%s key=%s", bucket, key)
	}
	return int(*out.ContentLength), nil
}

func (ph *S3ObjHandler) OpenRange(u *url.URL, start, end int) (*resource.Resource, error) {
	svc, bucket, key, err := ph.newS3(u)
	if err != nil {
		return nil, err
	}
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	})
	if err != nil {
		return nil, err
	}
	return resource.New(out.Body), nil
}

// S3ListHandler is AWS S3 list protocol handler
type S3ListHandler struct {
	Config *S3Config
}

var _ protocol.Protocol = (*S3ListHandler)(nil)

// Open opens a S3 URL.
func (ph *S3ListHandler) Open(u *url.URL) (*resource.Resource, error) {
	var (
		bucket = u.Host
		prefix = u.Path
	)
	conf := ph.Config.bucketConfig(bucket)
	sess, err := session.NewSession(conf.awsConfig())
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	if len(prefix) > 0 {
		prefix = prefix[1:]
	}
	in := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
		MaxKeys:   conf.maxKeys(),
	}
	// Setup continuation token of the request if available.
	if s := u.Query().Get(S3Token); len(s) > 0 {
		in.ContinuationToken = aws.String(s)
	}
	out, err := svc.ListObjectsV2(in)
	if err != nil {
		return nil, err
	}
	rc, err := ph.writeAsLTSV(out, bucket)
	if err != nil {
		return nil, err
	}
	rs := resource.New(rc)
	rs.Put(commonconst.LTSV, true)
	rs.Put(commonconst.ParsedKeys, []string{S3Token})
	// Add uplink if available.
	if prefix != "" {
		up := rxLastComponent.ReplaceAllString(prefix, "")
		link := fmt.Sprintf("/s3list://%s/%s?indexhtml", bucket, up)
		rs.Put(commonconst.UpLink, link)
	}
	// Embed next continuation token to rs if available.
	if out.NextContinuationToken != nil && *out.NextContinuationToken != "" {
		t := url.QueryEscape(*out.NextContinuationToken)
		link := fmt.Sprintf("/s3list://%s/%s?%s=%s",
			bucket, prefix, S3Token, t)
		rs.Put(commonconst.NextLink, link)
	}
	return rs, nil
}

func timeStr(ti time.Time) string {
	if s3config.UseUnixtime {
		return strconv.FormatInt(ti.Unix(), 10)
	}
	return ti.Format(time.RFC1123)
}

func (ph *S3ListHandler) writeAsLTSV(out *s3.ListObjectsV2Output, bucket string) (io.ReadCloser, error) {
	var (
		buf = &bytes.Buffer{}
		w   = ltsv.NewWriter(buf, "name", "type", "size", "modified_at", "link", "download")
	)
	// add prefixes
	for _, item := range out.CommonPrefixes {
		link := fmt.Sprintf("/s3list://%s/%s", bucket, *item.Prefix)
		err := w.Write(*item.Prefix, "prefix", "", "", link, "")
		if err != nil {
			return nil, err
		}
	}
	// show objects
	for _, obj := range out.Contents {
		link := fmt.Sprintf("/s3obj://%s/%s", bucket, *obj.Key)
		download := link + "?all&download"
		mtime := timeStr(obj.LastModified.In(ph.Config.location()))
		err := w.Write(*obj.Key, "object", strconv.FormatInt(*obj.Size, 10), mtime, link, download)
		if err != nil {
			return nil, err
		}
	}
	return io.NopCloser(buf), nil
}

// S3Config is configuration of S3 protocol handler.
type S3Config struct {

	// Timezone forces timezone of modified times or so.
	Timezone string `yaml:"timezone,omitempty"`

	// Default is default bucket configuration.
	Default S3BucketConfig `yaml:"default,omitempty"`

	// Buckets
	Buckets map[string]S3BucketConfig `yaml:"buckets,omitempty"`

	loc *time.Location

	// UseUnixtime makes times in UNIX format: modified_at or so.
	UseUnixtime bool `yaml:"use_unixtime"`
}

func (c *S3Config) location() *time.Location {
	if c.loc != nil {
		return c.loc
	}
	if c.Timezone != "" {
		l, err := time.LoadLocation(c.Timezone)
		if err == nil {
			c.loc = l
			return c.loc
		}
		// FIXME: use server's logger.
		log.Printf("unknown timezone %q: %s", c.Timezone, err)
	}
	c.loc = time.Local
	return c.loc
}

func (c *S3Config) bucketConfig(bucket string) *S3BucketConfig {
	bc, ok := c.Buckets[bucket]
	if !ok {
		return &c.Default
	}
	return &bc
}

// S3BucketConfig is AWS configuration for buckets.
type S3BucketConfig struct {

	// Region is AWS region.
	//
	// @see http://docs.aws.amazon.com/general/latest/gr/rande.html
	Region string `yaml:"region"`

	// AccessKeyID is AWS access key ID.
	AccessKeyID string `yaml:"access_key_id"`

	// SecrentAccessKey is AWS secrent access key.
	SecretAccessKey config.SecretString `yaml:"secret_access_key"`

	// SessionToken is AWS session token.
	SessionToken string `yaml:"session_token,omitempty"`

	// MaxKeys used for S3 object listing.
	MaxKeys int64 `yaml:"max_keys,omitempty"`

	// HTTPProxy used as HTTP proxy to access S3.
	HTTPProxy string `yaml:"http_proxy,omitempty"`
}

func (bc *S3BucketConfig) region() string {
	if bc.Region == "" {
		return "ap-northeast-1"
	}
	return bc.Region
}

func (bc *S3BucketConfig) creds() *credentials.Credentials {
	return credentials.NewStaticCredentials(bc.AccessKeyID, string(bc.SecretAccessKey), bc.SessionToken)
}

func (bc *S3BucketConfig) awsConfig() *aws.Config {
	conf := aws.NewConfig().
		WithRegion(bc.region()).
		WithCredentials(bc.creds())
	if cl := bc.httpClient(); cl != nil {
		conf = conf.WithHTTPClient(cl)
	}
	return conf
}

func (bc *S3BucketConfig) maxKeys() *int64 {
	if bc.MaxKeys <= 0 || bc.MaxKeys >= 1000 {
		return nil
	}
	return aws.Int64(bc.MaxKeys)
}

func (bc *S3BucketConfig) httpClient() *http.Client {
	if bc.HTTPProxy == "" {
		return nil
	}
	return NewProxyClient(bc.HTTPProxy)
}
