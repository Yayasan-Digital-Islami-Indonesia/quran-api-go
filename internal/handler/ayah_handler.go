package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/ayah"
	"quran-api-go/internal/domain/surah"
)

type AyahHandler struct {
	ayahService  ayah.AyahService
	surahService surah.SurahService
}

func NewAyahHandler(ayahService ayah.AyahService, surahService surah.SurahService) *AyahHandler {
	return &AyahHandler{
		ayahService:  ayahService,
		surahService: surahService,
	}
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
