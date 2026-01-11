package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHelp renders the help page with symbols and commands
func RenderHelp(width, height int, theme Theme, scrollOffset int, updateAvailable bool, updateDismissed bool, latestVersion string) string {
	// Build help content sections
	sections := []string{
		renderSymbolsSection(theme),
		"",
		renderCommandsSection(theme),
		"",
		renderNavigationSection(theme),
		"",
		renderDataStorageSection(theme),
	}

	fullContent := strings.Join(sections, "\n")

	// Split content into lines
	contentLines := strings.Split(fullContent, "\n")
	totalLines := len(contentLines)

	// Calculate available height for content
	availableHeight := height - 3 // header + footer + message line
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Apply scroll offset and viewport windowing
	var content string
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

		visibleLines := contentLines[scrollStart:scrollEnd]
		content = strings.Join(visibleLines, "\n")
	} else {
		// All content fits, no scrolling needed
		content = fullContent
		scrollStart = 0
		scrollEnd = totalLines
	}

	// Header with page navigation keys
	header := RenderHeader(width, theme,
		"f/b", "page",
		"esc", "back",
		"q", "quit",
	)

	// Footer with scroll indicators
	footerTitle := "Help"
	footerStats := ""
	if totalLines > availableHeight {
		footerStats = fmt.Sprintf("lines %d-%d of %d", scrollStart+1, scrollEnd, totalLines)
	}

	footer := RenderFooter(width, theme, footerTitle, footerStats)

	// Generate update notice if available
	updateNotice := RenderUpdateNotice(width, theme, updateAvailable, updateDismissed, latestVersion)

	// Manual assembly (since AssembleView doesn't support scrolling)
	// Add padding to fill remaining vertical space
	contentHeight := strings.Count(content, "\n") + 1
	reservedLines := 3 // header + footer + message line
	if updateNotice != "" {
		reservedLines += 1 // update notice line
	}
	availableContentHeight := height - reservedLines
	paddingNeeded := availableContentHeight - contentHeight
	if paddingNeeded > 0 {
		content += strings.Repeat("\n", paddingNeeded)
	}

	// Assemble: header + content + footer + update notice + empty message line
	result := header + "\n" + content + "\n" + footer
	if updateNotice != "" {
		result += "\n" + updateNotice
	}
	result += "\n"

	return result
}

func renderSymbolsSection(theme Theme) string {
	// Get plain text content
	content := RenderSymbolsSectionText()

	// Extract title and body
	parts := strings.SplitN(content, "\n\n", 2)
	title := lipgloss.NewStyle().Bold(true).Render(parts[0])
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	return title + "\n\n" + body
}

func renderCommandsSection(theme Theme) string {
	// Get plain text content
	content := RenderCommandsSectionText()

	// Extract title and body
	parts := strings.SplitN(content, "\n\n", 2)
	title := lipgloss.NewStyle().Bold(true).Render(parts[0])
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	return title + "\n\n" + body
}

func renderNavigationSection(theme Theme) string {
	// Get plain text content
	content := RenderNavigationSectionText()

	// Extract title and body
	parts := strings.SplitN(content, "\n\n", 2)
	title := lipgloss.NewStyle().Bold(true).Render(parts[0])
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	return title + "\n\n" + body
}

func renderDataStorageSection(theme Theme) string {
	// Get plain text content
	content := RenderDataStorageSectionText()

	// Extract title and body
	parts := strings.SplitN(content, "\n\n", 2)
	title := lipgloss.NewStyle().Bold(true).Render(parts[0])
	body := ""
	if len(parts) > 1 {
		body = parts[1]
	}

	return title + "\n\n" + body
}
