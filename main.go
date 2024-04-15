package main

import (
	"fmt"
	"net/http"

	."github.com/phsultan/simple_pubsub_server/server" // PubSubServer struct as type
)

func main() {
	pubSubServer := NewPubSubServer()

	http.HandleFunc("/publish/", pubSubServer.PublishHandler)
	http.HandleFunc("/subscribe/", pubSubServer.SubscribeHandler)

	fmt.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
