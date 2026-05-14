package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type ResultRepository struct {
	storageTime time.Duration
	ctx         context.Context
	cache       *redis.Client
}

func NewResultRepository(storageTime int) *ResultRepository {
	host := os.Getenv("REDIS_HOST")
	addr := fmt.Sprintf("%s:6379", host)
	ctx := context.Background()
	cache := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	_, err := cache.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &ResultRepository{
		storageTime: time.Duration(storageTime) * time.Hour,
		ctx:         ctx,
		cache:       cache,
	}
}

func (r *ResultRepository) Get(key string) (string, error) {
	return r.cache.Get(r.ctx, key).Result()
}

func (r *ResultRepository) Set(key, value string) error {
	return r.cache.Set(r.ctx, key, value, r.storageTime).Err()
}
