// This abstracts common-use functions from the database
// All functions here will return client-safe messages.
// That is, nothing internal will be exposed in these messages.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
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
* Picture functions *
*					*
********************/

// This will refresh the validURL if ExpirationTime < now + 10 minutes
// when refresh is set to true
func getPicture(user *User, pictureMask string, refresh bool) (*Picture, error) {
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
	if picture.ExpirationTime.Before(goodTime) && refresh {
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

func getPictures(user *User, pictureMasks []string, refresh bool) ([]Picture, error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	goodTime := time.Now().Add(10 * time.Minute)
	var (
		user2 User
	)
	db.Where("id = ?", user.ID).Preload("Pictures", "mask IN (?)", pictureMasks).First(&user2)

	if len(user2.Pictures) == 0 {
		return nil, fmt.Errorf("no picture found for your user session")
	}
	if refresh {
		for _, picture := range user2.Pictures {
			if picture.ExpirationTime.Before(goodTime) && refresh {
				url, err := refreshURL(&picture)
				if err != nil {
					panic(err)
				}
				picture.ValidURL = url
				picture.ExpirationTime = time.Now().Add(urlExpirationDuration)
				db.Save(&picture)
			}
		}
	}

	return user2.Pictures, nil
}

func deletePictures(user *User, pictures []Picture) error {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	pictureMasks := make([]string, len(pictures))
	i := 0
	for _, p := range pictures {
		pictureMasks[i] = p.Mask
		i++
	}
	db.Exec("DELETE FROM pictures WHERE mask IN (?) AND user_id = ?", pictureMasks, user.ID)
	return nil
}

func deletePicture(pictureMask string) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec("DELETE FROM pictures WHERE mask = ?", pictureMask)
}

// This will attempt to paginate on lines of multiples of limit
// If page * limit > length of result, then it will return the modulus
// of the final page and limit
// e.g. If you ask for page 2 of 50 results on limit 30, it will return the final 20 results
func getUsersPicturesAndRefreshURL(user *User, limit int, page int) (pictures []Picture, currentPage int, maxPages int) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	pictures = make([]Picture, 0)
	var (
		rowCount int
		offset   int
	)
	if limit <= 0 {
		limit = 1
	}
	if page <= 0 {
		page = 1
	}

	db.Model(&Picture{}).Where("user_id = ?", user.ID).Count(&rowCount)
	if (page-1)*limit >= rowCount {
		page = (rowCount / limit)
		if page == 0 {
			page++
		}
	}

	offset = (page - 1) * limit
	maxPages = (rowCount / limit)
	if limit*page > rowCount || (maxPages == 0 && rowCount > 0) {
		maxPages++
	}

	db.Limit(limit).Offset(offset).Where("user_id = ?", user.ID).Find(&pictures)
	goodTime := time.Now().Add(10 * time.Minute)
	for i, picture := range pictures {
		if picture.ExpirationTime.Before(goodTime) {
			url, err := refreshURL(&picture)
			if err != nil {
				panic(err)
			}
			picture.ValidURL = url
			picture.ExpirationTime = time.Now().Add(urlExpirationDuration)
			db.Save(&picture)
		}
		if tags, err := getTags(&picture); err == nil {
			pictures[i].Tags = tags
		}
	}

	return pictures, page, maxPages
}

/**********
*		  *
* Tagging *
*		  *
**********/

func createTags(picture *Picture, tags []Tag) error {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	picture.Tags = append(picture.Tags, tags...)
	return db.Save(picture).Error
}

func deleteTags(picture *Picture, tags []string) error {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Exec("DELETE FROM tags WHERE tag IN (?) AND picture_mask = ?", tags, picture.Mask).Error

	p, _ := getPicture(&User{ID: picture.UserID}, picture.Mask, false)
	*picture = *p

	return err
}

func getTags(picture *Picture) (tags []Tag, err error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	tags = make([]Tag, 0)
	err = db.Where("picture_mask = ?", picture.Mask).Find(&tags).Error
	return
}

func searchWithTag(u *User, term string, front, back, refresh bool) (pictures []Picture, err error) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	pictureMasks := make([]string, 0)
	var fuzz = term
	if front && back {
		fuzz = fmt.Sprintf("%%%s%%", term)
	} else if front {
		fuzz = fmt.Sprintf("%%%s", term)
	} else if back {
		fuzz = fmt.Sprintf("%s%%", term)
	}

	rows, err := db.Raw(`SELECT t.picture_mask FROM tags t LEFT JOIN pictures p ON t.picture_mask = p.mask LEFT JOIN users u ON p.user_id = u.id WHERE u.id = ? AND t.tag LIKE ?`, u.ID, fuzz).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mask string
		rows.Scan(&mask)
		pictureMasks = append(pictureMasks, mask)
	}

	return getPictures(u, pictureMasks, refresh)
}

/*********
*		 *
* Albums *
*		 *
*********/

// Precondition: album only has valid pictures for this user
func saveAlbum(album *Album) error {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	values := make([]string, len(album.Pictures))
	for i, pic := range album.Pictures {
		values[i] = fmt.Sprintf("(`%s`, `%s`)", album.Mask, pic.Mask)
	}
	db.Save(album)
	sql := fmt.Sprintf("INSERT INTO `album_has_pictures` (`album_mask`,`picture_mask`) VALUES %s", strings.Join(values, ", "))
	return db.Exec(sql).Error
}

/****************
*				*
* Miscellaneous *
*				*
****************/

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
