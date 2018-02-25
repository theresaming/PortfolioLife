package main

import (
	"testing"
)

func TestDeletePictureFromS3(t *testing.T) {
	deleteFromS3("users/3/buzz.png")
}

func TestDeletePicturesFromS3(t *testing.T) {
	pics := []Picture{
		Picture{
			ImagePath: "users/3/buzz.png",
		},
		Picture{
			ImagePath: "users/3/buzz.png",
		},
	}
	err := deleteMultipleFromS3(pics)
	if err != nil {
		panic(err)
	}
}
