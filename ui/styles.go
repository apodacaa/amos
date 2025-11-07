package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// Cyberpunk color palette (Matrix/Blade Runner vibes)
const (
	NeonGreen    = lipgloss.Color("#00FF41") // Primary - Matrix green (headers/footers)
	ElectricBlue = lipgloss.Color("#0080FF") // Selection highlight - Tron blue
	HotPink      = lipgloss.Color("#FF00FF") // Accent 1 - Magenta
	Cyan         = lipgloss.Color("#00FFFF") // Accent 2 - Cyan/Aqua
	Orange       = lipgloss.Color("#FF8C00") // Warning/dates
	DarkGray     = lipgloss.Color("#555555") // Dimmed
	Black        = lipgloss.Color("#000000") // For text on neon backgrounds
)

// GetFullScreenBox returns a box that fills most of the terminal with consistent styling
func GetFullScreenBox(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(width - 2).  // Minimal margin for border
		Height(height - 3) // Adjust to prevent top border clipping
}

// GetTitleStyle returns a title style sized to container width
func GetTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width - 8).
		Align(lipgloss.Center)
}

// Textarea style helpers - use terminal defaults
func GetTextareaStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func GetPlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Faint(true)
}

func GetPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func GetTextStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// getBarStyle returns the shared full-width neon bar style with bold text
func getBarStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(NeonGreen).
		Foreground(Black).
		Bold(true).
		Width(width)
}

// GetEmptyStateStyle returns centered faint style for empty lists
func GetEmptyStateStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Faint(true).
		Width(width - 4).
		Align(lipgloss.Center)
}

// GetSelectedItemStyle returns electric blue selection style for list items
func GetSelectedItemStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(ElectricBlue).
		Foreground(Black).
		Bold(true).
		Width(width)
}

// GetSelectedDoneStyle returns electric blue selection style for completed items
func GetSelectedDoneStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(ElectricBlue).
		Foreground(Black).
		Bold(true).
		Width(width)
}

// GetNormalItemStyle returns terminal default foreground for list items
func GetNormalItemStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// GetDimmedStyle returns dark gray style for completed items
func GetDimmedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(DarkGray)
}

// GetBoldStyle returns bold style for emphasis
func GetBoldStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

// HighlightTagsInText highlights @tags with cyan color
func HighlightTagsInText(text string) string {
	tagRegex := regexp.MustCompile(`@\w+`)
	tagStyle := lipgloss.NewStyle().Foreground(Cyan).Bold(true)

	return tagRegex.ReplaceAllStringFunc(text, func(tag string) string {
		return tagStyle.Render(tag)
	})
}

// StyleDate returns orange-styled date string
func StyleDate(date string) string {
	return lipgloss.NewStyle().Foreground(Orange).Render(date)
}

// StyleTodoStatus returns color-coded status indicator
func StyleTodoStatus(status, indicator string) string {
	var style lipgloss.Style
	switch status {
	case "open":
		style = lipgloss.NewStyle().Foreground(HotPink).Bold(true)
	case "next":
		style = lipgloss.NewStyle().Foreground(NeonGreen).Bold(true)
	case "done":
		style = lipgloss.NewStyle().Foreground(DarkGray)
	default:
		style = lipgloss.NewStyle()
	}
	return style.Render(indicator)
}

// StyleMarker returns color-coded marker (D for delete, + for todo/entry)
func StyleMarker(marker string, isDelete bool) string {
	if isDelete {
		return lipgloss.NewStyle().Foreground(HotPink).Bold(true).Render(marker)
	}
	return lipgloss.NewStyle().Foreground(Cyan).Render(marker)
}

// ApplyTextareaStyle applies consistent styling to a textarea
func ApplyTextareaStyle(ta *textarea.Model) {
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Faint(true)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Faint(true)
	ta.FocusedStyle.Prompt = lipgloss.NewStyle()
	ta.BlurredStyle.Prompt = lipgloss.NewStyle()
	ta.FocusedStyle.Text = lipgloss.NewStyle()
	ta.BlurredStyle.Text = lipgloss.NewStyle()
}

// CalculateViewport calculates start/end indices for viewport windowing
// Centers selectedIdx in available height, handles boundary conditions
func CalculateViewport(totalItems, selectedIdx, height int) (start, end int) {
	availableHeight := height - 3 // Reserve space for header + footer + message line
	if availableHeight < 5 {
		availableHeight = 5
	}

	start = 0
	end = totalItems

	if totalItems > availableHeight {
		// Center selected item in viewport
		half := availableHeight / 2
		start = selectedIdx - half
		end = selectedIdx + half + 1

		// Adjust if near beginning
		if start < 0 {
			start = 0
			end = availableHeight
		}

		// Adjust if near end
		if end > totalItems {
			end = totalItems
			start = end - availableHeight
			if start < 0 {
				start = 0
			}
		}
	}

	return start, end
}

