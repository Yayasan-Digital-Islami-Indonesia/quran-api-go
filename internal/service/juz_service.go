package service

import (
	"context"

	"quran-api-go/internal/domain/juz"
)

type juzService struct {
	repo juz.JuzRepository
}

func NewJuzService(repo juz.JuzRepository) juz.JuzService {
	return &juzService{repo: repo}
}

func (s *juzService) GetAll(ctx context.Context) ([]juz.Juz, error) {
	return s.repo.FindAll(ctx)
}

func (s *juzService) GetByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	if number < 1 || number > 30 {
		return nil, nil
	}
	return s.repo.FindByNumber(ctx, number)
}

func (s *juzService) GetAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	if juzNumber < 1 || juzNumber > 30 {
		return nil, nil
	}

	// Default pagination values
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.FindAyahsByJuz(ctx, juzNumber, limit, offset)
}

func (s *juzService) GetSurahsByJuz(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error) {
	if juzNumber < 1 || juzNumber > 30 {
		return nil, nil
	}
	return s.repo.FindSurahsByJuz(ctx, juzNumber)
}
