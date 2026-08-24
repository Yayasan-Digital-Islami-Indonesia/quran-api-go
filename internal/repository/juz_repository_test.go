package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/repository"
)

var createTableJuz = `
CREATE TABLE juzs (
    id INTEGER PRIMARY KEY,
    juz_number INTEGER NOT NULL,
    first_ayah_id INTEGER NOT NULL,
    last_ayah_id INTEGER NOT NULL
);`

var createTableSurahForJuz = `
CREATE TABLE surahs (
    id INTEGER PRIMARY KEY,
    number INTEGER NOT NULL,
    name_arabic TEXT NOT NULL,
    name_latin TEXT NOT NULL,
    name_transliteration TEXT NOT NULL,
    number_of_ayahs INTEGER NOT NULL,
    revelation_type TEXT NOT NULL
);`

var createTableAyahForJuz = `
CREATE TABLE ayahs (
    id INTEGER PRIMARY KEY,
    surah_id INTEGER NOT NULL,
    number_in_surah INTEGER NOT NULL,
    text_uthmani TEXT NOT NULL,
    translation_indo TEXT NOT NULL,
    translation_en TEXT NOT NULL,
    juz_number INTEGER NOT NULL,
    sajda_type TEXT,
    revelation_type TEXT NOT NULL
);`

var seedTableJuz = `
INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id) VALUES
(1, 1, 1, 7),
(2, 2, 8, 10);

INSERT INTO surahs (id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type) VALUES
(1, 1, 'الفاتحة', 'Al-Fatihah', 'Al-Fatihah', 7, 'meccan'),
(2, 2, 'البقرة', 'Al-Baqarah', 'Al-Baqarah', 286, 'medinan');

INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, sajda_type, revelation_type) VALUES
(1, 1, 1, 'بِسۡمِ', 'Dengan nama Allah', 'In the name of Allah', 1, NULL, 'meccan'),
(2, 1, 2, 'ٱلۡحَمۡدُ', 'Segala puji', 'All praise', 1, NULL, 'meccan'),
(3, 1, 3, 'ٱلرَّحۡمَٰنِ', 'Yang Maha Pengasih', 'The Entirely Merciful', 1, NULL, 'meccan'),
(4, 1, 4, 'مَٰلِكِ', 'Pemilik hari', 'Sovereign of the Day', 1, NULL, 'meccan'),
(5, 1, 5, 'إِيَّاكَ', 'Hanya kepada Engkau', 'It is You we worship', 1, NULL, 'meccan'),
(6, 1, 6, 'ٱهۡدِنَا', 'Tunjukilah kami', 'Guide us', 1, NULL, 'meccan'),
(7, 1, 7, 'صِرَٰطَ', 'Jalan orang-orang', 'The path of those', 1, NULL, 'meccan'),
(8, 2, 1, 'الۤمۤ', 'Alif Lam Mim', 'Alif, Lam, Meem', 2, NULL, 'medinan'),
(9, 2, 2, 'ذَٰلِكَ', 'Kitab itu', 'This is the Book', 2, NULL, 'medinan'),
(10, 2, 3, 'ٱلَّذِينَ', 'Orang-orang', 'Who believe', 2, NULL, 'medinan');`

func setupJuzTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range []string{createTableJuz, createTableSurahForJuz, createTableAyahForJuz} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("exec setup: %v", err)
		}
	}

	if _, err := db.Exec(seedTableJuz); err != nil {
		t.Fatalf("exec seed: %v", err)
	}

	return db
}

func TestJuzRepository_FindAll(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	juzs, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(juzs) != 2 {
		t.Fatalf("expected 2 juz, got %d", len(juzs))
	}

	if juzs[0].JuzNumber != 1 {
		t.Errorf("expected juz_number=1, got %d", juzs[0].JuzNumber)
	}
	if juzs[0].TotalAyahs != 7 {
		t.Errorf("expected total_ayahs=7 for juz 1, got %d", juzs[0].TotalAyahs)
	}
	if juzs[1].TotalAyahs != 3 {
		t.Errorf("expected total_ayahs=3 for juz 2, got %d", juzs[1].TotalAyahs)
	}
}

