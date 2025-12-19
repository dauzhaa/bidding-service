package bidding_service

import (
	"context"
	"github.com/students-api/bidding-service/internal/broker/kafka"
	"github.com/students-api/bidding-service/internal/models"
)

type BidRepository interface {
	CreateBid(ctx context.Context, bid models.Bid) error
}

type EventSender interface {
	SendBidPlaced(event kafka.BidEvent) error
}

type LockService interface {
	AcquireLock(ctx context.Context, auctionID int64) (bool, error)
	ReleaseLock(ctx context.Context, auctionID int64) error
}