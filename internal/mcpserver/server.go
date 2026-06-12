package mcpserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/domain/juz"
	"quran-api-go/internal/domain/search"
	"quran-api-go/internal/domain/surah"
)

type server struct {
	surahSvc  surah.SurahService
	ayahSvc   ayah.AyahService
	juzSvc    juz.JuzService
	searchSvc search.SearchService
}

// New builds a configured *mcp.Server with all Quran tools registered.
func New(
	version string,
	surahSvc surah.SurahService,
	ayahSvc ayah.AyahService,
	juzSvc juz.JuzService,
	searchSvc search.SearchService,
) *mcp.Server {
	s := &server{
		surahSvc:  surahSvc,
		ayahSvc:   ayahSvc,
		juzSvc:    juzSvc,
		searchSvc: searchSvc,
	}

	srv := mcp.NewServer(&mcp.Implementation{Name: "quran-api", Version: version}, nil)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_surahs",
		Description: "List all 114 surahs of the Quran with names in Arabic, Latin transliteration, revelation type (Meccan/Medinan), and total ayah count.",
	}, s.listSurahs)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_surah",
		Description: "Get details of a specific surah by its number (1-114).",
	}, s.getSurah)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_ayahs_by_surah",
		Description: "Get all ayahs of a surah in Arabic (Uthmani script) with both Indonesian and English translations. Optionally restrict to a range using from/to ayah numbers within the surah.",
	}, s.getAyahsBySurah)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_ayah",
		Description: "Get a single ayah by its global sequential ID (1-6236) with Arabic text and both translations.",
	}, s.getAyah)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_ayah_by_ref",
		Description: "Get a single ayah using surah number and position within that surah. Example: surah_id=2, number=255 returns Ayat al-Kursi.",
	}, s.getAyahByRef)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "random_ayah",
		Description: "Get a random ayah from the Quran. Set surah_id to 0 for any surah, or a specific number (1-114) to restrict the pick.",
	}, s.randomAyah)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_juz",
		Description: "List all 30 juz (parts) of the Quran with first/last ayah IDs and total ayah count per juz.",
	}, s.listJuz)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_juz",
		Description: "Get details of a specific juz by number (1-30).",
	}, s.getJuz)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_ayahs_by_juz",
		Description: "Get all ayahs within a juz, paginated. Returns Arabic text with both Indonesian and English translations.",
	}, s.getAyahsByJuz)

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search_quran",
		Description: "Full-text search across Quran translations using SQLite FTS5. Use lang='id' to search Indonesian (default) or lang='en' for English. Optionally filter by surah_id or juz number.",
	}, s.searchQuran)

	return srv
}

// ─── Input types ──────────────────────────────────────────────────────────────

type emptyInput struct{}

type getSurahInput struct {
	ID int `json:"id" jsonschema:"Surah number (1-114)"`
}

type getAyahsBySurahInput struct {
	SurahID int `json:"surah_id" jsonschema:"Surah number (1-114)"`
	From    int `json:"from"     jsonschema:"First ayah number to return within the surah; 0 means start from the beginning"`
	To      int `json:"to"       jsonschema:"Last ayah number to return within the surah; 0 means return until the last ayah"`
}

type getAyahInput struct {
	ID int `json:"id" jsonschema:"Global ayah ID (1-6236)"`
}

type getAyahByRefInput struct {
	SurahID int `json:"surah_id" jsonschema:"Surah number (1-114)"`
	Number  int `json:"number"   jsonschema:"Ayah position within the surah (starts at 1)"`
}

type randomAyahInput struct {
	SurahID int `json:"surah_id" jsonschema:"Restrict random pick to this surah number; 0 means any surah"`
}

type getJuzInput struct {
	Number int `json:"number" jsonschema:"Juz number (1-30)"`
}

type getAyahsByJuzInput struct {
	JuzNumber int `json:"juz_number" jsonschema:"Juz number (1-30)"`
	Page      int `json:"page"       jsonschema:"Page number (default 1)"`
	Limit     int `json:"limit"      jsonschema:"Items per page (default 20, max 100)"`
}

