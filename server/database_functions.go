// This abstracts common-use functions from the database
// All functions here will return client-safe messages.
// That is, nothing internal will be exposed in these messages.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	tokenLength = 32
	saltLength  = 10
)

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
	user.Token = ""

}

func registerUser(user *User) (token string, err error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	token = generateRandomString(tokenLength)
	u := &User{}
	db.Where("email = ?", user.Email).First(u)
	if len(u.Email) != 0 {
		return "", fmt.Errorf("this email is already registered")
	}
	user.Token = token
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
