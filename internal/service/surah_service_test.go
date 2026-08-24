package service_test

import (
	"context"
	"errors"
	"testing"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/internal/service"
)

type MockSurahRepository struct {
	FindAllFn                func(ctx context.Context) ([]surah.Surah, error)
	FindByIDFn               func(ctx context.Context, id int) (*surah.Surah, error)
	FindByRevelationTypeFn   func(ctx context.Context, revelationType string) ([]surah.Surah, error)
}

func (m *MockSurahRepository) FindAll(ctx context.Context) ([]surah.Surah, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, nil
}

func (m *MockSurahRepository) FindByID(ctx context.Context, id int) (*surah.Surah, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockSurahRepository) FindByRevelationType(ctx context.Context, revelationType string) ([]surah.Surah, error) {
	if m.FindByRevelationTypeFn != nil {
		return m.FindByRevelationTypeFn(ctx, revelationType)
	}
	return nil, nil
}

func TestSurahService_GetAll_Success(t *testing.T) {
	expected := []surah.Surah{
		{ID: 1, Number: 1, NameLatin: "Al-Fatihah", NumberOfAyahs: 7, RevelationType: "meccan"},
		{ID: 2, Number: 2, NameLatin: "Al-Baqarah", NumberOfAyahs: 286, RevelationType: "medinan"},
	}

	mockRepo := &MockSurahRepository{
		FindAllFn: func(_ context.Context) ([]surah.Surah, error) {
			return expected, nil
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 surahs, got %d", len(result))
	}
	if result[0].NameLatin != "Al-Fatihah" {
		t.Errorf("expected Al-Fatihah, got %s", result[0].NameLatin)
	}
}

func TestSurahService_GetAll_Error(t *testing.T) {
	mockRepo := &MockSurahRepository{
		FindAllFn: func(_ context.Context) ([]surah.Surah, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetAll(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestSurahService_GetByID_Success(t *testing.T) {
	expected := &surah.Surah{ID: 1, Number: 1, NameLatin: "Al-Fatihah", NumberOfAyahs: 7, RevelationType: "meccan"}

	mockRepo := &MockSurahRepository{
		FindByIDFn: func(_ context.Context, id int) (*surah.Surah, error) {
			if id != 1 {
				t.Errorf("expected id=1, got %d", id)
			}
			return expected, nil
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected surah, got nil")
	}
	if result.NameLatin != "Al-Fatihah" {
		t.Errorf("expected Al-Fatihah, got %s", result.NameLatin)
	}
}

func TestSurahService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockSurahRepository{
		FindByIDFn: func(_ context.Context, id int) (*surah.Surah, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetByID(context.Background(), 999)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestSurahService_GetByID_Error(t *testing.T) {
	mockRepo := &MockSurahRepository{
		FindByIDFn: func(_ context.Context, id int) (*surah.Surah, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetByID(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestSurahService_GetByRevelationType_Success(t *testing.T) {
	expected := []surah.Surah{
		{ID: 1, Number: 1, NameLatin: "Al-Fatihah", RevelationType: "meccan"},
	}

	mockRepo := &MockSurahRepository{
		FindByRevelationTypeFn: func(_ context.Context, revelationType string) ([]surah.Surah, error) {
			if revelationType != "meccan" {
				t.Errorf("expected meccan, got %s", revelationType)
			}
			return expected, nil
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetByRevelationType(context.Background(), "meccan")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 surah, got %d", len(result))
	}
}

func TestSurahService_GetByRevelationType_Error(t *testing.T) {
	mockRepo := &MockSurahRepository{
		FindByRevelationTypeFn: func(_ context.Context, revelationType string) ([]surah.Surah, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewSurahService(mockRepo)
	result, err := svc.GetByRevelationType(context.Background(), "meccan")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
