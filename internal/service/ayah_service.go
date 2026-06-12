package service

import (
	"context"

	"quran-api-go/internal/domain/ayah"
)

type ayahService struct {
	repo ayah.AyahRepository
}

// NewAyahService creates a new instance of AyahService
func NewAyahService(repo ayah.AyahRepository) ayah.AyahService {
	return &ayahService{
		repo: repo,
	}
}

func (s *ayahService) GetByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ayahService) GetBySurah(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
	return s.repo.FindBySurah(ctx, surahID, from, to)
}

func (s *ayahService) GetBySurahAndNumber(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
	return s.repo.FindBySurahAndNumber(ctx, surahID, number)
}

func (s *ayahService) GetRandom(ctx context.Context, surahID int) (*ayah.Ayah, error) {
	return s.repo.FindRandom(ctx, surahID)
}

func (s *ayahService) GetSajda(ctx context.Context) ([]ayah.SajdaAyah, error) {
	return s.repo.FindSajda(ctx)
}
