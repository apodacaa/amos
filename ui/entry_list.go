package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, statusMsg string, todos []models.Todo, filterTag string) string {
	container := GetFullScreenBox(width, height)

	// Update title to show filter if active
	titleText := "Entries"
	if filterTag != "" {
		titleText = fmt.Sprintf("Entries (filtered: %s)", filterTag)
	}
	title := GetTitleStyle(width).Render(titleText)

	// Filter entries by tag if filter is active
	filtered := helpers.FilterEntriesByTag(entries, filterTag)

	// Sort entries by timestamp (newest first)
	sorted := helpers.SortEntriesForDisplay(filtered)

	// Build entry list
	var listItems []string

	if len(sorted) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(width - 4).
			Align(lipgloss.Center)
		listItems = append(listItems, emptyStyle.Render("No entries yet"))
	} else {
		for i, entry := range sorted {
			// Format: 2025-10-18 14:32 | Meeting with team
			timestamp := entry.Timestamp.Format("2006-01-02 15:04")
			line := fmt.Sprintf("%s | %s", timestamp, entry.Title)

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Style selected item differently
			var styled string
			if i == selectedIdx {
				selectedStyle := lipgloss.NewStyle().Foreground(subtleColor)
				styled = selectedStyle.Render("> " + line)
			} else {
				normalStyle := lipgloss.NewStyle().Foreground(subtleColor)
				styled = normalStyle.Render("  " + line)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor)
		status = "\n" + statusStyle.Render(statusMsg) + "\n"
	}

	// Help text (changes based on filter state) with bold keys
	var help string
	if filterTag != "" {
		help = FormatHelpLeft(width,
			"n", "new entry",
			"a", "add todo",
			"j/k", "navigate",
			"enter", "view",
			"@", "clear filter",
			"t", "todos",
			"esc", "cancel",
			"q", "quit",
		)
	} else {
		help = FormatHelpLeft(width,
			"n", "new entry",
			"a", "add todo",
			"j/k", "navigate",
			"enter", "view",
			"@", "filter",
			"t", "todos",
			"esc", "cancel",
			"q", "quit",
		)
	}

	// Combine sections - help anchored to bottom
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
		status,
	)

	// Place content with help anchored to bottom
	fullContent := lipgloss.Place(
		width-4, height-4,
		lipgloss.Left, lipgloss.Top,
		mainContent,
	) + "\n" + help

	return container.Render(fullContent)
}
