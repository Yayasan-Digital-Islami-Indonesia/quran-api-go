package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/search"
	"quran-api-go/pkg/response"
	"quran-api-go/pkg/validator"
)

type SearchHandler struct {
	service search.SearchService
}

type SearchResponse struct {
	Query   string          `json:"query"`
	Results []search.Result `json:"results"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

func NewSearchHandler(service search.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Search godoc
// @Summary     Search ayahs
// @Description Full-text search across Quran ayahs (Arabic, Indonesian, English) using FTS5
// @Tags        Search
// @Produce     json
// @Param       q         query    string  true   "Search query"
// @Param       lang      query    string  false  "Translation language"  Enums(id, en)  default(id)
// @Param       surah_id  query    int     false  "Filter by surah ID"  minimum(1)  maximum(114)
// @Param       juz       query    int     false  "Filter by juz number"  minimum(1)  maximum(30)
// @Param       page      query    int     false  "Page number"  minimum(1)  default(1)
// @Param       limit     query    int     false  "Items per page"  minimum(1)  maximum(100)  default(20)
// @Success     200       {object} response.SuccessResponse{data=SearchResponse}
// @Failure     400       {object} response.ErrorResponse
// @Failure     500       {object} response.ErrorResponse
// @Router      /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "query parameter 'q' is required")
		return
	}

	lang, err := validator.ValidateLang(c.Query("lang"))
	if err != nil {
		response.BadRequest(c, "lang must be 'id' or 'en'")
		return
	}

	surahID, _ := strconv.Atoi(c.Query("surah_id"))
	juz, _ := strconv.Atoi(c.Query("juz"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	params := search.Params{
		Query:   query,
		Lang:    lang,
		SurahID: surahID,
		Juz:     juz,
		Page:    page,
		Limit:   limit,
	}

	results, total, err := h.service.Search(c.Request.Context(), params)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, SearchResponse{
		Query:   query,
		Results: results,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}
