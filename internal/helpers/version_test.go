package helpers

import "testing"

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input               string
		major, minor, patch int
	}{
		{"1.4.2", 1, 4, 2},
		{"v1.4.2", 1, 4, 2},
		{"2.0.0", 2, 0, 0},
		{"v3.5.7", 3, 5, 7},
		{"dev", 0, 0, 0},
		{"", 0, 0, 0},
		{"invalid", 0, 0, 0},
		{"1.4", 0, 0, 0},     // invalid format (only 2 parts)
		{"1.4.2.3", 0, 0, 0}, // invalid format (too many parts)
		{"a.b.c", 0, 0, 0},   // non-numeric
	}

	for _, tt := range tests {
		major, minor, patch := ParseVersion(tt.input)
		if major != tt.major || minor != tt.minor || patch != tt.patch {
			t.Errorf("ParseVersion(%q) = %d.%d.%d, want %d.%d.%d",
				tt.input, major, minor, patch, tt.major, tt.minor, tt.patch)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		current, latest string
		want            int
	}{
		{"1.4.2", "1.5.0", 1},   // latest is newer
		{"1.4.2", "1.4.3", 1},   // patch newer
		{"1.4.2", "2.0.0", 1},   // major newer
		{"1.4.2", "1.4.2", 0},   // equal
		{"1.5.0", "1.4.2", -1},  // current is newer (rare)
		{"2.0.0", "1.9.9", -1},  // current major is newer
		{"dev", "1.5.0", 1},     // dev is treated as 0.0.0
		{"1.4.2", "dev", -1},    // dev is older than any version
		{"dev", "dev", 0},       // both dev
		{"v1.4.2", "v1.5.0", 1}, // with 'v' prefix
	}

	for _, tt := range tests {
		got := CompareVersions(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d",
				tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestIsUpdateAvailable(t *testing.T) {
	tests := []struct {
		current, latest string
		want            bool
	}{
		{"1.4.2", "1.5.0", true},
		{"1.4.2", "1.4.3", true},
		{"1.4.2", "2.0.0", true},
		{"1.4.2", "1.4.2", false},
		{"1.5.0", "1.4.2", false},
		{"2.0.0", "1.9.9", false},
		{"dev", "1.5.0", true},
		{"1.4.2", "dev", false},
		{"dev", "dev", false},
	}

	for _, tt := range tests {
		got := IsUpdateAvailable(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("IsUpdateAvailable(%q, %q) = %v, want %v",
				tt.current, tt.latest, got, tt.want)
		}
	}
}
