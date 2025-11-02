# Release Process

This document describes how to release a new version of Amos.

## Prerequisites

- All changes committed and pushed to `main` branch
- All tests passing (`make ci`)
- Changelog/release notes prepared

## Step 1: Update Version

Update the version in `Makefile` (4 places):

```makefile
build: ## Build the binary
	go build -ldflags "-X main.Version=X.Y.Z" -o amos

build-windows: ## Build Windows binary (amd64)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=X.Y.Z" -o amos.exe

build-all: ## Build binaries for all platforms
	...
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=X.Y.Z" -o amos-windows-amd64.exe
	...
```

Replace `X.Y.Z` with the new version in all build targets (e.g., `1.1.0`, `2.0.0`).

## Step 2: Commit Version Change

```bash
git add Makefile
git commit -m "Bump version to vX.Y.Z"
git push origin main
```

## Step 3: Create and Push Git Tag

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z

- Feature 1
- Feature 2
- Bug fix 3
"

git push origin vX.Y.Z
```

## Step 4: Create GitHub Release

1. Go to: https://github.com/apodacaa/amos/releases/new
2. Select tag: `vX.Y.Z`
3. Release title: `vX.Y.Z - Brief Description`
4. Add release notes describing changes
5. Click **Publish release**

## Step 5: Build and Upload Windows Binary

Build the Windows binary and upload it to the GitHub release:

```bash
# Build Windows binary
make build-windows

# Rename to include version
mv amos.exe amos-windows-amd64.exe
```

Upload to GitHub Release:
1. Go to the release you just created
2. Click **Edit release**
3. Drag and drop `amos-windows-amd64.exe` to the **Attach binaries** section
4. Click **Update release**

## Step 6: Calculate SHA256 Hash

Calculate SHA256 for Homebrew:

```bash
# Download source tarball
curl -L https://github.com/apodacaa/amos/archive/refs/tags/vX.Y.Z.tar.gz -o /tmp/amos-X.Y.Z.tar.gz
sha256sum /tmp/amos-X.Y.Z.tar.gz
```

Copy the SHA256 hash (first part of output).

## Step 7: Update Homebrew Formula

In the `homebrew-amos` repository, update `Formula/amos.rb`:

```ruby
class Amos < Formula
  desc "Minimal TUI for journal + todo management with brutalist design"
  homepage "https://github.com/apodacaa/amos"
  url "https://github.com/apodacaa/amos/archive/refs/tags/vX.Y.Z.tar.gz"  # Update version
  sha256 "NEW_SHA256_HASH_HERE"  # Update hash from Step 6
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-X main.Version=X.Y.Z")  # Update version
  end

  test do
    assert_match "amos version X.Y.Z", shell_output("#{bin}/amos --version")  # Update version
  end
end
```

## Step 8: Commit and Push Homebrew Formula

```bash
cd ../homebrew-amos
git add Formula/amos.rb
git commit -m "Update to vX.Y.Z"
git push origin main
```

## Step 9: Test Installation

### Homebrew (macOS/Linux)

```bash
# Uninstall old version
brew uninstall amos

# Update tap
brew update

# Install new version
brew install amos

# Verify version
amos --version
```

Should output: `amos version X.Y.Z`

## Versioning Guidelines

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes, incompatible API changes
- **MINOR** (x.Y.0): New features, backward-compatible
- **PATCH** (x.y.Z): Bug fixes, backward-compatible

Examples:
- `1.0.0` → `1.1.0`: Added new filtering feature
- `1.1.0` → `1.1.1`: Fixed bug in date parsing
- `1.1.1` → `2.0.0`: Changed data storage format (breaking)

## Checklist

Use this checklist for each release:

- [ ] All changes committed and pushed
- [ ] Tests passing (`make ci`)
- [ ] Version updated in `Makefile` (4 places: build, build-windows, build-all targets)
- [ ] Version change committed
- [ ] Git tag created and pushed
- [ ] GitHub Release created with notes
- [ ] Windows binary built (`make build-windows`)
- [ ] Windows binary uploaded to GitHub Release
- [ ] SHA256 calculated for tarball (Homebrew)
- [ ] Homebrew formula updated (url, sha256, version in 3 places)
- [ ] Homebrew formula committed and pushed
- [ ] Installation tested with Homebrew (macOS/Linux)

## Troubleshooting

### Homebrew

**"Formula fails to install"**
- Verify SHA256 is correct
- Check that the tag exists on GitHub
- Ensure version strings match in all 3 places in formula

**"--version shows wrong version"**
- Verify ldflags in formula has correct version
- Uninstall and reinstall: `brew uninstall amos && brew install amos`

**"brew update doesn't show new version"**
- Formula changes must be pushed to `homebrew-amos` repo
- Run `brew update` to refresh tap
