package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"quran-api-go/internal/config"
	"quran-api-go/internal/database"
	"quran-api-go/internal/handler"
	"quran-api-go/internal/mcpserver"
	"quran-api-go/internal/middleware"
	"quran-api-go/internal/repository"
	"quran-api-go/internal/service"
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

	healthCheckRepo := repository.NewHealthCheckRepository(db)
	healthCheckService := service.NewHealthCheckService(healthCheckRepo)
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckService)
	surahRepo := repository.NewSurahRepository(db)
	surahService := service.NewSurahService(surahRepo)
	surahHandler := handler.NewSurahHandler(surahService)
	ayahRepo := repository.NewAyahRepository(db)
	ayahService := service.NewAyahService(ayahRepo)
	ayahHandler := handler.NewAyahHandler(ayahService, surahService)
	juzRepo := repository.NewJuzRepository(db)
	juzService := service.NewJuzService(juzRepo)
	juzHandler := handler.NewJuzHandler(juzService)
	searchRepo := repository.NewSearchRepository(db)
	searchService := service.NewSearchService(searchRepo)
	searchHandler := handler.NewSearchHandler(searchService)
	docsHandler := handler.NewDocsHandler()

	mcpSrv := mcpserver.New(cfg.AppVersion, surahService, ayahService, juzService, searchService)
	mcpHandler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return mcpSrv
	}, &mcp.StreamableHTTPOptions{Stateless: true})

	r.GET("/", func(c *gin.Context) { c.Redirect(301, "/docs") })
	r.GET("/health", healthCheckHandler.HealthCheck)
	r.GET("/health/ready", healthCheckHandler.ReadyCheck)
	r.GET("/surah", surahHandler.List)
	r.GET("/surah/:id", surahHandler.Detail)
	r.GET("/ayah/:id", ayahHandler.Detail)
	r.GET("/surah/:id/ayah", ayahHandler.BySurah)
	r.GET("/surah/:id/ayah/:number", ayahHandler.BySurahAndNumber)
	r.GET("/random", ayahHandler.RandomAyah)
	r.GET("/juz", juzHandler.List)
	r.GET("/juz/:number", juzHandler.Detail)
	r.GET("/juz/:number/ayah", juzHandler.Ayahs)
	r.GET("/search", searchHandler.Search)

	// MCP endpoint (Streamable HTTP, stateless)
	r.Any("/mcp", gin.WrapH(mcpHandler))

	// Documentation
	r.GET("/docs", docsHandler.ServeDocs)
	r.GET("/openapi.yaml", docsHandler.ServeOpenAPI)
	r.GET("/static/:filename", docsHandler.ServeStatic)

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
