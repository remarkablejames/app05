package entities

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleSuperUser  Role = "superuser"
	RoleAdmin      Role = "admin"
	RoleInstructor Role = "instructor"
	RoleStudent    Role = "student"
)

// Common errors
var (
	ErrSessionExpired = errors.New("session has expired")
	ErrInvalidSession = errors.New("invalid session")
	ErrSessionExists  = errors.New("active session already exists")
)

// SessionStatus represents the current state of a session
type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "active"
	SessionStatusExpired SessionStatus = "expired"
	SessionStatusRevoked SessionStatus = "revoked"
)

type User struct {
	ID                     uuid.UUID  `json:"id"`
	Email                  string     `json:"email"`
	HashedPassword         string     `json:"-"`
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	ProfilePictureURL      *string    `json:"profile_picture_url"`
	Title                  string     `json:"title"`
	Bio                    string     `json:"bio"`
	Role                   Role       `json:"role"`
	Active                 bool       `json:"active"`
	EmailVerified          bool       `json:"email_verified"`
	SubscribedToNewsletter bool       `json:"subscribed_to_newsletter"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	LastLoginAt            *time.Time `json:"last_login_at,omitempty"`
	CurrentSession         *Session   `json:"current_session,omitempty"`
	PasswordResetToken     string     `json:"-"`
	ResetTokenExpiresAt    time.Time  `json:"-"`
}

type Session struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	UserRole       Role
	Token          string        `json:"token"`
	RefreshToken   string        `json:"refresh_token"`
	Status         SessionStatus `json:"status"`
	DeviceInfo     DeviceInfo    `json:"device_info"`
	ExpiresAt      time.Time     `json:"expires_at"`
	LastActivityAt time.Time     `json:"last_activity_at"`
	CreatedAt      time.Time     `json:"created_at"`
	RevokedAt      *time.Time    `json:"revoked_at,omitempty"`
	RevokedReason  *string       `json:"revoked_reason,omitempty"`
}

type DeviceInfo struct {
	UserAgent      string `json:"user_agent"`
	IPAddress      string `json:"ip_address"`
	DeviceID       string `json:"device_id"`
	DeviceType     string `json:"device_type"`
	DeviceName     string `json:"device_name"`
	OSName         string `json:"os_name"`
	OSVersion      string `json:"os_version"`
	BrowserName    string `json:"browser_name"`
	BrowserVersion string `json:"browser_version"`
}

// Methods for Session management
func (s *Session) IsActive() bool {
	return s.Status == SessionStatusActive &&
		time.Now().Before(s.ExpiresAt) &&
		s.RevokedAt == nil
}

func (s *Session) Revoke(reason string) {
	now := time.Now()
	s.Status = SessionStatusRevoked
	s.RevokedAt = &now
	s.RevokedReason = &reason
}

func (s *Session) UpdateActivity() {
	s.LastActivityAt = time.Now()
}

// Methods for User
func (u *User) HasActiveSession() bool {
	return u.CurrentSession != nil && u.CurrentSession.IsActive()
}

func (u *User) SetCurrentSession(session *Session) {
	u.CurrentSession = session
	if session != nil {
		now := time.Now()
		u.LastLoginAt = &now
	}
}

func (r Role) String() string {
	switch r {
	case RoleSuperUser:
		return "SuperUser"
	case RoleAdmin:
		return "Admin"
	case RoleInstructor:
		return "Instructor"
	case RoleStudent:
		return "Student"
	default:
		return "Unknown"
	}
}
