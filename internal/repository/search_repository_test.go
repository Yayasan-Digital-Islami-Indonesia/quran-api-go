package repository_test

import (
	"context"
	"testing"

	"quran-api-go/internal/domain/search"
	"quran-api-go/internal/repository"
)

var createTablesSearch = `
CREATE TABLE surahs (
	id INTEGER PRIMARY KEY,
	number INTEGER NOT NULL,
	name_arabic TEXT NOT NULL,
	name_latin TEXT NOT NULL,
	name_transliteration TEXT NOT NULL,
	number_of_ayahs INTEGER NOT NULL,
	revelation_type TEXT NOT NULL
);
CREATE TABLE ayahs (
	id INTEGER PRIMARY KEY,
	surah_id INTEGER NOT NULL,
	number_in_surah INTEGER NOT NULL,
	text_uthmani TEXT NOT NULL,
	translation_indo TEXT NOT NULL,
	translation_en TEXT NOT NULL,
	juz_number INTEGER NOT NULL
);
CREATE VIRTUAL TABLE ayahs_fts USING fts5(
	text_uthmani,
	translation_indo,
	translation_en,
	content='ayahs',
	content_rowid='id'
);
`

var seedTablesSearch = `
INSERT INTO surahs (id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type)
VALUES (1, 1, 'الفاتحة', 'Pembukaan', 'Al-Fatihah', 7, 'meccan'),
       (2, 2, 'البقرة', 'Sapi Betina', 'Al-Baqarah', 286, 'medinan');

INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number)
VALUES (1, 1, 1, 'bismillah', 'dengan nama Allah Yang Maha Penyayang', 'in the name of Allah the Merciful', 1),
       (2, 1, 2, 'alhamdulillah', 'segala puji bagi Allah Tuhan semesta alam', 'all praise to Allah Lord of the worlds', 1),
       (3, 2, 1, 'alif lam mim', 'kitab Allah tidak ada keraguan', 'the book of Allah no doubt', 1);

INSERT INTO ayahs_fts (rowid, text_uthmani, translation_indo, translation_en)
SELECT id, text_uthmani, translation_indo, translation_en FROM ayahs;
`

func TestSearchRepository_Search_Match(t *testing.T) {
	db := setupTestDB(t, createTablesSearch, seedTablesSearch)
	repo := repository.NewSearchRepository(db)

	results, total, err := repo.Search(context.Background(), search.Params{
		Query: "Allah",
		Lang:  "id",
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected 3 total matches, got %d", total)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].SurahInfo.NameLatin != "Pembukaan" {
		t.Errorf("expected surah name 'Pembukaan', got %s", results[0].SurahInfo.NameLatin)
	}
	if results[0].Translation != "dengan nama Allah Yang Maha Penyayang" {
		t.Errorf("expected Indonesian translation, got %s", results[0].Translation)
	}
}

func TestSearchRepository_Search_EnglishLang(t *testing.T) {
	db := setupTestDB(t, createTablesSearch, seedTablesSearch)
	repo := repository.NewSearchRepository(db)

	results, _, err := repo.Search(context.Background(), search.Params{
		Query: "worlds",
		Lang:  "en",
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Translation != "all praise to Allah Lord of the worlds" {
		t.Errorf("expected English translation, got %s", results[0].Translation)
	}
}

func TestSearchRepository_Search_SurahFilter(t *testing.T) {
	db := setupTestDB(t, createTablesSearch, seedTablesSearch)
	repo := repository.NewSearchRepository(db)

	results, total, err := repo.Search(context.Background(), search.Params{
		Query:   "Allah",
		Lang:    "id",
		SurahID: 1,
		Page:    1,
		Limit:   20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 matches in surah 1, got %d", total)
	}
	for _, r := range results {
		if r.SurahID != 1 {
			t.Errorf("expected only surah 1 results, got surah %d", r.SurahID)
		}
	}
}

func TestSearchRepository_Search_JuzFilter(t *testing.T) {
	db := setupTestDB(t, createTablesSearch, seedTablesSearch)
	repo := repository.NewSearchRepository(db)

	_, total, err := repo.Search(context.Background(), search.Params{
		Query: "Allah",
		Lang:  "id",
		Juz:   1,
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected 3 matches in juz 1, got %d", total)
	}
}

func TestSearchRepository_Search_NoMatch(t *testing.T) {
	db := setupTestDB(t, createTablesSearch, seedTablesSearch)
	repo := repository.NewSearchRepository(db)

	results, total, err := repo.Search(context.Background(), search.Params{
		Query: "zzzznotfound",
		Lang:  "id",
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected 0 matches, got %d", total)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}
