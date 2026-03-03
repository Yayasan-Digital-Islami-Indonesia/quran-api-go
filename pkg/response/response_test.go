package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyTimestamp struct {
	Timestamp string `json:"timestamp"`
}

func testTimestamp(t *testing.T, timestamp string) {
	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Logf("invalid timestamp format: %v", err)
		t.Fail()
	}
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	msg := "success fetch"

	Success(ctx, msg)

	if w.Code != http.StatusOK {
		t.Logf("expected 200, got %d", w.Code)
		t.Fail()
	}

	var body struct {
		bodyTimestamp
	}

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Logf("invalid json %v", err)
	}

	testTimestamp(t, body.Timestamp)
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

	var body struct {
		Error string `json:"error"`
		bodyTimestamp
	}

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Logf("invalid json %v", err)
	}

	if body.Error != msg {
		t.Logf("expected %s,surah got %s", msg, body.Error)
		t.Fail()
	}

	testTimestamp(t, body.Timestamp)
}

func TestBadRequestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	msg := "invalid lang"

	BadRequest(ctx, msg)

	if w.Code != http.StatusBadRequest {
		t.Logf("expected 400, got %d", w.Code)
		t.Fail()
	}

	var body struct {
		Error string `json:"error"`
		bodyTimestamp
	}

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Logf("invalid json %v", err)
	}

	if body.Error != msg {
		t.Logf("expected %s,surah got %s", msg, body.Error)
		t.Fail()
	}

	testTimestamp(t, body.Timestamp)
}

func TestInternalErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	InternalError(ctx)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var body struct {
		bodyTimestamp
	}

	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Logf("invalid json %v", err)
	}

	testTimestamp(t, body.Timestamp)
}
