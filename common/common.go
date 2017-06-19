package common

import (
	"log"

	"github.com/salmander/github-notifications/config"
	"github.com/streadway/amqp"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetQueue(ch *amqp.Channel, c config.Config) amqp.Queue {
	q, err := ch.QueueDeclare(
		c.Queue.Name,             // name
		c.Queue.Durable,          // durable
		c.Queue.DeleteWhenUnused, // delete when unused
		c.Queue.Exclusive,        // exclusive
		c.Queue.NoWait,           // no-wait
		nil,                      // arguments
	)

	FailOnError(err, "Failed to declare a queue")
	return q
}
