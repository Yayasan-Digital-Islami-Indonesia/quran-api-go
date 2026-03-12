package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"quran-api-go/internal/repository"

	_ "modernc.org/sqlite"
)

// setupJuzTestDB creates an in-memory SQLite database seeded with test data
// for the juzs, surahs, and ayahs tables.
func setupJuzTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	schema := `
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
		juz_number INTEGER NOT NULL,
		sajda_type TEXT,
		revelation_type TEXT NOT NULL,
		FOREIGN KEY (surah_id) REFERENCES surahs(id)
	);
	CREATE INDEX idx_ayahs_juz_number ON ayahs (juz_number);
	CREATE TABLE juzs (
		id INTEGER PRIMARY KEY,
		juz_number INTEGER NOT NULL,
		first_ayah_id INTEGER NOT NULL,
		last_ayah_id INTEGER NOT NULL,
		FOREIGN KEY (first_ayah_id) REFERENCES ayahs(id),
		FOREIGN KEY (last_ayah_id) REFERENCES ayahs(id)
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	seed := `
	INSERT INTO surahs (id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type)
	VALUES
		(1, 1, 'الفاتحة', 'Pembukaan', 'Al-Fatihah', 7, 'meccan'),
		(2, 2, 'البقرة', 'Sapi Betina', 'Al-Baqarah', 286, 'medinan');

	INSERT INTO ayahs (id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, sajda_type, revelation_type)
	VALUES
		(1, 1, 1, 'بِسْمِ ٱللَّهِ', 'Dengan nama Allah', 'In the name of Allah', 1, NULL, 'meccan'),
		(2, 1, 2, 'ٱلْحَمْدُ لِلَّهِ', 'Segala puji bagi Allah', 'All praise is due to Allah', 1, NULL, 'meccan'),
		(3, 1, 3, 'ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ', 'Yang Maha Pengasih', 'The Entirely Merciful', 1, NULL, 'meccan'),
		(4, 2, 1, 'الٓمٓ', 'Alif Lam Mim', 'Alif Lam Mim', 1, NULL, 'medinan'),
		(5, 2, 2, 'ذَٰلِكَ ٱلْكِتَـٰبُ', 'Kitab itu', 'This is the Book', 1, NULL, 'medinan');

	INSERT INTO juzs (id, juz_number, first_ayah_id, last_ayah_id)
	VALUES
		(1, 1, 1, 5),
		(2, 2, 6, 10);`

	if _, err := db.Exec(seed); err != nil {
		t.Fatal(err)
	}

	return db
}

func TestJuzRepository_FindAll(t *testing.T) {
	db := setupJuzTestDB(t)
	defer db.Close()
	repo := repository.NewJuzRepository(db)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "returns all juz entries",
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			juzs, err := repo.FindAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(juzs) != tt.wantCount {
				t.Errorf("FindAll() returned %d juz(s), want %d", len(juzs), tt.wantCount)
			}
		})
	}
}

func TestJuzRepository_FindAll_Order(t *testing.T) {
	db := setupJuzTestDB(t)
	defer db.Close()
	repo := repository.NewJuzRepository(db)

	juzs, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(juzs) < 2 {
		t.Fatal("expected at least 2 juz entries")
	}
	if juzs[0].JuzNumber > juzs[1].JuzNumber {
		t.Errorf("expected ascending order, got juz %d before juz %d", juzs[0].JuzNumber, juzs[1].JuzNumber)
	}
}

func TestJuzRepository_FindByNumber(t *testing.T) {
	db := setupJuzTestDB(t)
	defer db.Close()
	repo := repository.NewJuzRepository(db)

	tests := []struct {
		name       string
		number     int
		wantNil    bool
		wantJuzNum int
		wantErr    bool
	}{
		{
			name:       "juz 1 found",
			number:     1,
			wantNil:    false,
			wantJuzNum: 1,
			wantErr:    false,
		},
		{
			name:       "juz 2 found",
			number:     2,
			wantNil:    false,
			wantJuzNum: 2,
			wantErr:    false,
		},
		{
			name:    "juz 99 not found",
			number:  99,
			wantNil: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := repo.FindByNumber(context.Background(), tt.number)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByNumber(%d) error = %v, wantErr %v", tt.number, err, tt.wantErr)
			}
			if tt.wantNil {
				if j != nil {
					t.Errorf("FindByNumber(%d) = %+v, want nil", tt.number, j)
				}
				return
			}
			if j == nil {
				t.Fatalf("FindByNumber(%d) = nil, want juz %d", tt.number, tt.wantJuzNum)
			}
			if j.JuzNumber != tt.wantJuzNum {
				t.Errorf("FindByNumber(%d).JuzNumber = %d, want %d", tt.number, j.JuzNumber, tt.wantJuzNum)
			}
		})
	}
}

func TestJuzRepository_FindAyahsByJuz(t *testing.T) {
	db := setupJuzTestDB(t)
	defer db.Close()
	repo := repository.NewJuzRepository(db)

	tests := []struct {
		name        string
		juzNumber   int
		limit       int
		offset      int
		wantCount   int
		wantFirstID int
		wantErr     bool
	}{
		{
			name:        "all ayahs in juz 1",
			juzNumber:   1,
			limit:       100,
			offset:      0,
			wantCount:   5,
			wantFirstID: 1,
			wantErr:     false,
		},
		{
			name:        "paginated — limit 2 offset 0",
			juzNumber:   1,
			limit:       2,
			offset:      0,
			wantCount:   2,
			wantFirstID: 1,
			wantErr:     false,
		},
		{
			name:        "paginated — limit 2 offset 3",
			juzNumber:   1,
			limit:       2,
			offset:      3,
			wantCount:   2,
			wantFirstID: 4,
			wantErr:     false,
		},
		{
			name:      "juz with no ayahs returns empty",
			juzNumber: 2,
			limit:     100,
			offset:    0,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ayahs, err := repo.FindAyahsByJuz(context.Background(), tt.juzNumber, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAyahsByJuz(%d, %d, %d) error = %v, wantErr %v",
					tt.juzNumber, tt.limit, tt.offset, err, tt.wantErr)
			}
			if len(ayahs) != tt.wantCount {
				t.Errorf("FindAyahsByJuz(%d, %d, %d) returned %d ayah(s), want %d",
					tt.juzNumber, tt.limit, tt.offset, len(ayahs), tt.wantCount)
			}
			if tt.wantCount > 0 && ayahs[0].AyahID != tt.wantFirstID {
				t.Errorf("first ayah ID = %d, want %d", ayahs[0].AyahID, tt.wantFirstID)
			}
		})
	}
}

func TestJuzRepository_FindAyahsByJuz_JoinFields(t *testing.T) {
	db := setupJuzTestDB(t)
	defer db.Close()
	repo := repository.NewJuzRepository(db)

	ayahs, err := repo.FindAyahsByJuz(context.Background(), 1, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ayahs) != 1 {
		t.Fatalf("expected 1 ayah, got %d", len(ayahs))
	}

	a := ayahs[0]
	if a.SurahNameLatin != "Pembukaan" {
		t.Errorf("SurahNameLatin = %q, want %q", a.SurahNameLatin, "Pembukaan")
	}
	if a.SurahID != 1 {
		t.Errorf("SurahID = %d, want 1", a.SurahID)
	}
	if a.JuzNumber != 1 {
		t.Errorf("JuzNumber = %d, want 1", a.JuzNumber)
	}
}
