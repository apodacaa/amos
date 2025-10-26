package helpers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

// Taggable is an interface for items that have tags
type Taggable interface {
	GetTags() []string
}

// filterByTags is a generic function that filters items by tags with AND logic
// Tags should be provided with @ prefix (e.g., ["@work", "@urgent"])
// Returns filtered list or original list if filterTags is empty
func filterByTags[T Taggable](items []T, filterTags []string) []T {
	if len(filterTags) == 0 {
		return items
	}

	// Normalize filter tags (remove @ prefix, lowercase)
	normalizedFilters := make([]string, len(filterTags))
	for i, tag := range filterTags {
		normalizedFilters[i] = strings.ToLower(strings.TrimPrefix(tag, "@"))
	}

	filtered := []T{}

	for _, item := range items {
		// Check if item has ALL filter tags (AND logic)
		hasAllTags := true
		for _, filterTag := range normalizedFilters {
			found := false
			for _, itemTag := range item.GetTags() {
				if strings.ToLower(itemTag) == filterTag {
					found = true
					break
				}
			}
			if !found {
				hasAllTags = false
				break
			}
		}

		if hasAllTags {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// extractUniqueTags is a generic function that collects all unique tags from a list of items
// Returns tags sorted alphabetically with @ prefix
func extractUniqueTags[T Taggable](items []T) []string {
	tagMap := make(map[string]int) // tag -> count

	for _, item := range items {
		for _, tag := range item.GetTags() {
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
	return extractUniqueTags(entries)
}

// ExtractUniqueTagsFromTodos collects all unique tags from a list of todos
// Returns tags sorted alphabetically with @ prefix
func ExtractUniqueTagsFromTodos(todos []models.Todo) []string {
	return extractUniqueTags(todos)
}

// ExtractUniqueTagsFromAll collects all unique tags from entries and todos
// Returns tags sorted alphabetically with @ prefix
func ExtractUniqueTagsFromAll(entries []models.Entry, todos []models.Todo) []string {
	tagMap := make(map[string]int) // tag -> count

	// Collect from entries
	for _, entry := range entries {
		for _, tag := range entry.GetTags() {
			tagMap[tag]++
		}
	}

	// Collect from todos
	for _, todo := range todos {
		for _, tag := range todo.GetTags() {
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

// FilterEntriesByTags filters entries to only those containing ALL specified tags (AND logic)
// Tags should be provided with @ prefix (e.g., ["@work", "@urgent"])
// Returns filtered list or original list if filterTags is empty
func FilterEntriesByTags(entries []models.Entry, filterTags []string) []models.Entry {
	return filterByTags(entries, filterTags)
}

// FilterTodosByTags filters todos to only those containing ALL specified tags (AND logic)
// Tags should be provided with @ prefix (e.g., ["@work", "@urgent"])
// Returns filtered list or original list if filterTags is empty
func FilterTodosByTags(todos []models.Todo, filterTags []string) []models.Todo {
	return filterByTags(todos, filterTags)
}
