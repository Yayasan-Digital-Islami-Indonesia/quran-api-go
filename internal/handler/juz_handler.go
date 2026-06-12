package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/juz"
	"quran-api-go/pkg/pagination"
	"quran-api-go/pkg/response"
	"quran-api-go/pkg/validator"
)

type JuzHandler struct {
	service juz.JuzService
}

type JuzAyahListItem struct {
	ID            int    `json:"id"`
	SurahID       int    `json:"surah_id"`
	SurahName     string `json:"surah_name"`
	NumberInSurah int    `json:"number_in_surah"`
	TextUthmani   string `json:"text_uthmani"`
	Translation   string `json:"translation"`
	JuzNumber     int    `json:"juz_number"`
}

type JuzAyahsResponse struct {
	Juz   JuzInfo           `json:"juz"`
	Ayahs []JuzAyahListItem `json:"ayahs"`
}

type JuzInfo struct {
	JuzNumber  int `json:"juz_number"`
	TotalAyahs int `json:"total_ayahs"`
}

func NewJuzHandler(service juz.JuzService) *JuzHandler {
	return &JuzHandler{service: service}
}

// List godoc
// @Summary     List all juz
// @Description Get a list of all 30 juz (parts) of the Quran
// @Tags        Juz
// @Produce     json
// @Success     200  {object} response.SuccessResponse{data=[]juz.Juz}
// @Failure     500  {object} response.ErrorResponse
// @Router      /juz [get]
func (h *JuzHandler) List(c *gin.Context) {
	juzs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Success(c, juzs)
}

// Detail godoc
// @Summary     Get juz by number
// @Description Get detailed information about a specific juz
// @Tags        Juz
// @Produce     json
// @Param       number  path     int  true  "Juz number (1-30)"  minimum(1)  maximum(30)
// @Success     200     {object} response.SuccessResponse{data=juz.Juz}
// @Failure     400     {object} response.ErrorResponse
// @Failure     404     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Router      /juz/{number} [get]
func (h *JuzHandler) Detail(c *gin.Context) {
	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		response.BadRequest(c, "invalid juz number")
		return
	}
	j, err := h.service.GetByNumber(c.Request.Context(), number)
	if err != nil {
		response.InternalError(c)
		return
	}
	if j == nil {
		response.NotFound(c, "juz not found")
		return
	}
	response.Success(c, j)
}

// Ayahs godoc
// @Summary     Get ayahs by juz
// @Description Get all ayahs from a specific juz with pagination
// @Tags        Juz
// @Produce     json
// @Param       number  path     int     true   "Juz number (1-30)"  minimum(1)  maximum(30)
// @Param       page    query    int     false  "Page number"  minimum(1)  default(1)
// @Param       limit   query    int     false  "Items per page"  minimum(1)  maximum(100)  default(50)
// @Param       lang    query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Success     200     {object} response.SuccessResponse{data=JuzAyahsResponse}
// @Failure     400     {object} response.ErrorResponse
// @Failure     404     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Router      /juz/{number}/ayah [get]
func (h *JuzHandler) Ayahs(c *gin.Context) {
	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		response.BadRequest(c, "invalid juz number")
		return
	}
	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}
	params := pagination.Parse(c.Query("page"), c.Query("limit"))
	j, err := h.service.GetByNumber(c.Request.Context(), number)
	if err != nil {
		response.InternalError(c)
		return
	}
	if j == nil {
		response.NotFound(c, "juz not found")
		return
	}
	ayahs, err := h.service.GetAyahsByJuz(c.Request.Context(), number, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "ayahs not found")
			return
		}
		response.InternalError(c)
		return
	}
	response.Success(c, JuzAyahsResponse{
		Juz:   JuzInfo{JuzNumber: j.JuzNumber, TotalAyahs: j.TotalAyahs},
		Ayahs: newJuzAyahsResponse(ayahs, lang),
	})
}

// Surahs godoc
// @Summary     Get surahs by juz
// @Description Get all surahs that appear in a specific juz
// @Tags        Juz
// @Produce     json
// @Param       number  path     int  true  "Juz number (1-30)"  minimum(1)  maximum(30)
// @Success     200     {object} response.SuccessResponse{data=[]juz.JuzSurah}
// @Failure     400     {object} response.ErrorResponse
// @Failure     404     {object} response.ErrorResponse
// @Failure     500     {object} response.ErrorResponse
// @Router      /juz/{number}/surah [get]
func (h *JuzHandler) Surahs(c *gin.Context) {
	number, err := strconv.Atoi(c.Param("number"))
	if err != nil || number < 1 || number > 30 {
		response.BadRequest(c, "invalid juz number")
		return
	}
	surahs, err := h.service.GetSurahsByJuz(c.Request.Context(), number)
	if err != nil {
		response.InternalError(c)
		return
	}
	if surahs == nil {
		response.NotFound(c, "juz not found")
		return
	}
	response.Success(c, surahs)
}

func newJuzAyahsResponse(ayahs []juz.JuzAyah, lang string) []JuzAyahListItem {
	result := make([]JuzAyahListItem, 0, len(ayahs))
	for _, item := range ayahs {
		translation := item.TranslationIdo
		if lang == "en" {
			translation = item.TranslationEn
		}
		result = append(result, JuzAyahListItem{
			ID:            item.AyahID,
			SurahID:       item.SurahID,
			SurahName:     item.SurahNameLatin,
			NumberInSurah: item.NumberInSurah,
			TextUthmani:   item.TextUthmani,
			Translation:   translation,
			JuzNumber:     item.JuzNumber,
		})
	}
	return result
}
