package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	biddingSvcAddr := getEnv("BIDDING_SVC_ADDR", "localhost:50051")
	auctionSvcAddr := getEnv("AUCTION_SVC_ADDR", "localhost:50052")

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := bidding_api.RegisterBiddingServiceHandlerFromEndpoint(ctx, mux, biddingSvcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register BiddingService gateway: %v", err)
	}
	log.Printf("Registered BiddingService gateway -> %s", biddingSvcAddr)

	err = auction_api.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, auctionSvcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register AuctionService gateway: %v", err)
	}
	log.Printf("Registered AuctionService gateway -> %s", auctionSvcAddr)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Println("Starting HTTP Gateway on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down Gateway...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}