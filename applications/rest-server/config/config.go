package config

import (
	envConfig "github.com/pdh9523/go-hexarch/shared/common/config"
)

type ServerConfig struct {
	Port string
	Mode string
}

func NewDefaultServerConfig() *ServerConfig {
	envConfig.LoadEnv()
	return &ServerConfig{
		Port: envConfig.GetEnv("SERVER_PORT", "8080"),
		Mode: envConfig.GetEnv("SERVER_MODE", "dev"),
	}
}
