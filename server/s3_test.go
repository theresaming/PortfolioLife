package main

import (
	"fmt"
	"testing"

	minio "github.com/minio/minio-go"
)

func TestBuckets(t *testing.T) {
	s3Client, err := minio.New(config.S3Endpoint, config.S3Key, config.S3Secret, true)
	if err != nil {
		panic(err)
	}

	for object := range s3Client.ListObjects("portfoliolife", "", true, nil) {
		fmt.Println(object.Key)
	}
}
