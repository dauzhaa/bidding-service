package auction_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/students-api/bidding-service/internal/pb/auction_api"
)

type AuctionService struct {
	auction_api.UnimplementedAuctionServiceServer
	repo   AuctionRepository
	APIURL string
}

func NewAuctionService(repo AuctionRepository) *AuctionService {
	return &AuctionService{
		repo:   repo,
		APIURL: "https://collectionapi.metmuseum.org/public/collection/v1/objects",
	}
}

type MetObject struct {
	ObjectID     int    `json:"objectID"`
	Title        string `json:"title"`
	Artist       string `json:"artistDisplayName"`
	PrimaryImage string `json:"primaryImage"`
}

func (s *AuctionService) CreateAuction(ctx context.Context, req *auction_api.CreateAuctionRequest) (*auction_api.CreateAuctionResponse, error) {
	url := fmt.Sprintf("%s/%d", s.APIURL, req.ObjectId)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Met API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Met API returned status: %d", resp.StatusCode)
	}

	var metObj MetObject
	if err := json.NewDecoder(resp.Body).Decode(&metObj); err != nil {
		return nil, fmt.Errorf("failed to decode Met API response: %v", err)
	}

	auction := &auction_api.Auction{
		Id:           int64(metObj.ObjectID),
		Title:        metObj.Title,
		Artist:       metObj.Artist,
		CurrentPrice: req.StartPrice,
		ImageUrl:     metObj.PrimaryImage,
	}

	if auction.Title == "" {
		auction.Title = "Unknown Title"
	}
	if auction.Artist == "" {
		auction.Artist = "Unknown Artist"
	}

	err = s.repo.CreateAuction(ctx, auction)
	if err != nil {
		log.Printf("Failed to create auction in DB: %v", err)
		return nil, err
	}

	return &auction_api.CreateAuctionResponse{
		AuctionId: auction.Id,
		Title:     auction.Title,
		Artist:    auction.Artist,
		ImageUrl:  auction.ImageUrl,
	}, nil
}

func (s *AuctionService) ListAuctions(ctx context.Context, req *auction_api.ListAuctionsRequest) (*auction_api.ListAuctionsResponse, error) {
	auctions, err := s.repo.ListAuctions(ctx)
	if err != nil {
		return nil, err
	}
	return &auction_api.ListAuctionsResponse{Auctions: auctions}, nil
}