// node.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Message struct {
	Message   string `json:"message"`
	Recipient int    `json:"recipient"`
	Sender    int    `json:"sender"`
	Number    int64  `json:"number"` // Number from which to start counting
}

var (
	nodeID       = flag.Int("id", 1, "Node ID")
	currentCount int64
)

func receiveMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the current count based on the received message
	atomic.StoreInt64(&currentCount, msg.Number)

	fmt.Printf("Node %d received a message to start counting from %d\n", *nodeID, msg.Number)
}

func count() {
	for {
		time.Sleep(1 * time.Second) // Slow down the counting for demonstration
		val := atomic.AddInt64(&currentCount, 1)
		fmt.Printf("Node %d counting: %d\n", *nodeID, val)
	}
}

func main() {
	flag.Parse()
	go count()

	http.HandleFunc("/receive", receiveMessage)
	port := 8080 + *nodeID
	log.Printf("Node %d listening on port %d\n", *nodeID, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
