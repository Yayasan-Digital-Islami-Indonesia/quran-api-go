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

func (h *SurahHandler) List(c *gin.Context) {
	surahs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, surahs)
}

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
