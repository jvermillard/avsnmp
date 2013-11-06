package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alouca/gosnmp"
	"net/http"
	"os"
)

const (
	_           = iota
	ECIO_OID    = ".1.3.6.1.4.1.20992.1.2.2.1.6.0"
	RSSI_OID    = ".1.3.6.1.4.1.20992.1.2.2.1.4.0"
	SERVICE_OID = ".1.3.6.1.4.1.20992.1.2.2.1.10.0"
	SN_OID      = ".1.3.6.1.4.1.20992.1.2.2.1.7.0"
	RM_FW_OID   = ".1.3.6.1.4.1.20992.1.2.2.1.8.0"
	FW_OID      = ".1.3.6.1.2.1.1.1.0"
)

func main() {
	fmt.Printf("AirVantage SNMP gateway\n")

	if len(os.Args) != 5 {
		fmt.Printf("Usage : avsnmp [device host] [community] [platform] [password]\n\n[device host] is the IP address or the hostname of the SNMP device\n[community] the SNMP community (ex: public)\n[platform] can be na,eu,edge,etc..\n[password] is the password associated with the device\n")
		os.Exit(1)
	}

	s, err := gosnmp.NewGoSNMP(os.Args[1], os.Args[2], gosnmp.Version1, 5)

	if err != nil {
		fmt.Printf("err : %s\n", err)
		return
	}

	var values = make(map[string]interface{})

	resp, err := s.Get(SN_OID)

	if err == nil {
		fmt.Printf("serial number: %s\n", string(resp.Variables[0].Value.([]uint8)))

		values["SERIAL_NUMBER"] = string(resp.Variables[0].Value.([]uint8))

	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}

	resp, err = s.Get(SERVICE_OID)
	if err == nil {
		fmt.Printf("service type: %s\n", string(resp.Variables[0].Value.([]uint8)))

		values["_NETWORK_SERVICE_TYPE"] = string(resp.Variables[0].Value.([]uint8))

	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}

	resp, err = s.Get(RSSI_OID)

	if err == nil {
		fmt.Printf("RSSI: %d\n", resp.Variables[0].Value.(int))
		values["_RSSI"] = resp.Variables[0].Value.(int)

	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}

	resp, err = s.Get(ECIO_OID)

	if err == nil {
		fmt.Printf("ECIO: %d\n", resp.Variables[0].Value.(int))
		values["_ECIO"] = resp.Variables[0].Value.(int)
	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}

	resp, err = s.Get(FW_OID)

	if err == nil {
		fmt.Printf("FW: %s\n", string(resp.Variables[0].Value.([]uint8)))
		values["FW"] = string(resp.Variables[0].Value.([]uint8))

	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}

	resp, err = s.Get(RM_FW_OID)

	if err == nil {
		//		fmt.Printf("response: %s\n", resp)
		fmt.Printf("Radio Module FW: %s\n", string(resp.Variables[0].Value.([]uint8)))
		values["RMFW"] = string(resp.Variables[0].Value.([]uint8))

	} else {
		fmt.Printf("SNMP err : %s\n", err)
	}
	if values["SERIAL_NUMBER"] != nil {
		dataPush(values, os.Args[3], values["SERIAL_NUMBER"].(string), os.Args[4])
	}
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
