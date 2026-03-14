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
	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/internal/handler"
)

const canonicalAyahPath = "/surah/1/ayah/1"

type MockAyahService struct {
	GetBySurahFunc          func(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error)
	GetBySurahAndNumberFunc func(ctx context.Context, surahID, number int) (*ayah.Ayah, error)
	GetRandomFunc           func(ctx context.Context, surahID int) (*ayah.Ayah, error)
}

func (m *MockAyahService) GetByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahService) GetBySurah(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
	if m.GetBySurahFunc != nil {
		return m.GetBySurahFunc(ctx, surahID, from, to)
	}
	return nil, nil
}

func (m *MockAyahService) GetBySurahAndNumber(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
	if m.GetBySurahAndNumberFunc != nil {
		return m.GetBySurahAndNumberFunc(ctx, surahID, number)
	}
	return nil, nil
}

func (m *MockAyahService) GetRandom(ctx context.Context, surahID int) (*ayah.Ayah, error) {
	if m.GetRandomFunc != nil {
		return m.GetRandomFunc(ctx, surahID)
	}
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
	r.GET("/surah/:id/ayah", h.BySurah)
	r.GET("/surah/:id/ayah/:number", h.BySurahAndNumber)
	r.GET("/random", h.RandomAyah)
	return r
}

