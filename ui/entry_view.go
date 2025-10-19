package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryView renders a read-only view of an entry
func RenderEntryView(width, height int, entry models.Entry, allTodos []models.Todo) string {
	container := GetFullScreenBox(width, height)

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Width(width - 8)
	title := titleStyle.Render(entry.Title)

	// Metadata line: date and tags
	metaStyle := lipgloss.NewStyle().
		Foreground(mutedColor)

	timestamp := entry.Timestamp.Format("2006-01-02 15:04")
	meta := timestamp

	if len(entry.Tags) > 0 {
		meta += " " + strings.Join(entry.Tags, " ")
	}

	metadata := metaStyle.Render(meta)

	// Body
	bodyStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Width(width - 8)
	body := bodyStyle.Render(entry.Body)

	// Todos section (if any)
	var todosSection string
	if len(entry.TodoIDs) > 0 {
		// Filter todos that belong to this entry
		entryTodos := helpers.FilterTodosByEntry(allTodos, entry.ID)

		if len(entryTodos) > 0 {
			// Count open todos
			openCount, totalCount := helpers.CountTodoStats(entryTodos)

			todosTitle := lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				Render(fmt.Sprintf("Todos (%d open, %d total)", openCount, totalCount))

			// Render each todo
			var todoLines []string
			for _, todo := range entryTodos {
				checkbox := "[ ]"
				if todo.Status == "done" {
					checkbox = "[x]"
				}

				todoLine := fmt.Sprintf("%s %s", checkbox, todo.Title)

				// Add tags if present
				if len(todo.Tags) > 0 {
					tagStr := ""
					for _, tag := range todo.Tags {
						tagStr += " @" + tag
					}
					todoLine += lipgloss.NewStyle().
						Foreground(mutedColor).
						Render(tagStr)
				}

				// Dim completed todos
				if todo.Status == "done" {
					todoLine = lipgloss.NewStyle().
						Foreground(mutedColor).
						Render(todoLine)
				} else {
					todoLine = lipgloss.NewStyle().
						Foreground(subtleColor).
						Render(todoLine)
				}

				todoLines = append(todoLines, "  "+todoLine)
			}

			todosContent := strings.Join(todoLines, "\n")
			todosSection = "\n\n" + todosTitle + "\n" + todosContent
		}
	}

	// Help text at bottom
	help := FormatHelpLeft(width,
		"n", "new entry",
		"a", "add todo",
		"j/k", "navigate",
		"t", "todos",
		"e", "entries",
		"d", "dashboard",
		"q", "quit",
	)

	// Combine main content - help anchored to bottom
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		metadata,
		"",
		body,
		todosSection,
	)

	// Place content with help anchored to bottom
	fullContent := lipgloss.Place(
		width-4, height-4,
		lipgloss.Left, lipgloss.Top,
		mainContent,
	) + "\n" + help

	return container.Render(fullContent)
}
