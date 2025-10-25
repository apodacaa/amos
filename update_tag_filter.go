package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// handleTagFilterKeys processes keyboard input (tag filter input view)
func (m Model) handleTagFilterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Cancel and return to filter view
		m.view = "filter_view"
		m.statusMsg = ""
		return m, nil
	case "tab":
		// Autocomplete last word with best match
		currentInput := m.tagFilterInput.Value()
		currentWord := getLastWord(currentInput)

		if currentWord != "" {
			// Find best matching tag
			match := findBestMatch(currentWord, m.availableTags)
			if match != "" {
				// Replace last word with match
				newInput := replaceLastWord(currentInput, match)
				m.tagFilterInput.SetValue(newInput)
			}
		}
		// Clear autocomplete hint after completion
		m.autocompleteTag = ""
		return m, nil
	case "enter":
		// Parse input and apply filters
		input := strings.TrimSpace(m.tagFilterInput.Value())
		if input == "" {
			// No input, clear tag filter and return to filter view
			m.filterTags = []string{}
			m.view = "filter_view"
			return m, nil
		}

		// Parse space-separated tags
		tags := parseTagInput(input)
		if len(tags) > 0 {
			m.filterTags = tags
			m.view = "filter_view"
			m.statusMsg = ""
			return m, nil
		}

		// No valid tags, show error
		m.statusMsg = "No valid tags entered"
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()
	}

	// Update textarea and recalculate autocomplete suggestion
	m.tagFilterInput, cmd = m.tagFilterInput.Update(msg)

	// Update autocomplete suggestion based on last word in input
	// Only show suggestion if user has typed at least one character after @
	currentInput := m.tagFilterInput.Value()
	currentWord := getLastWord(currentInput)

	// Check if word has content beyond just @ symbol
	wordWithoutAt := strings.TrimPrefix(currentWord, "@")
	if currentWord != "" && wordWithoutAt != "" && wordWithoutAt != currentWord {
		// User typed @ plus at least one character
		m.autocompleteTag = findBestMatch(currentWord, m.availableTags)
	} else {
		m.autocompleteTag = ""
	}

	return m, cmd
}

// getLastWord extracts the last word from the input (after the last space)
func getLastWord(input string) string {
	if input == "" {
		return ""
	}

	// Find last space
	lastSpaceIdx := strings.LastIndex(input, " ")
	if lastSpaceIdx == -1 {
		// No spaces, return entire input (trimmed)
		return strings.TrimSpace(input)
	}

	// Return everything after the last space (trimmed)
	return strings.TrimSpace(input[lastSpaceIdx+1:])
}

// findBestMatch finds the first tag that starts with the given prefix (case-insensitive)
func findBestMatch(prefix string, availableTags []string) string {
	if prefix == "" {
		return ""
	}

	// Normalize prefix (remove @ if present, lowercase)
	normalizedPrefix := strings.ToLower(strings.TrimPrefix(prefix, "@"))

	// Find first matching tag
	for _, tag := range availableTags {
		normalizedTag := strings.ToLower(strings.TrimPrefix(tag, "@"))
		if strings.HasPrefix(normalizedTag, normalizedPrefix) {
			return tag // Return with @ prefix
		}
	}

	return ""
}

// replaceLastWord replaces the last word in the input with the completion
func replaceLastWord(input string, completion string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return completion + " "
	}

	// Find last space
	lastSpaceIdx := strings.LastIndex(input, " ")
	if lastSpaceIdx == -1 {
		// No spaces, replace entire input
		return completion + " "
	}

	// Replace everything after last space with completion
	return input[:lastSpaceIdx+1] + completion + " "
}

// parseTagInput parses space-separated tag input into a slice of tags
// Handles @prefix automatically, deduplicates, and filters empty strings
func parseTagInput(input string) []string {
	parts := strings.Fields(input) // Split by whitespace and trim
	tagMap := make(map[string]bool)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Ensure @ prefix
		if !strings.HasPrefix(part, "@") {
			part = "@" + part
		}

		// Normalize to lowercase for deduplication
		tagMap[strings.ToLower(part)] = true
	}

	// Convert map to slice
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return tags
}
