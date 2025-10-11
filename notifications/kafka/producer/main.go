package main

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func setupKafkaProducer() *kafka.Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "voteNotifications",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Transport:    kafka.DefaultTransport,
	}
	return writer
}
