package repository_test

import (
	"context"
	"quran-api-go/internal/repository"
	"testing"
)

var createTableAyah = `
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
CREATE INDEX idx_ayahs_surah_id ON ayahs (surah_id);
CREATE INDEX idx_ayahs_juz_number ON ayahs (juz_number);
`

var seedTableAyah = `
INSERT INTO ayahs
(id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, sajda_type, revelation_type)
VALUES
(1, 1, 1,
'بِسۡمِ ٱللَّهِ ٱلرَّحۡمَٰنِ ٱلرَّحِيمِ',
'Dengan nama Allah Yang Maha Pengasih, Maha Penyayang',
'In the name of Allah, the Entirely Merciful, the Especially Merciful',
1, '' , 'meccan'),

(2, 1, 2,
'ٱلۡحَمۡدُ لِلَّهِ رَبِّ ٱلۡعَٰلَمِينَ',
'Segala puji bagi Allah, Tuhan seluruh alam',
'[All] praise is [due] to Allah, Lord of the worlds',
1, '' , 'meccan'),

(3, 1, 3,
'ٱلرَّحۡمَٰنِ ٱلرَّحِيمِ',
'Yang Maha Pengasih, Maha Penyayang',
'The Entirely Merciful, the Especially Merciful',
1, '' , 'meccan'),

(4, 1, 4,
'مَٰلِكِ يَوۡمِ ٱلدِّينِ',
'Pemilik hari pembalasan',
'Sovereign of the Day of Recompense',
1, '' , 'meccan'),

(5, 1, 5,
'إِيَّاكَ نَعۡبُدُ وَإِيَّاكَ نَسۡتَعِينُ',
'Hanya kepada Engkaulah kami menyembah dan hanya kepada Engkaulah kami mohon pertolongan',
'It is You we worship and You we ask for help',
1, '' , 'meccan'),

(6, 1, 6,
'ٱهۡدِنَا ٱلصِّرَٰطَ ٱلۡمُسۡتَقِيمَ',
'Tunjukilah kami jalan yang lurus',
'Guide us to the straight path',
1, '' , 'meccan'),

(7, 1, 7,
'صِرَٰطَ ٱلَّذِينَ أَنۡعَمۡتَ عَلَيۡهِمۡ غَيۡرِ ٱلۡمَغۡضُوبِ عَلَيۡهِمۡ وَلَا ٱلضَّآلِّينَ',
'(yaitu) jalan orang-orang yang telah Engkau beri nikmat kepadanya; bukan (jalan) mereka yang dimurkai, dan bukan (pula jalan) mereka yang sesat',
'The path of those upon whom You have bestowed favor, not of those who have evoked [Your] anger or of those who are astray',
1, '' , 'meccan');
`

func TestAyahRepository_FindByID_Success(t *testing.T) {
	db := setupTestDB(t, createTableAyah, seedTableAyah)
	repo := repository.NewAyahRepository(db)

	ctx := context.Background()

	ayah, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ayah == nil {
		t.Fatalf("expected surah, got nil")
	}

	if ayah.TranslationIdo != "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang" {
		t.Errorf("expected 'Dengan nama Allah Yang Maha Pengasih, Maha Penyayang', got %s", ayah.TranslationIdo)
	}
}

func TestAyahRepository_FindBySurah_Success(t *testing.T) {
	db := setupTestDB(t, createTableAyah, seedTableAyah)
	repo := repository.NewAyahRepository(db)

	ctx := context.Background()

	ayahs, err := repo.FindBySurah(ctx, 1, 1, 3)
	if err != nil {
		t.Fatalf("failed to get ayahs %v", err)
	}

	if len(ayahs) != 3 {
		t.Fatalf("expected 3 ayahs, got %d", len(ayahs))
	}

}

func TestAyahRepository_FindBySurahAndNumber_Success(t *testing.T) {
	db := setupTestDB(t, createTableAyah, seedTableAyah)
	repo := repository.NewAyahRepository(db)

	ctx := context.Background()

	ayah, err := repo.FindBySurahAndNumber(ctx, 1, 1)
	if err != nil {
		t.Fatalf("failed to get ayah %v", err)
	}

	if ayah == nil {
		t.Fatal("expected ayah, got nil")
	}

	if ayah.TranslationIdo != "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang" {
		t.Fatalf("expected 	'Dengan nama Allah Yang Maha Pengasih, Maha Penyayang', got %s", ayah.TranslationIdo)
	}
}

func TestAyahRepository_FindRandom_Success(t *testing.T) {
	db := setupTestDB(t, createTableAyah, seedTableAyah)
	repo := repository.NewAyahRepository(db)

	ctx := context.Background()

	ayah, err := repo.FindRandom(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get ayah %v", err)
	}

	if ayah == nil {
		t.Fatal("expected ayah, got nil")
	}

	if ayah.SurahID != 1 {
		t.Fatalf("expected surah_id 1, got %d", ayah.SurahID)
	}

	if ayah.NumberInSurah < 1 || ayah.NumberInSurah > 7 {
		t.Fatalf("unexpected ayah number %d", ayah.NumberInSurah)
	}
}

func TestAyahRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t, createTableAyah, seedTableAyah)
	repo := repository.NewAyahRepository(db)

	ctx := context.Background()

	ayah, err := repo.FindByID(ctx, 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ayah != nil {
		t.Fatal("expected nil, got data")
	}
}
