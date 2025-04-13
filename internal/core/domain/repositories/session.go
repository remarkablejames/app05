package repositories

import (
	"app05/internal/core/domain/entities"
	"context"
	"time"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entities.Session) error
	GetActiveSession(ctx context.Context, userID int64) (*entities.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*entities.Session, error)
	//UpdateSession(ctx context.Context, session *entities.Session) error
	RevokeUserSessions(ctx context.Context, userID int64, reason string) error
	GetActiveSessionByUserID(ctx context.Context, userID uint) (*entities.Session, error)
	UpdateSession(ctx context.Context, session *entities.Session) error
	DeleteRevokedSessions(ctx context.Context, olderThan time.Time) (int64, error)
	//CleanupExpiredSessions(ctx context.Context) error
}
