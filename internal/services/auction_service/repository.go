package auction_service

import (
	"context"

	"github.com/students-api/bidding-service/internal/pb/auction_api"
)

type AuctionRepository interface {
	CreateAuction(ctx context.Context, auction *auction_api.Auction) error
	ListAuctions(ctx context.Context) ([]*auction_api.Auction, error)
}