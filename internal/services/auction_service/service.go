package auction_service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/students-api/bidding-service/internal/integration/metmuseum"
	"github.com/students-api/bidding-service/internal/pb/auction_api"
)

type Implementation struct {
	auction_api.UnimplementedAuctionServiceServer
	metClient *metmuseum.Client
	dbShard1  *pgxpool.Pool
}

func NewAuctionService(db *pgxpool.Pool) *Implementation {
	return &Implementation{
		metClient: metmuseum.NewClient(),
		dbShard1:  db,
	}
}

func (s *Implementation) CreateAuction(ctx context.Context, req *auction_api.CreateAuctionRequest) (*auction_api.CreateAuctionResponse, error) {
	data, err := s.metClient.GetObjectData(req.ObjectId)
	if err != nil {
		log.Printf("Error fetching from Met Museum: %v", err)
		return nil, err
	}

	log.Printf("Creating auction for: %s by %s", data.Title, data.Artist)
	
	query := `INSERT INTO auctions (id, title, artist, start_price, image_url, status) VALUES ($1, $2, $3, $4, $5, 'active') RETURNING id`
	
	_, err = s.dbShard1.Exec(ctx, query, data.ObjectID, data.Title, data.Artist, req.StartPrice, data.PrimaryImage)
	if err != nil {
		log.Printf("Failed to insert auction: %v", err)
		return nil, err
	}

	return &auction_api.CreateAuctionResponse{
		AuctionId: data.ObjectID,
		Title:     data.Title,
		Artist:    data.Artist,
		ImageUrl:  data.PrimaryImage,
	}, nil
}

func (s *Implementation) ListAuctions(ctx context.Context, req *auction_api.ListAuctionsRequest) (*auction_api.ListAuctionsResponse, error) {
	return &auction_api.ListAuctionsResponse{}, nil
}