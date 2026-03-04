package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/internal/handler"
)

// mockSurahService is a test double for surah.SurahService.
type mockSurahService struct {
	getAllFn  func(ctx context.Context) ([]surah.Surah, error)
	getByIDFn func(ctx context.Context, id int) (*surah.Surah, error)
}

func (m *mockSurahService) GetAll(ctx context.Context) ([]surah.Surah, error) {
	return m.getAllFn(ctx)
}

func (m *mockSurahService) GetByID(ctx context.Context, id int) (*surah.Surah, error) {
	return m.getByIDFn(ctx, id)
}

func newTestRouter(h *handler.SurahHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/surah", h.List)
	r.GET("/surah/:id", h.Detail)
	return r
}

func TestSurahHandler_List_OK(t *testing.T) {
	svc := &mockSurahService{
		getAllFn: func(_ context.Context) ([]surah.Surah, error) {
			return []surah.Surah{
				{ID: 1, Number: 1, NameArabic: "الفاتحة", NameLatin: "Al-Fatihah", NameTransliteration: "Al-Fatihah", NumberOfAyahs: 7, RevelationType: "Meccan"},
			}, nil
		},
	}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	data, ok := body["data"].([]any)
	if !ok || len(data) != 1 {
		t.Fatalf("expected data array with 1 element, got %v", body["data"])
	}
}

func TestSurahHandler_List_InternalError(t *testing.T) {
	svc := &mockSurahService{
		getAllFn: func(_ context.Context) ([]surah.Surah, error) {
			return nil, errors.New("db error")
		},
	}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestSurahHandler_Detail_OK(t *testing.T) {
	svc := &mockSurahService{
		getByIDFn: func(_ context.Context, id int) (*surah.Surah, error) {
			return &surah.Surah{ID: id, Number: 1, NameLatin: "Al-Fatihah", NumberOfAyahs: 7, RevelationType: "Meccan"}, nil
		},
	}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSurahHandler_Detail_NotFound(t *testing.T) {
	svc := &mockSurahService{
		getByIDFn: func(_ context.Context, _ int) (*surah.Surah, error) {
			return nil, domain.ErrNotFound
		},
	}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/999", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSurahHandler_Detail_InvalidID(t *testing.T) {
	svc := &mockSurahService{}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSurahHandler_Detail_InternalError(t *testing.T) {
	svc := &mockSurahService{
		getByIDFn: func(_ context.Context, _ int) (*surah.Surah, error) {
			return nil, errors.New("db error")
		},
	}
	r := newTestRouter(handler.NewSurahHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
