package service

import (
	"context"
	"log"
	"time"

	"qflow/internal/domain"
)

type OTPCleanupJob struct {
	authRepo domain.AuthRepository
	interval time.Duration
	logger   *log.Logger
}

func NewOTPCleanupJob(authRepo domain.AuthRepository, interval time.Duration, logger *log.Logger) *OTPCleanupJob {
	return &OTPCleanupJob{
		authRepo: authRepo,
		interval: interval,
		logger:   logger,
	}
}

func (j *OTPCleanupJob) Start(ctx context.Context) {
	if j.interval <= 0 {
		j.logf("OTP cleanup job disabled")
		return
	}

	go j.Run(ctx)
}

func (j *OTPCleanupJob) Run(ctx context.Context) {
	j.cleanup(time.Now())

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			j.logf("OTP cleanup job stopped")
			return
		case now := <-ticker.C:
			j.cleanup(now)
		}
	}
}

func (j *OTPCleanupJob) CleanupExpiredOTPs(now time.Time) (int64, error) {
	return j.authRepo.DeleteExpiredOTPs(now)
}

func (j *OTPCleanupJob) cleanup(now time.Time) {
	deleted, err := j.CleanupExpiredOTPs(now)
	if err != nil {
		j.logf("failed to cleanup expired OTPs: %v", err)
		return
	}
	if deleted > 0 {
		j.logf("deleted %d expired OTPs", deleted)
	}
}

func (j *OTPCleanupJob) logf(format string, args ...any) {
	if j.logger != nil {
		j.logger.Printf(format, args...)
		return
	}
	log.Printf(format, args...)
}
