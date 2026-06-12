package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/pkg/response"
	"quran-api-go/pkg/validator"
)

type AyahHandler struct {
	ayahService  ayah.AyahService
	surahService surah.SurahService
}

type SurahAyahsResponse struct {
	Surah SurahSummaryResponse `json:"surah"`
	Ayahs []AyahListItem       `json:"ayahs"`
}

type SurahSummaryResponse struct {
	ID        int    `json:"id"`
	Number    int    `json:"number"`
	NameLatin string `json:"name_latin"`
}

type AyahListItem struct {
	Number        int     `json:"number"`
	NumberInSurah int     `json:"number_in_surah"`
	TextUthmani   string  `json:"text_uthmani"`
	Translation   string  `json:"translation"`
	Juz           int     `json:"juz"`
	Sajda         *string `json:"sajda"`
}

type AyahDetailResponse struct {
	ID             int                 `json:"id"`
	Number         int                 `json:"number"`
	SurahID        int                 `json:"surah_id"`
	NumberInSurah  int                 `json:"number_in_surah"`
	TextUthmani    string              `json:"text_uthmani"`
	Translation    string              `json:"translation"`
	SurahInfo      AyahDetailSurahInfo `json:"surah_info"`
	Juz            int                 `json:"juz"`
	Sajda          *string             `json:"sajda"`
	RevelationType *string             `json:"revelation_type"`
}

type AyahDetailSurahInfo struct {
	ID        int    `json:"id"`
	NameLatin string `json:"name_latin"`
}

type SajdaListItem struct {
	ID            int    `json:"id"`
	SurahID       int    `json:"surah_id"`
	SurahName     string `json:"surah_name"`
	NumberInSurah int    `json:"number_in_surah"`
	TextUthmani   string `json:"text_uthmani"`
	Translation   string `json:"translation"`
	Juz           int    `json:"juz"`
	SajdaType     string `json:"sajda_type"`
}

func NewAyahHandler(ayahService ayah.AyahService, surahService surah.SurahService) *AyahHandler {
	return &AyahHandler{ayahService: ayahService, surahService: surahService}
}

// BySurah godoc
// @Summary     Get ayahs by surah
// @Description Get ayahs from a specific surah with optional range filtering
// @Tags        Ayah
// @Produce     json
// @Param       id    path     int     true   "Surah ID (1-114)"  minimum(1)  maximum(114)
// @Param       from  query    int     false  "Start ayah number (must use with 'to')"
// @Param       to    query    int     false  "End ayah number (must use with 'from')"
// @Param       lang  query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200   {object} response.SuccessResponse{data=SurahAyahsResponse}
// @Failure     400   {object} response.ErrorResponse
// @Failure     404   {object} response.ErrorResponse
// @Failure     500   {object} response.ErrorResponse
// @Router      /surah/{id}/ayah [get]
func (h *AyahHandler) BySurah(c *gin.Context) {
	surahIDParam, err := validator.ValidateIDParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid surah id")
		return
	}
	surahID, _ := strconv.Atoi(surahIDParam)
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	sur, err := h.surahService.GetByID(c.Request.Context(), surahID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "surah not found")
			return
		}
		response.InternalError(c)
		return
	}
	if sur == nil {
		response.NotFound(c, "surah not found")
		return
	}
	from, to, err := parseAyahRange(c.Query("from"), c.Query("to"), sur.NumberOfAyahs)
	if err != nil {
		response.BadRequest(c, "invalid ayah range")
		return
	}
	ayahs, err := h.ayahService.GetBySurah(c.Request.Context(), surahID, from, to)
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, newSurahAyahsResponse(*sur, ayahs, lang))
}

// Detail godoc
// @Summary     Get ayah by global ID
// @Description Get a specific ayah by its global ID (1-6236)
// @Tags        Ayah
// @Produce     json
// @Param       id    path     int     true   "Global ayah ID (1-6236)"  minimum(1)  maximum(6236)
// @Param       lang  query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200   {object} response.SuccessResponse{data=AyahDetailResponse}
// @Failure     400   {object} response.ErrorResponse
// @Failure     404   {object} response.ErrorResponse
// @Failure     500   {object} response.ErrorResponse
// @Router      /ayah/{id} [get]
func (h *AyahHandler) Detail(c *gin.Context) {
	ayahID, err := parseIDParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid ayah id")
		return
	}
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	ay, err := h.ayahService.GetByID(c.Request.Context(), ayahID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "ayah not found")
			return
		}
		response.InternalError(c)
		return
	}
	if ay == nil {
		response.NotFound(c, "ayah not found")
		return
	}
	h.respondWithAyahDetail(c, *ay, lang)
}

// BySurahAndNumber godoc
// @Summary     Get ayah by surah and number
// @Description Get a specific ayah by its surah ID and number within that surah
// @Tags        Ayah
// @Produce     json
// @Param       id      path     int     true   "Surah ID (1-114)"  minimum(1)  maximum(114)
// @Param       number  path     int     true   "Ayah number within the surah"  minimum(1)
// @Param       lang    query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200     {object} response.SuccessResponse{data=AyahDetailResponse}
// @Failure     400     {object} response.ErrorResponse
// @Failure     404     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Router      /surah/{id}/ayah/{number} [get]
func (h *AyahHandler) BySurahAndNumber(c *gin.Context) {
	surahID, err := parseIDParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid surah id")
		return
	}
	number, err := parseIDParam(c.Param("number"))
	if err != nil {
		response.BadRequest(c, "invalid ayah number")
		return
	}
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	ay, err := h.ayahService.GetBySurahAndNumber(c.Request.Context(), surahID, number)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "ayah not found")
			return
		}
		response.InternalError(c)
		return
	}
	if ay == nil {
		response.NotFound(c, "ayah not found")
		return
	}
	h.respondWithAyahDetail(c, *ay, lang)
}

