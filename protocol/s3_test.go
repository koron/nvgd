package protocol

import (
	"reflect"
	"testing"

	"github.com/koron/nvgd/config"
)

func TestS3Config(t *testing.T) {
	_, err := config.LoadConfig("testdata/s3_test.yml")
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
	if !reflect.DeepEqual(act0, exp0) {
		t.Errorf("default config:%#v is not as expected %#v", act0, exp0)
	}

	var (
		act1 = s3config.bucketConfig("foo")
		exp1 = &S3BucketConfig{
			Region:          "aaa",
			AccessKeyID:     "bbb",
			SecretAccessKey: "ccc",
			SessionToken:    "ddd",
		}
	)
	if !reflect.DeepEqual(act1, exp1) {
		t.Errorf("foo %#v is not as expected %#v", act1, exp1)
	}

	var (
		act2 = s3config.bucketConfig("bar")
		exp2 = &S3BucketConfig{
			Region:          "eee",
			AccessKeyID:     "fff",
			SecretAccessKey: "ggg",
			SessionToken:    "hhh",
		}
	)
	if !reflect.DeepEqual(act2, exp2) {
		t.Errorf("bar %#v is not as expected %#v", act2, exp2)
	}
}
