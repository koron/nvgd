package protocol

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/ltsv"
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
	if len(key) > 0 {
		key = key[1:]
	}
	return ph.listObjects(svc, bucket, key)
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
		buf = &bytes.Buffer{}
		w   = ltsv.NewWriter(buf, "name", "type", "size", "modified_at", "link")
	)
	for _, item := range out.CommonPrefixes {
		link := fmt.Sprintf("/s3list://%s/%s?indexhtml", bucket, *item.Prefix)
		err := w.Write(*item.Prefix, "prefix", "", "", link)
		if err != nil {
			return nil, err
		}
	}
	for _, obj := range out.Contents {
		link := fmt.Sprintf("/s3obj://%s/%s", bucket, *obj.Key)
		t := obj.LastModified.In(ph.Config.location())
		err := w.Write(*obj.Key, "object", strconv.FormatInt(*obj.Size, 10),
			t.Format(time.RFC1123), link)
		if err != nil {
			return nil, err
		}
	}
	return ioutil.NopCloser(buf), nil
}

// S3Config is configuration of S3 protocol handler.
type S3Config struct {

	// Timezone forces timezone of modified times or so.
	Timezone string `yaml:"timezone,omitempty"`

	// Default is default bucket configuration.
	Default S3BucketConfig `yaml:"default,omitempty"`

	// Buckets
	Buckets map[string]S3BucketConfig `yaml:"buckets"`

	loc *time.Location
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
