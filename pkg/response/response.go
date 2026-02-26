package response

// TODO: implement — see issue #5
//
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

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
	panic("not implemented")
}

func NotFound(c *gin.Context, message string) {
	panic("not implemented")
}

func BadRequest(c *gin.Context, message string) {
	panic("not implemented")
}

func InternalError(c *gin.Context) {
	panic("not implemented")
}
