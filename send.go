package main

import (
	"log"

	"github.com/salmander/github-notifications/common"
	"github.com/streadway/amqp"
)

func main() {
	config := common.ReadFromConfig("config.yaml")
	conn, err := amqp.Dial(config.GetURL())
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	common.FailOnError(err, "Failed to declare a queue")

	body := "hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	common.FailOnError(err, "Failed to publish a message")
}
