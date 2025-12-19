package models

import "time"

type Bid struct {
	ID        int64
	AuctionID int64
	UserID    int64
	Amount    int64
	CreatedAt time.Time
}