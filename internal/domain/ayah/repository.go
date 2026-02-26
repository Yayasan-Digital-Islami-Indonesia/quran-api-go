package ayah

import "context"

// AyahRepository defines read-only access to ayah data.
// Implement this interface in internal/repository/ayah_repository.go.
type AyahRepository interface {
	FindByID(ctx context.Context, id int) (*Ayah, error)
	FindBySurah(ctx context.Context, surahID, from, to int) ([]Ayah, error)
	FindBySurahAndNumber(ctx context.Context, surahID, number int) (*Ayah, error)
	FindRandom(ctx context.Context, surahID int) (*Ayah, error) // surahID=0 means any surah
}
