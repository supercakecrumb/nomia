# Contributing to Affirm Name

Thank you for your interest in contributing to Affirm Name! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Changelog Management](#changelog-management)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Code Style](#code-style)

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive environment for all contributors.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/affirm-name.git
   cd affirm-name
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/original-owner/affirm-name.git
   ```
4. **Install dependencies**:
   ```bash
   make install-deps
   ```

## Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style guidelines

3. **Add tests** for new functionality

4. **Run tests** to ensure everything works:
   ```bash
   make test
   ```

5. **Create a changelog entry** (see below)

6. **Commit your changes** with a descriptive message

7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Open a Pull Request** on GitHub

## Changelog Management

We use [Changie](https://changie.dev/) to manage our changelog. Every contribution that affects users should include a changelog entry.

### When to Create a Changelog Entry

Create a changelog entry for:
- ✅ New features
- ✅ Bug fixes
- ✅ Breaking changes
- ✅ Deprecations
- ✅ Security fixes
- ✅ Performance improvements
- ✅ API changes

Skip changelog entries for:
- ❌ Documentation updates (unless they document new features)
- ❌ Internal refactoring (unless it affects performance)
- ❌ Test additions
- ❌ CI/CD changes

### Creating a Changelog Entry

Use the Makefile target:

```bash
make changelog-new
```

Or use Changie directly:

```bash
changie new
```

You'll be prompted for:

1. **Kind** - Select the type of change:
   - `Added` - New features
   - `Changed` - Changes in existing functionality
   - `Deprecated` - Soon-to-be removed features
   - `Removed` - Removed features
   - `Fixed` - Bug fixes
   - `Security` - Security improvements

2. **Body** - Describe the change from a user's perspective:
   - ✅ Good: "Add support for UK baby names dataset"
   - ✅ Good: "Fix race condition in worker pool shutdown"
   - ❌ Bad: "Update parser.go"
   - ❌ Bad: "Fix bug"

3. **Author** - Your GitHub username

4. **Issue** - Related GitHub issue number (if applicable)

### Changelog Entry Examples

#### Adding a Feature
```yaml
kind: Added
body: Add support for filtering names by decade
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "42"
```

#### Fixing a Bug
```yaml
kind: Fixed
body: Fix pagination not working correctly with more than 1000 results
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "56"
```

#### Breaking Change
```yaml
kind: Changed
body: Change API response format for /api/v1/names endpoint to include metadata
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "78"
  Breaking: true
```

### Manual Changelog Entry

If you prefer, create a file manually in `.changes/unreleased/`:

```bash
# File: .changes/unreleased/Added-uk-dataset-support.yaml
kind: Added
body: Add support for UK baby names dataset
time: 2025-11-12T15:00:00Z
custom:
  Author: your-github-username
  Issue: "42"
```

## Commit Messages

Write clear, descriptive commit messages:

### Format
```
<type>: <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat: add UK baby names parser

Implement parser for UK Office for National Statistics baby names data.
Includes validation and normalization for UK-specific formats.

Closes #42
```

```
fix: resolve race condition in worker pool

Add proper synchronization when shutting down worker pool to prevent
panic when jobs are still being processed.

Fixes #56
```

## Pull Request Process

1. **Update documentation** if you've changed APIs or added features

2. **Ensure all tests pass**:
   ```bash
   make test
   ```

3. **Run linters**:
   ```bash
   make lint
   ```

4. **Create a changelog entry** (if applicable)

5. **Fill out the PR template** with:
   - Description of changes
   - Related issues
   - Testing performed
   - Screenshots (if UI changes)

6. **Request review** from maintainers

7. **Address feedback** and update your PR as needed

8. **Squash commits** if requested before merging

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test -v ./internal/parser/...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

- Write unit tests for all new functions
- Write integration tests for API endpoints
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Aim for >80% code coverage

### Test Example

```go
func TestNameNormalizer(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"lowercase", "john", "John"},
        {"uppercase", "JOHN", "John"},
        {"mixed", "jOhN", "John"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NormalizeName(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

## Code Style

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write clear, self-documenting code
- Add comments for exported functions
- Keep functions small and focused

### Formatting

```bash
# Format code
make fmt

# Run linters
make lint
```

### Documentation

- Document all exported functions, types, and constants
- Use complete sentences in comments
- Include examples for complex functionality

```go
// NormalizeName converts a name to title case and trims whitespace.
// It handles Unicode characters correctly and preserves hyphens.
//
// Example:
//   NormalizeName("mary-jane") // Returns "Mary-Jane"
func NormalizeName(name string) string {
    // implementation
}
```

## Questions?

If you have questions or need help:

1. Check existing [documentation](docs/)
2. Search [existing issues](https://github.com/owner/affirm-name/issues)
3. Open a new issue with the `question` label
4. Reach out to maintainers

Thank you for contributing! 🎉