func TestBySurah(t *testing.T) {
	t.Run("Success default lang returns full surah", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahFunc: func(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
				if surahID != 1 || from != 1 || to != 7 {
					t.Fatalf("expected surahID=1 from=1 to=7, got surahID=%d from=%d to=%d", surahID, from, to)
				}

				return []ayah.Ayah{
					{
						ID:             1,
						SurahID:        1,
						NumberInSurah:  1,
						TextUthmani:    "Bismillah",
						TranslationIdo: "Dengan nama Allah",
						TranslationEn:  "In the name of Allah",
						JuzNumber:      1,
					},
					{
						ID:             2,
						SurahID:        1,
						NumberInSurah:  2,
						TextUthmani:    "Alhamdulillah",
						TranslationIdo: "Segala puji",
						TranslationEn:  "All praise",
						JuzNumber:      1,
					},
				}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{
					ID:            1,
					Number:        1,
					NameLatin:     "Al-Fatihah",
					NumberOfAyahs: 7,
				}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		data, ok := body["data"].(map[string]any)
		if !ok {
			t.Fatalf("expected data object, got %T", body["data"])
		}

		surahData, ok := data["surah"].(map[string]any)
		if !ok {
			t.Fatalf("expected surah object, got %T", data["surah"])
		}
		if surahData["id"] != float64(1) || surahData["number"] != float64(1) || surahData["name_latin"] != "Al-Fatihah" {
			t.Fatalf("unexpected surah payload: %v", surahData)
		}

		ayahsData, ok := data["ayahs"].([]any)
		if !ok || len(ayahsData) != 2 {
			t.Fatalf("expected 2 ayahs, got %v", data["ayahs"])
		}

		firstAyah, ok := ayahsData[0].(map[string]any)
		if !ok {
			t.Fatalf("expected ayah object, got %T", ayahsData[0])
		}
		if firstAyah["number"] != float64(1) || firstAyah["number_in_surah"] != float64(1) {
			t.Fatalf("unexpected ayah numbering payload: %v", firstAyah)
		}
		if firstAyah["translation"] != "Dengan nama Allah" {
			t.Fatalf("expected Indonesian translation, got %v", firstAyah["translation"])
		}
	})

	t.Run("Success english lang with explicit range", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahFunc: func(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
				if surahID != 1 || from != 2 || to != 3 {
					t.Fatalf("expected surahID=1 from=2 to=3, got surahID=%d from=%d to=%d", surahID, from, to)
				}

				return []ayah.Ayah{
					{
						ID:             2,
						SurahID:        1,
						NumberInSurah:  2,
						TextUthmani:    "Alhamdulillah",
						TranslationIdo: "Segala puji",
						TranslationEn:  "All praise",
						JuzNumber:      1,
					},
				}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{
					ID:            1,
					Number:        1,
					NameLatin:     "Al-Fatihah",
					NumberOfAyahs: 7,
				}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah?lang=en&from=2&to=3", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		data := body["data"].(map[string]any)
		ayahsData := data["ayahs"].([]any)
		firstAyah := ayahsData[0].(map[string]any)
		if firstAyah["translation"] != "All praise" {
			t.Fatalf("expected English translation, got %v", firstAyah["translation"])
		}
	})

	t.Run("Invalid lang", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah?lang=fr", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid range", func(t *testing.T) {
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return &surah.Surah{
					ID:            1,
					Number:        1,
					NameLatin:     "Al-Fatihah",
					NumberOfAyahs: 7,
				}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah?from=3", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Surah not found", func(t *testing.T) {
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return nil, domain.ErrNotFound
			},
		}

		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/999/ayah", nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})
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
}

func TestRandomAyah(t *testing.T) {
	t.Run("Success ID lang, surah_id=0 uses ayah's SurahID for surah_info", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				if surahID != 0 {
					t.Fatalf("expected surahID=0 passed to GetRandom, got %d", surahID)
				}
				return &ayah.Ayah{
					ID:             10,
					SurahID:        2,
					NumberInSurah:  5,
					TextUthmani:    "Test Uthmani",
					TranslationIdo: "Terjemah ID",
					TranslationEn:  "Translation EN",
					JuzNumber:      1,
					SajdaType:      nil,
					RevelationType: nil,
				}, nil
			},
		}

		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				if id != 2 {
					t.Fatalf("expected GetByID called with id=2 (ayah's SurahID), got %d", id)
				}
				return &surah.Surah{
					ID:        2,
					NameLatin: "Al-Baqarah",
				}, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/random", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["id"] != float64(10) {
			t.Errorf("expected id=10")
		}
		if response["surah_id"] != float64(2) {
			t.Errorf("expected surah_id=2")
		}
		if response["number_in_surah"] != float64(5) {
			t.Errorf("expected number_in_surah=5")
		}
		if response["text_uthmani"] != "Test Uthmani" {
			t.Errorf("expected text_uthmani")
		}
		if response["translation"] != "Terjemah ID" {
			t.Errorf("expected translation (id), got %v", response["translation"])
		}

		surahInfo, ok := response["surah_info"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected surah_info map")
		}
		if surahInfo["id"] != float64(2) || surahInfo["name_latin"] != "Al-Baqarah" {
			t.Errorf("invalid surah info")
		}
	})

	t.Run("Success EN lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				if surahID != 1 {
					t.Fatalf("expected surahID=1 passed to GetRandom, got %d", surahID)
				}
				return &ayah.Ayah{
					ID:             11,
					SurahID:        1,
					NumberInSurah:  1,
					TextUthmani:    "Bismillah",
					TranslationIdo: "Dengan nama Allah",
					TranslationEn:  "In the name of Allah",
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

		req, _ := http.NewRequest("GET", "/random?lang=en&surah_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["translation"] != "In the name of Allah" {
			t.Errorf("expected english translation, got %v", response["translation"])
		}
	})

	t.Run("Invalid Lang", func(t *testing.T) {
		h := handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/random?lang=fr", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Ayah Not Found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/random?surah_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Surah Not Found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return &ayah.Ayah{ID: 1, SurahID: 999}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return nil, nil
			},
		}

		h := handler.NewAyahHandler(mockAyahService, mockSurahService)
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/random", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Ayah Service Error", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return nil, errors.New("db error")
			},
		}

		h := handler.NewAyahHandler(mockAyahService, &MockSurahService{})
		r := setupRouter(h)

		req, _ := http.NewRequest("GET", "/random?surah_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
