package service_test

import (
	"context"
	"errors"
	"testing"

	"quran-api-go/internal/domain/search"
	"quran-api-go/internal/service"
)

type MockSearchRepository struct {
	SearchFn func(ctx context.Context, params search.Params) ([]search.Result, int, error)
}

func (m *MockSearchRepository) Search(ctx context.Context, params search.Params) ([]search.Result, int, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, params)
	}
	return []search.Result{}, 0, nil
}

func TestSearchService_Search_Success(t *testing.T) {
	expected := []search.Result{
		{ID: 1, SurahID: 1, Translation: "Dengan nama Allah"},
	}

	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Query != "bismillah" {
				t.Errorf("expected query='bismillah', got %s", p.Query)
			}
			if p.Lang != "id" {
				t.Errorf("expected lang='id', got %s", p.Lang)
			}
			if p.Page != 1 {
				t.Errorf("expected page=1, got %d", p.Page)
			}
			if p.Limit != 20 {
				t.Errorf("expected limit=20, got %d", p.Limit)
			}
			return expected, 1, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	results, total, err := svc.Search(context.Background(), search.Params{
		Query: "bismillah",
		Lang:  "id",
		Page:  1,
		Limit: 20,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchService_Search_EmptyQuery(t *testing.T) {
	svc := service.NewSearchService(&MockSearchRepository{})
	results, total, err := svc.Search(context.Background(), search.Params{
		Query: "",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total=0, got %d", total)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchService_Search_DefaultLang(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Lang != "id" {
				t.Errorf("expected default lang='id', got %s", p.Lang)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query: "test",
		Lang:  "",
	})
}

func TestSearchService_Search_InvalidLang_FallbackToID(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Lang != "id" {
				t.Errorf("expected fallback lang='id', got %s", p.Lang)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query: "test",
		Lang:  "fr",
	})
}

func TestSearchService_Search_DefaultPage(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Page != 1 {
				t.Errorf("expected default page=1, got %d", p.Page)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query: "test",
		Page:  0,
	})
}

func TestSearchService_Search_DefaultLimit(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Limit != 20 {
				t.Errorf("expected default limit=20, got %d", p.Limit)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query: "test",
		Limit: 0,
	})
}

func TestSearchService_Search_CapLimit(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.Limit != 100 {
				t.Errorf("expected capped limit=100, got %d", p.Limit)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query: "test",
		Limit: 200,
	})
}

func TestSearchService_Search_Error(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			return nil, 0, errors.New("fts5 error")
		},
	}

	svc := service.NewSearchService(mockRepo)
	results, total, err := svc.Search(context.Background(), search.Params{
		Query: "test",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
	if results != nil {
		t.Errorf("expected nil results, got %v", results)
	}
}

func TestSearchService_Search_WithFilters(t *testing.T) {
	mockRepo := &MockSearchRepository{
		SearchFn: func(_ context.Context, p search.Params) ([]search.Result, int, error) {
			if p.SurahID != 2 {
				t.Errorf("expected surahID=2, got %d", p.SurahID)
			}
			if p.Juz != 3 {
				t.Errorf("expected juz=3, got %d", p.Juz)
			}
			return []search.Result{}, 0, nil
		},
	}

	svc := service.NewSearchService(mockRepo)
	svc.Search(context.Background(), search.Params{
		Query:   "test",
		SurahID: 2,
		Juz:     3,
	})
}
