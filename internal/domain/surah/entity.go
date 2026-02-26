package surah

type Surah struct {
	ID                  int    `json:"id"`
	Number              int    `json:"number"`
	NameArabic          string `json:"name_arabic"`
	NameLatin           string `json:"name_latin"`
	NameTransliteration string `json:"name_transliteration"`
	NumberOfAyahs       int    `json:"number_of_ayahs"`
	RevelationType      string `json:"revelation_type"`
}
