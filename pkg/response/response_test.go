package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	msg := "success fetch"

	Success(ctx, msg)

	if w.Code != http.StatusOK {
		t.Logf("expected 200, got %d", w.Code)
		t.Fail()
	}
}

func TestNotFoundResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	msg := "surah not found"

	NotFound(ctx, msg)

	if w.Code != http.StatusNotFound {
		t.Logf("expected 404, got %d", w.Code)
		t.Fail()
	}
}

func TestBadRequestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	BadRequest(ctx, "invalid lang")

	if w.Code != http.StatusBadRequest {
		t.Logf("expected 400, got %d", w.Code)
		t.Fail()
	}
}

func TestInternalErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	InternalError(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
