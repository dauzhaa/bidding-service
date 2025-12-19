package auction_repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/students-api/bidding-service/internal/pb/auction_api"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateAuction(ctx context.Context, auction *auction_api.Auction) error {
	query := `
		INSERT INTO auctions (id, title, artist, start_price, image_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		auction.Id,
		auction.Title,
		auction.Artist,
		auction.CurrentPrice,
		auction.ImageUrl,
		"active",
	)
	return err
}

func (r *PostgresRepository) ListAuctions(ctx context.Context) ([]*auction_api.Auction, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, artist, start_price, image_url, status FROM auctions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []*auction_api.Auction
	for rows.Next() {
		var a auction_api.Auction
		var statusStr string
		
		if err := rows.Scan(&a.Id, &a.Title, &a.Artist, &a.CurrentPrice, &a.ImageUrl, &statusStr); err != nil {
			return nil, err
		}
		auctions = append(auctions, &a)
	}
	return auctions, nil
}