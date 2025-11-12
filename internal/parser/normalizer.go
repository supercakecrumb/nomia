package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// Normalizer provides methods for normalizing and validating parsed data
type Normalizer struct {
	// MinYear is the minimum valid year (default: 1880)
	MinYear int
	// MaxYear is the maximum valid year (default: current year + 1)
	MaxYear int
	// MinCount is the minimum valid count (default: 1)
	MinCount int
	// MaxNameLength is the maximum length for a name (default: 100)
	MaxNameLength int
}

// NewNormalizer creates a new normalizer with default settings
func NewNormalizer() *Normalizer {
	return &Normalizer{
		MinYear:       1880,
		MaxYear:       2100,
		MinCount:      1,
		MaxNameLength: 100,
	}
}

// NormalizeGender converts various gender representations to standard M/F format
// Supported inputs:
// - "Male", "male", "M", "m", "1" -> "M"
// - "Female", "female", "F", "f", "2" -> "F"
func (n *Normalizer) NormalizeGender(gender string) (string, error) {
	gender = strings.TrimSpace(gender)
	if gender == "" {
		return "", fmt.Errorf("gender cannot be empty")
	}

	// Convert to uppercase for comparison
	upper := strings.ToUpper(gender)

	// Map various formats to M/F
	switch upper {
	case "M", "MALE", "1":
		return "M", nil
	case "F", "FEMALE", "2":
		return "F", nil
	default:
		return "", fmt.Errorf("invalid gender value: %s (expected M, F, Male, Female, 1, or 2)", gender)
	}
}

// NormalizeName cleans and validates a name
// - Trims whitespace
// - Validates length
// - Ensures it contains only valid characters (letters, hyphens, apostrophes, spaces)
func (n *Normalizer) NormalizeName(name string) (string, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Check if empty
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	// Check length
	if len(name) > n.MaxNameLength {
		return "", fmt.Errorf("name too long: %d characters (max: %d)", len(name), n.MaxNameLength)
	}

	// Validate characters - allow letters, hyphens, apostrophes, and spaces
	for _, r := range name {
		if !unicode.IsLetter(r) && r != '-' && r != '\'' && r != ' ' && r != '.' {
			return "", fmt.Errorf("name contains invalid character: %c", r)
		}
	}

	// Capitalize first letter of each word for consistency
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			// Handle hyphenated names and apostrophes
			parts := strings.Split(word, "-")
			for j, part := range parts {
				if len(part) > 0 {
					// Handle apostrophes within parts
					apostropheParts := strings.Split(part, "'")
					for k, apPart := range apostropheParts {
						if len(apPart) > 0 {
							apostropheParts[k] = strings.ToUpper(string(apPart[0])) + strings.ToLower(apPart[1:])
						}
					}
					parts[j] = strings.Join(apostropheParts, "'")
				}
			}
			words[i] = strings.Join(parts, "-")
		}
	}
	name = strings.Join(words, " ")

	return name, nil
}

// ValidateYear checks if a year is within valid range
func (n *Normalizer) ValidateYear(year int) error {
	if year < n.MinYear {
		return fmt.Errorf("year %d is before minimum year %d", year, n.MinYear)
	}
	if year > n.MaxYear {
		return fmt.Errorf("year %d is after maximum year %d", year, n.MaxYear)
	}
	return nil
}

// ValidateCount checks if a count is valid (positive integer)
func (n *Normalizer) ValidateCount(count int) error {
	if count < n.MinCount {
		return fmt.Errorf("count %d is less than minimum count %d", count, n.MinCount)
	}
	return nil
}

// NormalizeRecord normalizes all fields in a record
func (n *Normalizer) NormalizeRecord(record *Record) error {
	// Validate year
	if err := n.ValidateYear(record.Year); err != nil {
		return fmt.Errorf("invalid year: %w", err)
	}

	// Normalize and validate name
	normalizedName, err := n.NormalizeName(record.Name)
	if err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	record.Name = normalizedName

	// Normalize and validate gender
	normalizedGender, err := n.NormalizeGender(record.Gender)
	if err != nil {
		return fmt.Errorf("invalid gender: %w", err)
	}
	record.Gender = normalizedGender

	// Validate count
	if err := n.ValidateCount(record.Count); err != nil {
		return fmt.Errorf("invalid count: %w", err)
	}

	return nil
}
