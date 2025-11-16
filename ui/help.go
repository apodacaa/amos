package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderHelp renders the help page with symbols and commands
func RenderHelp(width, height int, theme Theme, scrollOffset int) string {
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

	// Manual assembly (since AssembleView doesn't support scrolling)
	// Add padding to fill remaining vertical space
	contentHeight := strings.Count(content, "\n") + 1
	availableContentHeight := height - 3 // header + footer + message line
	paddingNeeded := availableContentHeight - contentHeight
	if paddingNeeded > 0 {
		content += strings.Repeat("\n", paddingNeeded)
	}

	// Assemble: header + content + footer + empty message line
	return header + "\n" + content + "\n" + footer + "\n"
}

func renderSymbolsSection(theme Theme) string {
	title := lipgloss.NewStyle().Bold(true).Render("SYMBOLS")

	items := []string{
		"Entry List Markers:",
		"  +  Entry has linked todos (color shows highest priority)",
		"     Green = next, Magenta = open, Dim = done",
		"  D  Entry marked for deletion",
		"",
		"Todo List Markers:",
		"  +  Todo is linked to an entry",
		"  D  Todo marked for deletion",
		"",
		"Todo Status Indicators:",
		"  [ ] Open todo",
		"  [>] Next todo (high priority)",
		"  [x] Done todo",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

func renderCommandsSection(theme Theme) string {
	title := lipgloss.NewStyle().Bold(true).Render("COMMANDS")

	items := []string{
		"Global:",
		"  n:new         Create new entry",
		"  a:todo        Add standalone todo",
		"  e:entries     Jump to entry list",
		"  t:todos       Jump to todo list",
		"  s:theme       Select theme",
		"  ?:help        Show this help page",
		"  q:quit        Quit application",
		"",
		"Entry List:",
		"  j/k:nav       Navigate up/down",
		"  enter:view    View entry detail",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Entry View:",
		"  j/k:nav       Navigate between entries",
		"  f/b:page      Page forward/backward in long entries",
		"  d:del         Mark/unmark for deletion",
		"  esc:back      Return to entry list",
		"",
		"Todo List:",
		"  j/k:nav       Navigate up/down",
		"  space:toggle  Cycle todo status (open → next → done)",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Forms (Entry/Todo):",
		"  ctrl+s:save   Save and close",
		"  esc:cancel    Cancel without saving",
		"",
		"Filter:",
		"  @tag          Filter by tag (tab to autocomplete)",
		"  today         Show items from today",
		"  yesterday     Show items from yesterday",
		"  last N days   Show items from last N days",
		"  YYYY-MM-DD    Show items from specific date",
		"  date to date  Show items in date range",
		"  enter:apply   Apply filter",
		"  esc:cancel    Cancel filter",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

func renderNavigationSection(theme Theme) string {
	title := lipgloss.NewStyle().Bold(true).Render("NAVIGATION")

	items := []string{
		"Entry list and todo list are peer views - use e/t to jump between them.",
		"Global shortcuts (n, a, s, ?) work from any read-only view.",
		"Press esc to exit forms and return to their natural home view.",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

func renderDataStorageSection(theme Theme) string {
	title := lipgloss.NewStyle().Bold(true).Render("DATA STORAGE")

	items := []string{
		"All data is stored in plain JSON files in ~/.amos/:",
		"  ~/.amos/entries.json   Journal entries",
		"  ~/.amos/todos.json     Todo items",
		"  ~/.amos/config.json    User preferences (theme)",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}
