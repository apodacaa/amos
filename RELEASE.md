# Release Process

This document describes how to release a new version of Amos.

## Prerequisites

- All changes committed and pushed to `main` branch
- All tests passing (`make ci`)
- gh CLI installed and authenticated (`~/.local/bin/gh`)
  - Install: Download from https://github.com/cli/cli/releases and copy to `~/.local/bin/gh`
  - Authenticate: `~/.local/bin/gh auth login`
- `homebrew-amos` repository cloned at `~/Github/homebrew-amos`
- `chocolatey-amos` repository cloned at `~/Github/chocolatey-amos`

## Automated Release (Recommended)

The release process is fully automated via `make release`.

### Step 1: Update Version in Makefile

Update the version in `Makefile` (all build targets):

```makefile
build: ## Build the binary
	go build -ldflags "-X main.Version=X.Y.Z" -o amos
```

Replace `X.Y.Z` with the new version in:
- `build`
- `build-windows`
- `build-all` (3 instances: linux, darwin amd64, darwin arm64, windows)

### Step 2: Test with Dry-Run

Test the release process without making any changes:

```bash
make release VERSION=X.Y.Z DRY_RUN=true
```

This will show all commands that would be executed in yellow `[DRY RUN]` format:
- Git operations (merge, tag, push)
- GitHub release creation
- Tarball download and SHA256 calculation
- Homebrew formula updates
- Git commits to homebrew-amos repo

Review the output to ensure everything looks correct.

### Step 3: Run Automated Release

```bash
make release VERSION=X.Y.Z
```

This single command automatically:
1. Merges current branch to main (if not on main)
2. Creates annotated git tag `vX.Y.Z`
3. Pushes tag to GitHub
4. Generates release notes from commits since last tag
5. Creates GitHub release with notes
6. Builds Windows binary (amd64) and uploads to GitHub release
7. Calculates SHA256 of Windows binary
8. Updates Chocolatey package in `chocolatey-amos` repo:
   - Version in .nuspec
   - Release notes URL in .nuspec
   - Version in chocolateyInstall.ps1
   - SHA256 checksum in chocolateyInstall.ps1
9. Commits and pushes Chocolatey package changes
10. Downloads release tarball and calculates SHA256
11. Updates Homebrew formula in `homebrew-amos` repo:
   - URL with new version
   - SHA256 hash
   - Version in ldflags
   - Version in test assertion
12. Commits and pushes Homebrew formula changes

### Step 4: Verify Release

1. Check GitHub release: https://github.com/apodacaa/amos/releases/tag/vX.Y.Z
   - Verify source tarball is present
   - Verify Windows binary (`amos-windows-amd64-X.Y.Z.exe`) is attached
2. Check Chocolatey package: https://github.com/apodacaa/chocolatey-amos
3. Check Homebrew formula: https://github.com/apodacaa/homebrew-amos/blob/main/Formula/amos.rb
4. Test Homebrew installation:
   ```bash
   brew update
   brew upgrade amos
   amos --version  # Should show: amos version X.Y.Z
   ```
5. Test Chocolatey installation (on Windows):
   ```powershell
   choco source add -n=amos -s="https://github.com/apodacaa/chocolatey-amos"
   choco upgrade amos
   amos --version  # Should show: amos version X.Y.Z
   ```

## Manual Release (Fallback)

Use these manual steps if the automated release fails or for troubleshooting.

### Step 1: Update Version

Update version in `Makefile` (4 places: build, build-windows, build-all targets).

### Step 2: Commit Version Change

```bash
git add Makefile
git commit -m "Bump version to vX.Y.Z"
git push origin main
```

### Step 3: Create and Push Git Tag

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

### Step 4: Create GitHub Release

```bash
~/.local/bin/gh release create vX.Y.Z \
  --repo apodacaa/amos \
  --title "Release vX.Y.Z" \
  --notes "Release notes here"
```

Or use GitHub web UI:
1. Go to: https://github.com/apodacaa/amos/releases/new
2. Select tag: `vX.Y.Z`
3. Add release notes
4. Click **Publish release**

### Step 5: Calculate SHA256 Hash

```bash
curl -L https://github.com/apodacaa/amos/archive/refs/tags/vX.Y.Z.tar.gz -o /tmp/amos-X.Y.Z.tar.gz
sha256sum /tmp/amos-X.Y.Z.tar.gz
```

### Step 6: Update Homebrew Formula

In the `homebrew-amos` repository, update `Formula/amos.rb`:

```ruby
class Amos < Formula
  url "https://github.com/apodacaa/amos/archive/refs/tags/vX.Y.Z.tar.gz"
  sha256 "NEW_SHA256_HASH_HERE"

  def install
    system "go", "build", *std_go_args(ldflags: "-X main.Version=X.Y.Z")
  end

  test do
    assert_match "amos version X.Y.Z", shell_output("#{bin}/amos --version")
  end
end
```

### Step 7: Commit and Push Homebrew Formula

```bash
cd ../homebrew-amos
git add Formula/amos.rb
git commit -m "Bump version to vX.Y.Z"
git push origin main
```

