package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"quran-api-go/internal/config"
	"quran-api-go/internal/database"
	"quran-api-go/internal/middleware"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logging())
	if cfg.AllowedOrigins != "" {
		r.Use(middleware.CORS(cfg.AllowedOrigins))
	}

	// TODO: wire repositories, services, and handlers as issues are resolved.
	//
	// Pattern for each domain:
	//   repo    := repository.NewSurahRepository(db)
	//   svc     := service.NewSurahService(repo)
	//   handler := handler.NewSurahHandler(svc)
	//
	// Endpoints to register (see docs/prd-quran-api-go-2nd.md):
	//   r.GET("/surah",                    surahHandler.List)        // #8
	//   r.GET("/surah/:id",                surahHandler.Detail)      // #8
	//   r.GET("/surah/:id/ayat",           ayahHandler.BySurah)      // #9
	//   r.GET("/surah/:id/ayat/:number",   ayahHandler.BySurahAndNumber) // #10
	//   r.GET("/ayah/:id",                 ayahHandler.ByGlobalID)   // #11
	//   r.GET("/juz",                      juzHandler.List)          // #15
	//   r.GET("/juz/:number",              juzHandler.Detail)        // #15
	//   r.GET("/search",                   searchHandler.Search)     // #16
	//   r.GET("/random",                   ayahHandler.Random)       // #17
	//   r.GET("/health",                   healthHandler.Health)     // #20
	//   r.GET("/health/ready",             healthHandler.Ready)      // #20

	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Info().Str("addr", addr).Msg("starting server")
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func setupLogger(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
