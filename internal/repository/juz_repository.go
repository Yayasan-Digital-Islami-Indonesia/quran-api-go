package repository

import (
	"context"
	"database/sql"
	"errors"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/juz"
)

type juzRepository struct {
	db *sql.DB
}

func NewJuzRepository(db *sql.DB) juz.JuzRepository {
	return &juzRepository{db: db}
}

func (r *juzRepository) FindAll(ctx context.Context) ([]juz.Juz, error) {
	query := `
		SELECT j.id, j.juz_number, j.first_ayah_id, j.last_ayah_id,
		       COUNT(a.id) as total_ayahs
		FROM juzs j
		LEFT JOIN ayahs a ON a.juz_number = j.juz_number
		GROUP BY j.id
		ORDER BY j.juz_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var juzs []juz.Juz
	for rows.Next() {
		var j juz.Juz
		if err := rows.Scan(&j.ID, &j.JuzNumber, &j.FirstAyahID, &j.LastAyahID, &j.TotalAyahs); err != nil {
			return nil, err
		}
		juzs = append(juzs, j)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return juzs, nil
}

func (r *juzRepository) FindByNumber(ctx context.Context, number int) (*juz.Juz, error) {
	query := `
		SELECT j.id, j.juz_number, j.first_ayah_id, j.last_ayah_id,
		       COUNT(a.id) as total_ayahs
		FROM juzs j
		LEFT JOIN ayahs a ON a.juz_number = j.juz_number
		WHERE j.juz_number = ?
		GROUP BY j.id
	`

	var j juz.Juz
	err := r.db.QueryRowContext(ctx, query, number).Scan(&j.ID, &j.JuzNumber, &j.FirstAyahID, &j.LastAyahID, &j.TotalAyahs)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (r *juzRepository) FindAyahsByJuz(ctx context.Context, juzNumber, limit, offset int) ([]juz.JuzAyah, error) {
	query := `
		SELECT a.id, a.surah_id, s.name_latin, a.number_in_surah,
			   a.text_uthmani, a.translation_indo, a.translation_en, a.juz_number
		FROM ayahs a
		INNER JOIN surahs s ON a.surah_id = s.id
		WHERE a.juz_number = ?
		ORDER BY a.id ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, juzNumber, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ayahs []juz.JuzAyah
	for rows.Next() {
		var a juz.JuzAyah
		if err := rows.Scan(
			&a.AyahID,
			&a.SurahID,
			&a.SurahNameLatin,
			&a.NumberInSurah,
			&a.TextUthmani,
			&a.TranslationIdo,
			&a.TranslationEn,
			&a.JuzNumber,
		); err != nil {
			return nil, err
		}
		ayahs = append(ayahs, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ayahs, nil
}

func (r *juzRepository) FindSurahsByJuz(ctx context.Context, juzNumber int) ([]juz.JuzSurah, error) {
	query := `
		SELECT DISTINCT s.id, s.number, s.name_arabic, s.name_latin, s.name_transliteration,
		       s.number_of_ayahs, s.revelation_type
		FROM ayahs a
		INNER JOIN surahs s ON a.surah_id = s.id
		WHERE a.juz_number = ?
		ORDER BY s.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, juzNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surahs []juz.JuzSurah
	for rows.Next() {
		var s juz.JuzSurah
		if err := rows.Scan(&s.ID, &s.Number, &s.NameArabic, &s.NameLatin, &s.NameTransliteration, &s.NumberOfAyahs, &s.RevelationType); err != nil {
			return nil, err
		}
		surahs = append(surahs, s)
	}
	return surahs, rows.Err()
}