type searchInput struct {
	Query   string `json:"query"    jsonschema:"Search keyword in the translation, e.g. 'sabar', 'patience', 'rahman'"`
	Lang    string `json:"lang"     jsonschema:"Translation language to search: 'id' for Indonesian (default) or 'en' for English"`
	SurahID int    `json:"surah_id" jsonschema:"Restrict search to this surah number; 0 means all surahs"`
	Juz     int    `json:"juz"      jsonschema:"Restrict search to this juz number; 0 means all juz"`
	Page    int    `json:"page"     jsonschema:"Page number (default 1)"`
	Limit   int    `json:"limit"    jsonschema:"Items per page (default 20, max 100)"`
}

// ─── Output types (must all be structs — SDK requires JSON schema type "object") ──

type listSurahsOutput struct {
	Surahs []surah.Surah `json:"surahs"`
}

type getSurahOutput struct {
	surah.Surah
}

type ayahsOutput struct {
	Ayahs []ayah.Ayah `json:"ayahs"`
}

type ayahOutput struct {
	ayah.Ayah
}

type listJuzOutput struct {
	Juz []juz.Juz `json:"juz"`
}

type getJuzOutput struct {
	juz.Juz
}

type juzAyahsOutput struct {
	JuzNumber  int           `json:"juz_number"`
	TotalAyahs int           `json:"total_ayahs"`
	Ayahs      []juz.JuzAyah `json:"ayahs"`
}

