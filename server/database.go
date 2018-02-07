// This contains mostly definitions useful for database functions. For more meat,
// checkout database_functions.go
package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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

	Pictures []Picture
}

// A Picture is data about an uploaded photo
type Picture struct {
	CreatedAt time.Time
	UserID    uint   `gorm:"primary_key;size:11;index"`
	ImagePath string `gorm:"type:varchar(512)" json:"path"`       // Path to image in S3 Bucket
	Mask      string `gorm:"primary_key;unique;type:varchar(32)"` // portfoliolife.com/picture/mask
	// The mask acts as a key of sorts

	ValidURL       string `gorm:"type:varchar(1024)"`
	ExpirationTime time.Time

	// TODO: more metadata here

	Albums []Album `gorm:"many2many:picture_in_album;AssociationForeignKey:albumID;ForeignKey:pictureID;"`

	Tags []Tag

	PictureShare PictureShare
}

// An Album is a collection of a users photos
type Album struct {
	CreatedAt time.Time
	AlbumID   uint   `gorm:"primary_key;AUTO_INCREMENT;size:11"`
	UserID    uint   `gorm:"size:11;"`
	Title     string `gorm:"type:varchar(256)"`
	Mask      string `gorm:"unique;type:varchar(512)"` // portfoliolife.com/album/mask

	Pictures []Picture

	AlbumShare AlbumShare
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

// A Tag is metadata on a photo
type Tag struct {
	CreatedAt time.Time
	TagID     uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Picture   int    `gorm:"size:11;"`
	Tag       string `gorm:"type:varchar(256)"`
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
	return gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Username, config.Password, config.SQLURL, config.TableName))
}
