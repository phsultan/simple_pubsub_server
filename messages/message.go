package simple_pubsub_server

// Message represents the data structure that will be sent between publishers and subscribers
type Message struct {
	Content string `json:"content"`
}

