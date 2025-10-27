package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, entries []models.Entry, selectedIdx int, filterTags []string, filterDate string) string {
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

			// Table format: checkbox  date  title (padded)  tags
			dateStr := todo.CreatedAt.Format("2006-01-02")

			// Pad title to fixed width for column alignment
			titleWidth := 35
			paddedTitle := todo.Title
			if len(paddedTitle) > titleWidth {
				paddedTitle = paddedTitle[:titleWidth]
			} else {
				paddedTitle = paddedTitle + strings.Repeat(" ", titleWidth-len(paddedTitle))
			}

			line := fmt.Sprintf("%s %s  %s", checkbox, dateStr, paddedTitle)

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Apply selection and completion styling with inverted colors (brutalist full-width bar)
			var styled string
			if i == selectedIdx {
				// Selected items with inverted colors - full width bar
				if todo.Status == "done" {
					styled = GetSelectedDoneStyle(width).Render(line)
				} else {
					styled = GetSelectedItemStyle(width).Render(line)
				}
			} else {
				// Dim completed todos, normal color for open
				if todo.Status == "done" {
					styled = GetDimmedStyle().Render(line)
				} else {
					styled = GetNormalItemStyle().Render(line)
				}
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "enter", "view", "j/k", "nav", "space", "cycle", "/", "filter", "c", "clear", "e", "entries", "q", "quit")

	// Footer
	footerTitle := "Todos"
	if len(filterTags) > 0 {
		footerTitle += " " + strings.Join(filterTags, " ")
	}
	if filterDate != "" {
		footerTitle += " " + filterDate
	}

	// Stats for footer
	openCount := 0
	nextCount := 0
	doneCount := 0
	for _, todo := range filtered {
		switch todo.Status {
		case "open":
			openCount++
		case "next":
			nextCount++
		case "done":
			doneCount++
		}
	}

	// Build stats with scroll info if needed
	var stats string
	if len(filtered) > 0 {
		start, end := CalculateViewport(len(filtered), selectedIdx, height)
		if end-start < len(filtered) {
			stats = fmt.Sprintf("%d-%d of %d | %d open, %d next, %d done", start+1, end, len(filtered), openCount, nextCount, doneCount)
		} else {
			stats = fmt.Sprintf("%d open, %d next, %d done", openCount, nextCount, doneCount)
		}
	}

	footer := RenderFooter(width, footerTitle, stats)

	return AssembleView(header, list, footer, height)
}
