package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleListDevice(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(Devices)
	if err != nil {
		fmt.Fprintf(w, "Whoops : %s", err)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

func RunHttpServer() {
	fmt.Printf("starting the HTTP server\n")

	http.HandleFunc("/", handleListDevice)
	http.ListenAndServe(":8080", nil)

}
