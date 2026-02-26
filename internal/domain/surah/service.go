package surah

import "context"

// SurahService defines the business operations for surah data.
// Implement this interface in internal/service/surah_service.go.
type SurahService interface {
	GetAll(ctx context.Context) ([]Surah, error)
	GetByID(ctx context.Context, id int) (*Surah, error)
}
