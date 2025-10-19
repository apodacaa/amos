package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryView renders a read-only view of an entry
func RenderEntryView(width, height int, entry models.Entry, allTodos []models.Todo, scrollOffset int) string {
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

	// Todos section (if any)
	var todosSection string
	var todoLineCount int
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
			todoLineCount = 2 + len(todoLines) // Title + blank + todo lines
		}
	}

	// Calculate available height for content
	// Reserve space for: title (1) + metadata (1) + blank (1) + help (2) + margins (4)
	availableHeight := height - 9 - todoLineCount
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Split body into lines
	bodyLines := strings.Split(entry.Body, "\n")
	totalLines := len(bodyLines)

	// Apply scroll offset
	var body string
	var scrollIndicator string

	if totalLines > availableHeight {
		// Clamp scrollOffset to valid range
		maxOffset := totalLines - availableHeight
		if scrollOffset > maxOffset {
			scrollOffset = maxOffset
		}
		if scrollOffset < 0 {
			scrollOffset = 0
		}

		// Show windowed content
		start := scrollOffset
		end := scrollOffset + availableHeight
		if end > totalLines {
			end = totalLines
		}

		visibleLines := bodyLines[start:end]
		body = bodyStyle.Render(strings.Join(visibleLines, "\n"))

		// Add scroll indicator
		continuationStyle := lipgloss.NewStyle().Foreground(mutedColor)
		scrollIndicator = "\n" + continuationStyle.Render(fmt.Sprintf("(showing lines %d-%d of %d)", start+1, end, totalLines))
	} else {
		body = bodyStyle.Render(entry.Body)
	}

	// Help text at bottom - always show scroll controls for consistency
	help := FormatHelpLeft(width,
		"n", "new entry",
		"a", "add todo",
		"j/k", "navigate",
		"u/i", "scroll",
		"e", "entries",
		"t", "todos",
		"esc", "cancel",
		"q", "quit",
	)

	// Build main content (everything except help)
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		metadata,
		"",
		body,
		scrollIndicator,
		todosSection,
	)

	// Calculate how much vertical space to add to push help to bottom
	// Count lines in main content
	mainLines := strings.Count(mainContent, "\n") + 1
	helpLines := 1               // Help is single line
	availableSpace := height - 4 // Account for container margins
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
