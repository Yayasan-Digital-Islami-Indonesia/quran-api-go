package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/internal/handler"
)

const (
	canonicalAyahPath = "/surah/1/ayah/1"
	legacyAyatPath    = "/surah/1/ayat/1"
)

type MockAyahService struct {
	GetBySurahAndNumberFunc func(ctx context.Context, surahID, number int) (*ayah.Ayah, error)
}

func (m *MockAyahService) GetByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahService) GetBySurah(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahService) GetBySurahAndNumber(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
	if m.GetBySurahAndNumberFunc != nil {
		return m.GetBySurahAndNumberFunc(ctx, surahID, number)
	}
	return nil, nil
}

func (m *MockAyahService) GetRandom(ctx context.Context, surahID int) (*ayah.Ayah, error) {
	return nil, nil
}

type MockSurahService struct {
	GetByIDFunc func(ctx context.Context, id int) (*surah.Surah, error)
}

func (m *MockSurahService) GetAll(ctx context.Context) ([]surah.Surah, error) {
	return nil, nil
}

func (m *MockSurahService) GetByID(ctx context.Context, id int) (*surah.Surah, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func setupRouter(h *handler.AyahHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/surah/:id/ayah/:number", h.BySurahAndNumber)
	r.GET("/surah/:id/ayat/:number", h.BySurahAndNumber)
	return r
}

func TestBySurahAndNumber(t *testing.T) {
	t.Run("Success ID lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return &ayah.Ayah{
					ID:             1,
					SurahID:        1,
					NumberInSurah:  1,
					TextUthmani:    "Bismillah",
					TranslationIdo: "Dengan nama Allah",
					TranslationEn:  "In the name of Allah",
					JuzNumber:      1,
				}, nil
			},
		}

		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{
					ID:        1,
					NameLatin: "Al-Fatihah",
				}, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", canonicalAyahPath, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["id"] != float64(1) {
			t.Errorf("expected id=1")
		}
		if response["number_in_surah"] != float64(1) {
			t.Errorf("expected number_in_surah=1")
		}
		if response["text_uthmani"] != "Bismillah" {
			t.Errorf("expected text_uthmani")
		}
		if response["translation"] != "Dengan nama Allah" {
			t.Errorf("expected translation")
		}

		surahInfo, ok := response["surah_info"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected surah_info map")
		}
		if surahInfo["id"] != float64(1) || surahInfo["name_latin"] != "Al-Fatihah" {
			t.Errorf("invalid surah info")
		}
	})

	t.Run("Success EN lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return &ayah.Ayah{
					ID:             2,
					SurahID:        1,
					NumberInSurah:  2,
					TextUthmani:    "Alhamdulillah",
					TranslationIdo: "Segala puji",
					TranslationEn:  "All praise",
				}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{ID: 1, NameLatin: "Al-Fatihah"}, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/surah/1/ayah/2?lang=en", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["translation"] != "All praise" {
			t.Errorf("expected english translation, got %v", response["translation"])
		}
	})

	t.Run("Invalid Lang", func(t *testing.T) {
		h := handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/surah/1/ayah/1?lang=fr", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Ayah Not Found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/surah/1/ayah/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Surah Not Found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return &ayah.Ayah{ID: 1}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return nil, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/surah/999/ayah/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Ayah Service Error", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, errors.New("db error")
			},
		}

		h := handler.NewAyahHandler(mockAyahService, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", canonicalAyahPath, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})

	t.Run("Legacy Ayat Alias", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return &ayah.Ayah{
					ID:             1,
					SurahID:        1,
					NumberInSurah:  1,
					TextUthmani:    "Bismillah",
					TranslationIdo: "Dengan nama Allah",
				}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{ID: 1, NameLatin: "Al-Fatihah"}, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", legacyAyatPath, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
