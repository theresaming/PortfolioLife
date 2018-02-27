package main

import (
	"testing"
)

func TestAddTag(t *testing.T) {
	u := &User{
		ID: 3,
	}
	pic, err := getPicture(u, "-gtPq6Jd8ZHwHF5aeMMFr8ux-Gc6mOsP", false)
	if err != nil {
		panic(err)
	}
	pic.Tags = append(pic.Tags, Tag{Tag: "testing"})
	savePicture(pic)
}
