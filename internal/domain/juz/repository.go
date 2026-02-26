package juz

import "context"

// JuzRepository defines read-only access to juz data.
// Implement this interface in internal/repository/juz_repository.go.
type JuzRepository interface {
	FindAll(ctx context.Context) ([]Juz, error)
	FindByNumber(ctx context.Context, number int) (*Juz, error)
	FindAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]JuzAyah, error)
}

// JuzAyah is an ayah row joined with its surah name, used in juz detail responses.
type JuzAyah struct {
	AyahID        int    `json:"id"`
	SurahID       int    `json:"surah_id"`
	SurahNameLatin string `json:"surah_name_latin"`
	NumberInSurah int    `json:"number_in_surah"`
	TextUthmani   string `json:"text_uthmani"`
	TranslationIdo string `json:"translation_indo"`
	TranslationEn  string `json:"translation_en"`
	JuzNumber     int    `json:"juz_number"`
}
