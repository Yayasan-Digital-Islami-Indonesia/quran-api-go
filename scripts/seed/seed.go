package seed

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Surah struct {
	ID                  int    `json:"id"`
	NameArabic          string `json:"name"`
	NameLatin           string `json:"translation"`
	NameTransliteration string `json:"transliteration"`
	RevelationType      string `json:"type"`
	TotalVerses         int    `json:"total_verses"`
	Verses              []Verse
}

type Verse struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

type JuzMeta struct {
	Code int      `json:"code"`
	Data MetaData `json:"data"`
}

type MetaData struct {
	Juzs JuzsMeta `json:"juzs"`
}

type JuzsMeta struct {
	Count      int            `json:"count"`
	References []JuzReference `json:"references"`
}

type JuzReference struct {
	Surah int `json:"surah"`
	Ayah  int `json:"ayah"`
}

type FlatAyah struct {
	ID             int
	SurahID        int
	NumberInSurah  int
	TextUthmani    string
	TranslationID  string
	TranslationEN  string
	JuzNumber      int
	SajdaType      sql.NullString
	RevelationType string
}

type Juz struct {
	ID        int
	JuzNumber int
	FirstAyah int
	LastAyah  int
}

func Run(ctx context.Context, db *sql.DB, dataDir string) error {
	idPath := filepath.Join(dataDir, "quran_id.json")
	enPath := filepath.Join(dataDir, "quran_en.json")
	metaPath := filepath.Join(dataDir, "meta")

	log.Info().Str("path", idPath).Msg("loading quran_id")
	idSurahs, err := loadSurahs(idPath)
	if err != nil {
		return err
	}
	log.Info().Int("count", len(idSurahs)).Msg("quran_id loaded")

	log.Info().Str("path", enPath).Msg("loading quran_en")
	enSurahs, err := loadSurahs(enPath)
	if err != nil {
		return err
	}
	log.Info().Int("count", len(enSurahs)).Msg("quran_en loaded")

	log.Info().Str("path", metaPath).Msg("loading meta")
	meta, err := loadMeta(metaPath)
	if err != nil {
		return err
	}
	log.Info().Int("count", len(meta.Data.Juzs.References)).Msg("meta loaded")

	if err := validateSurahAlignment(idSurahs, enSurahs); err != nil {
		return err
	}

	flatAyahs, juzs, err := buildAyahsAndJuzs(idSurahs, enSurahs, meta)
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err := seedSurahs(ctx, tx, idSurahs); err != nil {
		return err
	}
	if err := seedAyahs(ctx, tx, flatAyahs); err != nil {
		return err
	}
	if err := seedJuzs(ctx, tx, juzs); err != nil {
		return err
	}

	if err := validateCounts(ctx, tx, len(idSurahs), len(flatAyahs), len(juzs)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Info().Msg("seed completed")
	return nil
}

func loadSurahs(path string) ([]Surah, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []Surah
	dec := json.NewDecoder(file)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func loadMeta(path string) (JuzMeta, error) {
	file, err := os.Open(path)
	if err != nil {
		return JuzMeta{}, err
	}
	defer file.Close()

	var data JuzMeta
	dec := json.NewDecoder(file)
	if err := dec.Decode(&data); err != nil {
		return JuzMeta{}, err
	}
	return data, nil
}

func validateSurahAlignment(idSurahs, enSurahs []Surah) error {
	if len(idSurahs) != len(enSurahs) {
		return fmt.Errorf("surah count mismatch: id=%d en=%d", len(idSurahs), len(enSurahs))
	}

	for i := range idSurahs {
		idS := idSurahs[i]
		enS := enSurahs[i]
		if idS.ID != enS.ID || len(idS.Verses) != len(enS.Verses) {
			return fmt.Errorf("surah mismatch at index %d: id=%d en=%d verses=%d/%d", i, idS.ID, enS.ID, len(idS.Verses), len(enS.Verses))
		}
	}

	return nil
}

func buildAyahsAndJuzs(idSurahs, enSurahs []Surah, meta JuzMeta) ([]FlatAyah, []Juz, error) {
	refs := meta.Data.Juzs.References
	if len(refs) != 30 {
		return nil, nil, fmt.Errorf("unexpected juz references count: %d", len(refs))
	}

	startIndex := map[[2]int]int{}
	for i, ref := range refs {
		key := [2]int{ref.Surah, ref.Ayah}
		startIndex[key] = i + 1 // juz number
	}

	flat := make([]FlatAyah, 0, 6236)
	globalID := 0
	juzStarts := make([]int, 0, 30)

	for sIdx := range idSurahs {
		idS := idSurahs[sIdx]
		enS := enSurahs[sIdx]
		for vIdx := range idS.Verses {
			globalID++
			key := [2]int{idS.ID, vIdx + 1}
			if j, ok := startIndex[key]; ok {
				juzStarts = append(juzStarts, globalID)
				_ = j
			}

			flat = append(flat, FlatAyah{
				ID:             globalID,
				SurahID:        idS.ID,
				NumberInSurah:  vIdx + 1,
				TextUthmani:    idS.Verses[vIdx].Text,
				TranslationID:  idS.Verses[vIdx].Translation,
				TranslationEN:  enS.Verses[vIdx].Translation,
				RevelationType: idS.RevelationType,
				SajdaType:      sql.NullString{},
			})
		}
	}

	if len(flat) == 0 {
		return nil, nil, errors.New("no ayahs built")
	}

	if len(juzStarts) != 30 {
		return nil, nil, fmt.Errorf("juz starts mismatch: %d", len(juzStarts))
	}

	juzs := make([]Juz, 0, 30)
	for i := 0; i < 30; i++ {
		start := juzStarts[i]
		end := 0
		if i < 29 {
			end = juzStarts[i+1] - 1
		} else {
			end = globalID
		}
		juzs = append(juzs, Juz{
			ID:        i + 1,
			JuzNumber: i + 1,
			FirstAyah: start,
			LastAyah:  end,
		})
	}

	// assign juz_number to flat ayahs based on ranges
	juzIdx := 0
	for i := range flat {
		pos := i + 1
		for juzIdx < len(juzs)-1 && pos > juzs[juzIdx].LastAyah {
			juzIdx++
		}
		flat[i].JuzNumber = juzs[juzIdx].JuzNumber
	}

	return flat, juzs, nil
}

func seedSurahs(ctx context.Context, tx *sql.Tx, surahs []Surah) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO surahs (
			id, number, name_arabic, name_latin, name_transliteration, number_of_ayahs, revelation_type
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range surahs {
		if _, err := stmt.ExecContext(ctx, s.ID, s.ID, s.NameArabic, s.NameLatin, s.NameTransliteration, s.TotalVerses, s.RevelationType); err != nil {
			return err
		}
	}

	log.Info().Int("count", len(surahs)).Msg("surahs seeded")
	return nil
}

func seedAyahs(ctx context.Context, tx *sql.Tx, ayahs []FlatAyah) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO ayahs (
			id, surah_id, number_in_surah, text_uthmani, translation_indo, translation_en, juz_number, sajda_type, revelation_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ftsStmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO ayahs_fts (rowid, ayah_id, text_uthmani, translation_indo, translation_en)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer ftsStmt.Close()

	for _, a := range ayahs {
		var sajda any = nil
		if a.SajdaType.Valid {
			sajda = a.SajdaType.String
		}
		if _, err := stmt.ExecContext(ctx, a.ID, a.SurahID, a.NumberInSurah, a.TextUthmani, a.TranslationID, a.TranslationEN, a.JuzNumber, sajda, a.RevelationType); err != nil {
			return err
		}
		if _, err := ftsStmt.ExecContext(ctx, a.ID, a.ID, a.TextUthmani, a.TranslationID, a.TranslationEN); err != nil {
			return err
		}
	}

	log.Info().Int("count", len(ayahs)).Msg("ayahs seeded")
	return nil
}

func seedJuzs(ctx context.Context, tx *sql.Tx, juzs []Juz) error {
	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR REPLACE INTO juzs (id, juz_number, first_ayah_id, last_ayah_id)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, j := range juzs {
		if _, err := stmt.ExecContext(ctx, j.ID, j.JuzNumber, j.FirstAyah, j.LastAyah); err != nil {
			return err
		}
	}

	log.Info().Int("count", len(juzs)).Msg("juzs seeded")
	return nil
}

func validateCounts(ctx context.Context, tx *sql.Tx, surahCount, ayahCount, juzCount int) error {
	var gotSurah, gotAyah, gotJuz int

	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM surahs").Scan(&gotSurah); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM ayahs").Scan(&gotAyah); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM juzs").Scan(&gotJuz); err != nil {
		return err
	}

	if gotSurah != surahCount || gotAyah != ayahCount || gotJuz != juzCount {
		return fmt.Errorf("seed validation failed: surah=%d/%d ayah=%d/%d juz=%d/%d", gotSurah, surahCount, gotAyah, ayahCount, gotJuz, juzCount)
	}

	log.Info().Msg("seed validation passed")
	return nil
}
