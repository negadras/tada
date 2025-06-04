# .github/workflows/auto-release.yml
name: Auto Release

on:
  push:
    branches:
      - main

permissions:
  contents: write
  pull-requests: write
  issues: write

jobs:
  release:
    name: Create Release
    runs-on: ubuntu-latest
    # Only run on main branch push (not on release commits)
    if: "!contains(github.event.head_commit.message, 'chore(release)')"

    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install semantic-release plugins
        run: |
          npm install -g \
            semantic-release \
            @semantic-release/git \
            @semantic-release/github \
            @semantic-release/changelog \
            @semantic-release/exec \
            conventional-changelog-conventionalcommits

      # .releaserc.json is now in the repository root

      - name: Create update script
        run: |
          mkdir -p .github/scripts
          cat > .github/scripts/update-files.sh << 'EOF'
          #!/bin/bash
          set -e

          VERSION=$1
          echo "Updating files to version ${VERSION}"

          # Update Homebrew formula
          if [ -f "Formula/tada.rb" ]; then
            sed -i -E "s/tag:[[:space:]]+\"[^\"]+\"/tag:      \"${VERSION}\"/" Formula/tada.rb
            echo "✅ Updated Formula/tada.rb"
          fi

          # Update VERSION file
          echo "${VERSION}" > VERSION
          echo "✅ Updated VERSION file"

          # Update version in Go code
          if [ -f "cmd/version.go" ]; then
            sed -i -E "s/Version = \"[^\"]+\"/Version = \"${VERSION}\"/" cmd/version.go
            echo "✅ Updated cmd/version.go"
          fi

          # Update version in main.go if exists
          if [ -f "main.go" ] && grep -q 'var version = ' main.go; then
            sed -i -E "s/var version = \"[^\"]*\"/var version = \"${VERSION}\"/" main.go
            echo "✅ Updated main.go"
          fi
          EOF
          chmod +x .github/scripts/update-files.sh

      - name: Run semantic-release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          npx semantic-release