type searchOutput struct {
	Results []search.Result `json:"results"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (s *server) listSurahs(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, listSurahsOutput, error) {
	surahs, err := s.surahSvc.GetAll(ctx)
	if err != nil {
		return nil, listSurahsOutput{}, err
	}
	return nil, listSurahsOutput{Surahs: surahs}, nil
}

func (s *server) getSurah(ctx context.Context, _ *mcp.CallToolRequest, in getSurahInput) (*mcp.CallToolResult, getSurahOutput, error) {
	if in.ID < 1 || in.ID > 114 {
		return nil, getSurahOutput{}, fmt.Errorf("surah id must be between 1 and 114, got %d", in.ID)
	}
	result, err := s.surahSvc.GetByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, getSurahOutput{}, fmt.Errorf("surah %d not found", in.ID)
		}
		return nil, getSurahOutput{}, err
	}
	return nil, getSurahOutput{*result}, nil
}

func (s *server) getAyahsBySurah(ctx context.Context, _ *mcp.CallToolRequest, in getAyahsBySurahInput) (*mcp.CallToolResult, ayahsOutput, error) {
	if in.SurahID < 1 || in.SurahID > 114 {
		return nil, ayahsOutput{}, fmt.Errorf("surah_id must be between 1 and 114, got %d", in.SurahID)
	}
	ayahs, err := s.ayahSvc.GetBySurah(ctx, in.SurahID, in.From, in.To)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ayahsOutput{}, fmt.Errorf("surah %d not found", in.SurahID)
		}
		return nil, ayahsOutput{}, err
	}
	return nil, ayahsOutput{Ayahs: ayahs}, nil
}

func (s *server) getAyah(ctx context.Context, _ *mcp.CallToolRequest, in getAyahInput) (*mcp.CallToolResult, ayahOutput, error) {
	if in.ID < 1 || in.ID > 6236 {
		return nil, ayahOutput{}, fmt.Errorf("global ayah id must be between 1 and 6236, got %d", in.ID)
	}
	result, err := s.ayahSvc.GetByID(ctx, in.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ayahOutput{}, fmt.Errorf("ayah with id %d not found", in.ID)
		}
		return nil, ayahOutput{}, err
	}
	return nil, ayahOutput{*result}, nil
}

func (s *server) getAyahByRef(ctx context.Context, _ *mcp.CallToolRequest, in getAyahByRefInput) (*mcp.CallToolResult, ayahOutput, error) {
	if in.SurahID < 1 || in.SurahID > 114 {
		return nil, ayahOutput{}, fmt.Errorf("surah_id must be between 1 and 114, got %d", in.SurahID)
	}
	if in.Number < 1 {
		return nil, ayahOutput{}, fmt.Errorf("number must be >= 1, got %d", in.Number)
	}
	result, err := s.ayahSvc.GetBySurahAndNumber(ctx, in.SurahID, in.Number)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ayahOutput{}, fmt.Errorf("ayah %d:%d not found", in.SurahID, in.Number)
		}
		return nil, ayahOutput{}, err
	}
	return nil, ayahOutput{*result}, nil
}

func (s *server) randomAyah(ctx context.Context, _ *mcp.CallToolRequest, in randomAyahInput) (*mcp.CallToolResult, ayahOutput, error) {
	result, err := s.ayahSvc.GetRandom(ctx, in.SurahID)
	if err != nil {
		return nil, ayahOutput{}, err
	}
	return nil, ayahOutput{*result}, nil
}

func (s *server) listJuz(ctx context.Context, _ *mcp.CallToolRequest, _ emptyInput) (*mcp.CallToolResult, listJuzOutput, error) {
	juzs, err := s.juzSvc.GetAll(ctx)
	if err != nil {
		return nil, listJuzOutput{}, err
	}
	return nil, listJuzOutput{Juz: juzs}, nil
}

func (s *server) getJuz(ctx context.Context, _ *mcp.CallToolRequest, in getJuzInput) (*mcp.CallToolResult, getJuzOutput, error) {
	if in.Number < 1 || in.Number > 30 {
		return nil, getJuzOutput{}, fmt.Errorf("juz number must be between 1 and 30, got %d", in.Number)
	}
	result, err := s.juzSvc.GetByNumber(ctx, in.Number)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, getJuzOutput{}, fmt.Errorf("juz %d not found", in.Number)
		}
		return nil, getJuzOutput{}, err
	}
	return nil, getJuzOutput{*result}, nil
}

func (s *server) getAyahsByJuz(ctx context.Context, _ *mcp.CallToolRequest, in getAyahsByJuzInput) (*mcp.CallToolResult, juzAyahsOutput, error) {
	if in.JuzNumber < 1 || in.JuzNumber > 30 {
		return nil, juzAyahsOutput{}, fmt.Errorf("juz_number must be between 1 and 30, got %d", in.JuzNumber)
	}
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Limit < 1 {
		in.Limit = 20
	}
	if in.Limit > 100 {
		in.Limit = 100
	}
	offset := (in.Page - 1) * in.Limit

	j, err := s.juzSvc.GetByNumber(ctx, in.JuzNumber)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, juzAyahsOutput{}, fmt.Errorf("juz %d not found", in.JuzNumber)
		}
		return nil, juzAyahsOutput{}, err
	}

	ayahs, err := s.juzSvc.GetAyahsByJuz(ctx, in.JuzNumber, in.Limit, offset)
	if err != nil {
		return nil, juzAyahsOutput{}, err
	}

	return nil, juzAyahsOutput{
		JuzNumber:  j.JuzNumber,
		TotalAyahs: j.TotalAyahs,
		Ayahs:      ayahs,
	}, nil
}

func (s *server) searchQuran(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, searchOutput, error) {
	if in.Query == "" {
		return nil, searchOutput{}, fmt.Errorf("query must not be empty")
	}
	if in.Lang == "" {
		in.Lang = "id"
	}
	if in.Lang != "id" && in.Lang != "en" {
		return nil, searchOutput{}, fmt.Errorf("lang must be 'id' or 'en', got %q", in.Lang)
	}
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Limit < 1 {
		in.Limit = 20
	}
	if in.Limit > 100 {
		in.Limit = 100
	}

	results, total, err := s.searchSvc.Search(ctx, search.Params{
		Query:   in.Query,
		Lang:    in.Lang,
		SurahID: in.SurahID,
		Juz:     in.Juz,
		Page:    in.Page,
		Limit:   in.Limit,
	})
	if err != nil {
		return nil, searchOutput{}, err
	}

	return nil, searchOutput{
		Results: results,
		Total:   total,
		Page:    in.Page,
		Limit:   in.Limit,
	}, nil
}
