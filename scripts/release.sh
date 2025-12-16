#!/bin/bash
set -e

# Release automation script for amos
# Usage: ./scripts/release.sh VERSION [RELEASE_NOTES_FILE]
# Example: ./scripts/release.sh 1.2.1
# Example: ./scripts/release.sh 1.2.1 release-notes.md

VERSION=$1
RELEASE_NOTES_FILE=$2
DRY_RUN=$3
HOMEBREW_PATH="/home/anthonyapodaca/Github/homebrew-amos"
CHOCOLATEY_PATH="/home/anthonyapodaca/Github/chocolatey-amos"
AMOS_PATH="/home/anthonyapodaca/Github/amos"
GH_CLI="$HOME/.local/bin/gh"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configure git for non-interactive mode
export GIT_TERMINAL_PROMPT=0
export GIT_SSH_COMMAND="ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=accept-new"

# Cleanup function for failures
cleanup_on_error() {
    if [ "$DRY_RUN" != "true" ] && [ -n "$VERSION" ]; then
        echo -e "${RED}Release failed. Cleaning up partial changes...${NC}"
        cd "$AMOS_PATH" 2>/dev/null || true
        # Delete local tag if it exists
        git tag -d "v${VERSION}" 2>/dev/null && echo -e "${YELLOW}✓ Deleted local tag v${VERSION}${NC}" || true
        # Note: We don't delete the GitHub release as it may be intentional to keep for debugging
        # Clean up tarball and Windows binary
        rm -f "/tmp/amos-${VERSION}.tar.gz" 2>/dev/null || true
        rm -f "/tmp/amos-windows-amd64-${VERSION}.exe" 2>/dev/null || true
    fi
}
trap cleanup_on_error ERR

# Logging helper with timestamps
log_step() {
    echo -e "${GREEN}[$(date '+%H:%M:%S')]${NC} $*"
}

# Helper function to run commands (respects DRY_RUN)
run_cmd() {
    if [ "$DRY_RUN" = "true" ]; then
        echo -e "${YELLOW}[DRY RUN]${NC} $*"
    else
        "$@"
    fi
}

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Version required${NC}"
    echo "Usage: $0 VERSION [RELEASE_NOTES_FILE] [DRY_RUN]"
    echo "Example: $0 1.2.1"
    echo "Example: $0 1.2.1 release-notes.md"
    echo "Example: $0 1.2.1 \"\" true  # Dry run"
    exit 1
fi

# Validate version format (x.y.z)
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid version format. Expected x.y.z (e.g., 1.2.1)${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Starting release process for v${VERSION}${NC}"
if [ "$DRY_RUN" = "true" ]; then
    echo -e "${YELLOW}[DRY RUN MODE - No changes will be made]${NC}"
fi
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if we're in the amos directory
cd "$AMOS_PATH"

# Check for uncommitted changes (skip in dry-run mode)
if [ "$DRY_RUN" != "true" ]; then
    if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}Error: Uncommitted changes detected. Commit or stash them first.${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}[DRY RUN] Skipping uncommitted changes check${NC}"
fi

# Pre-flight checks
log_step "Running pre-flight checks..."

# Check gh CLI is installed and authenticated
if ! command -v "$GH_CLI" &> /dev/null; then
    echo -e "${RED}Error: gh CLI not found at $GH_CLI${NC}"
    exit 1
fi

if [ "$DRY_RUN" != "true" ]; then
    if ! "$GH_CLI" auth status &>/dev/null; then
        echo -e "${RED}Error: gh CLI not authenticated. Run: gh auth login${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✓ gh CLI authenticated${NC}"

# Check homebrew-amos repo exists and is clean
if [ ! -d "$HOMEBREW_PATH" ]; then
    echo -e "${RED}Error: Homebrew tap not found at $HOMEBREW_PATH${NC}"
    exit 1
fi

