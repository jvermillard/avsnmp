package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Device struct {
	Host          string
	Port          int
	SnmpCommunity string
	Password      string
}

var Devices = make(map[string]Device)

func LoadFromJson(path string) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	e = json.Unmarshal(file, &Devices)

	if e != nil {
		fmt.Printf("Json error: %v\n", e)
		os.Exit(1)
	}
}
