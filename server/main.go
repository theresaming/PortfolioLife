package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type conf struct {
	SQLURL    string `json:"sqlURL"`
	TableName string `json:"tableName"`
	Username  string `json:"username"`
	Password  string `json:"password"`

	S3Endpoint  string `json:"s3Endpoint"`
	S3SpaceName string `json:"s3SpaceName"`
	S3Location  string `json:"s3Location"`

	Port int `json:"port"`

	S3Key    string `json:"-"`
	S3Secret string `json:"-"`
}

var (
	l      *log.Logger
	config conf
)

func main() {
	override := flag.Bool("force", false, "force start even without image upload keys")
	flag.Parse()

	if !*override && (len(config.S3Key) == 0 || len(config.S3Secret) == 0) {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Hello, 世界!\n")

	initServer()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &config)

	config.S3Key = os.Getenv("S3_KEY")
	config.S3Secret = os.Getenv("S3_SECRET")

	l = log.New(os.Stderr, "[main]: ", log.LstdFlags|log.Lshortfile)
}
