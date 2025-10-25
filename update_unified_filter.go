package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/apodacaa/amos/internal/helpers"
)

// handleUnifiedFilterKeys processes keyboard input (unified filter input view)
func (m Model) handleUnifiedFilterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Cancel and return to originating list
		m.view = m.filterContext
		m.statusMsg = ""
		return m, nil
	case "tab":
		// Autocomplete last word with best match (tags or dates)
		currentInput := m.unifiedFilterInput.Value()
		currentWord := getLastWord(currentInput)

		if currentWord != "" {
			var match string

			// Try tag match first (if starts with @)
			if strings.HasPrefix(currentWord, "@") {
				match = findBestTagMatch(currentWord, m.availableTags)
			} else {
				// Try date phrase match
				match = findBestDateMatch(currentWord)
			}

			if match != "" {
				// Replace last word with match
				newInput := replaceLastWord(currentInput, match)
				m.unifiedFilterInput.SetValue(newInput)
			}
		}
		// Clear autocomplete hint after completion
		m.autocompleteTag = ""
		return m, nil
	case "enter":
		// Parse input and apply filters
		input := strings.TrimSpace(m.unifiedFilterInput.Value())
		if input == "" {
			// No input, clear all filters and return to list
			m.filterTags = []string{}
			m.filterDate = ""
			m.view = m.filterContext

			// Reset selection to first item
			if m.filterContext == "entries" {
				m.selectedEntry = 0
			} else if m.filterContext == "todos" {
				m.selectedTodo = 0
			}

			return m, nil
		}

		// Parse unified input for tags and dates
		result := helpers.ParseFilterInput(input)

		// Apply parsed filters
		m.filterTags = result.Tags
		m.filterDate = result.Date

		// Show errors if any
		if len(result.Errors) > 0 {
			m.statusMsg = strings.Join(result.Errors, "; ") + ". Try: " + helpers.GetFilterHint()
			return m, nil
		}

		// Success - return to list
		m.view = m.filterContext
		m.statusMsg = ""

		// Reset selection to first item in filtered list
		if m.filterContext == "entries" {
			m.selectedEntry = 0
		} else if m.filterContext == "todos" {
			m.selectedTodo = 0
		}

		return m, nil
	}

	// Update textarea and recalculate autocomplete suggestion
	m.unifiedFilterInput, cmd = m.unifiedFilterInput.Update(msg)

	// Update autocomplete suggestion based on last word in input
	currentInput := m.unifiedFilterInput.Value()
	currentWord := getLastWord(currentInput)

	if currentWord != "" {
		// Check if typing tag (starts with @)
		if strings.HasPrefix(currentWord, "@") {
			wordWithoutAt := strings.TrimPrefix(currentWord, "@")
			if wordWithoutAt != "" {
				// User typed @ plus at least one character
				m.autocompleteTag = findBestTagMatch(currentWord, m.availableTags)
			} else {
				m.autocompleteTag = ""
			}
		} else {
			// Typing date phrase
			m.autocompleteTag = findBestDateMatch(currentWord)
		}
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

// findBestTagMatch finds the first tag that starts with the given prefix (case-insensitive)
func findBestTagMatch(prefix string, availableTags []string) string {
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

// findBestDateMatch finds the first date phrase that starts with the given prefix
func findBestDateMatch(prefix string) string {
	if prefix == "" {
		return ""
	}

	normalizedPrefix := strings.ToLower(prefix)
	dateSuggestions := helpers.GetDateSuggestions()

	for _, suggestion := range dateSuggestions {
		if strings.HasPrefix(suggestion, normalizedPrefix) {
			return suggestion
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
