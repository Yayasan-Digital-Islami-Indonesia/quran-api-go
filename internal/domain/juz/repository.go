package juz

import "context"

// JuzRepository defines read-only access to juz data.
// Implement this interface in internal/repository/juz_repository.go.
type JuzRepository interface {
	FindAll(ctx context.Context) ([]Juz, error)
	FindByNumber(ctx context.Context, number int) (*Juz, error)
	FindAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]JuzAyah, error)
	FindSurahsByJuz(ctx context.Context, juzNumber int) ([]JuzSurah, error)
}
