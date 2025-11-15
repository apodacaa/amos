package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryView renders a read-only view of an entry
func RenderEntryView(width, height int, theme Theme, entry models.Entry, allTodos []models.Todo, scrollOffset int, markedForDeletion map[string]string, statusMsg string, currentIndex int, totalCount int) string {
	// Title at top with date
	dateStr := entry.Timestamp.Format("2006-01-02")
	styledDate := StyleDate(dateStr, theme)
	styledTitle := HighlightTagsInText(entry.Title, theme)
	title := fmt.Sprintf("%s: %s", styledDate, styledTitle)

	// Todos section (if any)
	var todosSection string
	var todoLineCount int
	if len(entry.TodoIDs) > 0 {
		// Filter todos that belong to this entry
		entryTodos := helpers.FilterTodosByEntry(allTodos, entry.ID)

		if len(entryTodos) > 0 {
			// Count open todos
			openCount, totalCount := helpers.CountTodoStats(entryTodos)

			todosTitle := GetNormalItemStyle().Render(fmt.Sprintf("Todos (%d open, %d total)", openCount, totalCount))

			// Render each todo
			var todoLines []string
			for _, todo := range entryTodos {
				checkbox := "[ ]"
				if todo.Status == "done" {
					checkbox = "[x]"
				} else if todo.Status == "next" {
					checkbox = "[>]"
				}

				// Style checkbox based on status
				styledCheckbox := StyleTodoStatus(todo.Status, checkbox, theme)

				// Apply tag highlighting to title
				highlightedTitle := HighlightTagsInText(todo.Title, theme)
				todoLine := fmt.Sprintf("%s %s", styledCheckbox, highlightedTitle)

				// Add tags if present (with bold highlighting)
				if len(todo.Tags) > 0 {
					var tagParts []string
					for _, tag := range todo.Tags {
						tagParts = append(tagParts, "@"+tag)
					}
					tagStr := " " + strings.Join(tagParts, " ")
					todoLine += HighlightTagsInText(tagStr, theme)
				}

				// Dim completed todos
				if todo.Status == "done" {
					todoLine = GetDimmedStyle(theme).Render(todoLine)
				} else {
					todoLine = GetNormalItemStyle().Render(todoLine)
				}

				todoLines = append(todoLines, "  "+todoLine)
			}

			todosContent := strings.Join(todoLines, "\n")
			todosSection = "\n\n" + todosTitle + "\n" + todosContent
			todoLineCount = 2 + len(todoLines) // Title + blank + todo lines
		}
	}

	// Calculate available height for content
	availableHeight := height - 3 - todoLineCount // header + footer + message line + todos
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Wrap body text to fit width, preserving original line breaks
	wrappedBody := wrapBodyText(entry.Body, width-8)

	// Split wrapped body into lines
	bodyLines := strings.Split(wrappedBody, "\n")
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
		body = HighlightTagsInText(strings.Join(visibleLines, "\n"), theme)
	} else {
		body = HighlightTagsInText(wrappedBody, theme)
		scrollStart = 0
		scrollEnd = totalLines
	}

	// Header
	header := RenderHeader(width, theme, "n", "new", "a", "todo", "j/k", "nav", "d", "del", "e", "entries", "t", "todos", "s", "theme", "q", "quit")

	// Footer: entry position + marked indicator + scroll info
	footerTitle := fmt.Sprintf("Entry %d of %d", currentIndex+1, totalCount)

	// Add marked indicator if entry is marked for deletion
	if _, isMarked := markedForDeletion[entry.ID]; isMarked {
		footerTitle += " [MARKED FOR DELETION]"
	}

	// Footer stats (statusMsg now in message line, not footer)
	footerStats := ""
	if totalLines > availableHeight {
		footerStats = fmt.Sprintf("lines %d-%d of %d", scrollStart+1, scrollEnd, totalLines)
	}

	footer := RenderFooter(width, theme, footerTitle, footerStats)

	// Build main content
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		body,
		todosSection,
	)

	// Calculate padding for content area
	contentHeight := height - 3 // header + footer + message line
	mainLines := strings.Count(mainContent, "\n") + 1
	padding := contentHeight - mainLines
	if padding < 0 {
		padding = 0
	}

	// Build message line (neomutt-style)
	messageLine := RenderMessageLine(width, statusMsg)

	// Build full view
	content := header + "\n" + mainContent
	if padding > 0 {
		content += strings.Repeat("\n", padding)
	}
	content += "\n" + footer + "\n" + messageLine

	return content
}
