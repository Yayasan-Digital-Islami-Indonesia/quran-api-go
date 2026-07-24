package service_test

import (
	"context"
	"errors"
	"testing"

	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/service"
)

type MockAyahRepository struct {
	FindBySurahAndNumberFunc func(ctx context.Context, surahID, number int) (*ayah.Ayah, error)
}

func (m *MockAyahRepository) FindByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahRepository) FindBySurah(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahRepository) FindBySurahAndNumber(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
	if m.FindBySurahAndNumberFunc != nil {
		return m.FindBySurahAndNumberFunc(ctx, surahID, number)
	}
	return nil, nil
}

func (m *MockAyahRepository) FindRandom(ctx context.Context, surahID int) (*ayah.Ayah, error) {
	return nil, nil
}

func (m *MockAyahRepository) FindSajda(ctx context.Context) ([]ayah.SajdaAyah, error) {
	return nil, nil
}

func TestAyahService_GetBySurahAndNumber(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedAyah := &ayah.Ayah{ID: 1, SurahID: 1, NumberInSurah: 1, TextUthmani: "Bismillah"}

		mockRepo := &MockAyahRepository{
			FindBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				if surahID != 1 || number != 1 {
					t.Errorf("expected surahID=1 and number=1, got %d and %d", surahID, number)
				}
				return expectedAyah, nil
			},
		}

		ayahService := service.NewAyahService(mockRepo)
		ay, err := ayahService.GetBySurahAndNumber(ctx, 1, 1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ay != expectedAyah {
			t.Errorf("expected ayah %v, got %v", expectedAyah, ay)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := &MockAyahRepository{
			FindBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, nil
			},
		}

		ayahService := service.NewAyahService(mockRepo)
		ay, err := ayahService.GetBySurahAndNumber(ctx, 115, 1)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ay != nil {
			t.Errorf("expected mil ayah, got %v", ay)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &MockAyahRepository{
			FindBySurahAndNumberFunc: func(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
				return nil, errors.New("db error")
			},
		}

		ayahService := service.NewAyahService(mockRepo)
		ay, err := ayahService.GetBySurahAndNumber(ctx, 1, 2)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "db error" {
			t.Errorf("expected error message 'db error', got '%s'", err.Error())
		}
		if ay != nil {
			t.Errorf("expected nil ayah, got %v", ay)
		}
	})
}
