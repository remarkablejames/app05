package repo_impl

import (
	"app05/internal/core/domain/entities"
	"app05/internal/infrastructure/storage/postgres/repo_impl/dbUtils"
	"app05/pkg/appErrors"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	query := `
        INSERT INTO users (email, hashed_password, first_name, last_name,subscribed_to_newsletter, role, active, email_verified)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.HashedPassword,
		user.FirstName,
		user.LastName,
		user.SubscribedToNewsletter,
		user.Role,
		user.Active,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return appErrors.New(appErrors.CodeBadRequest, "email already taken")
		}
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	query := `
        SELECT id, email, hashed_password, first_name, last_name, profile_picture_url, role,
               active, email_verified, created_at, updated_at, last_login_at
        FROM users
        WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePictureURL,
		&user.Role,
		&user.Active,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.New(appErrors.CodeNotFound, "user not found")
	}
	return user, err
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	query := `
		SELECT id, email, hashed_password, first_name, last_name, profile_picture_url, role,
			   active, email_verified, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePictureURL,
		&user.Role,
		&user.Active,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.New(appErrors.CodeNotFound, "user not found")
	}
	return user, err
}

func (r *UserRepositoryImpl) GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}

	// First query user information
	userQuery := `
        SELECT id, email, first_name, last_name, profile_picture_url, role, title, bio,
               active, email_verified, created_at, updated_at, last_login_at
        FROM users
        WHERE id = $1`

	var bio, title sql.NullString
	err := r.db.QueryRowContext(ctx, userQuery, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePictureURL,
		&user.Role,
		&title,
		&bio,
		&user.Active,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.New(appErrors.CodeNotFound, "user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	if title.Valid {
		user.Title = title.String
	}

	if bio.Valid {
		user.Bio = bio.String
	}

	if err == sql.ErrNoRows {
		return nil, appErrors.New(appErrors.CodeNotFound, "user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// Then query active session if exists
	sessionQuery := `
        SELECT id, user_id, token, refresh_token, status, 
               device_info, expires_at, last_activity_at, created_at, 
               revoked_at, revoked_reason
        FROM sessions 
        WHERE user_id = $1 
          AND status = 'active' 
          AND expires_at > NOW()
        ORDER BY created_at DESC 
        LIMIT 1`

	session := &entities.Session{}
	var deviceInfoBytes []byte
	err = r.db.QueryRowContext(ctx, sessionQuery, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.RefreshToken,
		&session.Status,
		&deviceInfoBytes,
		&session.ExpiresAt,
		&session.LastActivityAt,
		&session.CreatedAt,
		&session.RevokedAt,
		&session.RevokedReason,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, appErrors.New(appErrors.CodeNotFound, "session not found")
	}

	if err != sql.ErrNoRows {
		if err := json.Unmarshal(deviceInfoBytes, &session.DeviceInfo); err != nil {
			return nil, fmt.Errorf("error unmarshalling device info: %w", err)
		}
		user.CurrentSession = session
	}

	return user, nil
}

func (r *UserRepositoryImpl) GetUserByResetToken(ctx context.Context, hashedToken string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	var user entities.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name, last_name, role, active, 
         email_verified, password_reset_token, reset_token_expires_at 
         FROM users WHERE password_reset_token = $1`,
		hashedToken).Scan(
		&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName,
		&user.Role, &user.Active, &user.EmailVerified, &user.PasswordResetToken,
		&user.ResetTokenExpiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET 
         hashed_password = $1,
         password_reset_token = $2,
         reset_token_expires_at = $3
         WHERE id = $4`,
		user.HashedPassword,
		user.PasswordResetToken,
		user.ResetTokenExpiresAt,
		user.ID)
	return err
}

func (r *UserRepositoryImpl) UpdateProfilePicture(ctx context.Context, userID uuid.UUID, removeProfilePicture bool, profilePictureURL string) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()
	query := `
        UPDATE users 
        SET profile_picture_url = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`

	if removeProfilePicture {
		profilePictureURL = ""
	}

	result, err := r.db.ExecContext(ctx, query, profilePictureURL, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return appErrors.New(appErrors.CodeNotFound, "user not found")
	}

	return nil
}

func (r *UserRepositoryImpl) CreateVerificationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	// Delete existing codes for the user
	deleteQuery := `
        DELETE FROM email_verification_codes
        WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}

	// Insert the new code
	insertQuery := `
        INSERT INTO email_verification_codes (user_id, code, expires_at)
        VALUES ($1, $2, $3)`
	_, err = r.db.ExecContext(ctx, insertQuery, userID, code, expiresAt)
	return err
}

func (r *UserRepositoryImpl) VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error {
	ctx, cancel := context.WithTimeout(ctx, dbUtils.QueryTimeoutDuration)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify the code, delete it if valid and not expired
	verifyQuery := `
        DELETE FROM email_verification_codes
        WHERE user_id = $1 
        AND code = $2 
        AND used_at IS NULL 
        AND expires_at > CURRENT_TIMESTAMP`

	result, err := tx.ExecContext(ctx, verifyQuery, userID, code)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return appErrors.New(appErrors.CodeBadRequest, "invalid or expired verification code")
	}

	// Update user's email_verified status
	updateUserQuery := `
        UPDATE users 
        SET email_verified = true,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateUserQuery, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
