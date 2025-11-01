# Scoop Bucket Setup

This document describes how to set up the `scoop-amos` repository for Windows distribution.

## Overview

Scoop is a Windows package manager similar to Homebrew. It requires a separate "bucket" repository containing the app manifest.

## Repository Setup

1. Create a new GitHub repository named `scoop-amos`

2. Create the following directory structure:
```
scoop-amos/
├── bucket/
│   └── amos.json
└── README.md
```

3. Copy `scoop-manifest-example.json` to `bucket/amos.json` in the new repo

4. Update the manifest with the current version, URL, and SHA256 hash

## Manifest File

The `bucket/amos.json` file should contain:

```json
{
  "version": "X.Y.Z",
  "description": "Minimal TUI for journal + todo management with brutalist design",
  "homepage": "https://github.com/apodacaa/amos",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/apodacaa/amos/releases/download/vX.Y.Z/amos-windows-amd64.exe",
      "hash": "SHA256_HASH_HERE",
      "bin": [["amos-windows-amd64.exe", "amos"]]
    }
  },
  "checkver": {
    "github": "https://github.com/apodacaa/amos"
  },
  "autoupdate": {
    "architecture": {
      "64bit": {
        "url": "https://github.com/apodacaa/amos/releases/download/v$version/amos-windows-amd64.exe"
      }
    }
  }
}
```

## README for scoop-amos Repository

Create a `README.md` in the `scoop-amos` repository with:

```markdown
# Scoop Bucket for Amos

This is the official Scoop bucket for [Amos](https://github.com/apodacaa/amos), a minimal TUI for journal + todo management.

## Installation

```powershell
scoop bucket add amos https://github.com/apodacaa/scoop-amos
scoop install amos
```

## Updating

```powershell
scoop update amos
```

## Uninstalling

```powershell
scoop uninstall amos
```

## About Amos

Amos is a minimal Bubble Tea (Go) TUI for journal + todo management with a brutalist design philosophy.

See the [main repository](https://github.com/apodacaa/amos) for more information.
```

## Release Process Integration

When releasing a new version of Amos:

1. Build Windows binary: `make build-windows`
2. Upload `amos-windows-amd64.exe` to GitHub Release
3. Calculate SHA256: `sha256sum amos-windows-amd64.exe`
4. Update `bucket/amos.json` in `scoop-amos` repository:
   - Update `version` field
   - Update `url` with new version tag
   - Update `hash` with new SHA256
5. Commit and push changes

See `RELEASE.md` in the main repository for complete release instructions.

## Testing

Test installation locally:

```powershell
# Add bucket
scoop bucket add amos https://github.com/apodacaa/scoop-amos

# Install
scoop install amos

# Verify
amos --version
```

## Troubleshooting

**Installation fails**
- Verify SHA256 hash is correct
- Check that the Windows binary exists in the GitHub Release
- Ensure the URL in manifest points to the correct version

**Version check fails**
- Run `scoop update`
- Try removing and re-adding the bucket:
  ```powershell
  scoop bucket rm amos
  scoop bucket add amos https://github.com/apodacaa/scoop-amos
  ```

## Scoop Resources

- [Scoop Documentation](https://scoop.sh/)
- [Creating Buckets](https://github.com/ScoopInstaller/Scoop/wiki/Buckets)
- [App Manifest Reference](https://github.com/ScoopInstaller/Scoop/wiki/App-Manifests)
