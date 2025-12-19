package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/students-api/bidding-service/internal/pb/auction_api"
	"github.com/students-api/bidding-service/internal/services/auction_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	dsn := "postgres://user:password@localhost:5432/auction_db_1"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	svc := auction_service.NewAuctionService(pool)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	auction_api.RegisterAuctionServiceServer(grpcServer, svc)
	reflection.Register(grpcServer)

	go func() {
		log.Println("Starting Auction Service on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down Auction Service...")
	grpcServer.GracefulStop()
}