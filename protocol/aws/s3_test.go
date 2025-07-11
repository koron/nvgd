package aws

import (
	"testing"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/assert"
)

func TestS3Config(t *testing.T) {
	_, err := config.LoadConfig("s3_test.yml")
	if err != nil {
		t.Fatal(err)
	}

	var (
		act0 = s3config.bucketConfig("")
		exp0 = &S3BucketConfig{
			AccessKeyID:     "XXX",
			SecretAccessKey: "YYY",
		}
	)
	assert.Equal(t, exp0, act0, "default config")

	var (
		act1 = s3config.bucketConfig("foo")
		exp1 = &S3BucketConfig{
			Region:          "aaa",
			AccessKeyID:     "bbb",
			SecretAccessKey: "ccc",
			SessionToken:    "ddd",
		}
	)
	assert.Equal(t, exp1, act1, "%q config", "foo")

	var (
		act2 = s3config.bucketConfig("bar")
		exp2 = &S3BucketConfig{
			Region:          "eee",
			AccessKeyID:     "fff",
			SecretAccessKey: "ggg",
			SessionToken:    "hhh",
		}
	)
	assert.Equal(t, exp2, act2, "%q config", "bar")
}