## Versioning Guidelines

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes, incompatible API changes
- **MINOR** (x.Y.0): New features, backward-compatible
- **PATCH** (x.y.Z): Bug fixes, backward-compatible

Examples:
- `1.0.0` → `1.1.0`: Added new filtering feature
- `1.1.0` → `1.1.1`: Fixed bug in date parsing
- `1.1.1` → `2.0.0`: Changed data storage format (breaking)

## Release Checklist

### Automated Release
- [ ] All changes committed and pushed to `main`
- [ ] Tests passing (`make ci`)
- [ ] Version updated in `Makefile` (all build targets)
- [ ] Dry-run tested (`make release VERSION=X.Y.Z DRY_RUN=true`)
- [ ] Dry-run output reviewed (no errors, versions correct)
- [ ] Automated release executed (`make release VERSION=X.Y.Z`)
- [ ] GitHub release verified (source tarball + Windows binary present)
- [ ] Chocolatey package verified (version and SHA256 updated)
- [ ] Homebrew formula verified (version and SHA256 updated)
- [ ] Homebrew installation tested (`brew update && brew upgrade amos && amos --version`)
- [ ] Chocolatey installation tested on Windows (optional, if Windows available)

### Manual Release (if automation fails)
- [ ] All changes committed and pushed
- [ ] Tests passing (`make ci`)
- [ ] Version updated in `Makefile`
- [ ] Version change committed
- [ ] Git tag created and pushed
- [ ] GitHub Release created with notes
- [ ] SHA256 calculated for tarball
- [ ] Homebrew formula updated (url, sha256, version in 3 places)
- [ ] Homebrew formula committed and pushed
- [ ] Installation tested with Homebrew

## Troubleshooting

### Automated Release

**Dry-run shows wrong version:**
- Check that you updated all build targets in Makefile
- Verify VERSION parameter matches Makefile version

**"gh CLI not found":**
- Ensure gh is installed at `~/.local/bin/gh`
- Run `~/.local/bin/gh auth status` to verify authentication
- If not authenticated: `~/.local/bin/gh auth login`

**Release script fails mid-process:**
- Check script output for specific error
- Use manual steps to complete remaining steps
- Clean up partial release if needed:
  ```bash
  git tag -d vX.Y.Z              # Delete local tag
  git push origin :refs/tags/vX.Y.Z  # Delete remote tag
  gh release delete vX.Y.Z --yes     # Delete GitHub release
  ```

**Uncommitted changes detected:**
- Commit or stash your changes first in amos, homebrew-amos, or chocolatey-amos repos
- Or run in dry-run mode: `make release VERSION=X.Y.Z DRY_RUN=true`
  (dry-run skips uncommitted changes check)

**"Chocolatey package repo not found" error:**
- Clone the chocolatey-amos repo: `git clone https://github.com/apodacaa/chocolatey-amos ~/Github/chocolatey-amos`
- Or update CHOCOLATEY_PATH in scripts/release.sh to match your local path

### Homebrew

**Formula fails to install:**
- Verify SHA256 is correct
- Check that the tag exists on GitHub
- Ensure version strings match in all 3 places in formula

**--version shows wrong version:**
- Verify ldflags in formula has correct version
- Uninstall and reinstall: `brew uninstall amos && brew install amos`

**brew update doesn't show new version:**
- Formula changes must be pushed to `homebrew-amos` repo
- Run `brew update` to refresh tap
- Check formula on GitHub: https://github.com/apodacaa/homebrew-amos/blob/main/Formula/amos.rb

### Chocolatey

**Package fails to install:**
- Verify Windows binary exists in GitHub release
- Check SHA256 matches in chocolateyInstall.ps1
- Verify version strings match in .nuspec and chocolateyInstall.ps1

**--version shows wrong version:**
- Verify the Windows binary was built with correct ldflags
- Check that version in chocolateyInstall.ps1 matches release
- Uninstall and reinstall: `choco uninstall amos && choco install amos`

**choco upgrade doesn't find new version:**
- Package changes must be pushed to `chocolatey-amos` repo
- Re-add source: `choco source remove -n=amos && choco source add -n=amos -s="https://github.com/apodacaa/chocolatey-amos"`
- Check package on GitHub: https://github.com/apodacaa/chocolatey-amos

**Windows binary not in GitHub release:**
- Check that release.sh completed Windows binary upload step
- Manually upload: `gh release upload vX.Y.Z amos-windows-amd64-X.Y.Z.exe --repo apodacaa/amos`
- Update chocolatey-amos repo with correct SHA256

### GitHub Release

**Release notes are auto-generated from commits:**
- If you want custom notes, create a file (e.g., `release-notes.md`)
- Run: `make release VERSION=X.Y.Z NOTES=release-notes.md`

**Need to delete a bad release:**
```bash
# Delete tag locally and remotely
git tag -d vX.Y.Z
git push origin :refs/tags/vX.Y.Z

# Delete GitHub release
~/.local/bin/gh release delete vX.Y.Z --yes --repo apodacaa/amos

# Revert Homebrew formula (if already pushed)
cd ../homebrew-amos
git revert HEAD
git push origin main
```
