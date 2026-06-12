package response

// SuccessResponse is the envelope returned for all successful API responses.
type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// ErrorResponse is the envelope returned for all error API responses.
type ErrorResponse struct {
	Error     string `json:"error" example:"resource not found"`
	Code      string `json:"code" example:"not found"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}
