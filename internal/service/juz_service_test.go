package service_test

import (
	"context"
	"errors"
	"testing"

	"quran-api-go/internal/domain/juz"
	"quran-api-go/internal/service"
)

type MockJuzRepository struct {
	FindAllFn           func(ctx context.Context) ([]juz.Juz, error)
	FindByNumberFn      func(ctx context.Context, number int) (*juz.Juz, error)
	FindAyahsByJuzFn    func(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error)
	FindSurahsByJuzFn   func(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error)
}

func (m *MockJuzRepository) FindAll(ctx context.Context) ([]juz.Juz, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn(ctx)
	}
	return nil, nil
}

func (m *MockJuzRepository) FindByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	if m.FindByNumberFn != nil {
		return m.FindByNumberFn(ctx, number)
	}
	return nil, nil
}

func (m *MockJuzRepository) FindAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	if m.FindAyahsByJuzFn != nil {
		return m.FindAyahsByJuzFn(ctx, juzNumber, limit, offset)
	}
	return nil, nil
}

func (m *MockJuzRepository) FindSurahsByJuz(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error) {
	if m.FindSurahsByJuzFn != nil {
		return m.FindSurahsByJuzFn(ctx, juzNumber)
	}
	return nil, nil
}

func TestJuzService_GetAll_Success(t *testing.T) {
	expected := []juz.Juz{
		{ID: 1, JuzNumber: 1, TotalAyahs: 141},
		{ID: 2, JuzNumber: 2, TotalAyahs: 111},
	}

	mockRepo := &MockJuzRepository{
		FindAllFn: func(_ context.Context) ([]juz.Juz, error) {
			return expected, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 juz, got %d", len(result))
	}
}

func TestJuzService_GetAll_Error(t *testing.T) {
	mockRepo := &MockJuzRepository{
		FindAllFn: func(_ context.Context) ([]juz.Juz, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetAll(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestJuzService_GetByNumber_Success(t *testing.T) {
	expected := &juz.Juz{ID: 1, JuzNumber: 1, TotalAyahs: 141}

	mockRepo := &MockJuzRepository{
		FindByNumberFn: func(_ context.Context, number int) (*juz.Juz, error) {
			return expected, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetByNumber(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected juz, got nil")
	}
	if result.JuzNumber != 1 {
		t.Errorf("expected juz_number=1, got %d", result.JuzNumber)
	}
}

func TestJuzService_GetByNumber_OutOfRange_Low(t *testing.T) {
	svc := service.NewJuzService(&MockJuzRepository{})
	result, err := svc.GetByNumber(context.Background(), 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for out-of-range, got %v", result)
	}
}

func TestJuzService_GetByNumber_OutOfRange_High(t *testing.T) {
	svc := service.NewJuzService(&MockJuzRepository{})
	result, err := svc.GetByNumber(context.Background(), 31)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for out-of-range, got %v", result)
	}
}

func TestJuzService_GetAyahsByJuz_Success(t *testing.T) {
	expected := []juz.JuzAyah{
		{AyahID: 1, SurahID: 1, SurahNameLatin: "Al-Fatihah", NumberInSurah: 1, JuzNumber: 1},
	}

	mockRepo := &MockJuzRepository{
		FindAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			if juzNumber != 1 || limit != 50 || offset != 0 {
				t.Errorf("unexpected args: juzNumber=%d limit=%d offset=%d", juzNumber, limit, offset)
			}
			return expected, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetAyahsByJuz(context.Background(), 1, 50, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 ayah, got %d", len(result))
	}
}

func TestJuzService_GetAyahsByJuz_OutOfRange(t *testing.T) {
	svc := service.NewJuzService(&MockJuzRepository{})
	result, err := svc.GetAyahsByJuz(context.Background(), 0, 50, 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for out-of-range, got %v", result)
	}
}

func TestJuzService_GetAyahsByJuz_DefaultLimit(t *testing.T) {
	mockRepo := &MockJuzRepository{
		FindAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			if limit != 50 {
				t.Errorf("expected default limit=50, got %d", limit)
			}
			return []juz.JuzAyah{}, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	svc.GetAyahsByJuz(context.Background(), 1, 0, 0)
}

func TestJuzService_GetAyahsByJuz_CapLimit(t *testing.T) {
	mockRepo := &MockJuzRepository{
		FindAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			if limit != 100 {
				t.Errorf("expected capped limit=100, got %d", limit)
			}
			return []juz.JuzAyah{}, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	svc.GetAyahsByJuz(context.Background(), 1, 200, 0)
}

func TestJuzService_GetAyahsByJuz_Error(t *testing.T) {
	mockRepo := &MockJuzRepository{
		FindAyahsByJuzFn: func(_ context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetAyahsByJuz(context.Background(), 1, 50, 0)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestJuzService_GetSurahsByJuz_Success(t *testing.T) {
	expected := []juz.JuzSurah{
		{ID: 1, Number: 1, NameLatin: "Al-Fatihah"},
	}

	mockRepo := &MockJuzRepository{
		FindSurahsByJuzFn: func(_ context.Context, juzNumber int) ([]juz.JuzSurah, error) {
			return expected, nil
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetSurahsByJuz(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 surah, got %d", len(result))
	}
}

func TestJuzService_GetSurahsByJuz_OutOfRange(t *testing.T) {
	svc := service.NewJuzService(&MockJuzRepository{})
	result, err := svc.GetSurahsByJuz(context.Background(), 31)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil for out-of-range, got %v", result)
	}
}

func TestJuzService_GetSurahsByJuz_Error(t *testing.T) {
	mockRepo := &MockJuzRepository{
		FindSurahsByJuzFn: func(_ context.Context, juzNumber int) ([]juz.JuzSurah, error) {
			return nil, errors.New("db error")
		},
	}

	svc := service.NewJuzService(mockRepo)
	result, err := svc.GetSurahsByJuz(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
