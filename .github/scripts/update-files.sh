#!/bin/bash
set -e

VERSION=$1
echo "Updating files to version ${VERSION}"

# Update Homebrew formula
if [ -f "Formula/tada.rb" ]; then
  echo "Current formula content:"
  grep -n "tag:" Formula/tada.rb || true

  # Create a temporary file for the sed operation
  cp Formula/tada.rb Formula/tada.rb.tmp
  sed -E "s/tag:[[:space:]]+\"[^\"]+\"/tag:      \"${VERSION}\"/" Formula/tada.rb.tmp > Formula/tada.rb
  rm Formula/tada.rb.tmp

  echo "Updated formula content:"
  grep -n "tag:" Formula/tada.rb || true
  echo "✅ Updated Formula/tada.rb"
fi

# Update VERSION file
echo "${VERSION}" > VERSION
echo "✅ Updated VERSION file"

# Update version in Go code
if [ -f "cmd/version.go" ]; then
  cp cmd/version.go cmd/version.go.tmp
  sed -E "s/Version = \"[^\"]+\"/Version = \"${VERSION}\"/" cmd/version.go.tmp > cmd/version.go
  rm cmd/version.go.tmp
  echo "✅ Updated cmd/version.go"
fi

echo "📋 Files updated successfully for version ${VERSION}"
