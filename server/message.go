package simple_pubsub_server

type Message struct {
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

func NewMessage(msg string, topic string) *Message {
	// Returns the message object
	return &Message{
		Topic: topic,
		Body:  msg,
	}
}
