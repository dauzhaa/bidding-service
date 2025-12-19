package main

import (
	"log"

	"github.com/students-api/bidding-service/internal/services/notification"
)

func main() {
	brokers := []string{"localhost:9092"}

	consumer, err := notification.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	consumer.StartConsume("bids.placed")
}