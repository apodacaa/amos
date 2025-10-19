package helpers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

// ParseEntryContent splits entry content into title and body
// First line becomes title, rest becomes body
func ParseEntryContent(content string) (title, body string) {
	lines := strings.Split(content, "\n")

	if len(lines) > 0 {
		title = strings.TrimSpace(lines[0])
		if len(lines) > 1 {
			body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		}
	}

	return title, body
}

// ExtractTags finds all @word patterns in text and returns them lowercase
func ExtractTags(text string) []string {
	// Match @word patterns (letters, numbers, underscores, hyphens)
	re := regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	// Use map to deduplicate
	tagMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			tag := strings.ToLower(match[1]) // Case-insensitive
			tagMap[tag] = true
		}
	}

	// Convert to slice
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return tags
}

// ExtractUniqueTags collects all unique tags from a list of entries
// Returns tags sorted alphabetically with @ prefix
func ExtractUniqueTags(entries []models.Entry) []string {
	tagMap := make(map[string]int) // tag -> count

	for _, entry := range entries {
		for _, tag := range entry.Tags {
			tagMap[tag]++
		}
	}

	// Convert to slice with @ prefix and sort
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, "@"+tag)
	}

	// Sort alphabetically
	sort.Strings(tags)

	return tags
}

// FilterEntriesByTag filters entries to only those containing the specified tag
// Tag should be provided with @ prefix (e.g., "@client")
// Returns filtered list or original list if filterTag is empty
func FilterEntriesByTag(entries []models.Entry, filterTag string) []models.Entry {
	if filterTag == "" {
		return entries
	}

	filtered := []models.Entry{}
	tagWithoutAt := strings.TrimPrefix(filterTag, "@")

	for _, entry := range entries {
		for _, entryTag := range entry.Tags {
			if entryTag == tagWithoutAt {
				filtered = append(filtered, entry)
				break
			}
		}
	}

	return filtered
}
