package pagination

// TODO: implement — see issue #18
//
// Parse reads ?page and ?limit query params and returns safe, clamped values.
//
// Rules:
//   - Default page: 1
//   - Default limit: 20
//   - Max limit: 100
//   - page < 1 → 1
//   - limit < 1 → 20
//   - limit > 100 → 100
//
// Usage:
//   p := pagination.Parse(c.Query("page"), c.Query("limit"))
//   // use p.Limit and p.Offset in SQL queries

// Params holds the parsed and clamped pagination values.
type Params struct {
	Page   int
	Limit  int
	Offset int
}

func Parse(pageStr, limitStr string) Params {
	panic("not implemented")
}
