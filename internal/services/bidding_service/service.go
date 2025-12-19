package bidding_service

import (
	"context"
	"errors"

	"github.com/students-api/bidding-service/internal/pb/bidding_api"
)

type Implementation struct {
	bidding_api.UnimplementedBiddingServiceServer
}

func NewBiddingService() *Implementation {
	return &Implementation{}
}

func (s *Implementation) PlaceBid(ctx context.Context, req *bidding_api.PlaceBidRequest) (*bidding_api.PlaceBidResponse, error) {
	return &bidding_api.PlaceBidResponse{
		Success: true,
		Message: "Ставка принята (пока заглушка)",
	}, nil
}

func (s *Implementation) GetAuctionState(ctx context.Context, req *bidding_api.GetAuctionStateRequest) (*bidding_api.GetAuctionStateResponse, error) {
	return nil, errors.New("not implemented yet")
}