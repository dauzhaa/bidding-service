package bidding_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/students-api/bidding-service/internal/broker/kafka"
	"github.com/students-api/bidding-service/internal/models"
	"github.com/students-api/bidding-service/internal/pb/bidding_api"
)

type Implementation struct {
	bidding_api.UnimplementedBiddingServiceServer
	storage    BidRepository
	producer   EventSender
	redisStore LockService
}

func NewBiddingService(storage BidRepository, producer EventSender, redis LockService) *Implementation {
	return &Implementation{
		storage:    storage,
		producer:   producer,
		redisStore: redis,
	}
}

func (s *Implementation) PlaceBid(ctx context.Context, req *bidding_api.PlaceBidRequest) (*bidding_api.PlaceBidResponse, error) {
	if req.Amount <= 0 {
		return &bidding_api.PlaceBidResponse{Success: false, Message: "Сумма должна быть больше нуля"}, nil
	}

	locked := false
	for i := 0; i < 50; i++ {
		ok, err := s.redisStore.AcquireLock(ctx, req.AuctionId)
		if err != nil {
			log.Printf("Redis error: %v", err)
			break 
		}
		if ok {
			locked = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if !locked {
		return &bidding_api.PlaceBidResponse{
			Success: false, 
			Message: "Сервер перегружен (слишком много ставок на этот лот), попробуйте позже",
		}, nil
	}

	defer s.redisStore.ReleaseLock(ctx, req.AuctionId)
	
	newBid := models.Bid{
		AuctionID: req.AuctionId,
		UserID:    req.UserId,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	err := s.storage.CreateBid(ctx, newBid)
	if err != nil {
		log.Printf("Failed to create bid: %v", err)
		return &bidding_api.PlaceBidResponse{Success: false, Message: "Ошибка БД"}, nil
	}

	event := kafka.BidEvent{
		EventID:      fmt.Sprintf("%d-%d", newBid.AuctionID, time.Now().UnixNano()),
		AuctionID:    newBid.AuctionID,
		UserID:       newBid.UserID,
		Amount:       newBid.Amount,
		CurrencyCode: "USD",
		CreatedAt:    time.Now(),
	}

	if err := s.producer.SendBidPlaced(event); err != nil {
			log.Printf("Warning: Failed to send Kafka event: %v", err)
	}

	return &bidding_api.PlaceBidResponse{
		Success: true,
		Message: "Ставка принята",
	}, nil
}

func (s *Implementation) GetAuctionState(ctx context.Context, req *bidding_api.GetAuctionStateRequest) (*bidding_api.GetAuctionStateResponse, error) {
	return nil, errors.New("not implemented yet")
}