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
	Identifier    string
	Model         string
	Polling       int
}

var Devices = make(map[string]Device)

var Models = make(map[string]map[string]string)

func LoadModels(path string) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	e = json.Unmarshal(file, &Models)

	if e != nil {
		fmt.Printf("Models Json error: %v\n", e)
		os.Exit(1)
	}
}

func LoadDevices(path string) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	e = json.Unmarshal(file, &Devices)

	if e != nil {
		fmt.Printf("Devices Json error: %v\n", e)
		os.Exit(1)
	}

	// check if the linked model exists

	for k, v := range Devices {
		if Models[v.Model] == nil {
			fmt.Printf("model '%s' referenced by device '%s' not found\n", v.Model, k)
			os.Exit(1)
		}
	}
}
