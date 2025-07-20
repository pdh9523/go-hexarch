package config

import (
	envConfig "github.com/pdh9523/go-hexarch/shared/common/config"
	"time"
)

type TokenConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

type Config struct {
	Token TokenConfig
}

func DefaultConfig() Config {
	envConfig.LoadEnv()

	return Config{
		Token: TokenConfig{
			AccessTokenSecret:  envConfig.GetEnv("ACCESS_SECRET", "default-access-secret-change-in-production"),
			RefreshTokenSecret: envConfig.GetEnv("REFRESH_SECRET", "default-refresh-secret-change-in-production"),
			AccessTokenTTL:     envConfig.GetEnvDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL:    envConfig.GetEnvDuration("REFRESH_TOKEN_TTL", 24*time.Hour),
		},
	}
}
