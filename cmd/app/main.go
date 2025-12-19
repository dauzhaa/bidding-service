package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/students-api/bidding-service/internal/broker/kafka"
	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"github.com/students-api/bidding-service/internal/services/bidding_service"
	"github.com/students-api/bidding-service/internal/storage/bid_storage"
	"github.com/students-api/bidding-service/internal/storage/redis_storage" 
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	redisStore := redis_storage.NewRedisStorage("localhost:6379")
	log.Println("Connected to Redis")

	producer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()

	dsn1 := "postgres://user:password@localhost:5432/auction_db_1"
	pool1, err := pgxpool.New(ctx, dsn1)
	if err != nil {
		log.Fatalf("Unable to connect to shard 1: %v", err)
	}
	defer pool1.Close()

	dsn2 := "postgres://user:password@localhost:5433/auction_db_2"
	pool2, err := pgxpool.New(ctx, dsn2)
	if err != nil {
		log.Fatalf("Unable to connect to shard 2: %v", err)
	}
	defer pool2.Close()

	log.Println("Connected to Database Shards and Kafka successfully")

	storage := bid_storage.NewBidStorage(pool1, pool2)
	
	serviceImplementation := bidding_service.NewBiddingService(storage, producer, redisStore)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	bidding_api.RegisterBiddingServiceServer(grpcServer, serviceImplementation)
	reflection.Register(grpcServer)

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
}