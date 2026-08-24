package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/search"
	"quran-api-go/internal/handler"
)

type mockSearchServiceWithErrors struct {
	searchFn func(ctx context.Context, p search.Params) ([]search.Result, int, error)
}

func (m *mockSearchServiceWithErrors) Search(ctx context.Context, p search.Params) ([]search.Result, int, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, p)
	}
	return []search.Result{}, 0, nil
}

func setupSearchRouter(h *handler.SearchHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/search", h.Search)
	return r
}

func TestSearchHandler_EmptyQuery(t *testing.T) {
	svc := &mockSearchServiceWithErrors{}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSearchHandler_InvalidLang(t *testing.T) {
	svc := &mockSearchServiceWithErrors{}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=test&lang=fr", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSearchHandler_Success(t *testing.T) {
	svc := &mockSearchServiceWithErrors{
		searchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Query != "bismillah" {
				t.Fatalf("expected query='bismillah', got %s", p.Query)
			}
			if p.Lang != "id" {
				t.Fatalf("expected lang='id', got %s", p.Lang)
			}
			return []search.Result{
				{ID: 1, SurahID: 1, NumberInSurah: 1, TextUthmani: "Bismillah", Translation: "Dengan nama Allah", JuzNumber: 1},
			}, 1, nil
		},
	}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=bismillah", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}

	if data["query"] != "bismillah" {
		t.Fatalf("expected query='bismillah', got %v", data["query"])
	}
	if data["total"] != float64(1) {
		t.Fatalf("expected total=1, got %v", data["total"])
	}
}

func TestSearchHandler_WithFilters(t *testing.T) {
	svc := &mockSearchServiceWithErrors{
		searchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.SurahID != 2 {
				t.Fatalf("expected surahID=2, got %d", p.SurahID)
			}
			if p.Juz != 3 {
				t.Fatalf("expected juz=3, got %d", p.Juz)
			}
			return []search.Result{}, 0, nil
		},
	}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=test&surah_id=2&juz=3", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSearchHandler_ServiceError(t *testing.T) {
	svc := &mockSearchServiceWithErrors{
		searchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			return nil, 0, errors.New("db error")
		},
	}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=test", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestSearchHandler_EmptyResults(t *testing.T) {
	svc := &mockSearchServiceWithErrors{
		searchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			return []search.Result{}, 0, nil
		},
	}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=nonexistent", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)

	results := data["results"].([]any)
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
	if data["total"] != float64(0) {
		t.Fatalf("expected total=0, got %v", data["total"])
	}
}

func TestSearchHandler_EnglishLang(t *testing.T) {
	svc := &mockSearchServiceWithErrors{
		searchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Lang != "en" {
				t.Fatalf("expected lang='en', got %s", p.Lang)
			}
			return []search.Result{
				{ID: 1, Translation: "In the name of Allah"},
			}, 1, nil
		},
	}
	h := handler.NewSearchHandler(svc)
	r := setupSearchRouter(h)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=bismillah&lang=en", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
