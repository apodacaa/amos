package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, theme Theme, entries []models.Entry, selectedIdx int, todos []models.Todo, filterTags []string, filterDate string, markedForDeletion map[string]string, statusMsg string) string {
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
			// Table format: markers  date  title
			timestamp := entry.Timestamp.Format("2006-01-02")

			// Build 2-character marker prefix (D for deletion, + for todos)
			var markerText string
			_, isMarked := markedForDeletion[entry.ID]
			hasTodos := len(entry.TodoIDs) > 0

			if isMarked && hasTodos {
				markerText = "D+"
			} else if isMarked && !hasTodos {
				markerText = "D "
			} else if !isMarked && hasTodos {
				markerText = " +"
			} else {
				markerText = "  "
			}

			// Truncate title if too long
			titleText := entry.Title
			plainLen := len(markerText) + 1 + len(timestamp) + 2 + len(titleText)
			maxLen := width - 6
			if plainLen > maxLen {
				titleText = titleText[:len(titleText)-(plainLen-maxLen)-3] + "..."
			}

			// Build line - style only for non-selected items
			var styled string
			if i == selectedIdx {
				// Selected: plain text, let selection style handle it
				plainLine := fmt.Sprintf("%s %s  %s", markerText, timestamp, titleText)
				styled = GetSelectedItemStyle(width, theme).Render(plainLine)
			} else {
				// Not selected: apply color styling
				var styledMarker string
				if isMarked && hasTodos {
					styledMarker = StyleMarker("D", true, theme) + StyleMarker("+", false, theme)
				} else if isMarked && !hasTodos {
					styledMarker = StyleMarker("D", true, theme) + " "
				} else if !isMarked && hasTodos {
					styledMarker = " " + StyleMarker("+", false, theme)
				} else {
					styledMarker = "  "
				}
				styledDate := StyleDate(timestamp, theme)
				styledTitle := HighlightTagsInText(titleText, theme)
				line := fmt.Sprintf("%s %s  %s", styledMarker, styledDate, styledTitle)
				styled = GetNormalItemStyle().Width(width).Render(line)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	header := RenderHeader(width, theme, "n", "new", "a", "todo", "j/k", "nav", "enter", "view", "d", "del", "/", "filter", "t", "todos", "s", "theme", "q", "quit")

	// Footer (show only filter context, no view label)
	var footerTitle string
	if len(filterTags) > 0 {
		footerTitle = strings.Join(filterTags, " ")
	}
	if filterDate != "" {
		if footerTitle != "" {
			footerTitle += " "
		}
		footerTitle += filterDate
	}
	// Wrap filters in brackets if present
	if footerTitle != "" {
		footerTitle = "[" + footerTitle + "]"
	}

	// Build stats showing selected entry position
	var stats string
	if len(sorted) > 0 {
		stats = fmt.Sprintf("Entry %d of %d", selectedIdx+1, len(sorted))
	}

	footer := RenderFooter(width, theme, stats, footerTitle)

	return AssembleView(header, list, footer, width, height, statusMsg)
}
