#!/usr/bin/env bash
set -e

APP_NAME="snip"
OUTPUT_DIR="build"

mkdir -p "$OUTPUT_DIR"

platforms=(
  "windows/amd64"
  "windows/arm64"
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "android/arm64"
)


for platform in "${platforms[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$platform"
  output_name="${APP_NAME}-${GOOS}-${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    output_name+=".exe"
  fi

  echo "Building for $GOOS/$GOARCH ..."
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$output_name" ./main.go
done

echo "All builds complete. Binaries are in $OUTPUT_DIR/"