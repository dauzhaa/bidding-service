package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
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

	totalRequests := 100
	var successCount int64
	var failCount int64

	var wg sync.WaitGroup
	wg.Add(totalRequests)

	fmt.Printf("🚀 ЗАПУСК: Атака %d одновременных запросов...\n", totalRequests)
	startTime := time.Now()

	for i := 0; i < totalRequests; i++ {
		go func(userID int) {
			defer wg.Done()

			req := &bidding_api.PlaceBidRequest{
				AuctionId: 100,
				UserId:    int64(userID),
				Amount:    int64(1000 + userID),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			resp, err := client.PlaceBid(ctx, req)
			
			if err == nil && resp.Success {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("------------------------------------------------")
	fmt.Printf("УСПЕШНО: %d\n", successCount)
	fmt.Printf("ОТКЛОНЕНО: %d (сработал Redis Lock или таймаут)\n", failCount)
	fmt.Printf("ВРЕМЯ: %v\n", duration)
	fmt.Println("------------------------------------------------")
}