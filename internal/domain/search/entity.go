package search

// Params holds the inputs for a full-text search query.
type Params struct {
	Query   string
	Lang    string
	SurahID int // 0 = no filter
	Juz     int // 0 = no filter
	Page    int
	Limit   int
}

// Result is a single ayah match returned from a search query.
type Result struct {
	ID            int       `json:"id"`
	SurahID       int       `json:"surah_id"`
	SurahInfo     SurahInfo `json:"surah_info"`
	NumberInSurah int       `json:"number_in_surah"`
	TextUthmani   string    `json:"text_uthmani"`
	Translation   string    `json:"translation"`
	JuzNumber     int       `json:"juz_number"`
}

// SurahInfo is the minimal surah metadata embedded in a search result.
type SurahInfo struct {
	ID        int    `json:"id"`
	NameLatin string `json:"name_latin"`
}
