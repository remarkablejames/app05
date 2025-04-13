package scheduler

import (
	"app05/internal/core/application/contracts"
	"context"
	"time"
)

type Job interface {
	Name() string
	Run(ctx context.Context) error
	Interval() time.Duration
}

type Scheduler struct {
	jobs   []Job
	logger contracts.Logger
}

func NewScheduler(logger contracts.Logger) *Scheduler {
	return &Scheduler{
		jobs:   make([]Job, 0),
		logger: logger,
	}
}

func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		go s.runJob(ctx, job)
	}
}

func (s *Scheduler) runJob(ctx context.Context, job Job) {
	ticker := time.NewTicker(job.Interval())
	defer ticker.Stop()

	s.logger.Info("Starting job", "name", job.Name(), "interval", job.Interval())

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping job", "name", job.Name())
			return
		case <-ticker.C:
			if err := job.Run(ctx); err != nil {
				s.logger.Error("Job failed", "name", job.Name(), "error", err)
			}
		}
	}
}
