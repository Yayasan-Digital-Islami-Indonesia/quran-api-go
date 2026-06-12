package ayah

type Ayah struct {
	ID             int     `json:"id"`
	SurahID        int     `json:"surah_id"`
	NumberInSurah  int     `json:"number_in_surah"`
	TextUthmani    string  `json:"text_uthmani"`
	TranslationIdo string  `json:"translation_indo"`
	TranslationEn  string  `json:"translation_en"`
	JuzNumber      int     `json:"juz_number"`
	SajdaType      *string `json:"sajda_type"`
	RevelationType *string `json:"revelation_type"`
}

// SajdaAyah is an ayah row joined with its surah name, used in /sajda responses.
type SajdaAyah struct {
	AyahID         int    `json:"id"`
	SurahID        int    `json:"surah_id"`
	SurahNameLatin string `json:"surah_name_latin"`
	NumberInSurah  int    `json:"number_in_surah"`
	TextUthmani    string `json:"text_uthmani"`
	TranslationIdo string `json:"translation_indo"`
	TranslationEn  string `json:"translation_en"`
	JuzNumber      int    `json:"juz_number"`
	SajdaType      string `json:"sajda_type"`
}
