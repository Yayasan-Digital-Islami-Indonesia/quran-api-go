package surah

import "context"

// SurahRepository defines read-only access to surah data.
// Implement this interface in internal/repository/surah_repository.go.
type SurahRepository interface {
	FindAll(ctx context.Context) ([]Surah, error)
	FindByID(ctx context.Context, id int) (*Surah, error)
}
