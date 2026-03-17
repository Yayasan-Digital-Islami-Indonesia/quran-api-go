package repository

import (
	"context"
	"database/sql"
	"fmt"

	juz "quran-api-go/internal/domain/juz"
)

// JuzRepository is the SQLite implementation of juz.JuzRepository.
type JuzRepository struct {
	db *sql.DB
}

// NewJuzRepository creates a new JuzRepository backed by the given database.
func NewJuzRepository(db *sql.DB) juz.JuzRepository {
	return &JuzRepository{db: db}
}

// FindAll returns all 30 juz entries ordered by juz_number.
func (r *JuzRepository) FindAll(ctx context.Context) ([]juz.Juz, error) {
	query := `SELECT id, juz_number, first_ayah_id, last_ayah_id
		FROM juzs
		ORDER BY juz_number`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("juz FindAll: %w", err)
	}
	defer rows.Close()

	var juzs []juz.Juz
	for rows.Next() {
		var j juz.Juz
		if err := rows.Scan(&j.ID, &j.JuzNumber, &j.FirstAyahID, &j.LastAyahID); err != nil {
			return nil, fmt.Errorf("juz FindAll scan: %w", err)
		}
		juzs = append(juzs, j)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("juz FindAll rows: %w", err)
	}

	return juzs, nil
}

// FindByNumber returns a single juz by its number. Returns nil, nil when not found.
func (r *JuzRepository) FindByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	query := `SELECT id, juz_number, first_ayah_id, last_ayah_id
		FROM juzs
		WHERE juz_number = ?`

	var j juz.Juz
	err := r.db.QueryRowContext(ctx, query, number).
		Scan(&j.ID, &j.JuzNumber, &j.FirstAyahID, &j.LastAyahID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("juz FindByNumber: %w", err)
	}

	return &j, nil
}

// FindAyahsByJuz returns ayahs belonging to the given juz, joined with surah
// name_latin. Results are paginated via limit/offset and ordered by ayah id.
func (r *JuzRepository) FindAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	query := `SELECT
			a.id,
			a.surah_id,
			s.name_latin,
			a.number_in_surah,
			a.text_uthmani,
			a.translation_indo,
			a.translation_en,
			a.juz_number
		FROM ayahs a
		JOIN surahs s ON s.id = a.surah_id
		WHERE a.juz_number = ?
		ORDER BY a.id
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, juzNumber, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("juz FindAyahsByJuz: %w", err)
	}
	defer rows.Close()

	var ayahs []juz.JuzAyah
	for rows.Next() {
		var ja juz.JuzAyah
		if err := rows.Scan(
			&ja.AyahID,
			&ja.SurahID,
			&ja.SurahNameLatin,
			&ja.NumberInSurah,
			&ja.TextUthmani,
			&ja.TranslationIdo,
			&ja.TranslationEn,
			&ja.JuzNumber,
		); err != nil {
			return nil, fmt.Errorf("juz FindAyahsByJuz scan: %w", err)
		}
		ayahs = append(ayahs, ja)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("juz FindAyahsByJuz rows: %w", err)
	}

	return ayahs, nil
}
