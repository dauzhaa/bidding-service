package main

import (
	"log"
	"os"

	"github.com/students-api/bidding-service/internal/services/notification"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	kafkaAddr := getEnv("KAFKA_ADDR", "localhost:9092")
	brokers := []string{kafkaAddr}

	log.Printf("Connecting to Kafka at: %s", kafkaAddr)
	consumer, err := notification.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	consumer.StartConsume("bids.placed")
}