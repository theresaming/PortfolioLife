package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type conf struct {
	SQLURL    string `json:"sqlURL"`
	TableName string `json:"tableName"`
	Username  string `json:"username"`
	Password  string `json:"password"`

	Port int `json:"port"`
}

var config conf

func main() {
	fmt.Printf("Hello, 世界!\n")
	fmt.Printf("I'm going to connect to %s!\n", config.SQLURL)

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
}
