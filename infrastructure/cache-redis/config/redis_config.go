package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(config *RedisConfig) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.New("failed to connect to Redis")
	}
	return &RedisClient{client: client}, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
