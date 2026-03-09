package pagination

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		pageStr     string
		limitStr    string
		wantPage    int
		wantLimit   int
		wantOffset  int
	}{
		{
			name:       "empty strings use defaults",
			pageStr:    "",
			limitStr:   "",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "non-numeric strings use defaults",
			pageStr:    "abc",
			limitStr:   "xyz",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "page=0 clamped to 1",
			pageStr:    "0",
			limitStr:   "20",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "negative page clamped to 1",
			pageStr:    "-1",
			limitStr:   "20",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "limit=0 clamped to 20",
			pageStr:    "1",
			limitStr:   "0",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "negative limit clamped to 20",
			pageStr:    "1",
			limitStr:   "-5",
			wantPage:   1,
			wantLimit:  20,
			wantOffset: 0,
		},
		{
			name:       "limit > 100 clamped to 100",
			pageStr:    "1",
			limitStr:   "999",
			wantPage:   1,
			wantLimit:  100,
			wantOffset: 0,
		},
		{
			name:       "limit exactly 100 accepted",
			pageStr:    "1",
			limitStr:   "100",
			wantPage:   1,
			wantLimit:  100,
			wantOffset: 0,
		},
		{
			name:       "happy path: page=2, limit=10",
			pageStr:    "2",
			limitStr:   "10",
			wantPage:   2,
			wantLimit:  10,
			wantOffset: 10,
		},
		{
			name:       "page=5, limit=50 offset=200",
			pageStr:    "5",
			limitStr:   "50",
			wantPage:   5,
			wantLimit:  50,
			wantOffset: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Parse(tc.pageStr, tc.limitStr)

			if got.Page != tc.wantPage {
				t.Errorf("Page: got %d, want %d", got.Page, tc.wantPage)
			}
			if got.Limit != tc.wantLimit {
				t.Errorf("Limit: got %d, want %d", got.Limit, tc.wantLimit)
			}
			if got.Offset != tc.wantOffset {
				t.Errorf("Offset: got %d, want %d", got.Offset, tc.wantOffset)
			}
		})
	}
}
