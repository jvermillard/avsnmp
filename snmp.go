package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alouca/gosnmp"
	"net/http"
)

func main() {
	fmt.Printf("AirVantage SNMP gateway\n")

	s, err := gosnmp.NewGoSNMP("127.0.0.1", "public", gosnmp.Version1, 5)

	if err != nil {
		fmt.Printf("err : %s\n", err)
		return
	}
	resp, err := s.Get(".1.3.6.1.2.1.1.1.0")

	if err == nil {
		fmt.Printf("some response: %s\n", resp)

		fmt.Printf("data : %s\n", resp.Variables[0].Name)
		fmt.Printf("value : %s\n", resp.Variables[0].Type)
		fmt.Printf("as string : %s\n", string(resp.Variables[0].Value.([]uint8)))
		dataPush("datapath", string(resp.Variables[0].Value.([]uint8)), "na", "mySN", "myP4ssw0rd")

	} else {
		fmt.Printf("err : %s\n", err)
	}
}

// push a data to airvantage using HTTP REST for devices
func dataPush(path string, value string, platform string, username string, password string) {
	var container [1]map[string]interface{}
	var dataKey [1]map[string]interface{}

	container[0] = make(map[string]interface{})
	dataKey[0] = make(map[string]interface{})
	dataKey[0]["timestamp"] = nil

	dataKey[0]["value"] = value

	container[0][path] = dataKey
	b, err := json.Marshal(container)

	if err == nil {
		fmt.Printf("JSon :%s\n", b)
		//  build an HTTP request
		reader := bytes.NewBuffer(b)
		r, _ := http.NewRequest("POST", "http://"+platform+".airvantage.net/device/messages", reader)
		r.SetBasicAuth("myuser", "passw0rdz")

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