// AssembleView builds full view with header, content, footer, message line and padding
func AssembleView(header, content, footer string, width, height int, statusMsg string) string {
	contentHeight := height - 3 // Reserve space for header + footer + message line

	// Ensure content fits within available space
	listLines := strings.Split(content, "\n")
	if len(listLines) > contentHeight {
		listLines = listLines[:contentHeight]
	}

	// Calculate padding to fill remaining space
	padding := contentHeight - len(listLines)
	if padding < 0 {
		padding = 0
	}

	// Build message line (neomutt-style)
	messageLine := RenderMessageLine(width, statusMsg)

	// Assemble full view
	result := header + "\n" + strings.Join(listLines, "\n")
	if padding > 0 {
		result += strings.Repeat("\n", padding)
	}
	result += "\n" + footer + "\n" + messageLine

	return result
}

// RenderFormView renders a form with header, input, and message line
func RenderFormView(width, height int, inputView, statusMsg string, hasUnsaved bool, headerKeys ...string) string {
	// Build header with provided keys
	header := RenderHeader(width, headerKeys...)

	// Calculate padding for content area
	contentHeight := height - 2 // header + message line (no footer in forms)
	inputLines := lipgloss.Height(inputView)
	padding := contentHeight - inputLines
	if padding < 0 {
		padding = 0
	}

	// Build message line (neomutt-style)
	messageLine := RenderMessageLine(width, statusMsg)

	// Build full view (no footer)
	result := header + "\n" + inputView
	if padding > 0 {
		result += strings.Repeat("\n", padding)
	}
	result += "\n" + messageLine

	return result
}

// FormatHelp formats help text with reverse keys
// Centered alignment
// Example: FormatHelp(width, "n", "new entry", "a", "add todo")
func FormatHelp(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Center, keyDescPairs...)
}

// FormatHelpLeft formats help text with reverse keys
// Left-aligned for utility views (honest functional UI)
// Example: FormatHelpLeft(width, "n", "new entry", "a", "add todo")
func FormatHelpLeft(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Left, keyDescPairs...)
}

// formatHelpWithAlign is the shared implementation for help text formatting
func formatHelpWithAlign(width int, align lipgloss.Position, keyDescPairs ...string) string {
	var parts []string

	keyStyle := lipgloss.NewStyle().
		Reverse(true).
		Padding(0, 1). // Add small padding for readability
		Inline(true)   // Keep inline to prevent breaking

	descStyle := lipgloss.NewStyle().
		Inline(true) // Keep inline to prevent breaking

	for i := 0; i < len(keyDescPairs); i += 2 {
		if i+1 < len(keyDescPairs) {
			key := keyStyle.Render(keyDescPairs[i])
			desc := descStyle.Render("\u00A0" + keyDescPairs[i+1])
			// Combine key and desc as single unit with inline wrapper
			pair := lipgloss.NewStyle().Inline(true).Render(key + desc)
			parts = append(parts, pair)
		}
	}

	// Join parts with spacing
	result := ""
	for i, part := range parts {
		result += part
		if i < len(parts)-1 {
			result += "  " // Two spaces between items
		}
	}

	// Render with lipgloss (which may wrap based on width)
	rendered := lipgloss.NewStyle().
		Faint(true).
		Width(width - 8).
		Align(align).
		Render(result)

	// Add vertical spacing between wrapped lines
	lines := strings.Split(rendered, "\n")
	if len(lines) > 1 {
		rendered = strings.Join(lines, "\n\n") // Double newline for breathing room
	}

	return rendered
}

// RenderHeader renders the top bar with app name and help shortcuts
// Format: "n:new  a:todo  q:quit"
func RenderHeader(width int, keyDescPairs ...string) string {
	// Build shortcuts string
	var shortcuts []string
	for i := 0; i < len(keyDescPairs); i += 2 {
		if i+1 < len(keyDescPairs) {
			shortcuts = append(shortcuts, keyDescPairs[i]+":"+keyDescPairs[i+1])
		}
	}

	content := strings.Join(shortcuts, "  ")

	// Truncate if too long
	if len(content) > width {
		content = content[:width]
	}

	return getBarStyle(width).Render(content)
}

// RenderFooter renders the bottom bar with view context and stats
// Format: "Entries @work  15 items" or "15 items" (no leading space when title is empty)
func RenderFooter(width int, title string, stats string) string {
	content := title
	if stats != "" {
		if title != "" {
			content += "  "
		}
		content += stats
	}

	// Truncate if too long
	if len(content) > width {
		content = content[:width]
	}

	return getBarStyle(width).Render(content)
}

// RenderMessageLine renders the message line below footer (neomutt-style)
// Shows status messages, confirmations, errors. Blank when no message.
func RenderMessageLine(width int, statusMsg string) string {
	if statusMsg == "" {
		// Return blank line to maintain consistent layout
		return strings.Repeat(" ", width)
	}

	// Truncate if too long
	if len(statusMsg) > width {
		statusMsg = statusMsg[:width]
	}

	// Pad to full width for consistent appearance
	if len(statusMsg) < width {
		statusMsg = statusMsg + strings.Repeat(" ", width-len(statusMsg))
	}

	return statusMsg
}
