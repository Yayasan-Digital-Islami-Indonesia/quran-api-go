package repository

import (
	"context"
	"database/sql"
	"quran-api-go/internal/domain/healthcheck"
)

type HealthCheckRepository struct {
	db *sql.DB
}

func NewHealthCheckRepository(db *sql.DB) healthcheck.HealthCheckRepository {
	return &HealthCheckRepository{
		db: db,
	}
}

func (h *HealthCheckRepository) HealthCheck(ctx context.Context) error {
	return h.db.PingContext(ctx)
}
