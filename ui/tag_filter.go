package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// RenderTagFilter renders the tag filter input view
func RenderTagFilter(width, height int, ti textarea.Model, availableTags []string, autocompleteTag string) string {
	box := GetFullScreenBox(width, height)
	titleStyle := GetTitleStyle(width)

	title := titleStyle.Render("FILTER BY TAGS")

	// Input field
	input := ti.View()

	// Autocomplete hint (if available)
	autocompleteHint := ""
	if autocompleteTag != "" {
		hintStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(false) // Brutalist: no italics
		autocompleteHint = hintStyle.Render(fmt.Sprintf("Tab to complete: %s", autocompleteTag))
	}

	// Filtered matches section (brutalist: only show when typing)
	matchesSection := ""
	currentInput := ti.Value()
	currentWord := getLastWord(currentInput)

	// Only show matches if user is typing a word with content beyond @
	wordWithoutAt := strings.TrimPrefix(currentWord, "@")
	if currentWord != "" && wordWithoutAt != "" && wordWithoutAt != currentWord {
		// Filter tags that match the current word
		matches := filterMatchingTags(currentWord, availableTags)

		if len(matches) > 0 {
			labelStyle := lipgloss.NewStyle().
				Foreground(subtleColor).
				Bold(true)
			matchesLabel := labelStyle.Render("Matches:")

			tagsStyle := lipgloss.NewStyle().
				Foreground(mutedColor)
			matchesList := tagsStyle.Render(strings.Join(matches, " "))

			matchesSection = lipgloss.JoinVertical(
				lipgloss.Left,
				"",
				matchesLabel,
				matchesList,
			)
		}
	}

	// Help text
	help := FormatHelpLeft(width, "enter", "apply", "esc", "cancel")

	// Build main content (everything except help)
	mainParts := []string{title, "", input}
	if autocompleteHint != "" {
		mainParts = append(mainParts, "", autocompleteHint)
	}
	if matchesSection != "" {
		mainParts = append(mainParts, matchesSection)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Left, mainParts...)

	// Calculate how much vertical space to add to push help to bottom
	mainLines := lipgloss.Height(mainContent)
	helpLines := 1
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	// Add padding and help
	content := mainContent
	if padding > 0 {
		content += "\n" + lipgloss.NewStyle().Height(padding).Render("")
	}
	content += "\n" + help

	return box.Render(content)
}

// getLastWord extracts the last word from the input (after the last space)
func getLastWord(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// Find last space
	lastSpaceIdx := strings.LastIndex(input, " ")
	if lastSpaceIdx == -1 {
		// No spaces, entire input is the word
		return input
	}

	// Return everything after the last space
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
