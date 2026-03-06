package repository

import (
	"context"
	"database/sql"
	"quran-api-go/internal/domain/ayah"
)

type AyahRepository struct {
	db *sql.DB
}

func NewAyahRepository(db *sql.DB) ayah.AyahRepository {
	return &AyahRepository{
		db: db,
	}
}

func (a *AyahRepository) FindByID(ctx context.Context, id int) (*ayah.Ayah, error) {
	query := `SELECT 
	id, 
	surah_id, 
	number_in_surah, 
	text_uthmani,
	translation_indo, 
	translation_en,
	juz_number,
	sajda_type,
	revelation_type
	FROM ayahs WHERE id = ?`

	row := a.db.QueryRowContext(ctx, query, id)

	var ayah ayah.Ayah
	err := row.Scan(
		&ayah.ID,
		&ayah.SurahID,
		&ayah.NumberInSurah,
		&ayah.TextUthmani,
		&ayah.TranslationIdo,
		&ayah.TranslationEn,
		&ayah.JuzNumber,
		&ayah.SajdaType,
		&ayah.RevelationType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}
	return &ayah, nil
}

func (a *AyahRepository) FindBySurah(ctx context.Context, surahID, from, to int) ([]ayah.Ayah, error) {
	query := `SELECT 
	id, 
	surah_id, 
	number_in_surah, 
	text_uthmani,
	translation_indo, 
	translation_en,
	juz_number,
	sajda_type,
	revelation_type
	FROM ayahs WHERE surah_id = ? AND number_in_surah BETWEEN ? AND ?
	ORDER BY number_in_surah ASC`

	rows, err := a.db.QueryContext(ctx, query, surahID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ayahs []ayah.Ayah

	for rows.Next() {
		var ayah ayah.Ayah
		err := rows.Scan(
			&ayah.ID,
			&ayah.SurahID,
			&ayah.NumberInSurah,
			&ayah.TextUthmani,
			&ayah.TranslationIdo,
			&ayah.TranslationEn,
			&ayah.JuzNumber,
			&ayah.SajdaType,
			&ayah.RevelationType,
		)
		if err != nil {
			return nil, err
		}

		ayahs = append(ayahs, ayah)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ayahs, nil
}

func (a *AyahRepository) FindBySurahAndNumber(ctx context.Context, surahID, number int) (*ayah.Ayah, error) {
	query := `SELECT 
	id, 
	surah_id, 
	number_in_surah, 
	text_uthmani,
	translation_indo, 
	translation_en,
	juz_number,
	sajda_type,
	revelation_type
	FROM ayahs WHERE surah_id = ? AND number_in_surah = ?`

	row := a.db.QueryRowContext(ctx, query, surahID, number)

	var ayah ayah.Ayah
	err := row.Scan(
		&ayah.ID,
		&ayah.SurahID,
		&ayah.NumberInSurah,
		&ayah.TextUthmani,
		&ayah.TranslationIdo,
		&ayah.TranslationEn,
		&ayah.JuzNumber,
		&ayah.SajdaType,
		&ayah.RevelationType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ayah, nil
}

func (a *AyahRepository) FindRandom(ctx context.Context, surahID int) (*ayah.Ayah, error) {
	var query string
	var row *sql.Row

	if surahID == 0 {
		// random dari semua surah
		query = `SELECT 
		id,
		surah_id,
		number_in_surah,
		text_uthmani,
		translation_indo,
		translation_en,
		juz_number,
		sajda_type,
		revelation_type
		FROM ayahs
		ORDER BY RANDOM()
		LIMIT 1`

		row = a.db.QueryRowContext(ctx, query)

	} else {
		// random dari surah tertentu
		query = `SELECT 
		id,
		surah_id,
		number_in_surah,
		text_uthmani,
		translation_indo,
		translation_en,
		juz_number,
		sajda_type,
		revelation_type
		FROM ayahs
		WHERE surah_id = ?
		ORDER BY RANDOM()
		LIMIT 1`

		row = a.db.QueryRowContext(ctx, query, surahID)
	}

	var ay ayah.Ayah
	err := row.Scan(
		&ay.ID,
		&ay.SurahID,
		&ay.NumberInSurah,
		&ay.TextUthmani,
		&ay.TranslationIdo,
		&ay.TranslationEn,
		&ay.JuzNumber,
		&ay.SajdaType,
		&ay.RevelationType,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ay, nil
}
