package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

// highlightTodoInText returns the body text as-is
// Previously tried to highlight the todo, but rendering caused indentation issues
func highlightTodoInText(body, todoTitle string) string {
	return body
}

// RenderTodoView renders a read-only view of a todo
func RenderTodoView(width, height int, theme Theme, todo models.Todo, allEntries []models.Entry, scrollOffset int, markedForDeletion map[string]string, statusMsg string, currentIndex int, totalCount int, updateAvailable bool, updateDismissed bool, latestVersion string) string {
	// Status and date at top
	statusIcon := "[ ]" // open
	if todo.Status == "next" {
		statusIcon = "[>]"
	} else if todo.Status == "done" {
		statusIcon = "[x]"
	}

	// Style status icon
	styledStatus := StyleTodoStatus(todo.Status, statusIcon, theme)

	// Build timestamp line (single UpdatedAt, no label, keep status icon)
	timestampStr := todo.UpdatedAt.Format("2006-01-02")
	timestampLines := fmt.Sprintf("%s %s", styledStatus, StyleDate(timestampStr, theme))

	// Todo title (wrappable)
	titleStyle := GetNormalItemStyle().
		Width(width - 8)

	// Wrap title if long
	wrappedTitle := wordWrap(todo.Title, width-8)
	// Apply tag highlighting to wrapped title
	highlightedTitle := HighlightTagsInText(wrappedTitle, theme)
	titleLines := strings.Split(highlightedTitle, "\n")

	// Linked entry section (if EntryID is set)
	var linkedSection string
	var linkedLineCount int
	if todo.EntryID != nil {
		// Find the entry
		var linkedEntry *models.Entry
		for _, entry := range allEntries {
			if entry.ID == *todo.EntryID {
				linkedEntry = &entry
				break
			}
		}

		if linkedEntry != nil {
			// Build full entry context
			var linkedParts []string

			// Header line: "From the following entry..."
			linkedParts = append(linkedParts, GetDimmedStyle(theme).Render("From entry..."))

			// Date, title, tags line
			dateStr := linkedEntry.CreatedAt.Format("2006-01-02")
			styledDate := StyleDate(dateStr, theme)
			styledTitle := HighlightTagsInText(linkedEntry.Title, theme)
			tagStr := ""

			// Not necessary to repeate all the tags that are clearly visible in the text
			//if len(linkedEntry.Tags) > 0 {
			//	var tagStrings []string
			//	for _, tag := range linkedEntry.Tags {
			//		tagStrings = append(tagStrings, "@"+tag)
			//	}
			//	tagStr = " " + HighlightTagsInText(strings.Join(tagStrings, " "), theme)
			//}

			metaLine := fmt.Sprintf("%s %s%s", styledDate, styledTitle, tagStr)
			linkedParts = append(linkedParts, metaLine)

			// Body (full text) - wrap, highlight tags and todo reference
			if linkedEntry.Body != "" {
				// Wrap body text to fit width
				wrappedBody := wrapBodyText(linkedEntry.Body, width-8)
				// First highlight tags, then highlight todo reference
				bodyWithTags := HighlightTagsInText(wrappedBody, theme)
				styledBody := highlightTodoInText(bodyWithTags, todo.Title)
				linkedParts = append(linkedParts, styledBody)
			}

			linkedSection = "\n\n" + strings.Join(linkedParts, "\n\n")

			// Count lines (header + meta + body lines + blank lines)
			linkedLineCount = 2  // blank lines before section (\n\n prefix)
			linkedLineCount += 1 // header line
			linkedLineCount += 1 // blank line after header (\n\n in join)
			linkedLineCount += 1 // meta line
			if linkedEntry.Body != "" {
				linkedLineCount += 1 // blank line after meta (\n\n in join)
				bodyLines := strings.Split(linkedEntry.Body, "\n")
				linkedLineCount += len(bodyLines)
			}
		}
	}

	// Calculate available height for title content
	availableHeight := height - 3 - linkedLineCount // header + footer + message line + linked section
	if availableHeight < 5 {
		availableHeight = 5
	}

	totalLines := len(titleLines)

	// Apply scroll offset
	var renderedTitle string
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

		visibleLines := titleLines[scrollStart:scrollEnd]
		renderedTitle = titleStyle.Render(strings.Join(visibleLines, "\n"))
	} else {
		renderedTitle = titleStyle.Render(highlightedTitle)
		scrollStart = 0
		scrollEnd = totalLines
	}

	// Header
	header := RenderHeader(width, theme, "n", "new", "a", "todo", "i", "edit", "space", "cycle", "j/k", "nav", "f/b", "page", "d", "del", "e", "entries", "t", "todos", "?", "help", "q", "quit")

	// Footer: position + marked indicator + scroll info
	footerTitle := fmt.Sprintf("Todo %d of %d", currentIndex+1, totalCount)

	// Add marked indicator if todo is marked for deletion
	if _, isMarked := markedForDeletion[todo.ID]; isMarked {
		footerTitle += " [MARKED FOR DELETION]"
	}

	// Footer stats (statusMsg now in message line, not footer)
	footerStats := ""
	if totalLines > availableHeight {
		footerStats = fmt.Sprintf("lines %d-%d of %d", scrollStart+1, scrollEnd, totalLines)
	}

	footer := RenderFooter(width, theme, footerTitle, footerStats)

	// Build main content (status + timestamps at top, then blank line, then title)
	mainContent := timestampLines + "\n\n" + renderedTitle + linkedSection

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

// wrapBodyText wraps body text while preserving original line breaks
func wrapBodyText(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		if line == "" {
			// Preserve blank lines
			wrappedLines = append(wrappedLines, "")
		} else {
			// Wrap long lines
			wrapped := wordWrap(line, width)
			wrappedLines = append(wrappedLines, wrapped)
		}
	}

	return strings.Join(wrappedLines, "\n")
}

// wordWrap wraps text to fit within a given width
func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		// If word itself is longer than width, just add it
		if len(word) > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			lines = append(lines, word)
			continue
		}

		// Check if adding this word would exceed width
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	// Add last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
