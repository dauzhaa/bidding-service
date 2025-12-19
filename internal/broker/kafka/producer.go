package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
}


type BidEvent struct {
	EventID      string    `json:"event_id"`
	AuctionID    int64     `json:"auction_id"`
	UserID       int64     `json:"user_id"`
	Amount       int64     `json:"amount"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{syncProducer: producer}, nil
}

func (p *Producer) SendBidPlaced(event BidEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "bids.placed",
		Value: sarama.ByteEncoder(bytes),
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", event.AuctionID)),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}