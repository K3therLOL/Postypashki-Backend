package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ResultRepository struct {
	storageTime time.Duration
	ctx         context.Context
	cache       *redis.Client
}

func NewResultRepository(storageTime int) *ResultRepository {
	return &ResultRepository{
		storageTime: time.Duration(storageTime) * time.Hour,
		ctx:         context.TODO(),
		cache: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (r *ResultRepository) Get(key string) (string, error) {
	return r.cache.Get(r.ctx, key).Result()
}

func (r *ResultRepository) Set(key, value string) error {
	return r.cache.Set(r.ctx, key, value, r.storageTime).Err()
}
