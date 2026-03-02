package response

// All handlers must use these helpers. Never call c.JSON directly in handlers.
//
// Response formats:
//   Success:  { "data": any, "timestamp": string }
//   Error:    { "error": string, "code": string, "timestamp": string }
//
// Usage:
//   response.Success(c, data)
//   response.NotFound(c, "surah not found")
//   response.BadRequest(c, "invalid lang")
//   response.InternalError(c)

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	success := gin.H{
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339), // This Should Return ISO 8601 Timestamp Format
	}

	c.JSON(http.StatusOK, success)
}

func NotFound(c *gin.Context, message string) {
	error := gin.H{
		"error":     message,
		"code":      "not found",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusNotFound, error)
}

func BadRequest(c *gin.Context, message string) {
	error := gin.H{
		"error":     message,
		"code":      "bad request",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusBadRequest, error)
}

func InternalError(c *gin.Context) {
	error := gin.H{
		"code":      "internal server error",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusInternalServerError, error)
}
