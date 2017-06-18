package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/salmander/github-notifications/common"
	"github.com/streadway/amqp"
)

var config common.Config
var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue

type webhook struct {
	Header webhookHeader
	Body   struct{}
}

type webhookHeader struct {
	Event     string
	Delivery  string
	Signature string
}

func main() {
	config = common.ReadFromConfig("config.yaml")

	// Setup URL handler
	http.HandleFunc("/", webhookHandler)

	// Setup connection to the queue
	conn = setupQueue()
	defer conn.Close() // Close connection at the end

	// Setup channel to the queue
	q = setupChannelAndQueue()
	defer ch.Close() // Close the channel at the end

	// Start HTTP server
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received: param 1: %s!", r.URL.Path[1:])

	// Get the request body

	// Publish message to the queue
	message := r.URL.Path[1:]
	if publishMessage(message) {
		fmt.Fprintf(w, "[X] Message sent: %v", message)
	} else {
		fmt.Fprint(w, "[Error] Something went wrong in sending the message")
	}

}

// Configure and setup queue
func setupChannelAndQueue() amqp.Queue {
	var err error
	// Open channel
	ch, err = conn.Channel()
	common.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		config.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")
	log.Println("Connected to the exchange")
	return q
}

// Setup RabbitMQ connection
func setupQueue() *amqp.Connection {
	var err error

	// Dial to connect
	conn, err = amqp.Dial(config.GetURL())
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	log.Print("Connected to the RabbitMQ server")
	return conn
}

// Publish message to the queue
func publishMessage(message string) bool {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	common.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", message)
	return true
}
