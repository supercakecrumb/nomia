package parser

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockParser is a mock implementation of the Parser interface for testing
type mockParser struct {
	metadata ParserMetadata
}

func (m *mockParser) Parse(ctx context.Context, reader io.Reader) (<-chan Record, <-chan error) {
	records := make(chan Record)
	errors := make(chan error)
	close(records)
	close(errors)
	return records, errors
}

func (m *mockParser) Validate(reader io.Reader) error {
	return nil
}

func (m *mockParser) Metadata() ParserMetadata {
	return m.metadata
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	parser := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "TEST",
			Name:        "Test Parser",
			Description: "A test parser",
			Version:     "1.0.0",
		},
	}

	// Test successful registration
	err := registry.Register("TEST", parser)
	assert.NoError(t, err)

	// Test duplicate registration
	err = registry.Register("TEST", parser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")

	// Test empty country code
	err = registry.Register("", parser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Test nil parser
	err = registry.Register("TEST2", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()

	parser := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "TEST",
			Name:        "Test Parser",
			Description: "A test parser",
			Version:     "1.0.0",
		},
	}

	// Register parser
	err := registry.Register("TEST", parser)
	require.NoError(t, err)

	// Test successful retrieval
	retrieved, err := registry.Get("TEST")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "TEST", retrieved.Metadata().CountryCode)

	// Test non-existent parser
	_, err = registry.Get("NONEXISTENT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no parser registered")
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	// Empty registry
	codes := registry.List()
	assert.Empty(t, codes)

	// Add parsers
	parser1 := &mockParser{metadata: ParserMetadata{CountryCode: "US"}}
	parser2 := &mockParser{metadata: ParserMetadata{CountryCode: "GB"}}
	parser3 := &mockParser{metadata: ParserMetadata{CountryCode: "FR"}}

	registry.Register("US", parser1)
	registry.Register("GB", parser2)
	registry.Register("FR", parser3)

	// List all parsers
	codes = registry.List()
	assert.Len(t, codes, 3)
	assert.Contains(t, codes, "US")
	assert.Contains(t, codes, "GB")
	assert.Contains(t, codes, "FR")
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()

	parser := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "TEST",
			Name:        "Test Parser",
		},
	}

	// Register parser
	err := registry.Register("TEST", parser)
	require.NoError(t, err)

	// Verify it exists
	_, err = registry.Get("TEST")
	assert.NoError(t, err)

	// Unregister
	err = registry.Unregister("TEST")
	assert.NoError(t, err)

	// Verify it's gone
	_, err = registry.Get("TEST")
	assert.Error(t, err)

	// Try to unregister again
	err = registry.Unregister("TEST")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no parser registered")
}

func TestRegistry_GetMetadata(t *testing.T) {
	registry := NewRegistry()

	// Empty registry
	metadata := registry.GetMetadata()
	assert.Empty(t, metadata)

	// Add parsers
	parser1 := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "US",
			Name:        "US Parser",
			Description: "US SSA Parser",
			Version:     "1.0.0",
		},
	}
	parser2 := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "GB",
			Name:        "GB Parser",
			Description: "UK ONS Parser",
			Version:     "1.0.0",
		},
	}

	registry.Register("US", parser1)
	registry.Register("GB", parser2)

	// Get all metadata
	metadata = registry.GetMetadata()
	assert.Len(t, metadata, 2)

	// Check that both parsers' metadata is present
	countryCodes := make([]string, len(metadata))
	for i, m := range metadata {
		countryCodes[i] = m.CountryCode
	}
	assert.Contains(t, countryCodes, "US")
	assert.Contains(t, countryCodes, "GB")
}

func TestGlobalRegistry(t *testing.T) {
	// Note: This test uses the global registry, so it may interfere with other tests
	// In a real scenario, you might want to reset the global registry between tests

	parser := &mockParser{
		metadata: ParserMetadata{
			CountryCode: "GLOBAL_TEST",
			Name:        "Global Test Parser",
		},
	}

	// Test global Register
	err := Register("GLOBAL_TEST", parser)
	assert.NoError(t, err)

	// Test global Get
	retrieved, err := Get("GLOBAL_TEST")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Test global List
	codes := List()
	assert.Contains(t, codes, "GLOBAL_TEST")

	// Test global GetMetadata
	metadata := GetMetadata()
	found := false
	for _, m := range metadata {
		if m.CountryCode == "GLOBAL_TEST" {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test global Unregister
	err = Unregister("GLOBAL_TEST")
	assert.NoError(t, err)
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewRegistry()

	// Test concurrent registration
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			parser := &mockParser{
				metadata: ParserMetadata{
					CountryCode: string(rune('A' + id)),
				},
			}
			registry.Register(string(rune('A'+id)), parser)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all parsers were registered
	codes := registry.List()
	assert.Len(t, codes, 10)
}
