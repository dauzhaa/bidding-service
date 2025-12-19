package bid_storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/students-api/bidding-service/internal/models"
)

type BidStorage struct {
	shard1 *pgxpool.Pool
	shard2 *pgxpool.Pool
}

func NewBidStorage(shard1, shard2 *pgxpool.Pool) *BidStorage {
	return &BidStorage{
		shard1: shard1,
		shard2: shard2,
	}
}

func (s *BidStorage) getShard(auctionID int64) *pgxpool.Pool {
	if auctionID%2 == 0 {
		return s.shard1
	}
	return s.shard2
}

func (s *BidStorage) CreateBid(ctx context.Context, bid models.Bid) error {
	pool := s.getShard(bid.AuctionID)
	
	query := `INSERT INTO bids (auction_id, user_id, amount, created_at) VALUES ($1, $2, $3, $4)`
	
	_, err := pool.Exec(ctx, query, bid.AuctionID, bid.UserID, bid.Amount, time.Now())
	return err
}