if [ "$DRY_RUN" != "true" ]; then
    cd "$HOMEBREW_PATH"
    if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}Error: homebrew-amos has uncommitted changes${NC}"
        cd "$AMOS_PATH"
        exit 1
    fi
    cd "$AMOS_PATH"
fi
echo -e "${GREEN}✓ Homebrew tap is clean${NC}"

# Check chocolatey-amos repo exists and is clean
if [ ! -d "$CHOCOLATEY_PATH" ]; then
    echo -e "${RED}Error: Chocolatey package repo not found at $CHOCOLATEY_PATH${NC}"
    exit 1
fi

if [ "$DRY_RUN" != "true" ]; then
    cd "$CHOCOLATEY_PATH"
    if ! git diff-index --quiet HEAD --; then
        echo -e "${RED}Error: chocolatey-amos has uncommitted changes${NC}"
        cd "$AMOS_PATH"
        exit 1
    fi
    cd "$AMOS_PATH"
fi
echo -e "${GREEN}✓ Chocolatey package repo is clean${NC}"

# Check if tag already exists
if git rev-parse "v${VERSION}" &>/dev/null; then
    echo -e "${RED}Error: Tag v${VERSION} already exists${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Tag v${VERSION} does not exist${NC}"

echo ""

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${YELLOW}Current branch: $CURRENT_BRANCH${NC}"

# Merge to main if not already on main
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo -e "${YELLOW}Merging $CURRENT_BRANCH to main...${NC}"
    run_cmd git checkout main
    run_cmd git pull origin main
    run_cmd git merge "$CURRENT_BRANCH" --no-edit
    echo -e "${GREEN}✓ Merged to main${NC}"
else
    echo -e "${YELLOW}Already on main, pulling latest...${NC}"
    run_cmd git pull origin main
fi

echo ""

# Create annotated tag
log_step "Creating tag v${VERSION}..."
run_cmd git tag -a "v${VERSION}" -m "Release v${VERSION}"
echo -e "${GREEN}✓ Tag created${NC}"

# Push tag
log_step "Pushing tag to GitHub..."
run_cmd git push origin "v${VERSION}"
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
log_step "Creating GitHub release..."
run_cmd $GH_CLI release create "v${VERSION}" \
    --repo apodacaa/amos \
    --title "Release v${VERSION}" \
    --notes "$RELEASE_NOTES"
echo -e "${GREEN}✓ GitHub release created${NC}"
echo ""

# Build and upload Windows binary for Chocolatey
log_step "Building Windows binary..."
WINDOWS_BINARY_PATH="/tmp/amos-windows-amd64-${VERSION}.exe"

if [ "$DRY_RUN" != "true" ]; then
    GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o "$WINDOWS_BINARY_PATH"
    echo -e "${GREEN}✓ Windows binary built${NC}"
else
    echo -e "${YELLOW}[DRY RUN]${NC} GOOS=windows GOARCH=amd64 go build -ldflags \"-X main.Version=${VERSION}\" -o $WINDOWS_BINARY_PATH"
fi

log_step "Uploading Windows binary to GitHub release..."
run_cmd $GH_CLI release upload "v${VERSION}" "$WINDOWS_BINARY_PATH" \
    --repo apodacaa/amos \
    --clobber
echo -e "${GREEN}✓ Windows binary uploaded${NC}"
echo ""

# Wait for GitHub to process the binary
if [ "$DRY_RUN" != "true" ]; then
    log_step "Waiting for Windows binary to be available..."
    BINARY_URL="https://github.com/apodacaa/amos/releases/download/v${VERSION}/amos-windows-amd64-${VERSION}.exe"
    BINARY_READY=false
    for i in {1..30}; do
        if timeout 5 curl -L -I -f "$BINARY_URL" &>/dev/null; then
            echo -e "${GREEN}✓ Windows binary available (took $((i * 5)) seconds)${NC}"
            BINARY_READY=true
            break
        fi
        if [ $i -eq 30 ]; then
            echo -e "${RED}Error: Windows binary not available after 2.5 minutes${NC}"
            exit 1
        fi
        echo -n "."
        sleep 5
    done
    echo ""
