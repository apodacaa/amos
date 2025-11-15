#!/bin/bash
set -e

# Release automation script for amos
# Usage: ./scripts/release.sh VERSION [RELEASE_NOTES_FILE]
# Example: ./scripts/release.sh 1.2.1
# Example: ./scripts/release.sh 1.2.1 release-notes.md

VERSION=$1
RELEASE_NOTES_FILE=$2
HOMEBREW_PATH="/home/anthonyapodaca/Github/homebrew-amos"
AMOS_PATH="/home/anthonyapodaca/Github/amos"
GH_CLI="$HOME/.local/bin/gh"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Version required${NC}"
    echo "Usage: $0 VERSION [RELEASE_NOTES_FILE]"
    echo "Example: $0 1.2.1"
    exit 1
fi

# Validate version format (x.y.z)
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid version format. Expected x.y.z (e.g., 1.2.1)${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Starting release process for v${VERSION}${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if we're in the amos directory
cd "$AMOS_PATH"

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}Error: Uncommitted changes detected. Commit or stash them first.${NC}"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${YELLOW}Current branch: $CURRENT_BRANCH${NC}"

# Merge to main if not already on main
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Merging $CURRENT_BRANCH to main...${NC}"
    git checkout main
    git pull origin main
    git merge "$CURRENT_BRANCH" --no-edit
    echo -e "${GREEN}✓ Merged to main${NC}"
else
    echo -e "${YELLOW}Already on main, pulling latest...${NC}"
    git pull origin main
fi

echo ""

# Create annotated tag
echo -e "${YELLOW}Creating tag v${VERSION}...${NC}"
git tag -a "v${VERSION}" -m "Release v${VERSION}"
echo -e "${GREEN}✓ Tag created${NC}"

# Push tag
echo -e "${YELLOW}Pushing tag to GitHub...${NC}"
git push origin "v${VERSION}"
echo -e "${GREEN}✓ Tag pushed${NC}"

echo ""

# Generate or read release notes
RELEASE_NOTES=""
if [ -n "$RELEASE_NOTES_FILE" ] && [ -f "$RELEASE_NOTES_FILE" ]; then
    echo -e "${YELLOW}Reading release notes from $RELEASE_NOTES_FILE...${NC}"
    RELEASE_NOTES=$(cat "$RELEASE_NOTES_FILE")
else
    echo -e "${YELLOW}Generating release notes from git commits...${NC}"
    # Get previous tag
    PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")

    if [ -n "$PREV_TAG" ]; then
        echo -e "${YELLOW}Changes since $PREV_TAG:${NC}"
        COMMITS=$(git log "$PREV_TAG..HEAD" --pretty=format:"- %s" --no-merges)
        RELEASE_NOTES="## Changes

$COMMITS

## Installation

Via Homebrew:
\`\`\`bash
brew tap apodacaa/amos
brew install amos
# or upgrade existing installation
brew upgrade amos
\`\`\`"
    else
        RELEASE_NOTES="Release v${VERSION}

## Installation

Via Homebrew:
\`\`\`bash
brew tap apodacaa/amos
brew install amos
\`\`\`"
    fi
fi

# Create GitHub release
echo -e "${YELLOW}Creating GitHub release...${NC}"
$GH_CLI release create "v${VERSION}" \
    --repo apodacaa/amos \
    --title "Release v${VERSION}" \
    --notes "$RELEASE_NOTES"
echo -e "${GREEN}✓ GitHub release created${NC}"
echo ""

# Wait for GitHub to process release
echo -e "${YELLOW}Waiting for GitHub to process release...${NC}"
sleep 5

# Download tarball and calculate SHA256
echo -e "${YELLOW}Downloading release tarball...${NC}"
TARBALL_URL="https://github.com/apodacaa/amos/archive/refs/tags/v${VERSION}.tar.gz"
TARBALL_PATH="/tmp/amos-${VERSION}.tar.gz"
curl -L -o "$TARBALL_PATH" "$TARBALL_URL" 2>/dev/null

echo -e "${YELLOW}Calculating SHA256...${NC}"
SHA256=$(sha256sum "$TARBALL_PATH" | awk '{print $1}')
echo -e "${GREEN}✓ SHA256: $SHA256${NC}"
echo ""

# Update Homebrew formula
echo -e "${YELLOW}Updating Homebrew formula...${NC}"
cd "$HOMEBREW_PATH"

# Update URL
sed -i "s|url \"https://github.com/apodacaa/amos/archive/refs/tags/v[0-9.]*\.tar\.gz\"|url \"https://github.com/apodacaa/amos/archive/refs/tags/v${VERSION}.tar.gz\"|" Formula/amos.rb

# Update SHA256
sed -i "s|sha256 \"[a-f0-9]*\"|sha256 \"${SHA256}\"|" Formula/amos.rb

# Update ldflags version
sed -i "s|ldflags: \"-X main.Version=[0-9.]*\"|ldflags: \"-X main.Version=${VERSION}\"|" Formula/amos.rb

# Update test version
sed -i "s|amos version [0-9.]*|amos version ${VERSION}|" Formula/amos.rb

echo -e "${GREEN}✓ Formula updated${NC}"

# Show diff
echo -e "${YELLOW}Formula changes:${NC}"
git diff Formula/amos.rb

echo ""

# Commit and push formula
echo -e "${YELLOW}Committing formula changes...${NC}"
git add Formula/amos.rb
git commit -m "Bump version to v${VERSION}

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

echo -e "${YELLOW}Pushing formula to GitHub...${NC}"
git push origin main

cd "$AMOS_PATH"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Release v${VERSION} complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Release URL:${NC} https://github.com/apodacaa/amos/releases/tag/v${VERSION}"
echo ""
echo -e "${YELLOW}Users can install with:${NC}"
echo "  brew tap apodacaa/amos"
echo "  brew install amos"
echo "  # or upgrade:"
echo "  brew upgrade amos"
