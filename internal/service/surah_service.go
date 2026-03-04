package service

import (
	"context"

	"quran-api-go/internal/domain/surah"
)

type surahService struct {
	repo surah.SurahRepository
}

func NewSurahService(repo surah.SurahRepository) surah.SurahService {
	return &surahService{repo: repo}
}

func (s *surahService) GetAll(ctx context.Context) ([]surah.Surah, error) {
	return s.repo.FindAll(ctx)
}

func (s *surahService) GetByID(ctx context.Context, id int) (*surah.Surah, error) {
	return s.repo.FindByID(ctx, id)
}
