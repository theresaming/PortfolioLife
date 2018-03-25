// This contains mostly definitions useful for database functions. For more meat,
// checkout database_functions.go
package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	verboseDatabase = false
)

// A User is a client of our service
type User struct {
	CreatedAt time.Time
	ID        uint `gorm:"primary_key;AUTO_INCREMENT;size:11"` // default primary key but for clarity
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	Salt      string         `gorm:"type:varchar(10)"`
	Oauth     int            `gorm:"size:11"`
	Token     sql.NullString `gorm:"type:varchar(32);unique"`

	Pictures []Picture `gorm:"foreignkey:UserID"`
}

// A Picture is data about an uploaded photo
type Picture struct {
	CreatedAt time.Time
	UserID    uint   `gorm:"size:11;AUTO_INCREMENT:false"`
	ImagePath string `gorm:"type:varchar(512)" json:"path"`        // Path to image in S3 Bucket
	Mask      string `gorm:"unique;type:varchar(32);primary_key;"` // portfoliolife.com/picture/mask
	// The mask acts as a key of sorts

	ValidURL       string `gorm:"type:varchar(1024)"`
	ExpirationTime time.Time

	// No need for a backreference from picture->albums
	// Albums []Album `gorm:"many2many:picture_in_album;"`

	Tags []Tag `gorm:"ForeignKey:picture_mask;PRELOAD:true;"`

	PictureShare PictureShare
}

// An Album is a collection of a users photos
type Album struct {
	CreatedAt time.Time
	UserID    uint   `gorm:"size:11;"`
	Title     string `gorm:"type:varchar(256)"`
	Mask      string `gorm:"unique;type:varchar(512);primary_key;"` // portfoliolife.com/album/mask

	Pictures []Picture `gorm:"many2many:album_has_pictures;"`

	AlbumShare AlbumShare
}

// A Tag is metadata on a photo
type Tag struct {
	CreatedAt   time.Time
	PictureMask string `gorm:"primary_key;type:varchar(32);index"`
	Tag         string `gorm:"primary_key;type:varchar(256)"`
}

// A Registration is details on a user's pending registration
type Registration struct {
	CreatedAt time.Time
	UserID    uint   `gorm:"primary_key;size:11;"`
	EmailHash string `gorm:"type:varchar(32);"`
	PassHash  string `gorm:"type:varchar(32);"`
}

// TableName returns Registration's correct table name
func (Registration) TableName() string {
	return "registration"
}

// An AlbumShare is metadata on a shared album
type AlbumShare struct {
	CreatedAt    time.Time
	AlbumID      uint `gorm:"primary_key;size:11;"`
	ShareSetting int
	Hash         string `gorm:"type:varchar(256);UNIQUE"`
}

// A PictureShare is metadata on a shared picture
type PictureShare struct {
	CreatedAt    time.Time
	PictureID    uint `gorm:"primary_key;size:11"`
	ShareSetting int
	Hash         string `gorm:"type:varchar(256);UNIQUE"`
}

func wrapDb(fn func(db *gorm.DB)) {
	db, err := openConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fn(db)
}

func openConnection() (*gorm.DB, error) {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Username, config.Password, config.SQLURL, config.TableName))
	if err == nil {
		db.LogMode(verboseDatabase)
	}
	return db, err
}
