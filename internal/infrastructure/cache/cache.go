package cache

import (
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/entities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type SessionCache struct {
	client *redis.Client
	logger contracts.Logger
}

func NewSessionCache(redisURL string, logger contracts.Logger) (*SessionCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &SessionCache{client: client, logger: logger}, nil
}

// StoreSession stores a session in the cache
func (c *SessionCache) StoreSession(ctx context.Context, session *entities.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	pipe := c.client.Pipeline()
	expiration := time.Until(session.ExpiresAt)

	// Store session with TTL
	pipe.Set(ctx,
		"session:"+session.Token,
		sessionJSON,
		expiration,
	)

	// Store user's active session with same TTL
	pipe.Set(ctx,
		fmt.Sprintf("user_session:%d", session.UserID),
		session.Token,
		expiration,
	)

	// Store user role only if it exists
	if session.UserRole != "" {
		pipe.Set(ctx,
			fmt.Sprintf("user_role:%s", session.Token),
			string(session.UserRole),
			expiration,
		)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	return nil
}

func (c *SessionCache) GetSession(ctx context.Context, token string) (*entities.Session, error) {
	sessionJSON, err := c.client.Get(ctx, "session:"+token).Result()
	if err != nil {
		return nil, err
	}

	var session entities.Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (c *SessionCache) ValidateSession(ctx context.Context, token string) (bool, error) {
	exists, err := c.client.Exists(ctx, "session:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (c *SessionCache) InvalidateSession(ctx context.Context, token string, userID uuid.UUID) error {
	// Remove both session entries
	pipe := c.client.Pipeline()
	pipe.Del(ctx, "session:"+token)
	pipe.Del(ctx, fmt.Sprintf("user_session:%d", userID)) // Fix: proper formatting of userID
	_, err := pipe.Exec(ctx)
	return err
}

// GetUserActiveSession retrieves the active session token for a user
func (c *SessionCache) GetUserActiveSession(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := c.client.Get(ctx, fmt.Sprintf("user_session:%d", userID)).Result() // Fix: proper formatting of userID
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return token, nil
}

// GetSessionByToken retrieves a session by its token
func (c *SessionCache) GetSessionByToken(ctx context.Context, token string) (*entities.Session, error) {
	return c.GetSession(ctx, token) // Reuse existing GetSession method
}

func (c *SessionCache) IsSessionHealthy(ctx context.Context, token string) bool {
	// First check if session exists in Redis
	exists, err := c.client.Exists(ctx, "session:"+token).Result()
	if err != nil || exists == 0 {
		return false
	}

	session, err := c.GetSession(ctx, token)
	if err != nil || session == nil {
		return false
	}

	return session.Status == entities.SessionStatusActive &&
		time.Now().Before(session.ExpiresAt) &&
		session.RevokedAt == nil
}

func (c *SessionCache) DeleteSession(ctx context.Context, token string, userID uuid.UUID) error {
	// Invalidate session in Redis by removing session and user_session keys
	if err := c.InvalidateSession(ctx, token, userID); err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}
	return nil
}
