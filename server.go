// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	Message   string `json:"message"`
	Recipient int    `json:"recipient"`
	Sender    int    `json:"sender"`
}

var nodes = map[int]string{
	1: "http://localhost:8081",
	2: "http://localhost:8082",
	// Add other nodes as needed
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipientAddr, ok := nodes[msg.Recipient]
	if !ok {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", recipientAddr+"/receive", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "Message sent to node %d", msg.Recipient)
}

func main() {
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
