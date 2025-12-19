package main

import (
	"context"
	"log"
	"time"

	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := bidding_api.NewBiddingServiceClient(conn)


	req := &bidding_api.PlaceBidRequest{
		AuctionId: 2,
		UserId:    777,
		Amount:    1000, 
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.PlaceBid(ctx, req)
	if err != nil {
		log.Fatalf("Error calling PlaceBid: %v", err)
	}

	log.Printf("Response: Success=%v, Message=%s", res.Success, res.Message)
}