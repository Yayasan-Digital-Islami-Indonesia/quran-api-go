package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"quran-api-go/internal/database"
	"quran-api-go/internal/domain"
	"quran-api-go/internal/repository"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := database.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createTable := `
	CREATE TABLE surahs (
        id INTEGER PRIMARY KEY,
        number INTEGER NOT NULL,
        name_arabic TEXT NOT NULL,
        name_latin TEXT NOT NULL,
        name_transliteration TEXT NOT NULL,
        number_of_ayahs INTEGER NOT NULL,
        revelation_type TEXT NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatal(err)
	}

	// Seeding
	seed := `
	INSERT INTO surahs 
	(id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type)
	VALUES
		(1, 1, 'الفاتحة', 'Pembukaan', 'Al-Fatihah', 7, 'meccan'),
		(2, 2, 'البقرة', 'Sapi Betina', 'Al-Baqarah', 286, 'medinan'),
		(3, 3, 'آل عمران', 'Keluarga Imran', 'Ali ''Imran', 200, 'medinan'),
		(4, 4, 'النساء', 'Wanita', 'An-Nisa', 176, 'medinan');`

	_, err = db.Exec(seed)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestSurahRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSurahRepository(db)

	ctx := context.Background()

	users, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed to get users: %v", err)
	}

	if len(users) != 4 {
		t.Fatalf("expected 4 surahs, got %d", len(users))
	}
}

func TestSurahRepository_FindByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSurahRepository(db)

	ctx := context.Background()

	surah, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if surah == nil {
		t.Fatalf("expected surah, got nil")
	}

	if surah.NameTransliteration != "Al-Fatihah" {
		t.Errorf("expected Al-Fatihah, got %s", surah.NameTransliteration)
	}
}

func TestSurahRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSurahRepository(db)

	ctx := context.Background()

	surah, err := repo.FindByID(ctx, 999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if surah != nil {
		t.Fatalf("expected nil, got %+v", surah)
	}
}
