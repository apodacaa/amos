package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, entries []models.Entry, selectedIdx int, filterTags []string, filterDate string) string {
	// Apply filters: first date, then tags
	filtered := helpers.FilterTodosByDateRange(todos, filterDate)
	filtered = helpers.FilterTodosByTags(filtered, filterTags)

	// Build todo list
	var listItems []string

	if len(filtered) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(width - 4).
			Align(lipgloss.Center)
		if len(filterTags) > 0 {
			listItems = append(listItems, emptyStyle.Render("No todos match the filter."))
		} else {
			listItems = append(listItems, emptyStyle.Render("No todos yet. Create an entry with !todo lines."))
		}
	} else {
		// Todos are pre-sorted (displayTodos from model) and filtered
		sorted := filtered

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

			// Add tags if present (aligned after padded title)
			if len(todo.Tags) > 0 {
				tagStr := ""
				for _, tag := range todo.Tags {
					tagStr += " @" + tag
				}
				line += tagStr
			}

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
					selectedStyle := lipgloss.NewStyle().
						Foreground(mutedColor).
						Reverse(true).
						Width(width - 4)
					styled = selectedStyle.Render(line)
				} else {
					selectedStyle := lipgloss.NewStyle().
						Foreground(subtleColor).
						Reverse(true).
						Width(width - 4)
					styled = selectedStyle.Render(line)
				}
			} else {
				// Dim completed todos, normal color for open
				if todo.Status == "done" {
					dimStyle := lipgloss.NewStyle().Foreground(mutedColor)
					styled = dimStyle.Render(line)
				} else {
					normalStyle := lipgloss.NewStyle().Foreground(subtleColor)
					styled = normalStyle.Render(line)
				}
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Header
	hasFilters := len(filterTags) > 0 || filterDate != ""
	var header string
	if hasFilters {
		header = RenderHeader(width, "n", "new", "a", "todo", "j/k", "nav", "space", "cycle", "/", "clear", "e", "entries", "esc", "cancel", "q", "quit")
	} else {
		header = RenderHeader(width, "n", "new", "a", "todo", "j/k", "nav", "space", "cycle", "/", "filter", "e", "entries", "esc", "cancel", "q", "quit")
	}

	// Footer
	footerTitle := "Todos"
	if len(filterTags) > 0 {
		footerTitle += " " + strings.Join(filterTags, " ")
	}
	if filterDate != "" {
		dateLabel := helpers.FormatDatePreset(filterDate)
		if dateLabel != "" {
			footerTitle += " " + dateLabel
		}
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
		// Calculate viewport info
		availableHeight := height - 2
		if availableHeight < 5 {
			availableHeight = 5
		}

		if len(filtered) > availableHeight {
			// Showing windowed view - calculate same viewport as rendering
			half := availableHeight / 2
			start := selectedIdx - half
			end := selectedIdx + half + 1

			if start < 0 {
				start = 0
				end = availableHeight
			}

			if end > len(filtered) {
				end = len(filtered)
				start = end - availableHeight
				if start < 0 {
					start = 0
				}
			}

			stats = fmt.Sprintf("%d-%d of %d | %d open, %d next, %d done", start+1, end, len(filtered), openCount, nextCount, doneCount)
		} else {
			stats = fmt.Sprintf("%d open, %d next, %d done", openCount, nextCount, doneCount)
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
