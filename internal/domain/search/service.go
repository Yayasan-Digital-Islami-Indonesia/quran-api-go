package search

import "context"

// SearchService defines the business operations for full-text search.
// Implement this interface in internal/service/search_service.go.
type SearchService interface {
	Search(ctx context.Context, params Params) (results []Result, total int, err error)
}
