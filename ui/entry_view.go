package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
)

// RenderEntryView renders a read-only view of an entry
func RenderEntryView(width, height int, theme Theme, entry models.Entry, allTodos []models.Todo, scrollOffset int, markedForDeletion map[string]string, statusMsg string, currentIndex int, totalCount int, updateAvailable bool, updateDismissed bool, latestVersion string) string {
	// Build timestamp line (single UpdatedAt, no label)
	timestampStr := entry.UpdatedAt.Format("2006-01-02")
	timestampLines := StyleDate(timestampStr, theme)

	// Title on separate line (wrap to fit width)
	wrappedTitle := WrapBodyText(entry.Title, width-8)
	styledTitle := HighlightTagsInText(wrappedTitle, theme)
	titleLineCount := strings.Count(wrappedTitle, "\n") + 1

	// Todos section (if any)
	var todosSection string
	// Filter todos that belong to this entry (single source of truth: Todo.EntryID)
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
			//if len(todo.Tags) > 0 {
			//	var tagParts []string
			//	for _, tag := range todo.Tags {
			//		tagParts = append(tagParts, "@"+tag)
			//	}
			//	tagStr := " " + strings.Join(tagParts, " ")
			//	todoLine += HighlightTagsInText(tagStr, theme)
			//}

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
	}

	// Calculate available height for scrollable body content
	// Reserve space for all fixed elements:
	// - header (1)
	// - timestamp (1)
	// - blank line after timestamp (1)
	// - title (titleLineCount)
	// - blank line after title (1)
	// - footer (1)
	// - message line (1)
	// - update notice (1 if shown)
	baseReserved := 6 + titleLineCount // 6 fixed lines + title height
	if updateAvailable && !updateDismissed {
		baseReserved += 1 // update notice line
	}

	// Calculate todos section height (if present)
	todosHeight := 0
	if todosSection != "" {
		todosHeight = strings.Count(todosSection, "\n") + 1
	}

	availableHeight := height - baseReserved - todosHeight
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Wrap body text to fit width, preserving original line breaks
	wrappedBody := WrapBodyText(entry.Body, width-8)

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
	header := RenderHeader(width, theme, "n", "new", "a", "todo", "i", "edit", "j/k", "nav", "f/b", "page", "d", "del", "e", "entries", "t", "todos", "s", "theme", "?", "help", "q", "quit")

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

	// Build main content (simple string concatenation to avoid lipgloss adding extra formatting)
	var contentParts []string
	contentParts = append(contentParts, timestampLines) // timestamps first
	contentParts = append(contentParts, "")             // blank line
	contentParts = append(contentParts, styledTitle)    // title on new line
	contentParts = append(contentParts, "")             // blank line
	contentParts = append(contentParts, body)
	if todosSection != "" {
		contentParts = append(contentParts, todosSection)
	}

	mainContent := strings.Join(contentParts, "\n")

	// Generate update notice if available
	updateNotice := RenderUpdateNotice(width, theme, updateAvailable, updateDismissed, latestVersion)

	// Calculate padding to fill remaining vertical space
	contentHeight := strings.Count(mainContent, "\n") + 1
	reservedLines := 3 // header + footer + message line
	if updateNotice != "" {
		reservedLines += 1 // update notice line
	}
	availableContentHeight := height - reservedLines
	paddingNeeded := availableContentHeight - contentHeight
	if paddingNeeded > 0 {
		mainContent += strings.Repeat("\n", paddingNeeded)
	}

	// Build message line (neomutt-style)
	messageLine := RenderMessageLine(width, statusMsg)

	// Assemble: header + mainContent + footer + update notice + message line
	result := header + "\n" + mainContent + "\n" + footer
	if updateNotice != "" {
		result += "\n" + updateNotice
	}
	result += "\n" + messageLine

	return result
}
