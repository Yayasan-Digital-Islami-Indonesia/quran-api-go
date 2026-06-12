package repository

import (
	"context"
	"database/sql"
	"fmt"

	"quran-api-go/internal/domain/search"
)

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) search.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) Search(ctx context.Context, p search.Params) ([]search.Result, int, error) {
	// Use a subquery so FTS5 only returns rowids from its index — avoids the
	// content-table read that would fail on the missing `ayah_id` column.
	ftsSubquery := "a.id IN (SELECT rowid FROM ayahs_fts WHERE ayahs_fts MATCH ?)"
	ftsArgs := []interface{}{p.Query + "*"} // prefix match for partial terms

	// Build optional outer filters on the ayahs table directly
	outerFilters := ""
	outerArgs := []interface{}{}
	if p.SurahID > 0 {
		outerFilters += " AND a.surah_id = ?"
		outerArgs = append(outerArgs, p.SurahID)
	}
	if p.Juz > 0 {
		outerFilters += " AND a.juz_number = ?"
		outerArgs = append(outerArgs, p.Juz)
	}

	whereClause := ftsSubquery + outerFilters
	baseArgs := append(ftsArgs, outerArgs...)

	// Count total results
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM ayahs a WHERE %s
	`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, baseArgs...).Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, err
	}

	limit := p.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (p.Page - 1) * p.Limit
	if offset < 0 {
		offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT a.id, a.surah_id, s.name_latin, a.number_in_surah,
			   a.text_uthmani, a.translation_indo, a.translation_en, a.juz_number
		FROM ayahs a
		JOIN surahs s ON a.surah_id = s.id
		WHERE %s
		ORDER BY a.id ASC
		LIMIT ? OFFSET ?
	`, whereClause)

	dataArgs := append(baseArgs, limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []search.Result
	for rows.Next() {
		var r search.Result
		var translationIndo, translationEn string

		if err := rows.Scan(
			&r.ID,
			&r.SurahID,
			&r.SurahInfo.NameLatin,
			&r.NumberInSurah,
			&r.TextUthmani,
			&translationIndo,
			&translationEn,
			&r.JuzNumber,
		); err != nil {
			return nil, 0, err
		}

		r.SurahInfo.ID = r.SurahID

		// Set translation based on lang
		if p.Lang == "en" {
			r.Translation = translationEn
		} else {
			r.Translation = translationIndo
		}

		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
