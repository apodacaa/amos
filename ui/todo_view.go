package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

// highlightTodoInText highlights the todo title in the entry body
func highlightTodoInText(body, todoTitle string) string {
	// Find "!todo [title]" pattern
	todoPattern := "!todo " + todoTitle

	// Split by the todo, highlight the todo part
	parts := strings.Split(body, todoPattern)
	if len(parts) <= 1 {
		// Not found, return as-is in normal style
		return GetNormalItemStyle().Render(body)
	}

	// Build styled string: body text + highlighted todo
	var result strings.Builder
	for i, part := range parts {
		result.WriteString(GetNormalItemStyle().Render(part))
		if i < len(parts)-1 {
			// Add highlighted todo between parts
			result.WriteString(GetNormalItemStyle().Render(todoPattern))
		}
	}
	return result.String()
}

// RenderTodoView renders a read-only view of a todo
func RenderTodoView(width, height int, todo models.Todo, allEntries []models.Entry, scrollOffset int) string {
	// Todo title (wrappable)
	titleStyle := GetNormalItemStyle().
		Width(width - 8)

	// Wrap title if long
	wrappedTitle := wordWrap(todo.Title, width-8)
	titleLines := strings.Split(wrappedTitle, "\n")

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
			linkedParts = append(linkedParts, GetDimmedStyle().Render("From entry..."))

			// Date, title, tags line
			dateStr := linkedEntry.Timestamp.Format("2006-01-02")
			tagStr := ""
			if len(linkedEntry.Tags) > 0 {
				var tagStrings []string
				for _, tag := range linkedEntry.Tags {
					tagStrings = append(tagStrings, "@"+tag)
				}
				tagStr = " " + strings.Join(tagStrings, " ")
			}

			metaLine := fmt.Sprintf("%s %s%s", dateStr, linkedEntry.Title, tagStr)
			linkedParts = append(linkedParts, GetNormalItemStyle().Render(metaLine))

			// Body (full text) - highlight the todo reference
			if linkedEntry.Body != "" {
				styledBody := highlightTodoInText(linkedEntry.Body, todo.Title)
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
	availableHeight := height - 2 - linkedLineCount // header + footer + linked section
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
		renderedTitle = titleStyle.Render(wrappedTitle)
		scrollStart = 0
		scrollEnd = totalLines
	}

	// Header
	header := RenderHeader(width, "n", "new", "a", "todo", "space", "cycle", "j/k", "nav", "d/u", "scroll", "e", "entries", "t", "todos", "q", "quit")

	// Footer: status indicator + date + tags + scroll info
	statusIcon := "[ ]" // open
	if todo.Status == "next" {
		statusIcon = "[>]"
	} else if todo.Status == "done" {
		statusIcon = "[x]"
	}

	footerTitle := fmt.Sprintf("%s %s", statusIcon, todo.CreatedAt.Format("2006-01-02"))

	footerStats := ""
	if totalLines > availableHeight {
		footerStats = fmt.Sprintf("lines %d-%d of %d", scrollStart+1, scrollEnd, totalLines)
	}

	footer := RenderFooter(width, footerTitle, footerStats)

	// Build main content
	mainContent := renderedTitle + linkedSection

	// Calculate padding
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
