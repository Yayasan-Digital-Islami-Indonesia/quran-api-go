package healthcheck

import "context"

// HealthCheckRepository defines read-only access to healthcheck data.
// Implement this interface in internal/repository/healthcheck_repository.go.
type HealthCheckRepository interface {
	HealthCheck(ctx context.Context) error
}
