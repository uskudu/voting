package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"voting/notifications/email"
	"voting/notifications/rmq"
)

func main() {
	rmqClient := rmq.NewRMQ("amqp://localhost:5672")
	defer rmqClient.Close()

	msgs, err := rmqClient.Consume("voteNotifications")
	if err != nil {
		log.Fatal(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Consumer started...")

	for {
		select {
		case message := <-msgs:
			var notification rmq.VoteNotification
			if err := json.Unmarshal(message.Body, &notification); err != nil {
				log.Println("Invalid message:", err)
				continue
			}

			log.Printf("Received notification for %s: %s\n", notification.To, notification.Message)

			// отправка email асинхронно
			go func(n rmq.VoteNotification) {
				if err := email.SendMail(n); err != nil {
					log.Println("Failed to send email:", err)
				} else {
					log.Println("Email sent to", n.To)
				}
			}(notification)

		case <-sigchan:
			log.Println("Interrupt detected! Exiting...")
			os.Exit(0)
		}
	}
}
