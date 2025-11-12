#!/usr/bin/env bash
set -e

VERSION=$1
APP_NAME="snip"
OUTPUT_DIR="build"

if [ -z "$VERSION" ]; then
  echo "Usage: ./release.sh v1.0.0"
  exit 1
fi

# Run your cross-platform build
#bash build.sh

# Generate changelog from last tag
LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [ -n "$LAST_TAG" ]; then
  echo "Generating changelog since $LAST_TAG..."
  CHANGELOG=$(git log "$LAST_TAG"..HEAD --pretty=format:"- %s (%h)" --no-merges)
else
  echo "No previous tag found. Generating changelog from start..."
  CHANGELOG=$(git log --pretty=format:"- %s (%h)" --no-merges)
fi

# Save changelog to file
CHANGELOG_FILE="CHANGELOG.md"
{
  echo -e "## $VERSION - $(date +'%Y-%m-%d')\n"
  echo -e "$CHANGELOG\n"
  echo
} >> "$CHANGELOG_FILE"

# Commit, tag, and push
git add .
git commit -m "Release $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"
git push origin main --tags

# Create GitHub release with changelog content
gh release create "$VERSION" "$OUTPUT_DIR"/* \
  --title "$VERSION" \
  --notes-file "$CHANGELOG_FILE"

echo "✅ Release $VERSION created with changelog."
