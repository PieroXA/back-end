package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Env                string
	DatabaseURL        string
	RedisURL           string
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
	S3Endpoint         string
	S3AccessKey        string
	S3SecretKey        string
	S3Bucket           string
	S3Region           string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error if .env doesn't exist (production uses real env vars)

	accessDur, err := time.ParseDuration(getEnv("JWT_ACCESS_DURATION", "24h"))
	if err != nil {
		return nil, fmt.Errorf("config.Load: invalid JWT_ACCESS_DURATION: %w", err)
	}
	refreshDur, err := time.ParseDuration(getEnv("JWT_REFRESH_DURATION", "720h"))
	if err != nil {
		return nil, fmt.Errorf("config.Load: invalid JWT_REFRESH_DURATION: %w", err)
	}

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DatabaseURL:        mustGetEnv("DATABASE_URL"),
		RedisURL:           getEnv("REDIS_URL", ""),
		JWTSecret:          mustGetEnv("JWT_SECRET"),
		JWTAccessDuration:  accessDur,
		JWTRefreshDuration: refreshDur,
		S3Endpoint:         getEnv("S3_ENDPOINT", ""),
		S3AccessKey:        getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:        getEnv("S3_SECRET_KEY", ""),
		S3Bucket:           getEnv("S3_BUCKET", "mifinca-photos"),
		S3Region:           getEnv("S3_REGION", "us-east-1"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: required env var %s is not set", key))
	}
	return v
}
