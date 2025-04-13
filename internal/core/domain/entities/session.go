package entities

import "time"

// SessionManager handles business rules for session management
type SessionManager struct {
	MaxConcurrentSessions int
	SessionDuration       time.Duration
	RefreshTokenDuration  time.Duration
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		MaxConcurrentSessions: 1,                  // Only one active session allowed
		SessionDuration:       24 * time.Hour,     // Token valid for 24 hours
		RefreshTokenDuration:  7 * 24 * time.Hour, // Refresh token valid for 7 days
	}
}

func (sm *SessionManager) ValidateNewSession(user *User) error {
	if user.HasActiveSession() {
		// If user has an active session, it needs to be revoked before creating a new one
		return ErrSessionExists
	}
	return nil
}
