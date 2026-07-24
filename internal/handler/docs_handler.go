package handler

import (
	"embed"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

const canonicalProductionBaseURL = "https://quran.api.digitalislami.id"

//go:embed api-reference/openapi.yaml
var openapiSpec []byte

const scalarHTML = `<!doctype html>
<html>
  <head>
    <title>Quran API Go - Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body { margin: 0; padding: 0; }
    </style>
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.yaml"></script>
    <script src="/static/scalar.js"></script>
  </body>
</html>
`

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeDocs(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, scalarHTML)
}

func (h *DocsHandler) ServeOpenAPI(c *gin.Context) {
	c.Header("Content-Type", "text/yaml; charset=utf-8")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-store")

	productionURL := resolveDocsBaseURL(c)
	yaml := string(openapiSpec)

	// Replace localhost URLs with production/base URL
	// OpenAPI format uses either "host: localhost:8080" or server URLs
	if strings.Contains(yaml, "host: localhost:8080") {
		// Extract the host from the production URL
		hostOnly := productionURL
		if idx := strings.Index(hostOnly, "://"); idx >= 0 {
			hostOnly = hostOnly[idx+3:]
		}
		yaml = strings.ReplaceAll(yaml, "host: localhost:8080", "host: "+hostOnly)
	}

	// Also replace any full localhost server URLs
	yaml = strings.ReplaceAll(yaml, "http://localhost:8080", productionURL)
	yaml = strings.ReplaceAll(yaml, "https://localhost:8080", productionURL)

	c.String(http.StatusOK, yaml)
}

// ServeStatic serves embedded static files (Scalar JS)
func (h *DocsHandler) ServeStatic(c *gin.Context) {
	filename := filepath.Base(c.Param("filename"))
	if filename != "scalar.js" {
		c.Status(http.StatusNotFound)
		return
	}

	content, err := staticFiles.ReadFile("static/" + filename)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/javascript", content)
}

func resolveDocsBaseURL(c *gin.Context) string {
	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}

	if host == "" {
		return canonicalProductionBaseURL
	}

	if isLoopbackHost(host) && gin.Mode() == gin.ReleaseMode {
		return canonicalProductionBaseURL
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	} else if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return scheme + "://" + host
}

func isLoopbackHost(host string) bool {
	hostOnly := host
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		hostOnly = parsedHost
	}

	hostOnly = strings.Trim(hostOnly, "[]")
	switch hostOnly {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}
