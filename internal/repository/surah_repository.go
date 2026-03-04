package repository

import (
	"context"
	"database/sql"

	"quran-api-go/internal/domain"
	surah "quran-api-go/internal/domain/surah"
)

type SurahRepository struct {
	db *sql.DB
}

func NewSurahRepository(db *sql.DB) surah.SurahRepository {
	return &SurahRepository{
		db: db,
	}
}

func (s *SurahRepository) FindAll(ctx context.Context) ([]surah.Surah, error) {
	query := `SELECT id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type FROM surahs`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surahs []surah.Surah

	for rows.Next() {
		var surah surah.Surah
		if err := rows.Scan(
			&surah.ID,
			&surah.Number,
			&surah.NameArabic,
			&surah.NameLatin,
			&surah.NameTransliteration,
			&surah.NumberOfAyahs,
			&surah.RevelationType,
		); err != nil {
			return nil, err
		}

		surahs = append(surahs, surah)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return surahs, nil
}

func (s *SurahRepository) FindByID(ctx context.Context, id int) (*surah.Surah, error) {
	query := `SELECT 
	id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type 
	FROM surahs
	WHERE id = ?`

	row := s.db.QueryRowContext(ctx, query, id)

	var surah surah.Surah
	err := row.Scan(
		&surah.ID,
		&surah.Number,
		&surah.NameArabic,
		&surah.NameLatin,
		&surah.NameTransliteration,
		&surah.NumberOfAyahs,
		&surah.RevelationType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &surah, nil
}
