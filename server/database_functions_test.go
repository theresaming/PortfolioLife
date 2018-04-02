package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetUserWithEmail(t *testing.T) {
	u, err := getUserFromEmail("paul@paul.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if u.Email != "paul@paul.com" {
		t.FailNow()
	}
}

func TestSetToken(t *testing.T) {
	u, err := getUserFromEmail("paul@paul.com")
	if err != nil {
		panic(err)
	}

	setUserToken(u, generateRandomString(32))
}

func TestPassHash(t *testing.T) {
	salt := generateRandomString(saltLength)

	fmt.Printf("%s+%s\n", "paul", salt)
	fmt.Printf("%s\n", hash("paul", salt))
}

func TestPictureSet(t *testing.T) {
	u, _ := getUserFromEmail("paul@paul.com")
	picture := &Picture{
		ImagePath:      "/users/3/buzz.png",
		Mask:           "abcd",
		ExpirationTime: time.Now().Add(time.Minute * 20),
	}
	u.Pictures = append(u.Pictures, *picture)
	saveUser(u)
}

func TestGetPicture(t *testing.T) {
	u, _ := getUserFromEmail("paul@paul.com")
	mask := "mask3v2"

	pic, err := getPicture(u, mask, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(pic.ValidURL)
}

func TestDeletePicture(t *testing.T) {
	mask := "Z16WB5jVTAMmGrwUAVF7PlXA4ZaQIHAb"
	deletePicture(mask)
}

func TestDeletePictures(t *testing.T) {
	//u := &User{ID: 3}

	//deletePictures(u, []string{"2qPYFrubXq2lVEnbcrBKev1hwZb8jbZE", "AfcABMyb_73VDf-B9qeWrg61qaSOWvS6"})
}

func TestGetPicturesFromUser(t *testing.T) {
	testUser := &User{
		ID: 3,
	}
	var (
		limit = 30
		page  = 2
	)
	pictures, page, maxPages := getUsersPicturesAndRefreshURL(testUser, limit, page)
	fmt.Printf("Page [%d] out of [%d]\n", page, maxPages)
	for _, pic := range pictures {
		fmt.Println(pic.Mask)
	}
}

func TestGetPictures(t *testing.T) {
	testUser := &User{ID: 3}

	pictures, err := getPictures(testUser, []string{"hu8Lx3dbF5bqDRlDOcssmbXjzOoqF43p", "M4AXad_fjpynNSyS0oJN0Qp9QwPZVV0X"}, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(pictures)
}

func TestInsertAlbum(t *testing.T) {
	testUser := &User{ID: 3}

	p1, _ := getPicture(testUser, "UN5gdZlb2MCdrxgjlfVnTU6t-m3-Cbbr", false)
	p2, _ := getPicture(testUser, "-RPhn7n0J8PZWekTwoUdFldfV6cyF8G3", false)
	album := &Album{
		UserID:   testUser.ID,
		Title:    "Test Title",
		Mask:     generateRandomString(32),
		Pictures: []Picture{*p1, *p2},
	}
	err := saveAlbum(album)
	if err != nil {
		panic(err)
	}
}

func init() {
	verboseDatabase = true
}
