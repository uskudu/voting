package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Notification struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "voteNotifications",
		GroupID: "vote-consumers",
	})
	defer r.Close()

	log.Println("Kafka consumer started...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var notif Notification
		if err := json.Unmarshal(m.Value, &notif); err != nil {
			log.Println("Invalid message:", err)
			continue
		}
		log.Printf("New vote on your poll '%s'. to user %s", notif.Message, notif.To)
	}
}
