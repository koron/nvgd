package protocol

import (
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/koron/nvgd/config"
)

var s3handler = &S3Handler{
	Config: S3Config{
		Buckets: map[string]S3BucketConfig{},
	},
}

func init() {
	MustRegister("s3", s3handler)
	config.RegisterProtocol("s3", &s3handler.Config)
}

// S3Handler is AWS S3 protocol handler
type S3Handler struct {
	Config S3Config
}

// Open opens a S3 URL.
func (ph *S3Handler) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		bucket = u.Host
		key    = u.Path
	)
	conf := ph.Config.bucketConfig(bucket).awsConfig()
	sess := session.New(conf)
	svc := s3.New(sess)
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
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
