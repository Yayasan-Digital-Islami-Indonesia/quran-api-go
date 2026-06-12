package handler

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

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

	// Read the swaggo-generated spec
	content, err := os.ReadFile("./docs/api-reference/openapi.yaml")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Determine the correct scheme and host from the request (handles proxies)
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	} else if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}

	productionURL := scheme + "://" + host
	yaml := strings.ReplaceAll(string(content), "http://localhost:8080", productionURL)
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
