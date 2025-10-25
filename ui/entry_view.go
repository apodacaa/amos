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
	// Title at top
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor)
	title := titleStyle.Render(entry.Title)

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
	availableHeight := height - 2 - todoLineCount // header + footer + todos
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Split body into lines
	bodyLines := strings.Split(entry.Body, "\n")
	totalLines := len(bodyLines)

	// Apply scroll offset
	var body string
	var scrollStart, scrollEnd int

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
		scrollStart = scrollOffset
		scrollEnd = scrollOffset + availableHeight
		if scrollEnd > totalLines {
			scrollEnd = totalLines
		}

		visibleLines := bodyLines[scrollStart:scrollEnd]
		body = bodyStyle.Render(strings.Join(visibleLines, "\n"))
	} else {
		body = bodyStyle.Render(entry.Body)
		scrollStart = 0
		scrollEnd = totalLines
	}

	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "u/i", "scroll", "e", "entries", "t", "todos", "esc", "cancel", "q", "quit")

	// Footer: date (no time) + tags + scroll info
	footerTitle := entry.Timestamp.Format("2006-01-02")
	if len(entry.Tags) > 0 {
		// Add @ prefix to tags for clarity
		var tagStrings []string
		for _, tag := range entry.Tags {
			tagStrings = append(tagStrings, "@"+tag)
		}
		footerTitle += " " + strings.Join(tagStrings, " ")
	}

	footerStats := ""
	if totalLines > availableHeight {
		footerStats = fmt.Sprintf("lines %d-%d of %d", scrollStart+1, scrollEnd, totalLines)
	}

	footer := RenderFooter(width, footerTitle, footerStats)

	// Build main content
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		body,
		todosSection,
	)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	mainLines := strings.Count(mainContent, "\n") + 1
	padding := contentHeight - mainLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	content := header + "\n" + mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer

	return content
}
