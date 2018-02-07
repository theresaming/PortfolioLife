// This abstracts common-use functions from the database
// All functions here will return client-safe messages.
// That is, nothing internal will be exposed in these messages.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	tokenLength = 32
	saltLength  = 10
)

/*****************
*				 *
* User Functions *
*				 *
******************/

func getUserFromSession(sessionID string) (*User, error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var u User
	db.Where("token = ?", sessionID).First(&u)
	if len(u.Email) == 0 {
		return nil, fmt.Errorf("no user found for session")
	}
	return &u, nil
}

func getUserFromEmail(email string) (*User, error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var u User
	db.Where("email = ?", email).First(&u)
	if len(u.Email) == 0 {
		return nil, fmt.Errorf("no user found for given email %s", email)
	}
	return &u, nil
}

func emailExists(email string) bool {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var u User
	db.Where("email = ?", email).First(&u)
	return len(u.Email) != 0
}

func logoutUser(user *User) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Model(user).Where("email = ?", user.Email).Update("token", "")
	user.Token.String = ""
	user.Token.Valid = false

}

func registerUser(user *User) (token string, err error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// TODO: collision handling
	token = generateRandomString(tokenLength)
	u := &User{}
	db.Where("email = ?", user.Email).First(u)
	if len(u.Email) != 0 {
		return "", fmt.Errorf("this email is already registered")
	}
	user.Token.String = token
	user.Token.Valid = true
	db.Create(user)

	return
}

func setUserToken(user *User, token string) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Model(user).Update("token", token)
}

func saveUser(user *User) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Save(user)
}

/********************
*					*
* picture functions *
*					*
********************/

// This will refresh the validURL if ExpirationTime < now + 10 minutes
func getPicture(user *User, pictureMask string) (*Picture, error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goodTime := time.Now().Add(10 * time.Minute)

	var (
		picture Picture
		user2   User
	)

	//db.Where("mask", pictureMask).First(&picture)
	db.Where("id = ?", user.ID).Preload("Pictures", "mask = ?", pictureMask).First(&user2)

	if len(user2.Pictures) == 0 {
		return nil, fmt.Errorf("no picture found for your user session")
	}
	picture = user2.Pictures[0]
	if picture.ExpirationTime.Before(goodTime) {
		url, err := refreshURL(&picture)
		if err != nil {
			panic(err)
		}
		picture.ValidURL = url
		picture.ExpirationTime = time.Now().Add(urlExpirationDuration)
		db.Save(&picture)
	}

	return &picture, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) string {
	b, _ := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)[:s]
}
