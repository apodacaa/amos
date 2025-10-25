package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, todos []models.Todo, filterTags []string, filterDate string) string {
	// Apply filters: first date, then tags
	filtered := helpers.FilterEntriesByDateRange(entries, filterDate)
	filtered = helpers.FilterEntriesByTags(filtered, filterTags)

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
		availableHeight := height - 2 // header + footer
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
			// Table format: date  title (padded)  tags
			timestamp := entry.Timestamp.Format("2006-01-02")

			// Pad title to fixed width for column alignment
			titleWidth := 40
			paddedTitle := entry.Title
			if len(paddedTitle) > titleWidth {
				paddedTitle = paddedTitle[:titleWidth]
			} else {
				paddedTitle = paddedTitle + strings.Repeat(" ", titleWidth-len(paddedTitle))
			}

			line := fmt.Sprintf("%s  %s", timestamp, paddedTitle)

			// Add tags if present (aligned after padded title)
			if len(entry.Tags) > 0 {
				tagStr := ""
				for _, tag := range entry.Tags {
					tagStr += " @" + tag
				}
				line += tagStr
			}

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Style selected item with inverted colors (brutalist full-width bar)
			var styled string
			if i == selectedIdx {
				selectedStyle := lipgloss.NewStyle().
					Foreground(subtleColor).
					Reverse(true).
					Width(width - 4)
				styled = selectedStyle.Render(line)
			} else {
				normalStyle := lipgloss.NewStyle().Foreground(subtleColor)
				styled = normalStyle.Render(line)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	hasFilters := len(filterTags) > 0 || filterDate != ""
	var header string
	if hasFilters {
		header = RenderHeader(width, "n", "new", "a", "todo", "j/k", "nav", "enter", "view", "/", "clear", "t", "todos", "esc", "cancel", "q", "quit")
	} else {
		header = RenderHeader(width, "n", "new", "a", "todo", "j/k", "nav", "enter", "view", "/", "filter", "t", "todos", "esc", "cancel", "q", "quit")
	}

	// Footer
	footerTitle := "Entries"
	if len(filterTags) > 0 {
		footerTitle += " " + strings.Join(filterTags, " ")
	}
	if filterDate != "" {
		dateLabel := helpers.FormatDatePreset(filterDate)
		if dateLabel != "" {
			footerTitle += " " + dateLabel
		}
	}

	// Build stats with scroll info if needed
	var stats string
	if len(sorted) > 0 {
		// Calculate viewport info
		availableHeight := height - 2
		if availableHeight < 5 {
			availableHeight = 5
		}

		if len(sorted) > availableHeight {
			// Showing windowed view - calculate same viewport as rendering
			half := availableHeight / 2
			start := selectedIdx - half
			end := selectedIdx + half + 1

			if start < 0 {
				start = 0
				end = availableHeight
			}

			if end > len(sorted) {
				end = len(sorted)
				start = end - availableHeight
				if start < 0 {
					start = 0
				}
			}

			stats = fmt.Sprintf("%d-%d of %d items", start+1, end, len(sorted))
		} else {
			stats = fmt.Sprintf("%d items", len(sorted))
		}
	}

	footer := RenderFooter(width, footerTitle, stats)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	listLines := strings.Count(list, "\n") + 1
	padding := contentHeight - listLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	content := header + "\n" + list
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
}
