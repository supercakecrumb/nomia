package parser

import (
	"fmt"
	"sync"
)

// Registry manages available parsers for different countries
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
}

// Global registry instance
var globalRegistry = NewRegistry()

// NewRegistry creates a new parser registry
func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

// Register adds a parser to the registry for a specific country code
// Country code should be ISO 3166-1 alpha-2 format (e.g., "US", "GB")
func (r *Registry) Register(countryCode string, parser Parser) error {
	if countryCode == "" {
		return fmt.Errorf("country code cannot be empty")
	}
	if parser == nil {
		return fmt.Errorf("parser cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if parser already registered
	if _, exists := r.parsers[countryCode]; exists {
		return fmt.Errorf("parser for country code %s already registered", countryCode)
	}

	r.parsers[countryCode] = parser
	return nil
}

// Get retrieves a parser for a specific country code
func (r *Registry) Get(countryCode string) (Parser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parser, exists := r.parsers[countryCode]
	if !exists {
		return nil, fmt.Errorf("no parser registered for country code: %s", countryCode)
	}

	return parser, nil
}

// List returns all registered country codes
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	codes := make([]string, 0, len(r.parsers))
	for code := range r.parsers {
		codes = append(codes, code)
	}
	return codes
}

// Unregister removes a parser from the registry
func (r *Registry) Unregister(countryCode string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parsers[countryCode]; !exists {
		return fmt.Errorf("no parser registered for country code: %s", countryCode)
	}

	delete(r.parsers, countryCode)
	return nil
}

// GetMetadata returns metadata for all registered parsers
func (r *Registry) GetMetadata() []ParserMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metadata := make([]ParserMetadata, 0, len(r.parsers))
	for _, parser := range r.parsers {
		metadata = append(metadata, parser.Metadata())
	}
	return metadata
}

// Global registry functions for convenience

// Register adds a parser to the global registry
func Register(countryCode string, parser Parser) error {
	return globalRegistry.Register(countryCode, parser)
}

// Get retrieves a parser from the global registry
func Get(countryCode string) (Parser, error) {
	return globalRegistry.Get(countryCode)
}

// List returns all registered country codes from the global registry
func List() []string {
	return globalRegistry.List()
}

// Unregister removes a parser from the global registry
func Unregister(countryCode string) error {
	return globalRegistry.Unregister(countryCode)
}

// GetMetadata returns metadata for all parsers in the global registry
func GetMetadata() []ParserMetadata {
	return globalRegistry.GetMetadata()
}