else
    echo -e "${YELLOW}[DRY RUN] Skipping Windows binary availability check${NC}"
fi

# Calculate SHA256 of Windows binary
log_step "Calculating SHA256 of Windows binary..."
if [ "$DRY_RUN" = "true" ]; then
    WINDOWS_SHA256="<WINDOWS_SHA256_WOULD_BE_CALCULATED>"
    echo -e "${YELLOW}[DRY RUN]${NC} sha256sum $WINDOWS_BINARY_PATH"
else
    WINDOWS_SHA256=$(sha256sum "$WINDOWS_BINARY_PATH" | awk '{print $1}')
fi
echo -e "${GREEN}✓ Windows SHA256: $WINDOWS_SHA256${NC}"
echo ""

# Update Chocolatey package
log_step "Updating Chocolatey package..."
cd "$CHOCOLATEY_PATH"

# Update .nuspec version
run_cmd sed -i "s|<version>[0-9.]*</version>|<version>${VERSION}</version>|" amos.nuspec

# Update .nuspec release notes URL
run_cmd sed -i "s|<releaseNotes>https://github.com/apodacaa/amos/releases/tag/v[0-9.]*</releaseNotes>|<releaseNotes>https://github.com/apodacaa/amos/releases/tag/v${VERSION}</releaseNotes>|" amos.nuspec

# Update chocolateyInstall.ps1 version
run_cmd sed -i "s|\$version = '[0-9.]*'|\$version = '${VERSION}'|" tools/chocolateyInstall.ps1

# Update chocolateyInstall.ps1 SHA256
run_cmd sed -i "s|\$checksum64 = '[A-Za-z0-9_]*'|\$checksum64 = '${WINDOWS_SHA256}'|" tools/chocolateyInstall.ps1

echo -e "${GREEN}✓ Chocolatey package updated${NC}"

# Show diff
echo -e "${YELLOW}Chocolatey package changes:${NC}"
git --no-pager diff amos.nuspec tools/chocolateyInstall.ps1

echo ""

# Commit and push Chocolatey package
log_step "Committing Chocolatey package changes..."
run_cmd git add amos.nuspec tools/chocolateyInstall.ps1
run_cmd git commit -m "Bump version to v${VERSION}

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

log_step "Pushing Chocolatey package to GitHub..."
run_cmd git push origin main

cd "$AMOS_PATH"

# Clean up Windows binary
if [ "$DRY_RUN" != "true" ]; then
    rm -f "$WINDOWS_BINARY_PATH"
    echo -e "${GREEN}✓ Cleaned up Windows binary${NC}"
fi

echo ""

# Wait for GitHub to process release and make tarball available
TARBALL_URL="https://github.com/apodacaa/amos/archive/refs/tags/v${VERSION}.tar.gz"
TARBALL_PATH="/tmp/amos-${VERSION}.tar.gz"

if [ "$DRY_RUN" != "true" ]; then
    log_step "Waiting for GitHub to generate release tarball..."
    TARBALL_READY=false
    for i in {1..30}; do
        if timeout 5 curl -L -I -f "$TARBALL_URL" &>/dev/null; then
            echo -e "${GREEN}✓ Tarball available (took $((i * 5)) seconds)${NC}"
            TARBALL_READY=true
            break
        fi
        if [ $i -eq 30 ]; then
            echo -e "${RED}Error: Tarball not available after 2.5 minutes${NC}"
            exit 1
        fi
        echo -n "."
        sleep 5
    done
    echo ""
else
    echo -e "${YELLOW}[DRY RUN] Skipping tarball availability check${NC}"
fi

# Download tarball with retry logic
log_step "Downloading release tarball..."

