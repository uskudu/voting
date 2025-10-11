package main

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"voting/notifications/email"
	"voting/notifications/rmq"
)

const queueName = "voteNotifications"

func main() {
	conn, err := amqp091.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Consumer started...")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-msgs:
			var event rmq.VoteEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("JSON decode error:", err)
				continue
			}
			log.Printf("Received vote event for poll '%s' by user '%s'\n", event.PollTitle, event.UserID)

			// Формируем уведомление
			n := rmq.VoteNotification{
				To:      event.PollOwner,
				Message: "new vote on your poll '" + event.PollTitle + "'.",
			}

			// Отправляем асинхронно
			go func() {
				if err := email.SendMail(n); err != nil {
					log.Printf("Failed to send email: %v\n", err)
				} else {
					log.Printf("Email sent to %s\n", n.To)
				}
			}()

		case <-sigchan:
			log.Println("interrupt detected - shutting down.")
			return
		}
	}
}
