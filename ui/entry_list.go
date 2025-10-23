package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, todos []models.Todo, filterTags []string) string {
	container := GetFullScreenBox(width, height)

	// Update title to show filter if active
	titleText := "Entries"
	if len(filterTags) > 0 {
		titleText = fmt.Sprintf("Entries (filtered: %s)", strings.Join(filterTags, " "))
	}
	title := GetTitleStyle(width).Render(titleText)

	// Filter entries by tags if filter is active
	filtered := helpers.FilterEntriesByTags(entries, filterTags)

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
		// Calculate viewport (visible window of items)
		// Reserve space: title (1) + blank (1) + status (1) + blank (1) + help (1) = 5 lines minimum
		availableHeight := height - 10 // Conservative estimate for chrome
		if availableHeight < 5 {
			availableHeight = 5
		}

		// Calculate window start and end to keep selected item visible
		start := 0
		end := len(sorted)

		if len(sorted) > availableHeight {
			// Center selected item in viewport
			half := availableHeight / 2
			start = selectedIdx - half
			end = selectedIdx + half + 1

			// Adjust if near beginning
			if start < 0 {
				start = 0
				end = availableHeight
			}

			// Adjust if near end
			if end > len(sorted) {
				end = len(sorted)
				start = end - availableHeight
				if start < 0 {
					start = 0
				}
			}
		}

		// Render visible items
		for i := start; i < end; i++ {
			entry := sorted[i]
			// Format: 2006-01-02 15:04 | Meeting with team
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

		// Add scroll indicator if needed
		if len(sorted) > availableHeight {
			scrollInfo := fmt.Sprintf("(%d-%d of %d)", start+1, end, len(sorted))
			scrollStyle := lipgloss.NewStyle().Foreground(mutedColor)
			listItems = append(listItems, scrollStyle.Render(scrollInfo))
		}
	}

	list := strings.Join(listItems, "\n")

	// Help text (changes based on filter state) with bold keys
	var help string
	if len(filterTags) > 0 {
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

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
	)

	// Calculate how much vertical space to add to push help to bottom
	mainLines := strings.Count(mainContent, "\n") + 1
	helpLines := 1
	availableSpace := height - 4
	padding := availableSpace - mainLines - helpLines
	if padding < 0 {
		padding = 0
	}

	// Add padding and help
	content := mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + help

	return container.Render(content)
}
