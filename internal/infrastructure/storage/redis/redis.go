package redis

import (
	"app05/internal/core/domain/entities"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{client: client}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

// Session-related methods
func (c *Cache) StoreSession(ctx context.Context, session *entities.Session) error {
	key := sessionKey(session.Token)
	value, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Store session with expiration
	if err := c.client.Set(ctx, key, value, time.Until(session.ExpiresAt)).Err(); err != nil {
		return err
	}

	// Store user's active session token
	userKey := userSessionKey(uint(session.UserID))
	return c.client.Set(ctx, userKey, session.Token, time.Until(session.ExpiresAt)).Err()
}

func (c *Cache) GetSession(ctx context.Context, token string) (*entities.Session, error) {
	value, err := c.client.Get(ctx, sessionKey(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session entities.Session
	if err := json.Unmarshal([]byte(value), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (c *Cache) InvalidateSession(ctx context.Context, token string, userID uuid.UUID) error {
	pipe := c.client.Pipeline()
	pipe.Del(ctx, sessionKey(token))
	pipe.Del(ctx, userSessionKey(userID))
	_, err := pipe.Exec(ctx)
	return err
}

// Helper functions for key generation
func sessionKey(token string) string {
	return "session:" + token
}

func userSessionKey(userID uint) string {
	return "user_session:" + string(userID)
}
