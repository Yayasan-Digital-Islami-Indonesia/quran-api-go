package healthcheck

import "context"

// HealthCheckService defines the business operations for healthcheck data.
// Implement this interface in internal/service/healthcheck_service.go.
type HealthCheckService interface {
	HealthCheck(ctx context.Context) (HealthCheck, error)
	ReadyCheck(ctx context.Context) (HealthCheck, error)
}
