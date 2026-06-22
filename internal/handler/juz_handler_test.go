package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/juz"
	"quran-api-go/internal/handler"
	"quran-api-go/internal/repository"
	"quran-api-go/internal/service"
)

// TestJuzHandler_ListResponse_IncludesTotalAyahs verifies list response has total_ayahs
func TestJuzHandler_ListResponse_IncludesTotalAyahs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Setup schema - need both tables for the JOIN
	db.ExecContext(context.Background(), `
		CREATE TABLE juzs (id INTEGER PRIMARY KEY, juz_number INTEGER, first_ayah_id INTEGER, last_ayah_id INTEGER);
		CREATE TABLE ayahs (id INTEGER PRIMARY KEY, juz_number INTEGER);

		INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id) VALUES (1, 1, 1, 7);
		-- Insert 7 ayahs for juz 1
		INSERT INTO ayahs (id, juz_number) VALUES (1, 1), (2, 1), (3, 1), (4, 1), (5, 1), (6, 1), (7, 1);
	`)

	repo := repository.NewJuzRepository(db)
	svc := service.NewJuzService(repo)
	h := handler.NewJuzHandler(svc)
	r.GET("/juz", h.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/juz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	data, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("expected data array, got %T", body["data"])
	}

	if len(data) == 0 {
		t.Fatal("expected at least one juz")
	}

	juz := data[0].(map[string]any)
	if juz["total_ayahs"] == nil {
		t.Errorf("Juz list response missing 'total_ayahs' field. Got keys: %v", getMapKeys(juz))
	}

	// Verify the count is correct
	totalAyahs := int(juz["total_ayahs"].(float64))
	if totalAyahs != 7 {
		t.Errorf("Expected total_ayahs=7, got %d", totalAyahs)
	}
}

// TestJuzHandler_AyahsResponse_WrappedStructure verifies ayahs response has correct structure
func TestJuzHandler_AyahsResponse_WrappedStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.ExecContext(context.Background(), `
		CREATE TABLE surahs (id INTEGER PRIMARY KEY, name_latin TEXT);
		CREATE TABLE ayahs (id INTEGER PRIMARY KEY, surah_id INTEGER, number_in_surah INTEGER, text_uthmani TEXT, translation_indo TEXT, translation_en TEXT, juz_number INTEGER);
		CREATE TABLE juzs (id INTEGER PRIMARY KEY, juz_number INTEGER, first_ayah_id INTEGER, last_ayah_id INTEGER, total_ayahs INTEGER);

		INSERT INTO surahs (id, name_latin) VALUES (1, 'Al-Fatihah');
		INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id, total_ayahs) VALUES (1, 1, 1, 7, 7);
		INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number)
		VALUES (1, 1, 1, 'bismillah', 'bismillah', 'bismillah', 1);
	`)

	repo := repository.NewJuzRepository(db)
	svc := service.NewJuzService(repo)
	h := handler.NewJuzHandler(svc)
	r.GET("/juz/:number/ayah", h.Ayahs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/juz/1/ayah", nil)
	r.ServeHTTP(w, req)

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
		t.Errorf("Response missing 'juz' wrapper. Got keys: %v", getMapKeys(data))
	}
	if data["ayahs"] == nil {
		t.Errorf("Response missing 'ayahs' array. Got keys: %v", getMapKeys(data))
	}
}

func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

type mockJuzService struct {
	getAllFn       func(ctx context.Context) ([]juz.Juz, error)
	getByNumberFn  func(ctx context.Context, number int) (*juz.Juz, error)
	getAyahsByJuz  func(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error)
	getSurahsByJuz func(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error)
}

func (m *mockJuzService) GetAll(ctx context.Context) ([]juz.Juz, error) {
	return m.getAllFn(ctx)
}

func (m *mockJuzService) GetByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	return m.getByNumberFn(ctx, number)
}

func (m *mockJuzService) GetAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	return m.getAyahsByJuz(ctx, juzNumber, limit, offset)
}

func (m *mockJuzService) GetSurahsByJuz(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error) {
	return m.getSurahsByJuz(ctx, juzNumber)
}

func newJuzRouter(h *handler.JuzHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/juz/:number", h.Detail)
	r.GET("/juz/:number/surah", h.Surahs)
	return r
}

func TestJuzHandler_Detail_OK(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return &juz.Juz{ID: number, JuzNumber: number, TotalAyahs: 148}, nil
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_InvalidNumber(t *testing.T) {
	r := newJuzRouter(handler.NewJuzHandler(&mockJuzService{}))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_NotFound(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, _ int) (*juz.Juz, error) {
			return nil, nil
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/31", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestJuzHandler_Detail_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getByNumberFn: func(_ context.Context, _ int) (*juz.Juz, error) {
			return nil, errors.New("db error")
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_OK(t *testing.T) {
	svc := &mockJuzService{
		getSurahsByJuz: func(_ context.Context, _ int) ([]juz.JuzSurah, error) {
			return []juz.JuzSurah{{ID: 1, Number: 1, NameLatin: "Al-Fatihah"}}, nil
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/surah", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_InvalidNumber(t *testing.T) {
	r := newJuzRouter(handler.NewJuzHandler(&mockJuzService{}))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/99/surah", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_NotFound(t *testing.T) {
	svc := &mockJuzService{
		getSurahsByJuz: func(_ context.Context, _ int) ([]juz.JuzSurah, error) {
			return nil, nil
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/surah", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestJuzHandler_Surahs_InternalError(t *testing.T) {
	svc := &mockJuzService{
		getSurahsByJuz: func(_ context.Context, _ int) ([]juz.JuzSurah, error) {
			return nil, errors.New("db error")
		},
	}
	r := newJuzRouter(handler.NewJuzHandler(svc))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/juz/1/surah", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
