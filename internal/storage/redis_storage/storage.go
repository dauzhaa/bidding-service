package redis_storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{client: client}
}

func (r *RedisStorage) AcquireLock(ctx context.Context, auctionID int64) (bool, error) {
	key := fmt.Sprintf("auction_lock:%d", auctionID)
	
	success, err := r.client.SetNX(ctx, key, "locked", 5*time.Second).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *RedisStorage) ReleaseLock(ctx context.Context, auctionID int64) error {
	key := fmt.Sprintf("auction_lock:%d", auctionID)
	return r.client.Del(ctx, key).Err()
}