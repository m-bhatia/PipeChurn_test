// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"
)

type Message struct {
	Message   string `json:"message"`
	Recipient int    `json:"recipient"`
	Sender    int    `json:"sender"`
	Number    int64  `json:"number"`
}

var nodes = map[int]string{
	1: "http://sp24-cs525-0902.cs.illinois.edu:8081",
	2: "http://sp24-cs525-0903.cs.illinois.edu:8082",
	// Add other nodes as needed
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
    var msg Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        log.Printf("Error decoding message: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    recipientAddr, ok := nodes[msg.Recipient]
    if !ok {
        errMsg := fmt.Sprintf("Recipient %d not found", msg.Recipient)
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusNotFound)
        return
    }

    // Correctly marshal the msg object into JSON
    jsonValue, err := json.Marshal(msg)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    client := &http.Client{}
    req, err := http.NewRequest("POST", recipientAddr+"/receive", bytes.NewBuffer(jsonValue))
    if err != nil {
        log.Printf("Error creating request to recipient: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending request to recipient %d at %s: %v", msg.Recipient, recipientAddr, err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    log.Printf("Message sent to node %d successfully.", msg.Recipient)
    fmt.Fprintf(w, "Message sent to node %d", msg.Recipient)
}



func main() {
	http.HandleFunc("/send", sendMessage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
