package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jinzhu/gorm"
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
	wrapDbTest(func(db *gorm.DB) {
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Picture{})
		db.AutoMigrate(&Album{})
		db.AutoMigrate(&Registration{})
		db.AutoMigrate(&Tag{})
		db.AutoMigrate(&PictureShare{})
		db.AutoMigrate(&AlbumShare{})
	})
}

func TestTableExistence(t *testing.T) {
	assert := assert.New(t)
	wrapDbTest(func(db *gorm.DB) {
		assert.True(db.HasTable(&User{}), "no user table")
		assert.True(db.HasTable(&Picture{}), "no picture table")
		assert.True(db.HasTable(&Album{}), "no album table")
		assert.True(db.HasTable(&Registration{}), "no registration table")
		assert.True(db.HasTable(&Tag{}), "no tag table")
		assert.True(db.HasTable(&PictureShare{}), "no picture_share table")
		assert.True(db.HasTable(&AlbumShare{}), "no album_share table")
	})
}

func TestUserHasPictures(t *testing.T) {
	wrapDbTest(func(db *gorm.DB) {
		u := User{}
		db.First(&u).Related(&u.Pictures)
		fmt.Println(u.Pictures)
	})
}

func TestStandup(t *testing.T) {
	wrapDbTest(func(db *gorm.DB) {
		for _, u := range testConf.Users {
			db.Create(&u)
		}
		for _, p := range testConf.Pictures {
			db.Create(&p)
		}
	})
}

func TestDropAllTables(t *testing.T) {
	wrapDbTest(func(db *gorm.DB) {
		db.DropTable(&User{}, &Picture{}, &Album{}, &Registration{}, &Tag{}, &PictureShare{}, &AlbumShare{})
		db.Exec("DROP TABLE IF EXISTS album_shares;").
			Exec("DROP TABLE IF EXISTS picture_in_album;").
			Exec("DROP TABLE IF EXISTS picture_shares;")
	})
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

func wrapDbTest(test func(db *gorm.DB)) {
	db, err := openTestConnection()
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	defer db.Close()
	test(db)
}

func openTestConnection() (*gorm.DB, error) {
	return gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Username, config.Password, config.SQLURL, "jd_test"))
}
