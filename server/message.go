package simple_pubsub_server

type Message struct {
	Topic string `json:"topic"`
	Body  string `json:"body"`
}
