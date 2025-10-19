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

		// Render each todo
		for i, todo := range sorted {
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
	) // Combine sections - help anchored to bottom
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
		status,
	)

	// Place content with help anchored to bottom
	fullContent := lipgloss.Place(
		width-4, height-4,
		lipgloss.Left, lipgloss.Top,
		mainContent,
	) + "\n" + help

	return container.Render(fullContent)
}
