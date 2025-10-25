package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderUnifiedFilter renders the unified filter input view (tags + dates)
func RenderUnifiedFilter(width, height int, ti textarea.Model, availableTags []string, autocompleteTag string, statusMsg string) string {
	// Header
	header := RenderHeader(width, "tab", "complete", "enter", "apply", "esc", "cancel")

	// Footer with hint
	footerTitle := "Filter"
	footerHint := helpers.GetFilterHint()
	if statusMsg != "" {
		footerHint = statusMsg // Show error/status instead of hint
	}
	footer := RenderFooter(width, footerTitle, footerHint)

	// Input field
	input := ti.View()

	// Autocomplete hint (if available)
	autocompleteHint := ""
	if autocompleteTag != "" {
		hintStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(false) // Brutalist: no italics
		autocompleteHint = hintStyle.Render(fmt.Sprintf("Tab: %s", autocompleteTag))
	}

	// Suggestions section - show both tags and dates
	suggestionsSection := ""
	currentInput := ti.Value()
	currentWord := getLastWord(currentInput)

	if currentWord != "" {
		var matches []string

		// Check if typing tag (starts with @)
		if strings.HasPrefix(currentWord, "@") {
			wordWithoutAt := strings.TrimPrefix(currentWord, "@")
			if wordWithoutAt != "" {
				// Show matching tags
				matches = filterMatchingTags(currentWord, availableTags)
			}
		} else {
			// Show matching date phrases
			matches = filterMatchingDates(currentWord)
		}

		if len(matches) > 0 {
			labelStyle := lipgloss.NewStyle().
				Foreground(subtleColor).
				Bold(true)
			matchesLabel := labelStyle.Render("Suggestions:")

			matchesStyle := lipgloss.NewStyle().
				Foreground(mutedColor)
			matchesList := matchesStyle.Render(strings.Join(matches, ", "))

			suggestionsSection = lipgloss.JoinVertical(
				lipgloss.Left,
				"",
				matchesLabel,
				matchesList,
			)
		}
	}

	// Build main content
	mainParts := []string{input}
	if autocompleteHint != "" {
		mainParts = append(mainParts, "", autocompleteHint)
	}
	if suggestionsSection != "" {
		mainParts = append(mainParts, suggestionsSection)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Left, mainParts...)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	mainLines := lipgloss.Height(mainContent)
	padding := contentHeight - mainLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	content := header + "\n" + mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
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

// filterMatchingTags returns tags that start with the given prefix (case-insensitive)
func filterMatchingTags(prefix string, availableTags []string) []string {
	if prefix == "" {
		return nil
	}

	// Normalize prefix (remove @ if present, lowercase)
	normalizedPrefix := strings.ToLower(strings.TrimPrefix(prefix, "@"))

	var matches []string
	for _, tag := range availableTags {
		normalizedTag := strings.ToLower(strings.TrimPrefix(tag, "@"))
		if strings.HasPrefix(normalizedTag, normalizedPrefix) {
			matches = append(matches, tag)
		}
	}

	return matches
}

// filterMatchingDates returns date phrases that start with the given prefix
func filterMatchingDates(prefix string) []string {
	if prefix == "" {
		return nil
	}

	normalizedPrefix := strings.ToLower(prefix)
	dateSuggestions := helpers.GetDateSuggestions()

	var matches []string
	for _, suggestion := range dateSuggestions {
		if strings.HasPrefix(suggestion, normalizedPrefix) {
			matches = append(matches, suggestion)
		}
	}

	return matches
}
