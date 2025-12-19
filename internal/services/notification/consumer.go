package notification

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type BidEvent struct {
	EventID      string `json:"event_id"`
	AuctionID    int64  `json:"auction_id"`
	UserID       int64  `json:"user_id"`
	Amount       int64  `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) StartConsume(topic string) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer for partition 0: %v", err)
	}
	defer partitionConsumer.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Notification Service listening on topic: %s", topic)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var event BidEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			log.Println("---------------------------------------------------")
			log.Printf("НОВОЕ УВЕДОМЛЕНИЕ!")
			log.Printf("Кому: Пользователь ID %d", event.UserID)
			log.Printf("Текст: Поздравляем! Ваша ставка %d у.е. на аукцион #%d принята.", event.Amount, event.AuctionID)
			log.Printf("Тех. инфо: Offset %d, EventID %s", msg.Offset, event.EventID)
			log.Println("---------------------------------------------------")

		case <-sigchan:
			log.Println("Shutting down consumer...")
			return
		}
	}
}