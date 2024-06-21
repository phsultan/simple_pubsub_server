package simple_pubsub_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PubSubServer struct {
	broker *Broker
}

func NewPubSubServer() *PubSubServer {
	return &PubSubServer{
		broker: NewBroker(),
	}
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

	msg.Topic = topic

	fmt.Printf("[PublishHandler] topic : %s, sending to %d subscribers\n", topic, len(ps.broker.subscribers))
	fmt.Printf("[PublishHandler] msg : %s\n", msg)
	fmt.Printf("[PublishHandler] msg.Topic : %s\n", msg.Topic)
	fmt.Printf("[PublishHandler] msg.Body : %s\n", msg.Body)

	ps.broker.Publish(topic, msg.Body)

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

	fmt.Printf("[SubscribeHandler] subscribing to topic : %s\n", topic)

	subscriber := ps.broker.AddSubscriber()
	ps.broker.Subscribe(subscriber, topic)

	go func() {
		select {
		case <-r.Context().Done():
			fmt.Printf("Subscriber %s left\n", subscriber.id)
			ps.broker.RemoveSubscriber(subscriber)
		}
	}()

	subscriber.Listen(w)

	// What do we reply to the HTTP client?
}
