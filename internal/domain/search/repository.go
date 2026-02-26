package search

import "context"

// SearchRepository defines full-text search over ayah content via FTS5.
// Implement this interface in internal/repository/search_repository.go.
type SearchRepository interface {
	Search(ctx context.Context, params Params) (results []Result, total int, err error)
}
