package main

import (
	"testing"
	"time"
)

func TestDeletePictureFromS3(t *testing.T) {
	err := deleteFromS3("users/3/buzz.png")
	if err != nil {
		panic(err)
	}
}

func TestDeletePicturesFromS3(t *testing.T) {
	pics := []Picture{
		Picture{
			ImagePath: "users/3/buzz.png",
		},
	}
	err := deleteMultipleFromS3(pics)
	if err != nil {
		panic(err)
	}

	<-time.After(10 * time.Second)
}
