package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// Brutalist color palette - Pure monochrome concrete
var (
	// Primary text and borders - Dark gray (light) / Light gray (dark)
	subtleColor = lipgloss.AdaptiveColor{Light: "#404040", Dark: "#CCCCCC"}

	// Accent/highlight - Pure black/white for maximum contrast
	accentColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

	// Muted/help text - Mid gray
	mutedColor = lipgloss.AdaptiveColor{Light: "#808080", Dark: "#666666"}

	// Header/footer bar colors (full-width top/bottom bars)
	barFg = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}
	barBg = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}

	// Key highlight colors (help text badges)
	highlightFg = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}
	highlightBg = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
)

// GetFullScreenBox returns a box that fills most of the terminal with consistent styling
func GetFullScreenBox(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Width(width - 2).  // Minimal margin for border
		Height(height - 3) // Adjust to prevent top border clipping
}

// GetTitleStyle returns a title style sized to container width
func GetTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		Width(width - 8).
		Align(lipgloss.Center)
}

// Textarea style helpers
func GetTextareaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(subtleColor)
}

func GetPlaceholderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor)
}

func GetPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accentColor)
}

func GetTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(subtleColor)
}

// getBarStyle returns the shared full-width inverted bar style
func getBarStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(barFg).
		Background(barBg).
		Width(width).
		Bold(true)
}

// GetEmptyStateStyle returns centered muted style for empty lists
func GetEmptyStateStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Width(width - 4).
		Align(lipgloss.Center)
}

// GetSelectedItemStyle returns inverted accent style for selected list items
func GetSelectedItemStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(accentColor).
		Reverse(true).
		Width(width)
}

// GetSelectedDoneStyle returns inverted muted style for selected completed items
func GetSelectedDoneStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(mutedColor).
		Reverse(true).
		Width(width)
}

// GetNormalItemStyle returns normal accent foreground for list items
func GetNormalItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(accentColor)
}

// GetDimmedStyle returns muted foreground for completed items
func GetDimmedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(mutedColor)
}

// ApplyTextareaStyle applies consistent styling to a textarea
func ApplyTextareaStyle(ta *textarea.Model) {
	ta.FocusedStyle.CursorLine = GetTextareaStyle()
	ta.BlurredStyle.CursorLine = GetTextareaStyle()
	ta.FocusedStyle.Placeholder = GetPlaceholderStyle()
	ta.BlurredStyle.Placeholder = GetPlaceholderStyle()
	ta.FocusedStyle.Prompt = GetPromptStyle()
	ta.BlurredStyle.Prompt = GetPromptStyle()
	ta.FocusedStyle.Text = GetTextStyle()
	ta.BlurredStyle.Text = GetTextStyle()
}

// CalculateViewport calculates start/end indices for viewport windowing
// Centers selectedIdx in available height, handles boundary conditions
func CalculateViewport(totalItems, selectedIdx, height int) (start, end int) {
	availableHeight := height - 2 // Reserve space for header + footer
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

// AssembleView builds full view with header, content, footer and padding
func AssembleView(header, content, footer string, height int) string {
	contentHeight := height - 2 // Reserve space for header + footer

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

	// Assemble full view
	result := header + "\n" + strings.Join(listLines, "\n")
	if padding > 0 {
		result += strings.Repeat("\n", padding)
	}
	result += "\n" + footer

	return result
}

// RenderFormView renders a form with header, input, footer, and unsaved indicator
func RenderFormView(width, height int, inputView, footerTitle, statusMsg string, hasUnsaved bool, headerKeys ...string) string {
	// Build header with provided keys
	header := RenderHeader(width, headerKeys...)

	// Build footer with unsaved indicator in statusMsg
	footerStatus := statusMsg
	if footerStatus == "" && hasUnsaved {
		footerStatus = "Unsaved"
	}
	footer := RenderFooter(width, footerTitle, footerStatus)

	// Calculate padding for content area
	contentHeight := height - 2 // header + footer
	inputLines := lipgloss.Height(inputView)
	padding := contentHeight - inputLines
	if padding < 0 {
		padding = 0
	}

	// Build full view
	result := header + "\n" + inputView
	if padding > 0 {
		result += strings.Repeat("\n", padding)
	}
	result += "\n" + footer

	return result
}

// FormatHelp formats help text with bold keys (reverse colors for impact)
// Centered alignment
// Example: FormatHelp(width, "n", "new entry", "a", "add todo")
func FormatHelp(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Center, keyDescPairs...)
}

// FormatHelpLeft formats help text with bold keys (reverse colors for impact)
// Left-aligned for utility views (honest functional UI)
// Example: FormatHelpLeft(width, "n", "new entry", "a", "add todo")
func FormatHelpLeft(width int, keyDescPairs ...string) string {
	return formatHelpWithAlign(width, lipgloss.Left, keyDescPairs...)
}

// formatHelpWithAlign is the shared implementation for help text formatting
func formatHelpWithAlign(width int, align lipgloss.Position, keyDescPairs ...string) string {
	var parts []string

	keyStyle := lipgloss.NewStyle().
		Foreground(highlightFg).
		Background(highlightBg).
		Bold(true).
		Padding(0, 1). // Add small padding for readability
		Inline(true)   // Keep inline to prevent breaking

	descStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
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
		Foreground(mutedColor).
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
// Format: "Entries @work  15 items"
func RenderFooter(width int, title string, stats string) string {
	content := title
	if stats != "" {
		content += "  " + stats
	}

	// Truncate if too long
	if len(content) > width {
		content = content[:width]
	}

	return getBarStyle(width).Render(content)
}
