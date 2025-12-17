package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, theme Theme, todos []models.Todo, entries []models.Entry, selectedIdx int, filterTags []string, filterDate string, markedForDeletion map[string]string, statusMsg string) string {
	// Apply filters: first date, then tags
	filtered := helpers.ApplyTodoFilters(todos, filterDate, filterTags)

	// Build todo list
	var listItems []string

	if len(filtered) == 0 {
		if len(filterTags) > 0 {
			listItems = append(listItems, GetEmptyStateStyle(width).Render("No todos match the filter."))
		} else {
			listItems = append(listItems, GetEmptyStateStyle(width).Render("No todos yet. Create an entry with !todo lines."))
		}
	} else {
		// Todos are pre-sorted (displayTodos from model) and filtered
		sorted := filtered

		// Calculate viewport (visible window of items)
		start, end := CalculateViewport(len(sorted), selectedIdx, height)

		// Render visible todos
		for i := start; i < end; i++ {
			todo := sorted[i]
			// Checkbox based on status
			checkbox := "[ ]" // open
			if todo.Status == "next" {
				checkbox = "[>]" // next (brutalist arrow)
			} else if todo.Status == "done" {
				checkbox = "[x]" // done
			}

			// Table format: markers  checkbox  date  title
			dateStr := todo.CreatedAt.Format("2006-01-02")

			// Build 2-character marker prefix (D for deletion, + for linked entry)
			var markerText string
			_, isMarked := markedForDeletion[todo.ID]
			hasEntry := todo.EntryID != nil

			if isMarked && hasEntry {
				markerText = "D+"
			} else if isMarked && !hasEntry {
				markerText = "D "
			} else if !isMarked && hasEntry {
				markerText = " +"
			} else {
				markerText = "  "
			}

			// Truncate title if too long
			titleText := todo.Title
			plainLen := len(markerText) + 1 + len(checkbox) + 1 + len(dateStr) + 2 + len(titleText)
			maxLen := width - 6
			if plainLen > maxLen {
				titleText = titleText[:len(titleText)-(plainLen-maxLen)-3] + "..."
			}

			// Build line - style differently for selected vs non-selected
			var styled string
			if i == selectedIdx {
				// Selected: plain text on neon green background (no color styling)
				plainLine := fmt.Sprintf("%s %s %s  %s", markerText, checkbox, dateStr, titleText)
				if todo.Status == "done" {
					styled = GetSelectedDoneStyle(width, theme).Render(plainLine)
				} else {
					styled = GetSelectedItemStyle(width, theme).Render(plainLine)
				}
			} else {
				// Not selected: check status first
				if todo.Status == "done" {
					// Done: plain text, dimmed
					plainLine := fmt.Sprintf("%s %s %s  %s", markerText, checkbox, dateStr, titleText)
					styled = GetDimmedStyle(theme).Width(width).Render(plainLine)
				} else {
					// Open/next: apply full color styling
					var styledMarker string
					if isMarked && hasEntry {
						styledMarker = StyleMarker("D", true, theme) + StyleMarker("+", false, theme)
					} else if isMarked && !hasEntry {
						styledMarker = StyleMarker("D", true, theme) + " "
					} else if !isMarked && hasEntry {
						styledMarker = " " + StyleMarker("+", false, theme)
					} else {
						styledMarker = "  "
					}
					styledCheckbox := StyleTodoStatus(todo.Status, checkbox, theme)
					styledDate := StyleDate(dateStr, theme)
					styledTitle := HighlightTagsInText(titleText, theme)
					line := fmt.Sprintf("%s %s %s  %s", styledMarker, styledCheckbox, styledDate, styledTitle)
					styled = GetNormalItemStyle().Width(width).Render(line)
				}
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	header := RenderHeader(width, theme, "n", "new", "a", "todo", "enter", "view", "j/k", "nav", "space", "cycle", "d", "del", "/", "filter", "e", "entries", "?", "help", "q", "quit")

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

	// Build stats showing selected todo position
	var stats string
	if len(filtered) > 0 {
		stats = fmt.Sprintf("Todo %d of %d", selectedIdx+1, len(filtered))
	}

	footer := RenderFooter(width, theme, stats, footerTitle)

	return AssembleView(header, list, footer, width, height, statusMsg)
}
