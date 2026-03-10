package service_test

import (
	"context"
	"quran-api-go/internal/service"
	"testing"
)

// mockHealthCheckRepository is a test double for healthcheck.HealthCheckRepository.
type mockHealthCheckRepository struct {
	healthCheck func(ctx context.Context) error
}

func (m *mockHealthCheckRepository) HealthCheck(ctx context.Context) error {
	return m.healthCheck(ctx)
}

func TestHealthCheckService(t *testing.T) {
	service := service.NewHealthCheckService(nil)
	health, err := service.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if health.Status != "OK" {
		t.Fatalf("expected status OK, got %s", health.Status)
	}
}

func TestReadyCheckService(t *testing.T) {
	repo := &mockHealthCheckRepository{
		healthCheck: func(_ context.Context) error {
			return nil
		},
	}
	service := service.NewHealthCheckService(repo)
	ready, err := service.ReadyCheck(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ready.DBStatus != "OK" {
		t.Fatalf("expected status OK, got %s", ready.DBStatus)
	}
}
