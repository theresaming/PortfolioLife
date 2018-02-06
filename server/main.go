package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type conf struct {
	SQLURL    string `json:"sqlURL"`
	TableName string `json:"tableName"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var config conf

func main() {
	fmt.Printf("Hello, 世界!\n")
	fmt.Printf("I'm going to connect to %s!\n", config.SQLURL)
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &config)
}
