package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quran-api-go/internal/domain/healthcheck"
	"quran-api-go/internal/handler"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// mockHealthCheckService is a test double for healthcheck.HealthCheckService.
type mockHealthCheckService struct {
	healthCheck func(ctx context.Context) (healthcheck.HealthCheck, error)
	readyCheck  func(ctx context.Context) (healthcheck.HealthCheck, error)
}

func (m *mockHealthCheckService) HealthCheck(ctx context.Context) (healthcheck.HealthCheck, error) {
	return m.healthCheck(ctx)
}

func (m *mockHealthCheckService) ReadyCheck(ctx context.Context) (healthcheck.HealthCheck, error) {
	return m.readyCheck(ctx)
}

func newHealthCheckTestRouter(h *handler.HealthCheckHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", h.HealthCheck)
	r.GET("/health/ready", h.ReadyCheck)
	return r
}

func TestHealthCheckHandler_HealthCheck_Success(t *testing.T) {
	svc := &mockHealthCheckService{
		healthCheck: func(_ context.Context) (healthcheck.HealthCheck, error) {
			return healthcheck.HealthCheck{
				Status:    "OK",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Version:   "v0.0.1",
			}, nil
		},
	}
	r := newHealthCheckTestRouter(handler.NewHealthCheckHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	status := body["status"]
	if status != "OK" {
		t.Fatalf("expected status field, got %v", body["status"])
	}
}

func TestHealthCheckHandler_ReadyCheck_Success(t *testing.T) {
	svc := &mockHealthCheckService{
		readyCheck: func(_ context.Context) (healthcheck.HealthCheck, error) {
			return healthcheck.HealthCheck{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				DBStatus:  "OK",
			}, nil
		},
	}
	r := newHealthCheckTestRouter(handler.NewHealthCheckHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health/ready", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	status := body["db_status"]
	if status != "OK" {
		t.Fatalf("expected status field, got %v", body["status"])
	}
}
