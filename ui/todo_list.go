package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderTodoList renders the todo list view
func RenderTodoList(width, height int, todos []models.Todo, selectedIdx int, statusMsg string) string {
	container := GetContainerStyle(width, height)
	title := GetTitleStyle(width).Render("Todos")

	// Build todo list
	var listItems []string

	if len(todos) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Width(width - 4).
			Align(lipgloss.Center)
		listItems = append(listItems, emptyStyle.Render("No todos yet. Create an entry with !todo lines."))
	} else {
		// Sort todos: open first, then by created date descending
		sorted := make([]models.Todo, len(todos))
		copy(sorted, todos)

		// Simple bubble sort: open before done, then newest first within each group
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				// Open todos come before done todos
				if sorted[i].Status == "done" && sorted[j].Status == "open" {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].Status == sorted[j].Status {
					// Within same status, newest first
					if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
		}

		// Render each todo
		for i, todo := range sorted {
			// Checkbox
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

			// Entry link indicator
			if todo.EntryID != nil {
				linkStyle := lipgloss.NewStyle().Foreground(mutedColor)
				titleText += linkStyle.Render(" (from entry)")
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
				selectedStyle := lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true)
				styled = selectedStyle.Render("► " + line)
			} else {
				normalStyle := lipgloss.NewStyle().
					Foreground(subtleColor)
				styled = normalStyle.Render("  " + line)
			}

			// Dim completed todos
			if todo.Status == "done" {
				dimStyle := lipgloss.NewStyle().Foreground(mutedColor)
				styled = dimStyle.Render(styled)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)
		status = "\n" + statusStyle.Render(statusMsg) + "\n"
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)
	help := helpStyle.Render("j/k: navigate • space: toggle • esc: back • q: quit")

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
		status,
		help,
	)

	return container.Render(content)
}
