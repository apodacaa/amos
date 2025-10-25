package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, entries []models.Entry, selectedIdx int, filterTags []string) string {
	container := GetFullScreenBox(width, height)

	// Apply filter if active
	filtered := helpers.FilterTodosByTags(todos, filterTags)

	// Update title to show filter status
	titleText := "Todos"
	if len(filterTags) > 0 {
		titleText += " " + strings.Join(filterTags, " ")
	}
	title := GetTitleStyle(width).Render(titleText)

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

			// Title with tags
			titleText := todo.Title
			if len(todo.Tags) > 0 {
				tagStyle := lipgloss.NewStyle().Foreground(mutedColor)
				tagStr := ""
				for _, tag := range todo.Tags {
					tagStr += " @" + tag
				}
				titleText += tagStyle.Render(tagStr)
			}

			line := fmt.Sprintf("%s %s", checkbox, titleText)

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Apply selection and completion styling
			var styled string
			if i == selectedIdx {
				// Selected items use subtle color (unless done)
				if todo.Status == "done" {
					selectedStyle := lipgloss.NewStyle().Foreground(mutedColor)
					styled = selectedStyle.Render("> " + line)
				} else {
					selectedStyle := lipgloss.NewStyle().Foreground(subtleColor)
					styled = selectedStyle.Render("> " + line)
				}
			} else {
				// Dim completed todos, normal color for open
				if todo.Status == "done" {
					dimStyle := lipgloss.NewStyle().Foreground(mutedColor)
					styled = dimStyle.Render("  " + line)
				} else {
					normalStyle := lipgloss.NewStyle().Foreground(subtleColor)
					styled = normalStyle.Render("  " + line)
				}
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

	// Help text at bottom
	var help string
	if len(filterTags) > 0 {
		help = FormatHelpLeft(width,
			"n", "new entry",
			"a", "add todo",
			"j/k", "navigate",
			"space", "cycle",
			"r", "refresh",
			"e", "entries",
			"@", "clear filter",
			"esc", "cancel",
			"q", "quit",
		)
	} else {
		help = FormatHelpLeft(width,
			"n", "new entry",
			"a", "add todo",
			"j/k", "navigate",
			"space", "cycle",
			"r", "refresh",
			"e", "entries",
			"@", "filter",
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
