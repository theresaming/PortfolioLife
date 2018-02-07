package main

import (
	"fmt"
	"testing"
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

func TestRegistration(t *testing.T) {
	u := testConf.Users[0]
	u.Email = "paul7@paul.com"

	token, err := registerUser(&u)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}

func TestPassHash(t *testing.T) {
	salt := generateRandomString(saltLength)

	fmt.Printf("%s+%s\n", "paul", salt)
	fmt.Printf("%s\n", hash("paul", salt))
}

func TestPictureSet(t *testing.T) {
	u, _ := getUserFromEmail("paul@paul.com")
	picture := &Picture{
		ImagePath: "/users/3/buzz.png",
		Mask:      "a",
	}
	u.Pictures = append(u.Pictures, *picture)
	saveUser(u)
}

func TestGetPicture(t *testing.T) {
	u, _ := getUserFromEmail("paul@paul.com")
	mask := "mask3v2"

	pic, err := getPicture(u, mask)
	if err != nil {
		panic(err)
	}
	fmt.Println(pic.ValidURL)
}
