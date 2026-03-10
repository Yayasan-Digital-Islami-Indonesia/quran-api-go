package service

import (
	"context"
	"fmt"
	"quran-api-go/internal/domain/healthcheck"
	"time"
)

type healthCheckService struct {
	repo healthcheck.HealthCheckRepository
}

func NewHealthCheckService(repo healthcheck.HealthCheckRepository) healthcheck.HealthCheckService {
	return &healthCheckService{
		repo,
	}
}

func (s *healthCheckService) HealthCheck(ctx context.Context) (healthcheck.HealthCheck, error) {
	return healthcheck.HealthCheck{
		Status:    "OK",
		Version:   "-",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *healthCheckService) ReadyCheck(ctx context.Context) (healthcheck.HealthCheck, error) {
	err := s.repo.HealthCheck(ctx)
	if err != nil {
		return healthcheck.HealthCheck{}, fmt.Errorf("error check db status, err: %v", err)
	}

	return healthcheck.HealthCheck{
		DBStatus:  "OK",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