// RandomAyah godoc
// @Summary     Get random ayah
// @Description Get a random ayah, optionally filtered by surah
// @Tags        Ayah
// @Produce     json
// @Param       surah_id  query    int     false  "Filter by surah ID (0 = any)"  minimum(0)  default(0)
// @Param       lang      query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200       {object} response.SuccessResponse{data=AyahDetailResponse}
// @Failure     400       {object} response.ErrorResponse
// @Failure     404       {object} response.ErrorResponse
// @Failure     500       {object} response.ErrorResponse
// @Router      /random [get]
func (h *AyahHandler) RandomAyah(c *gin.Context) {
	surahIDParam := c.DefaultQuery("surah_id", "0")
	surahID, err := strconv.Atoi(surahIDParam)
	if err != nil {
		surahID = 0
	}
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	ay, err := h.ayahService.GetRandom(c.Request.Context(), surahID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "ayah not found")
			return
		}
		response.InternalError(c)
		return
	}
	if ay == nil {
		response.NotFound(c, "ayah not found")
		return
	}
	h.respondWithAyahDetail(c, *ay, lang)
}

// Sajda godoc
// @Summary     List sajda ayahs
// @Description Get all 15 sajda tilawah ayahs in the Quran
// @Tags        Ayah
// @Produce     json
// @Param       lang  query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200   {object} response.SuccessResponse{data=[]SajdaListItem}
// @Failure     400   {object} response.ErrorResponse
// @Failure     500   {object} response.ErrorResponse
// @Router      /sajda [get]
func (h *AyahHandler) Sajda(c *gin.Context) {
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	ayahs, err := h.ayahService.GetSajda(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}
	result := make([]SajdaListItem, 0, len(ayahs))
	for _, a := range ayahs {
		translation := a.TranslationIdo
		if lang == "en" {
			translation = a.TranslationEn
		}
		result = append(result, SajdaListItem{
			ID:            a.AyahID,
			SurahID:       a.SurahID,
			SurahName:     a.SurahNameLatin,
			NumberInSurah: a.NumberInSurah,
			TextUthmani:   a.TextUthmani,
			Translation:   translation,
			Juz:           a.JuzNumber,
			SajdaType:     a.SajdaType,
		})
	}
	response.Success(c, result)
}

func parseAyahRange(fromParam, toParam string, maxAyahs int) (int, int, error) {
	if fromParam == "" && toParam == "" {
		return 1, maxAyahs, nil
	}
	if fromParam == "" || toParam == "" {
		return 0, 0, domain.ErrInvalidRangeParam
	}
	if err := validator.ValidateRangeParam(fromParam, toParam); err != nil {
		return 0, 0, err
	}
	from, _ := strconv.Atoi(fromParam)
	to, _ := strconv.Atoi(toParam)
	return from, to, nil
}

func newSurahAyahsResponse(sur surah.Surah, ayahs []ayah.Ayah, lang string) SurahAyahsResponse {
	responseAyahs := make([]AyahListItem, 0, len(ayahs))
	for _, item := range ayahs {
		responseAyahs = append(responseAyahs, AyahListItem{
			Number:        item.ID,
			NumberInSurah: item.NumberInSurah,
			TextUthmani:   item.TextUthmani,
			Translation:   translationByLang(item, lang),
			Juz:           item.JuzNumber,
			Sajda:         item.SajdaType,
		})
	}
	return SurahAyahsResponse{
		Surah: SurahSummaryResponse{ID: sur.ID, Number: sur.Number, NameLatin: sur.NameLatin},
		Ayahs: responseAyahs,
	}
}

func translationByLang(item ayah.Ayah, lang string) string {
	if lang == "en" {
		return item.TranslationEn
	}
	return item.TranslationIdo
}

func (h *AyahHandler) respondWithAyahDetail(c *gin.Context, ay ayah.Ayah, lang string) {
	sur, err := h.surahService.GetByID(c.Request.Context(), ay.SurahID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "surah not found")
			return
		}
		response.InternalError(c)
		return
	}
	if sur == nil {
		response.NotFound(c, "surah not found")
		return
	}
	response.Success(c, newAyahDetailResponse(ay, *sur, lang))
}

func newAyahDetailResponse(item ayah.Ayah, sur surah.Surah, lang string) AyahDetailResponse {
	return AyahDetailResponse{
		ID:             item.ID,
		Number:         item.ID,
		SurahID:        item.SurahID,
		NumberInSurah:  item.NumberInSurah,
		TextUthmani:    item.TextUthmani,
		Translation:    translationByLang(item, lang),
		SurahInfo:      AyahDetailSurahInfo{ID: sur.ID, NameLatin: sur.NameLatin},
		Juz:            item.JuzNumber,
		Sajda:          item.SajdaType,
		RevelationType: item.RevelationType,
	}
}

func parseIDParam(raw string) (int, error) {
	validated, err := validator.ValidateIDParam(raw)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(validated)
}
