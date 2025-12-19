package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/students-api/bidding-service/internal/pb/bidding_api"
	"github.com/students-api/bidding-service/internal/services/bidding_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	serviceImplementation := bidding_service.NewBiddingService()

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