package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type conf struct {
	SQLURL string `json:"sqlURL"`
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
