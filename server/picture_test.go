package main

import "testing"

func TestDeletePictureFromS3(t *testing.T) {
	deleteFromS3("users/3/buzz.png")
}
