package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/salmander/github-notifications/common"
	"github.com/salmander/github-notifications/config"
	"github.com/streadway/amqp"
)

var c config.Config
var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue

type webhook struct {
	Header webhookHeader `json:"Header"`
	Body   string        `json:"Body"`
}

type webhookHeader struct {
	Event    string `json:"GitHubEvent"`
	Delivery string `json:"GitHubDelivery"`
}

func main() {
	c = config.ReadFromConfig("config.yaml")

	// Setup URL handler
	http.HandleFunc("/", webhookHandler)

	// Setup connection to the queue
	conn = setupQueue()
	defer conn.Close() // Close connection at the end

	// Setup channel to the queue
	q = setupChannelAndQueue()
	defer ch.Close() // Close the channel at the end

	// Start HTTP server
	log.Println("Listening on port", c.Http.Port)
	err := http.ListenAndServe(":"+c.Http.Port, nil)
	common.FailOnError(err, "Error listening on port "+c.Http.Port)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path[1:]
	var message string = endpoint
	log.Printf("Request received for endpoint: %s!", endpoint)

	// Check if the endpoint is 'payload'
	if endpoint == "payload" {
		log.Println("Endpoint == payload")
		response := webhook{}

		// Get the request body
		bodyBuffer, err := ioutil.ReadAll(r.Body)
		common.FailOnError(err, "Error reading request body")
		response.Body = fmt.Sprintf("%s", bodyBuffer)

		// Read the headers
		response.Header.Delivery = r.Header.Get("X-GitHub-Delivery")
		response.Header.Event = r.Header.Get("X-GitHub-Event")

		message = convertToJsonBody(response)
	}

	// Publish message to the queue
	if publishMessage(message) {
		fmt.Fprint(w, "[X] Message sent")
	} else {
		fmt.Fprint(w, "[Error] Something went wrong in sending the message")
	}

}

func convertToJsonBody(resp webhook) string {
	str, err := json.MarshalIndent(resp, "", "	")

	common.FailOnError(err, "Error converting struct to JSON response")

	return string(str)
}

// Configure and setup queue
func setupChannelAndQueue() amqp.Queue {
	var err error
	// Open channel
	ch, err = conn.Channel()
	common.FailOnError(err, "Failed to open a channel")

	q := common.GetQueue(ch, c)
	log.Println("Connected to the exchange")
	return q
}

// Setup RabbitMQ connection
func setupQueue() *amqp.Connection {
	var err error

	// Dial to connect
	conn, err = amqp.Dial(c.GetURL())
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
