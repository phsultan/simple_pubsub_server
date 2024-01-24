# PubSub Server in Go

This Go program implements a simple Publish-Subscribe (PubSub) server that allows remote communication between publishers and subscribers over HTTP. Publishers can send messages to specific topics, and subscribers can subscribe to topics to receive messages.

## Features

- Publish messages to specific topics
- Subscribe to topics to receive messages
- Use of HTTP for communication
- JSON encoding for message payload

## Directory Structure

pubsub/
|-- main.go
|-- server/
| |-- pubsub_server.go
| |-- handlers/
| |-- publish_handler.go
| |-- subscribe_handler.go
|-- messages/
| |-- message.go


## Running the Server

1. Clone the repository:

```bash
git clone https://github.com/yourusername/pubsub-server-go.git
```

2. Navigate to the project directory:

```
cd pubsub-server-go
```

3. Run the server:

```
go run main.go
``` 

# Usage

## Publishing Messages

To publish a message to a specific topic, make an HTTP POST request to /publish/{topic}:

```
curl -X POST -H "content-type: application/json" -d '{"Content":"Your message content"}' http://localhost:8080/publish/{topic}
```

## Subscribing to Topics

To subscribe to a topic and receive messages, make an HTTP GET request to /subscribe/{topic}:

```
curl http://localhost:8080/subscribe/{topic}
```
