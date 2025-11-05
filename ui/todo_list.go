package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, entries []models.Entry, selectedIdx int, filterTags []string, filterDate string, markedForDeletion map[string]string, statusMsg string) string {
	// Apply filters: first date, then tags
	filtered := helpers.FilterTodosByDateRange(todos, filterDate)
	filtered = helpers.FilterTodosByTags(filtered, filterTags)

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
			var marker string
			_, isMarked := markedForDeletion[todo.ID]
			hasEntry := todo.EntryID != nil

			if isMarked && hasEntry {
				marker = "D+"
			} else if isMarked && !hasEntry {
				marker = "D "
			} else if !isMarked && hasEntry {
				marker = " +"
			} else {
				marker = "  "
			}

			line := fmt.Sprintf("%s %s %s  %s", marker, checkbox, dateStr, todo.Title)

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Apply selection styling (no tag highlighting in lists)
			var styled string
			if i == selectedIdx {
				// Selected items
				if todo.Status == "done" {
					styled = GetSelectedDoneStyle(width).Render(line)
				} else {
					styled = GetSelectedItemStyle(width).Render(line)
				}
			} else {
				// Dim completed todos, normal color for open
				if todo.Status == "done" {
					styled = GetDimmedStyle().Width(width).Render(line)
				} else {
					styled = GetNormalItemStyle().Width(width).Render(line)
				}
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "enter", "view", "j/k", "nav", "space", "cycle", "d", "del", "/", "filter", "e", "entries", "q", "quit")

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

	footer := RenderFooter(width, stats, footerTitle)

	return AssembleView(header, list, footer, width, height, statusMsg)
}
