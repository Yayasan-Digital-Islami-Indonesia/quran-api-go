package handler

import (
	"errors"
	"log"
	"net/http"
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

func NewAyahHandler(ayahService ayah.AyahService, surahService surah.SurahService) *AyahHandler {
	return &AyahHandler{
		ayahService:  ayahService,
		surahService: surahService,
	}
}

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

func (h *AyahHandler) BySurahAndNumber(c *gin.Context) {
	surahID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid surah id"})
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ayah number"})
		return
	}

	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lang parameter, must be 'id' or 'en'"})
		return
	}

	ay, err := h.ayahService.GetBySurahAndNumber(c.Request.Context(), surahID, number)
	if err != nil {
		log.Printf("error fetching ayah: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ayah"})
		return
	}
	if ay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ayah not found"})
		return
	}

	sur, err := h.surahService.GetByID(c.Request.Context(), surahID)
	if err != nil {
		log.Printf("error fetching surah info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch surah info"})
		return
	}
	if sur == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "surah not found"})
		return
	}

	translation := ay.TranslationIdo
	if lang == "en" {
		translation = ay.TranslationEn
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              ay.ID,
		"surah_id":        ay.SurahID,
		"number_in_surah": ay.NumberInSurah,
		"text_uthmani":    ay.TextUthmani,
		"translation":     translation,
		"surah_info": gin.H{
			"id":         sur.ID,
			"name_latin": sur.NameLatin,
		},
		"juz":             ay.JuzNumber,
		"sajda":           ay.SajdaType,
		"revelation_type": ay.RevelationType,
	})
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
		Surah: SurahSummaryResponse{
			ID:        sur.ID,
			Number:    sur.Number,
			NameLatin: sur.NameLatin,
		},
		Ayahs: responseAyahs,
	}
}

func translationByLang(item ayah.Ayah, lang string) string {
	if lang == "en" {
		return item.TranslationEn
	}

	return item.TranslationIdo
}

func (h *AyahHandler) RandomAyah(c *gin.Context) {
	surahIDParam := c.DefaultQuery("surah_id", "0")
	surahID, err := strconv.Atoi(surahIDParam)
	if err != nil {
		surahID = 0
	}

	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lang parameter, must be 'id' or 'en'"})
		return
	}

	ay, err := h.ayahService.GetRandom(c.Request.Context(), surahID)
	if err != nil {
		log.Printf("error fetching ayah: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ayah"})
		return
	}
	if ay == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ayah not found"})
		return
	}

	if surahID == 0 {
		surahID = ay.SurahID
	}

	sur, err := h.surahService.GetByID(c.Request.Context(), surahID)
	if err != nil {
		log.Printf("error fetching surah info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch surah info"})
		return
	}
	if sur == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "surah not found"})
		return
	}

	translation := ay.TranslationIdo
	if lang == "en" {
		translation = ay.TranslationEn
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              ay.ID,
		"surah_id":        ay.SurahID,
		"number_in_surah": ay.NumberInSurah,
		"text_uthmani":    ay.TextUthmani,
		"translation":     translation,
		"surah_info": gin.H{
			"id":         sur.ID,
			"name_latin": sur.NameLatin,
		},
		"juz":             ay.JuzNumber,
		"sajda":           ay.SajdaType,
		"revelation_type": ay.RevelationType,
	})
}
