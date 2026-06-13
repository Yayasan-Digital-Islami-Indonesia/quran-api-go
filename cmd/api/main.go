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
	_ "quran-api-go/docs"
)

// @title           Quran API Go
// @version         1.0.0
// @description     Internal RESTful API serving Al-Quran data (Arabic text, Indonesian & English translations) for the Ilmunara super app.
// @description
// @description     ---
// @description
// @description     ## MCP Server
// @description
// @description     API ini dilengkapi dengan **MCP (Model Context Protocol) server** sehingga bisa digunakan langsung dari AI assistant seperti Claude, Cursor, dan tools lainnya.
// @description
// @description     ### Koneksi
// @description
// @description     | | |
// @description     |---|---|
// @description     | **URL** | `https://quran.api.digitalislami.id/mcp` |
// @description     | **Transport** | Streamable HTTP |
// @description     | **Mode** | Stateless (tidak perlu session) |
// @description
// @description     ### Setup Claude Desktop
// @description
// @description     Tambahkan ke file `claude_desktop_config.json`:
// @description
// @description     ```json
// @description     {
// @description       "mcpServers": {
// @description         "quran": {
// @description           "type": "http",
// @description           "url": "https://quran.api.digitalislami.id/mcp"
// @description         }
// @description       }
// @description     }
// @description     ```
// @description
// @description     ### Tools yang Tersedia
// @description
// @description     | Tool | Deskripsi |
// @description     |---|---|
// @description     | `list_surahs` | Daftar semua 114 surah |
// @description     | `get_surah` | Detail surah berdasarkan ID |
// @description     | `get_ayahs_by_surah` | Ayat-ayat dalam surah tertentu |
// @description     | `get_ayah` | Ayat berdasarkan ID global |
// @description     | `get_ayah_by_ref` | Ayat berdasarkan nomor surah dan ayat |
// @description     | `random_ayah` | Ayat acak |
// @description     | `list_juz` | Daftar semua 30 juz |
// @description     | `get_juz` | Detail juz tertentu |
// @description     | `get_ayahs_by_juz` | Ayat-ayat dalam juz tertentu |
// @description     | `search_quran` | Pencarian full-text (Arab, Indonesia, Inggris) |
// @description
// @description     ### Contoh Penggunaan
// @description
// @description     Setelah terhubung, kamu bisa langsung tanya ke AI:
// @description
// @description     > *"Tampilkan ayat pertama surah Al-Baqarah beserta terjemahannya"*
// @description
// @description     > *"Cari ayat yang mengandung kata 'sabar' dalam terjemahan Indonesia"*
// @description
// @description     > *"Surah apa saja yang ada di Juz 30?"*
// @host            localhost:8080
// @BasePath        /
// @schemes         http https
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
	r.GET("/sajda", ayahHandler.Sajda)
	r.GET("/juz", juzHandler.List)
	r.GET("/juz/:number", juzHandler.Detail)
	r.GET("/juz/:number/ayah", juzHandler.Ayahs)
	r.GET("/juz/:number/surah", juzHandler.Surahs)
	r.GET("/search", searchHandler.Search)

	// MCP endpoint with per-route CORS so browser-based clients (MCP Inspector,
	// Claude.ai web, etc.) work regardless of the global ALLOWED_ORIGINS value.
	mcpOrigin := cfg.AllowedOrigins
	if mcpOrigin == "" {
		mcpOrigin = "*" // MCP is a public read-only endpoint
	}
	mcpCORS := middleware.CORS(mcpOrigin)
	r.OPTIONS("/mcp", mcpCORS)
	r.POST("/mcp", mcpCORS, gin.WrapH(mcpHandler))
	r.GET("/mcp", mcpCORS, gin.WrapH(mcpHandler))

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
