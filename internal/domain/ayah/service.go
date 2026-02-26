package ayah

import "context"

// AyahService defines the business operations for ayah data.
// Implement this interface in internal/service/ayah_service.go.
type AyahService interface {
	GetByID(ctx context.Context, id int) (*Ayah, error)
	GetBySurah(ctx context.Context, surahID, from, to int) ([]Ayah, error)
	GetBySurahAndNumber(ctx context.Context, surahID, number int) (*Ayah, error)
	GetRandom(ctx context.Context, surahID int) (*Ayah, error)
}
