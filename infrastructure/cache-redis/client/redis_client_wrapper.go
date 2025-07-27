package client

import (
	"context"
	"errors"
	"github.com/pdh9523/go-hexarch/infrastructure/cache-redis/config"
	"github.com/pdh9523/go-hexarch/infrastructure/cache-redis/serializer"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClientWrapper struct {
	client     *config.RedisClient
	serializer *serializer.JSONSerializer
}

func NewRedisClientWrapper(client *config.RedisClient) *RedisClientWrapper {
	return &RedisClientWrapper{
		client:     client,
		serializer: serializer.NewJSONSerializer(),
	}
}

func (w *RedisClientWrapper) GetClient() *redis.Client {
	return w.client.GetClient()
}

func (w *RedisClientWrapper) SetDataWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := w.GetClient().Set(ctx, key, value, ttl).Err(); err != nil {
		return errors.New("failed to set value in cache" + err.Error())
	}
	return nil
}

func (w *RedisClientWrapper) SetData(ctx context.Context, key string, value interface{}) error {
	return w.SetDataWithTTL(ctx, key, value, 0)
}

func (w *RedisClientWrapper) SetJSONWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := w.serializer.Serialize(value)
	if err != nil {
		return errors.New("failed to serialize to cache" + err.Error())
	}
	return w.SetDataWithTTL(ctx, key, data, ttl)
}

func (w *RedisClientWrapper) SetJSON(ctx context.Context, key string, value interface{}) error {
	return w.SetJSONWithTTL(ctx, key, value, 0)
}

func (w *RedisClientWrapper) Get(ctx context.Context, key string) (string, error) {
	result, err := w.GetClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", errors.New("failed to get from cache" + err.Error())
	}
	return result, nil
}

func (w *RedisClientWrapper) GetJSON(ctx context.Context, key string, result interface{}) error {
	data, err := w.GetClient().Get(ctx, key).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
	}

	if err := w.serializer.Deserialize(data, result); err != nil {
		return errors.New("failed to deserialize from cache" + err.Error())
	}

	return nil
}

func (w *RedisClientWrapper) Delete(ctx context.Context, key string) error {
	if err := w.GetClient().Del(ctx, key).Err(); err != nil {
		return errors.New("failed to delete value from cache" + err.Error())
	}

	return nil
}

func (w *RedisClientWrapper) DeletePattern(ctx context.Context, pattern string) error {
	client := w.GetClient()

	const batchSize = 100
	var totalDeleted int64

	for {
		var keys []string
		iter := client.Scan(ctx, 0, pattern, batchSize).Iterator()
		for iter.Next(ctx) && len(keys) < batchSize {
			keys = append(keys, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return errors.New("failed to iterate over redis" + err.Error())
		}

		if len(keys) == 0 {
			break
		}

		pipe := client.Pipeline()
		for _, key := range keys {
			pipe.Del(ctx, key)
		}

		results, err := pipe.Exec(ctx)
		if err != nil {
			return errors.New("failed to delete from redis" + err.Error())
		}

		for _, result := range results {
			if delResult, ok := result.(*redis.IntCmd); ok {
				totalDeleted += delResult.Val()
			}
		}

		if len(keys) < batchSize {
			break
		}
	}
	return nil
}

func (w *RedisClientWrapper) Exists(ctx context.Context, key string) (bool, error) {
	result, err := w.GetClient().Exists(ctx, key).Result()
	if err != nil {
		return false, errors.New("failed to get value from cache" + err.Error())
	}
	return result > 0, nil
}

func (w *RedisClientWrapper) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := w.GetClient().Expire(ctx, key, ttl).Err(); err != nil {
		return errors.New("failed to set value in cache" + err.Error())
	}
	return nil
}

func (w *RedisClientWrapper) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := w.GetClient().TTL(ctx, key).Result()
	if err != nil {
		return time.Duration(0), errors.New("failed to set value in cache" + err.Error())
	}
	return ttl, nil
}
