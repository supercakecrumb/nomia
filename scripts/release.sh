#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if changie is installed
if ! command -v changie &> /dev/null; then
    print_error "changie is not installed. Install it with:"
    echo "  go install github.com/miniscruff/changie@latest"
    exit 1
fi

# Check if git is installed
if ! command -v git &> /dev/null; then
    print_error "git is not installed"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
    print_error "You have uncommitted changes. Please commit or stash them first."
    git status -s
    exit 1
fi

# Get version bump type
VERSION_TYPE=${1:-patch}

if [[ ! "$VERSION_TYPE" =~ ^(major|minor|patch)$ ]]; then
    print_error "Invalid version type: $VERSION_TYPE"
    echo "Usage: $0 [major|minor|patch]"
    exit 1
fi

print_info "Starting release process with version bump: $VERSION_TYPE"

# Get current version from CHANGELOG.md
CURRENT_VERSION=$(grep -m 1 "^## " CHANGELOG.md | sed 's/## \([0-9.]*\).*/\1/')

if [[ -z "$CURRENT_VERSION" ]]; then
    print_error "Could not determine current version from CHANGELOG.md"
    exit 1
fi

print_info "Current version: $CURRENT_VERSION"

# Calculate next version
IFS='.' read -r -a version_parts <<< "$CURRENT_VERSION"
major="${version_parts[0]}"
minor="${version_parts[1]}"
patch="${version_parts[2]}"

case $VERSION_TYPE in
    major)
        major=$((major + 1))
        minor=0
        patch=0
        ;;
    minor)
        minor=$((minor + 1))
        patch=0
        ;;
    patch)
        patch=$((patch + 1))
        ;;
esac

NEW_VERSION="$major.$minor.$patch"
print_info "New version: $NEW_VERSION"

# Check if there are unreleased changes
UNRELEASED_COUNT=$(find .changes/unreleased -name "*.yaml" 2>/dev/null | wc -l)

if [[ $UNRELEASED_COUNT -eq 0 ]]; then
    print_warn "No unreleased changes found in .changes/unreleased/"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Release cancelled"
        exit 0
    fi
fi

# Batch changes
print_info "Batching unreleased changes..."
changie batch "$NEW_VERSION"

# Merge changes into CHANGELOG
print_info "Merging changes into CHANGELOG.md..."
changie merge

# Update version in go.mod if needed (optional)
# This is a placeholder - adjust based on your versioning needs
# sed -i.bak "s/version v.*/version v$NEW_VERSION/" go.mod && rm go.mod.bak

# Commit changes
print_info "Committing changes..."
git add CHANGELOG.md .changes/
git commit -m "Release v$NEW_VERSION"

# Create git tag
print_info "Creating git tag v$NEW_VERSION..."
git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"

print_info "Release v$NEW_VERSION created successfully!"
echo ""
print_info "Next steps:"
echo "  1. Review the changes: git show v$NEW_VERSION"
echo "  2. Push the changes: git push origin main"
echo "  3. Push the tag: git push origin v$NEW_VERSION"
echo ""
print_warn "To undo this release (before pushing):"
echo "  git tag -d v$NEW_VERSION"
echo "  git reset --hard HEAD~1"
echo ""

# Optional: Create GitHub release
read -p "Create GitHub release? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if command -v gh &> /dev/null; then
        print_info "Creating GitHub release..."
        
        # Extract release notes from CHANGELOG
        RELEASE_NOTES=$(awk "/^## $NEW_VERSION/,/^## [0-9]/" CHANGELOG.md | sed '1d;$d')
        
        gh release create "v$NEW_VERSION" \
            --title "v$NEW_VERSION" \
            --notes "$RELEASE_NOTES"
        
        print_info "GitHub release created!"
    else
        print_warn "GitHub CLI (gh) not installed. Skipping GitHub release."
        echo "Install it from: https://cli.github.com/"
    fi
fi

print_info "Release process complete!"