package parsers

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUSSSAParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		year        int
		expectCount int
		expectError bool
	}{
		{
			name: "valid data",
			input: `Emma,F,20355
Liam,M,20456
Olivia,F,18256`,
			year:        2023,
			expectCount: 3,
			expectError: false,
		},
		{
			name:        "empty file",
			input:       "",
			year:        2023,
			expectCount: 0,
			expectError: false,
		},
		{
			name: "invalid count",
			input: `Emma,F,invalid
Liam,M,20456`,
			year:        2023,
			expectCount: 0,
			expectError: true,
		},
		{
			name: "wrong number of fields",
			input: `Emma,F
Liam,M,20456`,
			year:        2023,
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewUSSSAParser(tt.year)
			reader := strings.NewReader(tt.input)
			ctx := context.Background()

			records, errors := parser.Parse(ctx, reader)

			// Collect all records and errors
			var recordList []string
			var errorList []error

			done := false
			for !done {
				select {
				case record, ok := <-records:
					if !ok {
						records = nil
					} else {
						recordList = append(recordList, record.Name)
					}
				case err, ok := <-errors:
					if !ok {
						errors = nil
					} else {
						errorList = append(errorList, err)
					}
				}
				if records == nil && errors == nil {
					done = true
				}
			}

			if tt.expectError {
				assert.NotEmpty(t, errorList, "Expected errors but got none")
			} else {
				assert.Empty(t, errorList, "Expected no errors but got: %v", errorList)
				assert.Equal(t, tt.expectCount, len(recordList), "Record count mismatch")
			}
		})
	}
}

func TestUSSSAParser_Validate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name: "valid format",
			input: `Emma,F,20355
Liam,M,20456
Olivia,F,18256`,
			expectError: false,
		},
		{
			name: "invalid count",
			input: `Emma,F,invalid
Liam,M,20456`,
			expectError: true,
		},
		{
			name: "wrong number of fields",
			input: `Emma,F
Liam,M,20456`,
			expectError: true,
		},
		{
			name:        "empty file",
			input:       "",
			expectError: true,
		},
		{
			name: "invalid gender",
			input: `Emma,X,20355
Liam,M,20456`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewUSSSAParser(2023)
			reader := strings.NewReader(tt.input)

			err := parser.Validate(reader)

			if tt.expectError {
				assert.Error(t, err, "Expected validation error but got none")
			} else {
				assert.NoError(t, err, "Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestUSSSAParser_Metadata(t *testing.T) {
	parser := NewUSSSAParser(2023)
	metadata := parser.Metadata()

	assert.Equal(t, "US", metadata.CountryCode)
	assert.Equal(t, "US Social Security Administration", metadata.Name)
	assert.NotEmpty(t, metadata.Description)
	assert.NotEmpty(t, metadata.Version)
}

func TestUSSSAParser_ParseRow(t *testing.T) {
	parser := NewUSSSAParser(2023)

	tests := []struct {
		name         string
		row          []string
		expectName   string
		expectGender string
		expectCount  int
		expectError  bool
	}{
		{
			name:         "valid row",
			row:          []string{"Emma", "F", "20355"},
			expectName:   "Emma",
			expectGender: "F",
			expectCount:  20355,
			expectError:  false,
		},
		{
			name:         "male name",
			row:          []string{"Liam", "M", "20456"},
			expectName:   "Liam",
			expectGender: "M",
			expectCount:  20456,
			expectError:  false,
		},
		{
			name:        "invalid count",
			row:         []string{"Emma", "F", "invalid"},
			expectError: true,
		},
		{
			name:        "wrong field count",
			row:         []string{"Emma", "F"},
			expectError: true,
		},
		{
			name:        "invalid gender",
			row:         []string{"Emma", "X", "20355"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := parser.parseRow(tt.row, 1)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectName, record.Name)
				assert.Equal(t, tt.expectGender, record.Gender)
				assert.Equal(t, tt.expectCount, record.Count)
				assert.Equal(t, 2023, record.Year)
			}
		})
	}
}
