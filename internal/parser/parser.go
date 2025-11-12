package parser

import (
	"context"
	"io"
)

// Record represents a single parsed name record from a CSV file
type Record struct {
	Year   int    `json:"year"`
	Name   string `json:"name"`
	Gender string `json:"gender"` // M or F
	Count  int    `json:"count"`
}

// ParserMetadata contains information about a parser
type ParserMetadata struct {
	CountryCode string `json:"country_code"` // ISO 3166-1 alpha-2 code (e.g., "US")
	Name        string `json:"name"`         // Human-readable name (e.g., "US Social Security Administration")
	Description string `json:"description"`  // Description of the data format
	Version     string `json:"version"`      // Parser version
}

// Parser defines the interface for parsing country-specific baby name CSV files
// Each country may have a different CSV format, so implementations handle
// the specific format for their country
type Parser interface {
	// Parse reads CSV data from the reader and returns channels for streaming records
	// The records channel emits parsed records, and the errors channel emits any parsing errors
	// Both channels are closed when parsing is complete or an error occurs
	// The context can be used to cancel the parsing operation
	Parse(ctx context.Context, reader io.Reader) (<-chan Record, <-chan error)

	// Validate checks if the reader contains valid data for this parser
	// This is a quick check before full parsing to detect format issues early
	// Returns an error if the data format is invalid
	Validate(reader io.Reader) error

	// Metadata returns information about this parser
	Metadata() ParserMetadata
}

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Line    int    // Line number where error occurred (0 if not applicable)
	Field   string // Field name where error occurred (empty if not applicable)
	Message string // Error message
	Err     error  // Underlying error (may be nil)
}

// Error implements the error interface
func (e *ParseError) Error() string {
	if e.Line > 0 {
		if e.Field != "" {
			return e.Message + " at line " + string(rune(e.Line)) + ", field " + e.Field
		}
		return e.Message + " at line " + string(rune(e.Line))
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *ParseError) Unwrap() error {
	return e.Err
}
