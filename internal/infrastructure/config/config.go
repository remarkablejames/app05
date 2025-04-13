// Package config provides configuration management for the application.
// It handles loading environment variables, database configuration,
// authentication settings, and rate limiting parameters.
package config

import (
	"app05/internal/infrastructure/config/env"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

// AppConfig holds all configuration parameters for the application.
// It centralizes all settings to make configuration management easier.
type AppConfig struct {
	DatabaseURL string
	ServerPort  string
	ServerHost  string
	AppVersion  string
	FrontendURL string
	Env         string
	LogLevel    string
	Auth        AuthConfig
	RateLimiter LimiterConfig
	Redis       RedisConfig
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	Token TokenConfig
}

// TokenConfig contains JWT token configuration parameters.
type TokenConfig struct {
	Secret string        // Secret key used for signing JWT tokens
	Exp    time.Duration // Token expiration duration
	Aud    string        // Token audience claim
	Iss    string        // Token issuer claim
}

type LimiterConfig struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// Database connection constants
const (
	maxOpenConns    = 25
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
	connTimeout     = 10 * time.Second
	defaultTokenExp = 30 * 24 * time.Hour // 30 days
	defaultPort     = "8081"
	defaultHost     = "localhost:8080"
	defaultVersion  = "1.0.0"
	defaultEnv      = "development"
	defaultLogLevel = "info"
)

// LoadConfig initializes and returns the application configuration by reading
// environment variables and setting default values where necessary.
func LoadConfig() *AppConfig {
	dbURL := buildDatabaseURL(
		env.GetString("POSTGRES_USER", ""),
		env.GetString("POSTGRES_PASSWORD", ""),
		env.GetString("POSTGRES_DB", ""),
		env.GetString("POSTGRES_DB_HOST", ""),
		env.GetString("POSTGRES_DB_PORT", ""),
	)

	return &AppConfig{
		DatabaseURL: env.GetString("DB_URL", dbURL),
		ServerPort:  fmt.Sprintf(":%s", env.GetString("SERVER_PORT", defaultPort)),
		ServerHost:  env.GetString("SERVER_HOST", defaultHost),
		AppVersion:  env.GetString("APP_VERSION", defaultVersion),
		FrontendURL: env.GetString("FRONTEND_URL", ""),
		Env:         env.GetString("ENV", defaultEnv),
		LogLevel:    env.GetString("LOG_LEVEL", defaultLogLevel),
		Auth: AuthConfig{
			Token: TokenConfig{
				Secret: env.GetString("TOKEN_SECRET", "MySecret"),
				Exp:    defaultTokenExp,
				Aud:    env.GetString("TOKEN_AUD", "SomoLabs"),
				Iss:    env.GetString("TOKEN_ISS", "SomoLabs"),
			},
		},
		RateLimiter: LimiterConfig{
			RequestPerTimeFrame: env.GetInt("RATE_LIMITER_REQUEST_PER_TIME_FRAME", 20),
			TimeFrame:           env.GetDuration("RATE_LIMITER_TIME_FRAME", 5*time.Second),
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
		Redis: RedisConfig{
			URL:      env.GetString("REDIS_URL", "redis://localhost:6379"),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
		},
	}
}

// buildDatabaseURL constructs a PostgreSQL connection string from individual components.
func buildDatabaseURL(user, password, dbName, host, port string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		user, password, host, port, dbName,
	)
}

// InitDB initializes and returns a configured database connection pool.
// It handles connection setup, pool configuration, and connection testing.
func InitDB(cfg *AppConfig) (*sql.DB, error) {
	log.Printf("\033[1;33mConnecting to the database...\033[0m")

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			return nil, fmt.Errorf("postgres connection error - Code: %s, Message: %s, Detail: %s",
				pqErr.Code, pqErr.Message, pqErr.Detail)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("\033[1;33mSuccessfully connected to the database\033[0m")
	return db, nil
}