func TestJuzRepository_FindByNumber_Success(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	j, err := repo.FindByNumber(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if j == nil {
		t.Fatal("expected juz, got nil")
	}

	if j.JuzNumber != 1 {
		t.Errorf("expected juz_number=1, got %d", j.JuzNumber)
	}
	if j.TotalAyahs != 7 {
		t.Errorf("expected total_ayahs=7, got %d", j.TotalAyahs)
	}
	if j.FirstAyahID != 1 {
		t.Errorf("expected first_ayah_id=1, got %d", j.FirstAyahID)
	}
	if j.LastAyahID != 7 {
		t.Errorf("expected last_ayah_id=7, got %d", j.LastAyahID)
	}
}

func TestJuzRepository_FindByNumber_NotFound(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	j, err := repo.FindByNumber(ctx, 99)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if j != nil {
		t.Fatalf("expected nil, got %+v", j)
	}
}

func TestJuzRepository_FindAyahsByJuz_Success(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	ayahs, err := repo.FindAyahsByJuz(ctx, 1, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ayahs) != 5 {
		t.Fatalf("expected 5 ayahs, got %d", len(ayahs))
	}

	if ayahs[0].AyahID != 1 {
		t.Errorf("expected first ayah_id=1, got %d", ayahs[0].AyahID)
	}
	if ayahs[0].SurahNameLatin != "Al-Fatihah" {
		t.Errorf("expected surah_name_latin='Al-Fatihah', got %s", ayahs[0].SurahNameLatin)
	}
	if ayahs[0].TranslationIdo != "Dengan nama Allah" {
		t.Errorf("expected Indonesian translation, got %s", ayahs[0].TranslationIdo)
	}
	if ayahs[0].TranslationEn != "In the name of Allah" {
		t.Errorf("expected English translation, got %s", ayahs[0].TranslationEn)
	}
}

func TestJuzRepository_FindAyahsByJuz_WithOffset(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	ayahs, err := repo.FindAyahsByJuz(ctx, 1, 3, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ayahs) != 3 {
		t.Fatalf("expected 3 ayahs, got %d", len(ayahs))
	}

	if ayahs[0].AyahID != 4 {
		t.Errorf("expected first ayah_id=4 with offset=3, got %d", ayahs[0].AyahID)
	}
}

func TestJuzRepository_FindAyahsByJuz_Empty(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	ayahs, err := repo.FindAyahsByJuz(ctx, 30, 100, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ayahs) != 0 {
		t.Fatalf("expected 0 ayahs for non-existent juz, got %d", len(ayahs))
	}
}

func TestJuzRepository_FindSurahsByJuz_Success(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	surahs, err := repo.FindSurahsByJuz(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(surahs) != 1 {
		t.Fatalf("expected 1 surah in juz 1, got %d", len(surahs))
	}

	if surahs[0].NameLatin != "Al-Fatihah" {
		t.Errorf("expected name_latin='Al-Fatihah', got %s", surahs[0].NameLatin)
	}
}

func TestJuzRepository_FindSurahsByJuz_MultipleSurahs(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	surahs, err := repo.FindSurahsByJuz(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(surahs) != 1 {
		t.Fatalf("expected 1 surah in juz 2, got %d", len(surahs))
	}

	if surahs[0].NameLatin != "Al-Baqarah" {
		t.Errorf("expected name_latin='Al-Baqarah', got %s", surahs[0].NameLatin)
	}
}

func TestJuzRepository_FindSurahsByJuz_Empty(t *testing.T) {
	db := setupJuzTestDB(t)
	repo := repository.NewJuzRepository(db)
	ctx := context.Background()

	surahs, err := repo.FindSurahsByJuz(ctx, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(surahs) != 0 {
		t.Fatalf("expected 0 surahs for non-existent juz, got %d", len(surahs))
	}
}
