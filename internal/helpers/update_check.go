package helpers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	githubAPIURL   = "https://api.github.com/repos/apodacaa/amos/releases/latest"
	requestTimeout = 3 * time.Second
)

// GitHubRelease represents the GitHub API response for latest release
type GitHubRelease struct {
	TagName string `json:"tag_name"` // e.g., "v1.5.0"
}

// FetchLatestVersion fetches the latest release version from GitHub
// Returns empty string on any error (silent failure)
func FetchLatestVersion() string {
	// Create request with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", githubAPIURL, nil)
	if err != nil {
		return "" // silent failure
	}

	// Set User-Agent header (GitHub API requires this)
	req.Header.Set("User-Agent", "amos-tui")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "" // network error, timeout, etc.
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "" // API error, rate limit, etc.
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return ""
	}

	return release.TagName // e.g., "v1.5.0"
}
