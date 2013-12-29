package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alouca/gosnmp"
	"net/http"
	"os"
	"sync"
	"time"
)

// latch for
var w sync.WaitGroup

func pollDevice(name string, device *Device, platform string) {
	fmt.Printf("community: %s\n", device.SnmpCommunity)

	s, err := gosnmp.NewGoSNMP(device.Host, device.SnmpCommunity, gosnmp.Version1, 5)
	if err != nil {
		fmt.Printf("err : %s\n", err)
		return
	}

	for {
		var values = make(map[string]interface{})

		fmt.Printf("Polling device with model %s\n", device.Model)

		// get the model

		model := Models[device.Model]

		for k, v := range model {
			fmt.Printf("polling '%s'\n", k)
			resp, err := s.Get(k)
			if err != nil {
				fmt.Printf("SNMP err : %s\n", err)
			} else {
				pdu := resp.Variables[0]
				var t = pdu.Type
				switch t {
				case gosnmp.Integer:
					fmt.Println("Integer!")
					values[v] = pdu.Value.(int)
				case gosnmp.OctetString:
					fmt.Println("OctetString!")
					values[v] = string(pdu.Value.([]uint8))
				default:
					fmt.Printf("unknown type: %d (%s)\n", t, v)
				}
			}
		}

		fmt.Println("push...")
		dataPush(values, platform, device.Identifier, device.Password)
		time.Sleep(time.Duration(device.Polling) * time.Second)
	}
	w.Done()
}

func main() {
	fmt.Printf("AirVantage SNMP gateway\n")

	if len(os.Args) != 2 {
		fmt.Printf("Usage : avsnmp [platform]\n\n platform: edge,na,eu,etc\n")
		os.Exit(1)
	}
	// load configuration from json

	LoadModels("conf/model.json")
	LoadDevices("conf/devices.json")

	w.Add(len(Devices))
	// start the web server
	go RunHttpServer()

	for k, v := range Devices {
		go pollDevice(k, &v, os.Args[1])
	}

	w.Wait()
}

// push a data to airvantage using HTTP REST for devices
func dataPush(values map[string]interface{}, platform string, username string, password string) {
	var container [1]map[string]interface{}
	container[0] = make(map[string]interface{})

	for k, v := range values {
		var dataKey [1]map[string]interface{}

		dataKey[0] = make(map[string]interface{})

		dataKey[0]["timestamp"] = nil
		dataKey[0]["value"] = v

		container[0][k] = dataKey
	}

	b, err := json.Marshal(container)

	if err == nil {
		fmt.Printf("JSon :%s\n", b)
		//  build an HTTP request
		reader := bytes.NewBuffer(b)
		r, _ := http.NewRequest("POST", "http://"+platform+".airvantage.net/device/messages", reader)
		r.SetBasicAuth(username, password)

		client := &http.Client{}
		resp, err := client.Do(r)

		if err == nil {
			fmt.Printf("HTTP status %d\n", resp.StatusCode)
			if resp.StatusCode == 200 {
				// happy
				fmt.Printf("HTTP server Ok\n")
			} else {
				// unhappy
				fmt.Printf("HTTP server error : %s\n", resp.Status)
			}
		} else {
			fmt.Printf("HTTP POST error : %s\n", err.Error())
		}
	} else {
		fmt.Printf("fail : %s\n", err)
	}
}
