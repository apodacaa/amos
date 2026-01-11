package helpers

import (
	"strconv"
	"strings"
)

// ParseVersion parses a semantic version string (e.g., "1.4.2" or "v1.4.2")
// Returns major, minor, patch as integers
// Returns 0, 0, 0 for "dev" or invalid versions
func ParseVersion(version string) (major, minor, patch int) {
	// Strip leading 'v' if present
	version = strings.TrimPrefix(version, "v")

	// Handle "dev" version or empty string
	if version == "dev" || version == "" {
		return 0, 0, 0
	}

	// Split on dots
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0
	}

	// Parse integers, ignoring errors (defaults to 0)
	major, _ = strconv.Atoi(parts[0])
	minor, _ = strconv.Atoi(parts[1])
	patch, _ = strconv.Atoi(parts[2])

	return major, minor, patch
}

// CompareVersions compares two version strings
// Returns:
//
//	 1 if latest > current (latest is newer)
//	 0 if latest == current (same version)
//	-1 if latest < current (current is newer)
func CompareVersions(current, latest string) int {
	cMajor, cMinor, cPatch := ParseVersion(current)
	lMajor, lMinor, lPatch := ParseVersion(latest)

	// Compare major
	if cMajor < lMajor {
		return 1 // latest is newer
	} else if cMajor > lMajor {
		return -1 // current is newer
	}

	// Compare minor
	if cMinor < lMinor {
		return 1
	} else if cMinor > lMinor {
		return -1
	}

	// Compare patch
	if cPatch < lPatch {
		return 1
	} else if cPatch > lPatch {
		return -1
	}

	return 0 // equal
}

// IsUpdateAvailable returns true if latest version is newer than current
func IsUpdateAvailable(currentVersion, latestVersion string) bool {
	return CompareVersions(currentVersion, latestVersion) == 1
}
