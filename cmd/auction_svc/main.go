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
	"github.com/students-api/bidding-service/internal/storage/auction_repo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	ctx := context.Background()

	dbHost := getEnv("DB_HOST", "localhost")
	dsn := "postgres://user:password@" + dbHost + ":5432/auction_db_1"
	
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	repo := auction_repo.NewPostgresRepository(pool)
	
	svc := auction_service.NewAuctionService(repo)

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