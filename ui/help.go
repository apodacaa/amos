package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHelp renders the help page with symbols and commands
func RenderHelp(width, height int, theme Theme, scrollOffset int, updateAvailable bool, updateDismissed bool, latestVersion string) string {
	// Get plain content and apply styling
	plainContent := GetHelpContent()
	fullContent := styleHelpContent(plainContent)

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

// styleHelpContent applies styling to help content (bold section titles)
func styleHelpContent(content string) string {
	// Split into sections (separated by double newlines)
	sections := strings.Split(content, "\n\n")
	boldStyle := lipgloss.NewStyle().Bold(true)

	for i, section := range sections {
		// Check if this section starts with a title (all caps, single line)
		lines := strings.SplitN(section, "\n", 2)
		if len(lines) > 0 && isTitle(lines[0]) {
			// Apply bold to the title
			lines[0] = boldStyle.Render(lines[0])
			sections[i] = strings.Join(lines, "\n")
		}
	}

	return strings.Join(sections, "\n\n")
}

// isTitle checks if a line is a section title (all uppercase letters/spaces)
func isTitle(line string) bool {
	if line == "" {
		return false
	}
	for _, r := range line {
		if r != ' ' && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
