package parsers

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/supercakecrumb/affirm-name/internal/parser"
)

// USSSAParser implements the Parser interface for US Social Security Administration data
// Format: name,gender,count (no header row)
// Example: Emma,F,20355
type USSSAParser struct {
	normalizer *parser.Normalizer
	year       int // Year is extracted from filename or provided separately
}

// NewUSSSAParser creates a new US SSA parser
func NewUSSSAParser(year int) *USSSAParser {
	return &USSSAParser{
		normalizer: parser.NewNormalizer(),
		year:       year,
	}
}

// Parse reads CSV data and returns channels for streaming records
func (p *USSSAParser) Parse(ctx context.Context, reader io.Reader) (<-chan parser.Record, <-chan error) {
	records := make(chan parser.Record, 100)
	errors := make(chan error, 1)

	go func() {
		defer close(records)
		defer close(errors)

		csvReader := csv.NewReader(reader)
		csvReader.FieldsPerRecord = 3 // Expect exactly 3 fields
		csvReader.TrimLeadingSpace = true

		lineNum := 0
		for {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			default:
			}

			lineNum++
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errors <- &parser.ParseError{
					Line:    lineNum,
					Message: "failed to read CSV row",
					Err:     err,
				}
				return
			}

			// Parse the row
			record, err := p.parseRow(row, lineNum)
			if err != nil {
				errors <- err
				return
			}

			// Send record to channel
			select {
			case records <- *record:
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			}
		}
	}()

	return records, errors
}

// parseRow parses a single CSV row into a Record
func (p *USSSAParser) parseRow(row []string, lineNum int) (*parser.Record, error) {
	if len(row) != 3 {
		return nil, &parser.ParseError{
			Line:    lineNum,
			Message: fmt.Sprintf("expected 3 fields, got %d", len(row)),
		}
	}

	name := strings.TrimSpace(row[0])
	gender := strings.TrimSpace(row[1])
	countStr := strings.TrimSpace(row[2])

	// Parse count
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, &parser.ParseError{
			Line:    lineNum,
			Field:   "count",
			Message: fmt.Sprintf("invalid count value: %s", countStr),
			Err:     err,
		}
	}

	// Create record
	record := &parser.Record{
		Year:   p.year,
		Name:   name,
		Gender: gender,
		Count:  count,
	}

	// Normalize the record
	if err := p.normalizer.NormalizeRecord(record); err != nil {
		return nil, &parser.ParseError{
			Line:    lineNum,
			Message: "normalization failed",
			Err:     err,
		}
	}

	return record, nil
}

// Validate checks if the reader contains valid US SSA format data
func (p *USSSAParser) Validate(reader io.Reader) error {
	// Read first few lines to validate format
	scanner := bufio.NewScanner(reader)
	lineNum := 0
	const maxLinesToCheck = 10

	for scanner.Scan() && lineNum < maxLinesToCheck {
		lineNum++
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Check if line has 3 comma-separated fields
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			return fmt.Errorf("line %d: expected 3 fields, got %d", lineNum, len(fields))
		}

		// Validate that third field is a number
		countStr := strings.TrimSpace(fields[2])
		if _, err := strconv.Atoi(countStr); err != nil {
			return fmt.Errorf("line %d: invalid count value: %s", lineNum, countStr)
		}

		// Validate gender field
		gender := strings.TrimSpace(fields[1])
		if _, err := p.normalizer.NormalizeGender(gender); err != nil {
			return fmt.Errorf("line %d: invalid gender: %s", lineNum, gender)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	if lineNum == 0 {
		return fmt.Errorf("file is empty")
	}

	return nil
}

// Metadata returns information about this parser
func (p *USSSAParser) Metadata() parser.ParserMetadata {
	return parser.ParserMetadata{
		CountryCode: "US",
		Name:        "US Social Security Administration",
		Description: "Parses US SSA baby names data in CSV format (name,gender,count)",
		Version:     "1.0.0",
	}
}

// init registers the US SSA parser in the global registry
func init() {
	// Register a factory function that creates parsers with the appropriate year
	// Note: The actual year will be set when the parser is created for a specific file
	// For now, we register a default parser
	defaultParser := NewUSSSAParser(2023)
	if err := parser.Register("US", defaultParser); err != nil {
		// This should not happen in normal operation
		panic(fmt.Sprintf("failed to register US SSA parser: %v", err))
	}
}
