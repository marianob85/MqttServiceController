package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Server   string `json:"Server"`
	Port     int    `json:"Port"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	ClientID string `json:"ClientID"`
	Topic    string `json:"Topic"`
}

func readConfiguration(filePath string) Config {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)
	return config
}
