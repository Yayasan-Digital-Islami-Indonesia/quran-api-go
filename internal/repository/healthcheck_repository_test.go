package repository_test

import (
	"context"
	"quran-api-go/internal/repository"
	"testing"
)

func TestHealthCheckRepository_HealthCheck_Success(t *testing.T) {
	db := setupTestDB(t, "", "")
	repo := repository.NewHealthCheckRepository(db)

	ctx := context.Background()

	err := repo.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
