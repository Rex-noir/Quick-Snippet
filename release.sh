#!/usr/bin/env bash
set -e

VERSION=$1
APP_NAME="snip"
OUTPUT_DIR="build"

if [ -z "$VERSION" ]; then
  echo "Usage: ./release.sh v1.0.0"
  exit 1
fi

bash build.sh

git add .
git commit -m "Release $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"
git push origin main --tags

gh release create "$VERSION" "$OUTPUT_DIR"/* --title "$VERSION" --notes "Automated release for $APP_NAME $VERSION"