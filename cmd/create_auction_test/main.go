package main

import (
	"context"
	"log"
	"time"

	"github.com/students-api/bidding-service/internal/pb/auction_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := auction_api.NewAuctionServiceClient(conn)

	log.Println("Запрос на создание аукциона (Object ID 436535)...")
	
	req := &auction_api.CreateAuctionRequest{
		ObjectId:   436535, 
		StartPrice: 1000000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.CreateAuction(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка создания: %v", err)
	}

	log.Println("---------------------------------------------------")
	log.Printf("АУКЦИОН СОЗДАН УСПЕШНО!")
	log.Printf("ID: %d", res.AuctionId)
	log.Printf("Название: %s", res.Title)
	log.Printf("Художник: %s", res.Artist)
	log.Printf("Картинка: %s", res.ImageUrl)
	log.Println("---------------------------------------------------")
}