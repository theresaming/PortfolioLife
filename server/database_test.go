package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConf testConfig

func TestDropCreateStandup(t *testing.T) {
	TestDropAllTables(t)
	TestMigration(t)
	TestStandup(t)
}

func TestDatabaseConnection(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.Close()
}

func TestMigration(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Picture{})
	db.AutoMigrate(&Album{})
	db.AutoMigrate(&Registration{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&PictureShare{})
	db.AutoMigrate(&AlbumShare{})
}

func TestTableExistence(t *testing.T) {
	assert := assert.New(t)
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()

	assert.True(db.HasTable(&User{}), "no user table")
	assert.True(db.HasTable(&Picture{}), "no picture table")
	assert.True(db.HasTable(&Album{}), "no album table")
	assert.True(db.HasTable(&Registration{}), "no registration table")
	assert.True(db.HasTable(&Tag{}), "no tag table")
	assert.True(db.HasTable(&PictureShare{}), "no picture_share table")
	assert.True(db.HasTable(&AlbumShare{}), "no album_share table")
}

func TestUserHasPictures(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()

	u := User{}
	db.First(&u).Related(&u.Pictures)
	fmt.Println(u.Pictures)
}

func TestStandup(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()

	for _, u := range testConf.Users {
		db.Create(&u)
	}
	for _, p := range testConf.Pictures {
		db.Create(&p)
	}
}

func TestDropAllTables(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()
	db.DropTable(&User{}, &Picture{}, &Album{}, &Registration{}, &Tag{}, &PictureShare{}, &AlbumShare{})
	db.Exec("DROP TABLE IF EXISTS album_shares;").
		Exec("DROP TABLE IF EXISTS picture_in_album;").
		Exec("DROP TABLE IF EXISTS picture_shares;")
}

func TestLol(t *testing.T) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(false)
	defer db.Close()

	pics, _ := getPictures(&User{ID: 3}, []string{"CNO9tTq4I7H0S-fLBF6wmlE6MeQv6VxP", "66mB0HLm-QuNH7o2xip_QVx04F5_OvLA"}, false)

	// u := &User{
	// 	ID: 3,
	// }
	a := &Album{
		Mask:   "album_mask",
		UserID: 3,
		Title:  "test",
		// Pictures: picturePtrs,
	}
	db.LogMode(true)
	if err := db.Create(a).Error; err != nil {
		panic(err)
	}
	if err := db.Model(a).Association("Pictures").Append(pics).Error; err != nil {
		panic(err)
	}
	if err := db.Save(a).Error; err != nil {
		panic(err)
	}
	// if err := db.Save(a).Error; err != nil {
	// 	panic(err)
	// }
	// db.Model(&Album{}).Related(&Picture{}, "Pictures")
}

type testConfig struct {
	Users    []User    `json:"users"`
	Pictures []Picture `json:"pictures"`
}

func init() {
	data, err := ioutil.ReadFile("db_standup.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &testConf)
	if err != nil {
		panic(err)
	}
}
