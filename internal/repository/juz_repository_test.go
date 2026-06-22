package repository_test

import (
	"context"
	"errors"
	"testing"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/repository"
)

var createTablesJuz = `
CREATE TABLE surahs (
	id INTEGER PRIMARY KEY,
	number INTEGER NOT NULL,
	name_arabic TEXT NOT NULL,
	name_latin TEXT NOT NULL,
	name_transliteration TEXT NOT NULL,
	number_of_ayahs INTEGER NOT NULL,
	revelation_type TEXT NOT NULL
);
CREATE TABLE juzs (
	id INTEGER PRIMARY KEY,
	juz_number INTEGER NOT NULL,
	first_ayah_id INTEGER NOT NULL,
	last_ayah_id INTEGER NOT NULL
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
`

var seedTablesJuz = `
INSERT INTO surahs (id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type)
VALUES
	(1, 1, 'الفاتحة', 'Pembukaan', 'Al-Fatihah', 7, 'meccan'),
	(2, 2, 'البقرة', 'Sapi Betina', 'Al-Baqarah', 286, 'medinan');

INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id) VALUES
	(1, 1, 1, 9),
	(2, 2, 10, 12);

INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number)
VALUES
	(1, 1, 1, 'bismillah', 'dengan nama Allah', 'in the name of Allah', 1),
	(2, 1, 2, 'alhamdulillah', 'segala puji', 'all praise', 1),
	(3, 2, 1, 'alif lam mim', 'alif lam mim', 'alif lam mim', 1),
	(10, 2, 142, 'sayaqulu', 'orang bodoh akan berkata', 'the foolish will say', 2);
`

func TestJuzRepository_FindAll(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	juzs, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(juzs) != 2 {
		t.Fatalf("expected 2 juzs, got %d", len(juzs))
	}
	if juzs[0].JuzNumber != 1 || juzs[0].TotalAyahs != 3 {
		t.Errorf("expected juz 1 with 3 ayahs, got number=%d total=%d", juzs[0].JuzNumber, juzs[0].TotalAyahs)
	}
}

func TestJuzRepository_FindByNumber_Success(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	j, err := repo.FindByNumber(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected juz, got nil")
	}
	if j.JuzNumber != 1 || j.TotalAyahs != 3 {
		t.Errorf("expected juz 1 with 3 ayahs, got number=%d total=%d", j.JuzNumber, j.TotalAyahs)
	}
}

func TestJuzRepository_FindByNumber_NotFound(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	j, err := repo.FindByNumber(context.Background(), 99)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if j != nil {
		t.Fatalf("expected nil, got %+v", j)
	}
}

func TestJuzRepository_FindAyahsByJuz(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	ayahs, err := repo.FindAyahsByJuz(context.Background(), 1, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ayahs) != 3 {
		t.Fatalf("expected 3 ayahs, got %d", len(ayahs))
	}
	if ayahs[0].SurahNameLatin != "Pembukaan" {
		t.Errorf("expected surah name 'Pembukaan', got %s", ayahs[0].SurahNameLatin)
	}
}

func TestJuzRepository_FindAyahsByJuz_Pagination(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	ayahs, err := repo.FindAyahsByJuz(context.Background(), 1, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ayahs) != 1 {
		t.Fatalf("expected 1 ayah after offset, got %d", len(ayahs))
	}
}

func TestJuzRepository_FindSurahsByJuz(t *testing.T) {
	db := setupTestDB(t, createTablesJuz, seedTablesJuz)
	repo := repository.NewJuzRepository(db)

	surahs, err := repo.FindSurahsByJuz(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(surahs) != 2 {
		t.Fatalf("expected 2 distinct surahs in juz 1, got %d", len(surahs))
	}
}
