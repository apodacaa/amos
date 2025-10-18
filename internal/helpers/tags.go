package helpers

import (
	"regexp"
	"strings"
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
