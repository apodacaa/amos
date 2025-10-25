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
// Supports mixed input like "@client last 30 days" or "yesterday @work @urgent"
// Returns FilterResult with parsed tags, date preset, and any errors
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
	words := strings.Fields(input)

	// Track what we've consumed for error detection
	consumedIndices := make(map[int]bool)

	// Pass 1: Extract tags (words starting with @)
	tagMap := make(map[string]bool)
	for i, word := range words {
		if strings.HasPrefix(word, "@") {
			tag := strings.ToLower(word)
			tagMap[tag] = true
			consumedIndices[i] = true
		}
	}

	// Convert tag map to slice
	for tag := range tagMap {
		result.Tags = append(result.Tags, tag)
	}

	// Pass 2: Extract date phrases
	// Try to match multi-word date phrases first, then single words
	dateFound := false

	// Multi-word phrases
	multiWordPhrases := map[string]string{
		"last 7 days":   DateFilterLast7Days,
		"last 30 days":  DateFilterLast30Days,
		"last 60 days":  DateFilterLast60Days,
		"last 90 days":  DateFilterLast90Days,
		"last 365 days": DateFilterLast365Days,
	}

	// Check for multi-word phrases
	inputLower := strings.ToLower(input)
	for phrase, preset := range multiWordPhrases {
		if strings.Contains(inputLower, phrase) {
			result.Date = preset
			dateFound = true
			// Mark words as consumed
			phraseWords := strings.Fields(phrase)
			// Find phrase location in words array
			for i := 0; i <= len(words)-len(phraseWords); i++ {
				match := true
				for j, pw := range phraseWords {
					if i+j >= len(words) || strings.ToLower(words[i+j]) != pw {
						match = false
						break
					}
				}
				if match {
					for j := 0; j < len(phraseWords); j++ {
						consumedIndices[i+j] = true
					}
					break
				}
			}
			break
		}
	}

	// Single-word dates (if no multi-word found)
	if !dateFound {
		singleWordDates := map[string]string{
			"today":     DateFilterToday,
			"yesterday": DateFilterYesterday,
		}

		for i, word := range words {
			wordLower := strings.ToLower(word)
			if preset, exists := singleWordDates[wordLower]; exists {
				result.Date = preset
				consumedIndices[i] = true
				dateFound = true
				break
			}
		}
	}

	// Pass 3: Check for unconsumed words (errors)
	var unconsumed []string
	for i, word := range words {
		if !consumedIndices[i] {
			unconsumed = append(unconsumed, word)
		}
	}

	if len(unconsumed) > 0 {
		result.Errors = append(result.Errors, "Unrecognized: "+strings.Join(unconsumed, " "))
	}

	return result
}

// GetFilterHint returns a usage hint for the filter input
func GetFilterHint() string {
	return "e.g. @work yesterday, last 30 days @client"
}

// GetDateSuggestions returns available date filter options for autocomplete
func GetDateSuggestions() []string {
	return []string{
		"today",
		"yesterday",
		"last 7 days",
		"last 30 days",
		"last 60 days",
		"last 90 days",
		"last 365 days",
	}
}
