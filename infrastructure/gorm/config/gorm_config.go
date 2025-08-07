package config

import (
	"fmt"
	envConfig "github.com/pdh9523/go-hexarch/shared/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel
}

func NewGormConfig() *GormConfig {
	envConfig.LoadEnv()
	return &GormConfig{
		Host:            envConfig.GetEnv("DB_HOST", "localhost"),
		Port:            envConfig.GetEnvInt("DB_PORT", 5432),
		Username:        envConfig.GetEnv("DB_USERNAME", "username"),
		Password:        envConfig.GetEnv("DB_PASSWORD", "password"),
		Database:        envConfig.GetEnv("DB_NAME", "database_name"),
		SSLMode:         envConfig.GetEnv("DB_SSL_MODE", "disable"),
		MaxIdleConns:    envConfig.GetEnvInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    envConfig.GetEnvInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: envConfig.GetEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		LogLevel:        parseLogLevel(envConfig.GetEnv("DB_LOG_LEVEL", "info")),
	}
}

func (g *GormConfig) NewGormDB() (*gorm.DB, error) {
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(g.LogLevel),
	}

	db, err := gorm.Open(postgres.Open(g.DSN()), &gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(g.MaxIdleConns)
	sqlDB.SetMaxOpenConns(g.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(g.ConnMaxLifetime)

	return db, nil
}

func (g *GormConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Seoul",
		g.Host, g.Username, g.Password, g.Database, g.Port, g.SSLMode,
	)
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