if [ "$DRY_RUN" != "true" ]; then
    DOWNLOAD_SUCCESS=false
    for attempt in {1..3}; do
        if curl -L -f --max-time 60 --connect-timeout 10 -o "$TARBALL_PATH" "$TARBALL_URL"; then
            DOWNLOAD_SUCCESS=true
            break
        fi
        if [ $attempt -lt 3 ]; then
            echo -e "${YELLOW}Download failed (attempt $attempt/3), retrying in 5 seconds...${NC}"
            sleep 5
        fi
    done

    if [ "$DOWNLOAD_SUCCESS" = false ]; then
        echo -e "${RED}Error: Failed to download tarball after 3 attempts${NC}"
        exit 1
    fi

    # Verify tarball was actually downloaded
    if [ ! -f "$TARBALL_PATH" ] || [ ! -s "$TARBALL_PATH" ]; then
        echo -e "${RED}Error: Tarball file not created or is empty${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Tarball downloaded successfully${NC}"
else
    echo -e "${YELLOW}[DRY RUN]${NC} curl -L -f --max-time 60 --connect-timeout 10 -o $TARBALL_PATH $TARBALL_URL"
fi

echo -e "${YELLOW}Calculating SHA256...${NC}"
if [ "$DRY_RUN" = "true" ]; then
    SHA256="<SHA256_WOULD_BE_CALCULATED>"
    echo -e "${YELLOW}[DRY RUN]${NC} sha256sum $TARBALL_PATH"
else
    SHA256=$(sha256sum "$TARBALL_PATH" | awk '{print $1}')
fi
echo -e "${GREEN}✓ SHA256: $SHA256${NC}"
echo ""

# Update Homebrew formula
log_step "Updating Homebrew formula..."
cd "$HOMEBREW_PATH"

# Update URL
run_cmd sed -i "s|url \"https://github.com/apodacaa/amos/archive/refs/tags/v[0-9.]*\.tar\.gz\"|url \"https://github.com/apodacaa/amos/archive/refs/tags/v${VERSION}.tar.gz\"|" Formula/amos.rb

# Update SHA256
run_cmd sed -i "s|sha256 \"[a-f0-9]*\"|sha256 \"${SHA256}\"|" Formula/amos.rb

# Update ldflags version
run_cmd sed -i "s|ldflags: \"-X main.Version=[0-9.]*\"|ldflags: \"-X main.Version=${VERSION}\"|" Formula/amos.rb

# Update test version
run_cmd sed -i "s|amos version [0-9.]*|amos version ${VERSION}|" Formula/amos.rb

echo -e "${GREEN}✓ Formula updated${NC}"

# Show diff
echo -e "${YELLOW}Formula changes:${NC}"
git --no-pager diff Formula/amos.rb

echo ""

# Commit and push formula
log_step "Committing formula changes..."
run_cmd git add Formula/amos.rb
run_cmd git commit -m "Bump version to v${VERSION}

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

log_step "Pushing formula to GitHub..."
run_cmd git push origin main

cd "$AMOS_PATH"

# Clean up temporary tarball
if [ "$DRY_RUN" != "true" ]; then
    rm -f "$TARBALL_PATH"
    echo -e "${GREEN}✓ Cleaned up temporary files${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Release v${VERSION} complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Release URL:${NC} https://github.com/apodacaa/amos/releases/tag/v${VERSION}"
echo ""
echo -e "${YELLOW}Users can install with:${NC}"
echo ""
echo "  Homebrew (macOS/Linux):"
echo "    brew tap apodacaa/amos"
echo "    brew install amos"
echo "    # or upgrade:"
echo "    brew upgrade amos"
echo ""
echo "  Chocolatey (Windows):"
echo "    choco source add -n=amos -s=\"https://github.com/apodacaa/chocolatey-amos\""
echo "    choco install amos"
echo "    # or upgrade:"
echo "    choco upgrade amos"
