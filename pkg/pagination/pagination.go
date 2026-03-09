// Package pagination provides shared pagination helpers for all HTTP handlers.
// Use Parse to convert raw query-string values into safe, clamped Params.
package pagination

import "strconv"

// Params holds the parsed and clamped pagination values.
type Params struct {
	Page   int
	Limit  int
	Offset int
}

// Parse reads raw ?page and ?limit query strings and returns safe, clamped Params.
//
// Defaults:  page=1, limit=20
// Clamping:  page < 1 → 1 | limit < 1 → 20 | limit > 100 → 100
// Offset:    (page - 1) * limit
//
// Usage (in a Gin handler):
//
//	p := pagination.Parse(c.Query("page"), c.Query("limit"))
//	// pass p.Limit and p.Offset to the repository
func Parse(pageStr, limitStr string) Params {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}