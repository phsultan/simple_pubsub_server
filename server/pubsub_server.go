package simple_pubsub_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type PubSubServer struct {
	mu            sync.Mutex
	Subscribers   map[chan Message]map[string]struct{} // Map of channels with subscribed topics
	messageBuffer []Message
}

// PublishHandler handles the HTTP POST request for message publishing
func (ps *PubSubServer) PublishHandler(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.Path, "/")
	if len(tokens) < 3 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	topic := tokens[2]

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	fmt.Printf("[PublishHandler] topic : %s, sending to %d subscribers\n", topic, len(ps.Subscribers))
	fmt.Printf("[PublishHandler] msg : %s\n", msg)

	// Add the message to the buffer
	ps.messageBuffer = append(ps.messageBuffer, msg)

	// Broadcast the message to subscribers interested in the topic
	for ch, topics := range ps.Subscribers {
		if _, ok := topics[topic]; ok {
			go func(ch chan Message) {
				select {
				case ch <- msg:
				default:
					// If the channel is full, skip the message to prevent blocking
				}
			}(ch)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// SubscribeHandler handles the HTTP GET request for subscriber registration
func (ps *PubSubServer) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.Path, "/")
	if len(tokens) < 3 {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	topic := tokens[2]

	ps.mu.Lock()
	defer ps.mu.Unlock()

	fmt.Printf("[SubscribeHandler] subscribing to topic : %s\n", topic)
	fmt.Printf("[SubscribeHandler] len(ps.messageBuffer) : %d\n", len(ps.messageBuffer))

	// Create a new subscriber channel
	ch := make(chan Message)

	// Add the subscriber channel to the map with subscribed topics
	if _, exists := ps.Subscribers[ch]; !exists {
		ps.Subscribers[ch] = make(map[string]struct{})
	}
	ps.Subscribers[ch][topic] = struct{}{}

	// Send the buffered messages to the new subscriber
	for _, msg := range ps.messageBuffer {
		if _, ok := ps.Subscribers[ch][topic]; ok {
			go func(ch chan Message, msg Message) {
				select {
				case ch <- msg:
				default:
					// If the channel is full, skip the message to prevent blocking
				}
			}(ch, msg)
		}
	}

	// Return the channel as a response
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ch); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

