package common

import (
	"fmt"
	"log"
)

const USERNAME = "guest1"
const PASSWORD = "guest1"
const PORT = "5672"
const QUEUE_NAME = "test-queue"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetURL() string {
	return fmt.Sprintf("amqp://%v:%v@localhost:%v/", USERNAME, PASSWORD, PORT)
}
