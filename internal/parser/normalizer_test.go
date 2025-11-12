package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizer_NormalizeGender(t *testing.T) {
	normalizer := NewNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"male uppercase", "M", "M", false},
		{"male lowercase", "m", "M", false},
		{"male full", "Male", "M", false},
		{"male full lowercase", "male", "M", false},
		{"male numeric", "1", "M", false},
		{"female uppercase", "F", "F", false},
		{"female lowercase", "f", "F", false},
		{"female full", "Female", "F", false},
		{"female full lowercase", "female", "F", false},
		{"female numeric", "2", "F", false},
		{"empty string", "", "", true},
		{"invalid", "X", "", true},
		{"invalid number", "3", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizer.NormalizeGender(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNormalizer_NormalizeName(t *testing.T) {
	normalizer := NewNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"simple name", "emma", "Emma", false},
		{"uppercase name", "EMMA", "Emma", false},
		{"mixed case", "EmMa", "Emma", false},
		{"hyphenated name", "mary-jane", "Mary-Jane", false},
		{"name with apostrophe", "o'brien", "O'Brien", false},
		{"two word name", "mary jane", "Mary Jane", false},
		{"name with period", "st. john", "St. John", false},
		{"empty string", "", "", true},
		{"whitespace only", "   ", "", true},
		{"too long", string(make([]byte, 101)), "", true},
		{"invalid character", "emma123", "", true},
		{"invalid character @", "emma@test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizer.NormalizeName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNormalizer_ValidateYear(t *testing.T) {
	normalizer := NewNormalizer()

	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{"valid year 2023", 2023, false},
		{"valid year 1880", 1880, false},
		{"valid year 2100", 2100, false},
		{"too early", 1879, true},
		{"too late", 2101, true},
		{"way too early", 1000, true},
		{"way too late", 3000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := normalizer.ValidateYear(tt.year)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizer_ValidateCount(t *testing.T) {
	normalizer := NewNormalizer()

	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{"valid count 1", 1, false},
		{"valid count 100", 100, false},
		{"valid count 20000", 20000, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"large negative", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := normalizer.ValidateCount(tt.count)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizer_NormalizeRecord(t *testing.T) {
	normalizer := NewNormalizer()

	tests := []struct {
		name    string
		record  Record
		wantErr bool
	}{
		{
			name: "valid record",
			record: Record{
				Year:   2023,
				Name:   "emma",
				Gender: "F",
				Count:  20000,
			},
			wantErr: false,
		},
		{
			name: "invalid year",
			record: Record{
				Year:   1800,
				Name:   "emma",
				Gender: "F",
				Count:  20000,
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			record: Record{
				Year:   2023,
				Name:   "",
				Gender: "F",
				Count:  20000,
			},
			wantErr: true,
		},
		{
			name: "invalid gender",
			record: Record{
				Year:   2023,
				Name:   "emma",
				Gender: "X",
				Count:  20000,
			},
			wantErr: true,
		},
		{
			name: "invalid count",
			record: Record{
				Year:   2023,
				Name:   "emma",
				Gender: "F",
				Count:  0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := tt.record
			err := normalizer.NormalizeRecord(&record)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// Check that name was normalized
				assert.Equal(t, "Emma", record.Name)
				assert.Equal(t, "F", record.Gender)
			}
		})
	}
}

func TestNormalizer_CustomSettings(t *testing.T) {
	normalizer := &Normalizer{
		MinYear:       2000,
		MaxYear:       2050,
		MinCount:      5,
		MaxNameLength: 50,
	}

	// Test custom year range
	assert.NoError(t, normalizer.ValidateYear(2000))
	assert.NoError(t, normalizer.ValidateYear(2050))
	assert.Error(t, normalizer.ValidateYear(1999))
	assert.Error(t, normalizer.ValidateYear(2051))

	// Test custom min count
	assert.NoError(t, normalizer.ValidateCount(5))
	assert.Error(t, normalizer.ValidateCount(4))

	// Test custom max name length
	longName := string(make([]byte, 51))
	for i := range longName {
		longName = string(append([]byte(longName[:i]), 'a'))
	}
	_, err := normalizer.NormalizeName(longName)
	assert.Error(t, err)
}
