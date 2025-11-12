# Changelog Management with Changie

This directory contains unreleased changelog entries managed by [Changie](https://changie.dev/).

## Directory Structure

```
.changes/
├── unreleased/          # Unreleased change files
│   ├── Added-*.yaml     # New features
│   ├── Changed-*.yaml   # Changes in existing functionality
│   ├── Fixed-*.yaml     # Bug fixes
│   └── ...
└── v*.md                # Released version files
```

## Creating a New Change Entry

Use the Makefile target to create a new changelog entry:

```bash
make changelog-new
```

Or use Changie directly:

```bash
changie new
```

This will prompt you for:
- **Kind**: Type of change (Added, Changed, Fixed, etc.)
- **Body**: Description of the change
- **Author**: Your GitHub username
- **Issue**: Related GitHub issue number (optional)

## Change Types

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security improvements

## Manual Change File Creation

You can also manually create change files in `.changes/unreleased/`:

```yaml
kind: Added
body: Implement new feature X
time: 2025-11-12T15:00:00Z
custom:
  Author: your-github-username
  Issue: "123"
```

File naming convention: `{Kind}-{description}.yaml`

Example: `Added-new-api-endpoint.yaml`

## Batching Changes

When ready to release, batch all unreleased changes:

```bash
make changelog-batch
```

This creates a version file in `.changes/` with all unreleased changes.

## Merging to CHANGELOG

Merge the batched changes into CHANGELOG.md:

```bash
make changelog-merge
```

## Full Release Process

Use the automated release script:

```bash
# Patch release (0.1.0 -> 0.1.1)
make release VERSION=patch

# Minor release (0.1.0 -> 0.2.0)
make release VERSION=minor

# Major release (0.1.0 -> 1.0.0)
make release VERSION=major
```

## Best Practices

1. **Create entries as you work**: Don't wait until release time
2. **Be descriptive**: Write clear, user-focused descriptions
3. **Reference issues**: Link to GitHub issues when applicable
4. **One change per file**: Keep changes atomic and focused
5. **Use proper categories**: Choose the most appropriate change type

## Examples

### Adding a Feature
```yaml
kind: Added
body: Add support for UK baby names dataset
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "42"
```

### Fixing a Bug
```yaml
kind: Fixed
body: Fix race condition in worker pool shutdown
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "56"
```

### Breaking Change
```yaml
kind: Changed
body: Change API response format for /api/v1/names endpoint
time: 2025-11-12T15:00:00Z
custom:
  Author: username
  Issue: "78"
  Breaking: true
```

## Resources

- [Changie Documentation](https://changie.dev/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)