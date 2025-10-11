package rmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type VoteNotification struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type VoteEvent struct {
	Type      string `json:"type"`
	PollID    string `json:"pollID"`
	OptionID  string `json:"optionID"`
	UserID    string `json:"userID"`
	PollOwner string `json:"pollOwner"`
	PollTitle string `json:"pollTitle"`
	Timestamp string `json:"timestamp"`
}

type RMQ struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRMQ(url string) *RMQ {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatal("RMQ connection error:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("RMQ channel error:", err)
	}

	_, err = ch.QueueDeclare("voteNotifications", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	return &RMQ{conn: conn, ch: ch}
}

func (r *RMQ) Publish(body []byte) error {
	return r.ch.Publish(
		"",
		"voteNotifications",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RMQ) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	return r.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}
