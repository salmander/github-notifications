package main

import (
	"log"

	"github.com/salmander/github-notifications/common"
	"github.com/salmander/github-notifications/config"
	"github.com/streadway/amqp"
)

func main() {
	c := config.ReadFromConfig("config.yaml")
	conn, err := amqp.Dial(c.GetURL())
	common.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	common.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q := common.GetQueue(ch, c)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	common.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
