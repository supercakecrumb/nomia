package handlers

import (
	"net/url"
	"testing"
)

func TestParseNameTrendParams(t *testing.T) {
	tests := []struct {
		name      string
		query     url.Values
		dbStart   int
		dbEnd     int
		wantErr   bool
		errMsg    string
		checkFunc func(*testing.T, *NameTrendParams)
	}{
		{
			name:    "missing name parameter",
			query:   url.Values{},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "name parameter is required",
		},
		{
			name: "valid name parameter",
			query: url.Values{
				"name": []string{"Oliver"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NameTrendParams) {
				if p.Name != "Oliver" {
					t.Errorf("Name = %s, want Oliver", p.Name)
				}
				if p.YearFrom != 2020 {
					t.Errorf("YearFrom = %d, want 2020", p.YearFrom)
				}
				if p.YearTo != 2024 {
					t.Errorf("YearTo = %d, want 2024", p.YearTo)
				}
			},
		},
		{
			name: "custom year range",
			query: url.Values{
				"name":      []string{"Emma"},
				"year_from": []string{"2022"},
				"year_to":   []string{"2023"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NameTrendParams) {
				if p.Name != "Emma" {
					t.Errorf("Name = %s, want Emma", p.Name)
				}
				if p.YearFrom != 2022 {
					t.Errorf("YearFrom = %d, want 2022", p.YearFrom)
				}
				if p.YearTo != 2023 {
					t.Errorf("YearTo = %d, want 2023", p.YearTo)
				}
			},
		},
		{
			name: "invalid year range",
			query: url.Values{
				"name":      []string{"Test"},
				"year_from": []string{"2024"},
				"year_to":   []string{"2022"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "year_from must be <= year_to",
		},
		{
			name: "countries filter",
			query: url.Values{
				"name":      []string{"John"},
				"countries": []string{"US,UK"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: false,
			checkFunc: func(t *testing.T, p *NameTrendParams) {
				if len(p.Countries) != 2 {
					t.Errorf("Countries length = %d, want 2", len(p.Countries))
				}
			},
		},
		{
			name: "invalid year_from",
			query: url.Values{
				"name":      []string{"Test"},
				"year_from": []string{"invalid"},
			},
			dbStart: 2020,
			dbEnd:   2024,
			wantErr: true,
			errMsg:  "year_from must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := ParseNameTrendParams(tt.query, tt.dbStart, tt.dbEnd)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseNameTrendParams() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseNameTrendParams() error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseNameTrendParams() unexpected error = %v", err)
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, params)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
