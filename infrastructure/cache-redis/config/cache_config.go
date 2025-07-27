package config

import "time"

type CacheConfig struct {
	DefaultTTL      time.Duration
	MaxEntries      int
	CleanupInterval time.Duration
}

func NewDefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL:      5 * time.Minute,
		MaxEntries:      10000,
		CleanupInterval: 10 * time.Minute,
	}
}
