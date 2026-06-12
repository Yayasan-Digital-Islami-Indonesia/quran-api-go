package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain"
	"quran-api-go/internal/domain/surah"
	"quran-api-go/pkg/response"
)

type SurahHandler struct {
	service surah.SurahService
}

func NewSurahHandler(service surah.SurahService) *SurahHandler {
	return &SurahHandler{service: service}
}

// List godoc
// @Summary     List all surahs
// @Description Get a list of all 114 surahs, optionally filtered by revelation type
// @Tags        Surah
// @Produce     json
// @Param       type  query    string  false  "Filter by revelation type"  Enums(meccan, medinan)
// @Success     200   {object} response.SuccessResponse{data=[]surah.Surah}
// @Failure     400   {object} response.ErrorResponse
// @Failure     500   {object} response.ErrorResponse
// @Router      /surah [get]
func (h *SurahHandler) List(c *gin.Context) {
	revelationType := c.Query("type")
	if revelationType != "" {
		if revelationType != "meccan" && revelationType != "medinan" {
			response.BadRequest(c, "type must be 'meccan' or 'medinan'")
			return
		}
		surahs, err := h.service.GetByRevelationType(c.Request.Context(), revelationType)
		if err != nil {
			response.InternalError(c)
			return
		}
		response.Success(c, surahs)
		return
	}

	surahs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, surahs)
}

// Detail godoc
// @Summary     Get surah by ID
// @Description Get detailed information about a specific surah
// @Tags        Surah
// @Produce     json
// @Param       id   path     int  true  "Surah ID (1-114)"  minimum(1)  maximum(114)
// @Success     200  {object} response.SuccessResponse{data=surah.Surah}
// @Failure     400  {object} response.ErrorResponse
// @Failure     404  {object} response.ErrorResponse
// @Failure     500  {object} response.ErrorResponse
// @Router      /surah/{id} [get]
func (h *SurahHandler) Detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid surah id")
		return
	}

	s, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, "surah not found")
			return
		}
		response.InternalError(c)
		return
	}

	response.Success(c, s)
}
