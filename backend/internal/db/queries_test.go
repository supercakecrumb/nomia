package db

import (
	"net/url"
	"testing"
)

func TestParseNamesListParams(t *testing.T) {
	tests := []struct {
		name      string
		query     url.Values
		dbStart   int
		dbEnd     int
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, *NamesListParams)
	}{
		{
			name:    "default parameters",
			query:   url.Values{},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.YearFrom != 2020 {
					t.Errorf("YearFrom = %d, want 2020", p.YearFrom)
				}
				if p.YearTo != 2024 {
					t.Errorf("YearTo = %d, want 2024", p.YearTo)
				}
				if p.Page != 1 {
					t.Errorf("Page = %d, want 1", p.Page)
				}
				if p.PageSize != 50 {
					t.Errorf("PageSize = %d, want 50", p.PageSize)
				}
			},
		},
		{
			name: "custom year range",
			query: url.Values{
				"year_from": []string{"2022"},
				"year_to":   []string{"2023"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.YearFrom != 2022 {
					t.Errorf("YearFrom = %d, want 2022", p.YearFrom)
				}
				if p.YearTo != 2023 {
					t.Errorf("YearTo = %d, want 2023", p.YearTo)
				}
			},
		},
		{
			name: "invalid year_from not integer",
			query: url.Values{
				"year_from": []string{"invalid"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "year_from must be an integer",
		},
		{
			name: "year_from > year_to",
			query: url.Values{
				"year_from": []string{"2024"},
				"year_to":   []string{"2022"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "year_from must be <= year_to",
		},
		{
			name: "gender balance filters",
			query: url.Values{
				"gender_balance_min": []string{"40"},
				"gender_balance_max": []string{"60"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.GenderBalanceMin != 40 {
					t.Errorf("GenderBalanceMin = %d, want 40", p.GenderBalanceMin)
				}
				if p.GenderBalanceMax != 60 {
					t.Errorf("GenderBalanceMax = %d, want 60", p.GenderBalanceMax)
				}
			},
		},
		{
			name: "invalid gender balance range",
			query: url.Values{
				"gender_balance_min": []string{"150"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "gender_balance_min must be between 0 and 100",
		},
		{
			name: "page out of range",
			query: url.Values{
				"page": []string{"101"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "page must be between 1 and 100",
		},
		{
			name: "page_size out of range",
			query: url.Values{
				"page_size": []string{"5"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "page_size must be between 10 and 100",
		},
		{
			name: "invalid sort_key",
			query: url.Values{
				"sort_key": []string{"invalid"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "sort_key must be one of",
		},
		{
			name: "invalid sort_order",
			query: url.Values{
				"sort_order": []string{"invalid"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "sort_order must be either 'asc' or 'desc'",
		},
		{
			name: "top_n filter",
			query: url.Values{
				"top_n": []string{"10"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.TopN != 10 {
					t.Errorf("TopN = %d, want 10", p.TopN)
				}
				if p.GetActivePopularityFilter() != "top_n" {
					t.Errorf("ActivePopularityFilter = %s, want top_n", p.GetActivePopularityFilter())
				}
			},
		},
		{
			name: "coverage_percent filter",
			query: url.Values{
				"coverage_percent": []string{"95.5"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.CoveragePercent != 95.5 {
					t.Errorf("CoveragePercent = %f, want 95.5", p.CoveragePercent)
				}
				if p.GetActivePopularityFilter() != "coverage_percent" {
					t.Errorf("ActivePopularityFilter = %s, want coverage_percent", p.GetActivePopularityFilter())
				}
			},
		},
		{
			name: "name_glob pattern",
			query: url.Values{
				"name_glob": []string{"Alex*"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if p.NameGlob != "Alex*" {
					t.Errorf("NameGlob = %s, want Alex*", p.NameGlob)
				}
			},
		},
		{
			name: "countries filter",
			query: url.Values{
				"countries": []string{"US,UK,CA"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NamesListParams) {
				if len(p.Countries) != 3 {
					t.Errorf("Countries length = %d, want 3", len(p.Countries))
				}
				if p.Countries[0] != "US" || p.Countries[1] != "UK" || p.Countries[2] != "CA" {
					t.Errorf("Countries = %v, want [US UK CA]", p.Countries)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := ParseNamesListParams(tt.query, tt.dbStart, tt.dbEnd)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseNamesListParams() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseNamesListParams() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseNamesListParams() unexpected error = %v", err)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, params)
			}
		})
	}
}

func TestGetActivePopularityFilter(t *testing.T) {
	tests := []struct {
		name   string
		params *NamesListParams
		want   string
	}{
		{
			name: "coverage_percent has priority",
			params: &NamesListParams{
				CoveragePercent: 90,
				TopN:            100,
				MinCount:        50,
			},
			want: "coverage_percent",
		},
		{
			name: "top_n when no coverage_percent",
			params: &NamesListParams{
				TopN:     100,
				MinCount: 50,
			},
			want: "top_n",
		},
		{
			name: "min_count when no other filters",
			params: &NamesListParams{
				MinCount: 50,
			},
			want: "min_count",
		},
		{
			name:   "none when no filters",
			params: &NamesListParams{},
			want:   "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.params.GetActivePopularityFilter()
			if got != tt.want {
				t.Errorf("GetActivePopularityFilter() = %s, want %s", got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
