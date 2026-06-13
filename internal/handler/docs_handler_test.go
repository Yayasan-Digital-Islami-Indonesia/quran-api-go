package handler_test

import (
	"net/http"
	"net/http/httptest"
	"quran-api-go/internal/handler"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestDocsHandler_Standalone verifies docs HTML doesn't use external CDN
func TestDocsHandler_Standalone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := handler.NewDocsHandler()
	r.GET("/docs", h.ServeDocs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	html := w.Body.String()

	// Should NOT contain CDN references
	if strings.Contains(html, "cdn.jsdelivr.net") {
		t.Error("Documentation HTML uses external CDN, violates standalone requirement")
	}
	if strings.Contains(html, "https://cdn.") {
		t.Error("Documentation HTML uses external CDN, violates standalone requirement")
	}

	// Should contain Scalar reference (local or embedded)
	if !strings.Contains(html, "scalar") && !strings.Contains(html, "api-reference") {
		t.Error("Documentation HTML doesn't contain Scalar references")
	}
}

func TestDocsHandler_OpenAPIRetainsLocalhostInTestMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := handler.NewDocsHandler()
	r.GET("/openapi.yaml", h.ServeOpenAPI)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	req.Host = "localhost:8080"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "http://localhost:8080") {
		t.Fatal("expected localhost URL to remain in test mode")
	}
}

func TestDocsHandler_OpenAPIUsesProductionURLInReleaseMode(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.ReleaseMode)
	t.Cleanup(func() {
		gin.SetMode(previousMode)
	})

	r := gin.New()

	h := handler.NewDocsHandler()
	r.GET("/openapi.yaml", h.ServeOpenAPI)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	req.Host = "localhost:8080"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if strings.Contains(body, "http://localhost:8080") {
		t.Fatal("expected localhost URL to be replaced in release mode")
	}

	if !strings.Contains(body, "https://quran.api.digitalislami.id") {
		t.Fatal("expected production URL in openapi output")
	}
}
