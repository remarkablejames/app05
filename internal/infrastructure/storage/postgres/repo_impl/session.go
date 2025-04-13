package repo_impl

import (
	"app05/internal/core/domain/entities"
	"app05/internal/infrastructure/storage/postgres/repo_impl/dbUtils"
	"app05/pkg/appErrors"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type SessionRepositoryImpl struct {
	db *sql.DB
}

func (r *SessionRepositoryImpl) GetActiveSessionByUserID(ctx context.Context, userID uint) (*entities.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (r *SessionRepositoryImpl) GetSessionByToken(ctx context.Context, token string) (*entities.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	session := &entities.Session{}
	query := `
        SELECT id, user_id, token, refresh_token, status, device_info,
               expires_at, last_activity_at, created_at, revoked_at, revoked_reason
        FROM sessions 
        WHERE token = $1`

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.RefreshToken,
		&session.Status,
		&session.DeviceInfo,
		&session.ExpiresAt,
		&session.LastActivityAt,
		&session.CreatedAt,
		&session.RevokedAt,
		&session.RevokedReason,
	)

	if err == sql.ErrNoRows {
		return nil,
			appErrors.New(appErrors.CodeNotFound, "session not found")
	}
	return session, err
}

func (r *SessionRepositoryImpl) UpdateSession(ctx context.Context, session *entities.Session) error {
	//TODO implement me
	panic("implement me")
}

func NewSessionRepository(db *sql.DB) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{db: db}
}

func (r *SessionRepositoryImpl) CreateSession(ctx context.Context, session *entities.Session) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, revoke any existing active sessions
	_, err = tx.ExecContext(ctx, `
        UPDATE sessions 
        SET status = $1, revoked_at = CURRENT_TIMESTAMP, revoked_reason = $2
        WHERE user_id = $3 AND status = $4`,
		entities.SessionStatusRevoked,
		"New login from another device",
		session.UserID,
		entities.SessionStatusActive,
	)
	if err != nil {
		return err
	}

	// Marshal DeviceInfo to JSON
	deviceInfoJSON, err := json.Marshal(session.DeviceInfo)
	if err != nil {
		return err
	}

	// Then create the new session
	query := `
        INSERT INTO sessions (
            user_id, token, refresh_token, status, expires_at,
            device_info, last_activity_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
        RETURNING id, created_at`

	err = tx.QueryRowContext(
		ctx,
		query,
		session.UserID,
		session.Token,
		session.RefreshToken,
		session.Status,
		session.ExpiresAt,
		deviceInfoJSON,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return err
	}

	// Update user's last login timestamp
	_, err = tx.ExecContext(ctx, `
        UPDATE users 
        SET last_login_at = CURRENT_TIMESTAMP
        WHERE id = $1`,
		session.UserID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SessionRepositoryImpl) GetActiveSession(ctx context.Context, userID int64) (*entities.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	session := &entities.Session{}
	query := `
        SELECT id, user_id, token, refresh_token, status, device_info,
               expires_at, last_activity_at, created_at, revoked_at, revoked_reason
        FROM sessions
        WHERE user_id = $1 AND status = $2 AND expires_at > CURRENT_TIMESTAMP
        ORDER BY created_at DESC
        LIMIT 1`

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		entities.SessionStatusActive,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.RefreshToken,
		&session.Status,
		&session.DeviceInfo,
		&session.ExpiresAt,
		&session.LastActivityAt,
		&session.CreatedAt,
		&session.RevokedAt,
		&session.RevokedReason,
	)

	if err == sql.ErrNoRows {
		return nil,
			appErrors.New(appErrors.CodeNotFound, " no active session found")
	}
	return session, err
}

func (r *SessionRepositoryImpl) RevokeUserSessions(ctx context.Context, userID int64, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
        UPDATE sessions 
        SET status = $1, revoked_at = CURRENT_TIMESTAMP, revoked_reason = $2
        WHERE user_id = $3 AND status = $4`,
		entities.SessionStatusRevoked,
		reason,
		userID,
		entities.SessionStatusActive,
	)
	return err
}

func (r *SessionRepositoryImpl) DeleteRevokedSessions(ctx context.Context, olderThan time.Time) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	// Delete in batches to avoid long-running transactions
	const batchSize = 1000
	var totalDeleted int64

	for {
		query := `
            WITH batch AS (
                SELECT id FROM sessions 
                WHERE status = 'revoked' 
				AND revoked_at < $1 
                LIMIT $2
                FOR UPDATE SKIP LOCKED
            )
            DELETE FROM sessions 
            WHERE id IN (SELECT id FROM batch)
            RETURNING id`

		result, err := r.db.ExecContext(ctx, query, olderThan, batchSize)
		if err != nil {
			return totalDeleted, err
		}

		deleted, _ := result.RowsAffected()
		totalDeleted += deleted

		if deleted < batchSize {
			break
		}
	}

	return totalDeleted, nil
}
