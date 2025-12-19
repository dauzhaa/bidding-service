package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/students-api/bidding-service/internal/pb/auction_api"
	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	biddingAddr := getEnv("BIDDING_SVC_ADDR", "localhost:50051")
	auctionAddr := getEnv("AUCTION_SVC_ADDR", "localhost:50052")

	log.Printf("Connecting to Bidding Service at %s", biddingAddr)
	err := bidding_api.RegisterBiddingServiceHandlerFromEndpoint(ctx, mux, biddingAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register Bidding Service: %v", err)
	}

	log.Printf("Connecting to Auction Service at %s", auctionAddr)
	err = auction_api.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, auctionAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register Auction Service: %v", err)
	}

	log.Println("HTTP Gateway listening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}