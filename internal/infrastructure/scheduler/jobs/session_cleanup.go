package jobs

import (
	"app05/internal/core/application/contracts"
	"app05/internal/core/domain/repositories"
	"context"
	"time"
)

type SessionCleanupJob struct {
	sessionRepo repositories.SessionRepository
	logger      contracts.Logger
	interval    time.Duration
}

func NewSessionCleanupJob(
	sessionRepo repositories.SessionRepository,
	logger contracts.Logger,
	interval time.Duration,
) *SessionCleanupJob {
	return &SessionCleanupJob{
		sessionRepo: sessionRepo,
		logger:      logger,
		interval:    interval,
	}
}

func (j *SessionCleanupJob) Name() string {
	return "session_cleanup"
}

func (j *SessionCleanupJob) Interval() time.Duration {
	return j.interval
}

func (j *SessionCleanupJob) Run(ctx context.Context) error {
	cutoffTime := time.Now().Add(-7 * 24 * time.Hour)
	deleted, err := j.sessionRepo.DeleteRevokedSessions(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if deleted > 0 {
		j.logger.Info("Cleaned up revoked sessions", "count", deleted)
	}
	return nil
}
