package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"quran-api-go/internal/config"
	"quran-api-go/internal/database"
	"quran-api-go/internal/mcpserver"
	"quran-api-go/internal/repository"
	"quran-api-go/internal/service"
)

func main() {
	cfg := config.Load()

	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	// Logs must go to stderr — stdout is reserved for the MCP stdio transport.
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// When DB_PATH is not set, resolve relative to the binary itself so the
	// server works regardless of what working directory the MCP client uses.
	if os.Getenv("DB_PATH") == "" {
		if exe, err := os.Executable(); err == nil {
			cfg.DBPath = filepath.Join(filepath.Dir(exe), "data", "quran.db")
		}
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	surahRepo := repository.NewSurahRepository(db)
	surahSvc := service.NewSurahService(surahRepo)
	ayahRepo := repository.NewAyahRepository(db)
	ayahSvc := service.NewAyahService(ayahRepo)
	juzRepo := repository.NewJuzRepository(db)
	juzSvc := service.NewJuzService(juzRepo)
	searchRepo := repository.NewSearchRepository(db)
	searchSvc := service.NewSearchService(searchRepo)

	srv := mcpserver.New(cfg.AppVersion, surahSvc, ayahSvc, juzSvc, searchSvc)

	log.Info().Msg("starting MCP server on stdio")
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal().Err(err).Msg("MCP server stopped")
	}
}
