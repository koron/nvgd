package protocol

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/koron/nvgd/config"
)

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
	MustRegister("s3obj", s3ObjHandler)
	MustRegister("s3list", s3ListHandler)
	config.RegisterProtocol("s3", &s3config)
}

// S3ObjHandler is AWS S3 object protocol handler
type S3ObjHandler struct {
	Config *S3Config
}

// Open opens a S3 URL.
func (ph *S3ObjHandler) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		bucket = u.Host
		key    = u.Path
	)
	conf := ph.Config.bucketConfig(bucket).awsConfig()
	sess := session.New(conf)
	svc := s3.New(sess)
	return ph.getObject(svc, bucket, key)
}

func (ph *S3ObjHandler) getObject(svc *s3.S3, bucket, key string) (io.ReadCloser, error) {
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// S3ListHandler is AWS S3 list protocol handler
type S3ListHandler struct {
	Config *S3Config
}

// Open opens a S3 URL.
func (ph *S3ListHandler) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		bucket = u.Host
		key    = u.Path
	)
	conf := ph.Config.bucketConfig(bucket).awsConfig()
	sess := session.New(conf)
	svc := s3.New(sess)
	return ph.listObjects(svc, bucket, key[1:])
}

func (ph *S3ListHandler) listObjects(svc *s3.S3, bucket, prefix string) (io.ReadCloser, error) {
	out, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}
	var (
		buf  = &bytes.Buffer{}
		cols = make([]string, 0, 4)
	)
	for _, obj := range out.Contents {
		cols := append(cols, *obj.Key, "obj", strconv.FormatInt(*obj.Size, 10),
			obj.LastModified.Format(time.RFC1123))
		_, err := buf.WriteString(strings.Join(cols, "\t") + "\n")
		if err != nil {
			return nil, err
		}
		cols = cols[0:0]
	}
	return ioutil.NopCloser(buf), nil
}

// S3Config is configuration of S3 protocol handler.
type S3Config struct {

	// Default is default bucket configuration.
	Default S3BucketConfig `yaml:"default,omitempty"`

	// Buckets
	Buckets map[string]S3BucketConfig `yaml:"buckets"`
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
	SecretAccessKey string `yaml:"secret_access_key"`

	// SessionToken is AWS session token.
	SessionToken string `yaml:"session_token"`
}

func (bc *S3BucketConfig) region() string {
	if bc.Region == "" {
		return "ap-northeast-1"
	}
	return bc.Region
}

func (bc *S3BucketConfig) creds() *credentials.Credentials {
	return credentials.NewStaticCredentials(bc.AccessKeyID, bc.SecretAccessKey, bc.SessionToken)
}

func (bc *S3BucketConfig) awsConfig() *aws.Config {
	return aws.NewConfig().WithRegion(bc.region()).WithCredentials(bc.creds())
}
