package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, entries []models.Entry, selectedIdx int, statusMsg string) string {
	container := GetFullScreenBox(width, height)
	title := GetTitleStyle(width).Render("Todos")

	// Build todo list
	var listItems []string

	if len(todos) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(width - 4).
			Align(lipgloss.Center)
		listItems = append(listItems, emptyStyle.Render("No todos yet. Create an entry with !todo lines."))
	} else {
		// Sort todos using helper (same logic as commands)
		sorted := helpers.SortTodosForDisplay(todos)

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
			// Checkbox (simple - status is always current)
			checkbox := "[ ]"
			if todo.Status == "done" {
				checkbox = "[x]"
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

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor)
		status = "\n" + statusStyle.Render(statusMsg) + "\n"
	}

	// Help text at bottom
	help := FormatHelpLeft(width,
		"n", "new entry",
		"a", "add todo",
		"j/k", "navigate",
		"u/i", "move",
		"space", "toggle",
		"e", "entries",
		"esc", "cancel",
		"q", "quit",
	)

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
		status,
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
