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
	"quran-api-go/internal/domain/juz"
	"quran-api-go/internal/handler"
)

type mockJuzService struct {
	getAllFn         func(ctx context.Context) ([]juz.Juz, error)
	getByNumberFn    func(ctx context.Context, number int) (*juz.Juz, error)
	getAyahsByJuzFn  func(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error)
	getSurahsByJuzFn func(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error)
}

func (m *mockJuzService) GetAll(ctx context.Context) ([]juz.Juz, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *mockJuzService) GetByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	if m.getByNumberFn != nil {
		return m.getByNumberFn(ctx, number)
	}
	return nil, nil
}

func (m *mockJuzService) GetAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	if m.getAyahsByJuzFn != nil {
		return m.getAyahsByJuzFn(ctx, juzNumber, limit, offset)
	}
	return nil, nil
}

func (m *mockJuzService) GetSurahsByJuz(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error) {
	if m.getSurahsByJuzFn != nil {
		return m.getSurahsByJuzFn(ctx, juzNumber)
	}
	return nil, nil
}

func setupJuzRouter(h *handler.JuzHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/juz", h.List)
	r.GET("/juz/:number", h.Detail)
	r.GET("/juz/:number/ayah", h.Ayahs)
	r.GET("/juz/:number/surah", h.Surahs)
	return r
}

func TestJuzHandler_List_Success(t *testing.T) {
	svc := &mockJuzService{
		getAllFn: func(_ context.Context) ([]juz.Juz, error) {
			return []juz.Juz{
				{ID: 1, JuzNumber: 1, FirstAyahID: 1, LastAyahID: 141, TotalAyahs: 141},
				{ID: 2, JuzNumber: 2, FirstAyahID: 142, LastAyahID: 252, TotalAyahs: 111},
			}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	data, ok := body["data"].([]any)
	if !ok || len(data) != 2 {
		t.Fatalf("expected data array with 2 elements, got %v", body["data"])
	}
}

func TestJuzHandler_List_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getAllFn: func(_ context.Context) ([]juz.Juz, error) {
			return nil, errors.New("db error")
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_Success(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: number, FirstAyahID: 1, LastAyahID: 141, TotalAyahs: 141}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_InvalidNumber(t *testing.T) {
	svc := &mockJuzService{}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_NotFound(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return nil, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/99", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return nil, errors.New("db error")
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_Success(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 7}, nil
		},
		getAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			return []juz.JuzAyah{
				{AyahID: 1, SurahID: 1, SurahNameLatin: "Al-Fatihah", NumberInSurah: 1, TextUthmani: "Bismillah", TranslationIdo: "Dengan nama Allah", TranslationEn: "In the name of Allah", JuzNumber: 1},
			}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/ayah", nil))

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

	if data["juz"] == nil {
		t.Errorf("response missing 'juz' wrapper")
	}
	if data["ayahs"] == nil {
		t.Errorf("response missing 'ayahs' array")
	}
}

func TestJuzHandler_Ayahs_InvalidNumber(t *testing.T) {
	svc := &mockJuzService{}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/abc/ayah", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_InvalidLang(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 7}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/ayah?lang=fr", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_NotFound(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return nil, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/99/ayah", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 7}, nil
		},
		getAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			return nil, errors.New("db error")
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/ayah", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_ErrNotFound(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 7}, nil
		},
		getAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			return nil, domain.ErrNotFound
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/ayah", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestJuzHandler_Ayahs_EnglishLang(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 7}, nil
		},
		getAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			return []juz.JuzAyah{
				{AyahID: 1, SurahID: 1, SurahNameLatin: "Al-Fatihah", NumberInSurah: 1, TextUthmani: "Bismillah", TranslationIdo: "Dengan nama Allah", TranslationEn: "In the name of Allah", JuzNumber: 1},
			}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/ayah?lang=en", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	data := body["data"].(map[string]any)
	ayahs := data["ayahs"].([]any)
	first := ayahs[0].(map[string]any)

	if first["translation"] != "In the name of Allah" {
		t.Fatalf("expected English translation, got %v", first["translation"])
	}
}

func TestJuzHandler_Surahs_Success(t *testing.T) {
	svc := &mockJuzService{
		getSurahsByJuzFn: func(_ context.Context, juzNumber int) ([]juz.JuzSurah, error) {
			return []juz.JuzSurah{
				{ID: 1, Number: 1, NameArabic: "الفاتحة", NameLatin: "Al-Fatihah", NameTransliteration: "Al-Fatihah", NumberOfAyahs: 7, RevelationType: "meccan"},
			}, nil
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/surah", nil))

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

func TestJuzHandler_Surahs_InvalidNumber(t *testing.T) {
	svc := &mockJuzService{}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/abc/surah", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_OutOfRange(t *testing.T) {
	svc := &mockJuzService{}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/31/surah", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getSurahsByJuzFn: func(_ context.Context, juzNumber int) ([]juz.JuzSurah, error) {
			return nil, errors.New("db error")
		},
	}
	r := setupJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/surah", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
