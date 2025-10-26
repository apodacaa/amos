package helpers

import (
	"strings"
)

// FilterResult holds parsed filter components from user input
type FilterResult struct {
	Tags     []string
	Date     string
	Errors   []string
	Warnings []string
}

// ParseFilterInput parses a unified filter input string
// Supports mixed input like "@client last 14 days" or "2025-10-23 @work"
// Returns FilterResult with parsed tags, date string, and validation errors
func ParseFilterInput(input string) FilterResult {
	result := FilterResult{
		Tags:     []string{},
		Date:     "",
		Errors:   []string{},
		Warnings: []string{},
	}

	if input == "" {
		return result
	}

	input = strings.TrimSpace(input)

	// Extract tags (words starting with @)
	tagMap := make(map[string]bool)
	words := strings.Fields(input)
	var nonTagWords []string

	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			tag := strings.ToLower(word)
			tagMap[tag] = true
		} else {
			nonTagWords = append(nonTagWords, word)
		}
	}

	// Convert tag map to slice
	for tag := range tagMap {
		result.Tags = append(result.Tags, tag)
	}

	// Everything that's not a tag is potentially part of the date filter
	if len(nonTagWords) > 0 {
		dateString := strings.Join(nonTagWords, " ")

		// Validate the date filter by attempting to parse it
		start, end := ParseDateFilter(dateString)
		if start.IsZero() || end.IsZero() {
			// Invalid date format
			result.Errors = append(result.Errors, "Invalid date format: '"+dateString+"'")
		} else {
			// Valid date - store it
			result.Date = dateString
		}
	}

	return result
}

// GetFilterHint returns a usage hint for the filter input
func GetFilterHint() string {
	return "e.g., @work yesterday, last 14 days, 2025-10-23, 2025-10-23 to 2025-11-30"
}

// GetDateSuggestions returns format examples for date filters (no autocomplete)
func GetDateSuggestions() []string {
	return []string{
		"e.g., today, last 14 days, 2025-10-23, 2025-10-23 to 2025-11-30",
	}
}
