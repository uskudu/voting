package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbitmq/amqp091-go"
)

const queueName = "voteNotifications"

func main() {
	// new conn
	conn, err := amqp091.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// open chan
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	// subscribe to get messages from the queue
	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// wait for messages
	for {
		select {
		case message := <-messages:
			log.Printf("message: %s\n", message.Body)
		case <-sigchan:
			log.Println("interrupt detected!")
			os.Exit(0)
		}
	}
}
