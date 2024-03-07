// client.go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"log"
)

func main() {
    // Create the message
    msg := map[string]interface{}{
        "message":   "start counting",
        "recipient": 2,
        "sender":    1,
        "number":    13,
    }

    // Marshal the message into a JSON body
    requestBody, err := json.Marshal(msg)
    if err != nil {
        fmt.Printf("Could not marshal request body: %v\n", err)
        return
    }

    // Send the request
    response, err := http.Post("http://sp24-cs525-0901:8080/send", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        fmt.Printf("Request failed: %v\n", err)
        return
    }
    defer response.Body.Close()

	// Inside client.go, after receiving the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("Couldn't read response body", err)
	}
	fmt.Printf("Response from server: %s\n", body)


    fmt.Println("Message sent successfully")
}
