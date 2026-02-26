package juz

import "context"

// JuzService defines the business operations for juz data.
// Implement this interface in internal/service/juz_service.go.
type JuzService interface {
	GetAll(ctx context.Context) ([]Juz, error)
	GetByNumber(ctx context.Context, number int) (*Juz, error)
	GetAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]JuzAyah, error)
}
