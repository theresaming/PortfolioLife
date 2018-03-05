package main

import (
	"fmt"
	"testing"
)

func TestAddTag(t *testing.T) {
	u := &User{
		ID: 3,
	}
	pic, err := getPicture(u, "N0JHtI2jjdSaHZI86cxwMIHqaCzoJPOs", false)
	if err != nil {
		panic(err)
	}
	tag := Tag{Tag: "testing2"}
	createTags(pic, []Tag{tag})
}

func TestDeleteTag(t *testing.T) {
	u := &User{
		ID: 3,
	}
	pic, err := getPicture(u, "N0JHtI2jjdSaHZI86cxwMIHqaCzoJPOs", false)

	if err != nil {
		panic(err)
	}

	err = deleteTags(pic, []string{"test3", "testing"})
	if err != nil {
		panic(err)
	}
	fmt.Println(pic.Tags)
}

func TestGetTags(t *testing.T) {
	u := &User{
		ID: 3,
	}
	pic, err := getPicture(u, "N0JHtI2jjdSaHZI86cxwMIHqaCzoJPOs", false)

	if err != nil {
		panic(err)
	}
	tags, err := getTags(pic)
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}