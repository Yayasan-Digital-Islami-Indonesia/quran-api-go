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

const (
	canonicalAyahPath       = "/surah/1/ayah/1"
	canonicalGlobalAyahPath = "/ayah/1"
)

type MockAyahService struct {
	GetByIDFunc             func(ctx context.Context, id int) (*ayah.Ayah, error)
	GetBySurahFunc          func(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error)
	GetBySurahAndNumberFunc func(ctx context.Context, surahID, number int) (*ayah.Ayah, error)
	GetRandomFunc           func(ctx context.Context, surahID int) (*ayah.Ayah, error)
}

func (m *MockAyahService) GetByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

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
	r.GET("/ayah/:id", h.Detail)
	r.GET("/surah/:id/ayah", h.BySurah)
	r.GET("/surah/:id/ayah/:number", h.BySurahAndNumber)
	r.GET("/random", h.RandomAyah)
	return r
}

func decodeBody(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return decoded
}

func decodeData(t *testing.T, body []byte) map[string]any {
	t.Helper()

	decoded := decodeBody(t, body)
	data, ok := decoded["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", decoded["data"])
	}

	return data
}

func TestAyahHandler_BySurah(t *testing.T) {
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

		data := decodeData(t, w.Body.Bytes())

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

		data := decodeData(t, w.Body.Bytes())
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

func TestAyahHandler_Detail(t *testing.T) {
	t.Run("Success default lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetByIDFunc: func(ctx context.Context, id int) (*ayah.Ayah, error) {
				if id != 1 {
					t.Fatalf("expected ayah id 1, got %d", id)
				}

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
				if id != 1 {
					t.Fatalf("expected surah id 1, got %d", id)
				}

				return &surah.Surah{ID: 1, NameLatin: "Al-Fatihah"}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, canonicalGlobalAyahPath, nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["id"] != float64(1) || data["number"] != float64(1) {
			t.Fatalf("unexpected ayah identifiers: %v", data)
		}
		if data["surah_id"] != float64(1) || data["number_in_surah"] != float64(1) {
			t.Fatalf("unexpected surah numbering: %v", data)
		}
		if data["translation"] != "Dengan nama Allah" {
			t.Fatalf("expected Indonesian translation, got %v", data["translation"])
		}

		surahInfo, ok := data["surah_info"].(map[string]any)
		if !ok {
			t.Fatalf("expected surah_info object, got %T", data["surah_info"])
		}
		if surahInfo["id"] != float64(1) || surahInfo["name_latin"] != "Al-Fatihah" {
			t.Fatalf("unexpected surah info: %v", surahInfo)
		}
	})

	t.Run("Success english lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetByIDFunc: func(ctx context.Context, id int) (*ayah.Ayah, error) {
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

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ayah/2?lang=en", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["translation"] != "All praise" {
			t.Fatalf("expected english translation, got %v", data["translation"])
		}
	})

	t.Run("Invalid ayah id", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ayah/abc", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid lang", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ayah/1?lang=fr", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Ayah not found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetByIDFunc: func(ctx context.Context, id int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ayah/999", nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Surah not found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetByIDFunc: func(ctx context.Context, id int) (*ayah.Ayah, error) {
				return &ayah.Ayah{ID: 1, SurahID: 999}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return nil, domain.ErrNotFound
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, canonicalGlobalAyahPath, nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})
}

func TestAyahHandler_BySurahAndNumber(t *testing.T) {
	t.Run("Success default lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				if surahID != 1 || number != 1 {
					t.Fatalf("unexpected arguments: surahID=%d number=%d", surahID, number)
				}

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
				return &surah.Surah{ID: 1, NameLatin: "Al-Fatihah"}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, canonicalAyahPath, nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["id"] != float64(1) || data["number"] != float64(1) {
			t.Fatalf("unexpected ayah identifiers: %v", data)
		}
		if data["translation"] != "Dengan nama Allah" {
			t.Fatalf("expected translation, got %v", data["translation"])
		}
	})

	t.Run("Success english lang", func(t *testing.T) {
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

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah/2?lang=en", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["translation"] != "All praise" {
			t.Fatalf("expected english translation, got %v", data["translation"])
		}
	})

	t.Run("Invalid ayah number", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah/abc", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid lang", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah/1?lang=fr", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Ayah not found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/surah/1/ayah/999", nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Ayah service error", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, errors.New("db error")
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, canonicalAyahPath, nil))

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", w.Code)
		}
	})
}

func TestAyahHandler_RandomAyah(t *testing.T) {
	t.Run("Success default lang uses ayah surah info", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				if surahID != 0 {
					t.Fatalf("expected surahID=0, got %d", surahID)
				}

				return &ayah.Ayah{
					ID:             10,
					SurahID:        2,
					NumberInSurah:  5,
					TextUthmani:    "Test Uthmani",
					TranslationIdo: "Terjemah ID",
					TranslationEn:  "Translation EN",
					JuzNumber:      1,
				}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				if id != 2 {
					t.Fatalf("expected surah id 2, got %d", id)
				}

				return &surah.Surah{ID: 2, NameLatin: "Al-Baqarah"}, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["id"] != float64(10) || data["number"] != float64(10) {
			t.Fatalf("unexpected ayah identifiers: %v", data)
		}
		if data["surah_id"] != float64(2) || data["translation"] != "Terjemah ID" {
			t.Fatalf("unexpected ayah detail: %v", data)
		}

		surahInfo, ok := data["surah_info"].(map[string]any)
		if !ok {
			t.Fatalf("expected surah_info object, got %T", data["surah_info"])
		}
		if surahInfo["id"] != float64(2) || surahInfo["name_latin"] != "Al-Baqarah" {
			t.Fatalf("unexpected surah info: %v", surahInfo)
		}
	})

	t.Run("Success english lang", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				if surahID != 1 {
					t.Fatalf("expected surahID=1, got %d", surahID)
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

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random?lang=en&surah_id=1", nil))

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		data := decodeData(t, w.Body.Bytes())
		if data["translation"] != "In the name of Allah" {
			t.Fatalf("expected english translation, got %v", data["translation"])
		}
	})

	t.Run("Invalid lang", func(t *testing.T) {
		r := setupRouter(handler.NewAyahHandler(&MockAyahService{}, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random?lang=fr", nil))

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("Ayah not found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random?surah_id=1", nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Surah not found", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return &ayah.Ayah{ID: 1, SurahID: 999}, nil
			},
		}
		mockSurahService := &MockSurahService{
			GetByIDFunc: func(ctx context.Context, id int) (*surah.Surah, error) {
				return nil, domain.ErrNotFound
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, mockSurahService))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random", nil))

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("Ayah service error", func(t *testing.T) {
		mockAyahService := &MockAyahService{
			GetRandomFunc: func(ctx context.Context, surahID int) (*ayah.Ayah, error) {
				return nil, errors.New("db error")
			},
		}

		r := setupRouter(handler.NewAyahHandler(mockAyahService, &MockSurahService{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/random?surah_id=1", nil))

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", w.Code)
		}
	})
}
