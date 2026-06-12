package juz

type Juz struct {
	ID          int `json:"id"`
	JuzNumber   int `json:"juz_number"`
	FirstAyahID int `json:"first_ayah_id"`
	LastAyahID  int `json:"last_ayah_id"`
	TotalAyahs  int `json:"total_ayahs"`
}

// JuzAyah is an ayah row joined with its surah name, used in juz detail responses.
type JuzAyah struct {
	AyahID         int    `json:"id"`
	SurahID        int    `json:"surah_id"`
	SurahNameLatin string `json:"surah_name_latin"`
	NumberInSurah  int    `json:"number_in_surah"`
	TextUthmani    string `json:"text_uthmani"`
	TranslationIdo string `json:"translation_indo"`
	TranslationEn  string `json:"translation_en"`
	JuzNumber      int    `json:"juz_number"`
}

// JuzSurah represents a surah that appears within a given juz.
type JuzSurah struct {
	ID                  int    `json:"id"`
	Number              int    `json:"number"`
	NameArabic          string `json:"name_arabic"`
	NameLatin           string `json:"name_latin"`
	NameTransliteration string `json:"name_transliteration"`
	NumberOfAyahs       int    `json:"number_of_ayahs"`
	RevelationType      string `json:"revelation_type"`
}
