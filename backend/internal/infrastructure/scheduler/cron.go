package scheduler

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"
	"haberbot/internal/usecase"
)

// Scheduler manages scheduled tasks for the HaberBot application.
type Scheduler struct {
	cron       *cron.Cron
	pipelineUC *usecase.DailyPipelineUseCase
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(pipelineUC *usecase.DailyPipelineUseCase) *Scheduler {
	return &Scheduler{
		pipelineUC: pipelineUC,
	}
}

// Start starts the background scheduler with the given schedules for fetching and processing.
func (s *Scheduler) Start(fetchSchedule, processSchedule string) error {
	s.cron = cron.New() // robfig/cron/v3 uses standard 5-field cron by default (minute, hour, dom, month, dow)

	// Register Fetch task
	_, err := s.cron.AddFunc(fetchSchedule, func() {
		slog.Info("cron: starting scheduled fetch job")
		fetched, err := s.pipelineUC.ExecuteFetch(context.Background())
		if err != nil {
			slog.Error("cron: scheduled fetch job failed", "error", err)
		} else {
			slog.Info("cron: scheduled fetch job completed successfully", "fetched", fetched)
		}
	})
	if err != nil {
		return err
	}

	// Register Process task
	_, err = s.cron.AddFunc(processSchedule, func() {
		slog.Info("cron: starting scheduled process job")
		processed, err := s.pipelineUC.ExecuteProcess(context.Background())
		if err != nil {
			slog.Error("cron: scheduled process job failed", "error", err)
		} else {
			slog.Info("cron: scheduled process job completed successfully", "processed", processed)
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	slog.Info("cron: scheduler started", "fetch_schedule", fetchSchedule, "process_schedule", processSchedule)
	return nil
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		slog.Info("cron: scheduler stopped")
	}
}
