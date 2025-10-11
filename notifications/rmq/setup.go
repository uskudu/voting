package rmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func SetupRMQ() *amqp091.Channel {
	conn, err := amqp091.Dial("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare("voteNotifications", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ch
}
