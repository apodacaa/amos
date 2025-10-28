package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
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
		listItems = append(listItems, GetEmptyStateStyle(width).Render("No entries yet"))
	} else {
		// Calculate viewport (visible window of items)
		start, end := CalculateViewport(len(sorted), selectedIdx, height)

		// Render visible items
		for i := start; i < end; i++ {
			entry := sorted[i]
			// Table format: date  title (padded)  tags
			timestamp := entry.Timestamp.Format("2006-01-02")

			// Pad title to fixed width for column alignment
			titleWidth := 30
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

			// Add selection marker
			var prefix string
			if i == selectedIdx {
				prefix = "> "
			} else {
				prefix = "  " // Indent for alignment
			}
			line = prefix + line

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Style selected item
			var styled string
			if i == selectedIdx {
				styled = GetSelectedItemStyle(width).Render(line)
			} else {
				styled = GetNormalItemStyle().Render(line)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "j/k", "nav", "enter", "view", "/", "filter", "c", "clear", "t", "todos", "q", "quit")

	// Footer
	footerTitle := "Entries"
	if len(filterTags) > 0 {
		footerTitle += " " + strings.Join(filterTags, " ")
	}
	if filterDate != "" {
		footerTitle += " " + filterDate
	}

	// Build stats with scroll info if needed
	var stats string
	if len(sorted) > 0 {
		start, end := CalculateViewport(len(sorted), selectedIdx, height)
		if end-start < len(sorted) {
			stats = fmt.Sprintf("%d-%d of %d items", start+1, end, len(sorted))
		} else {
			stats = fmt.Sprintf("%d items", len(sorted))
		}
	}

	footer := RenderFooter(width, footerTitle, stats)

	return AssembleView(header, list, footer, height)
}
