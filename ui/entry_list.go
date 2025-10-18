package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, statusMsg string) string {
	container := GetContainerStyle(width, height)
	title := GetTitleStyle(width).Render("Entries")

	// Sort entries by timestamp (newest first)
	sorted := make([]models.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})

	// Build entry list
	var listItems []string

	if len(sorted) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Width(width - 4).
			Align(lipgloss.Center)
		listItems = append(listItems, emptyStyle.Render("No entries yet"))
	} else {
		for i, entry := range sorted {
			// Format: 2025-10-18 14:32 | Meeting with team
			timestamp := entry.Timestamp.Format("2006-01-02 15:04")
			line := fmt.Sprintf("%s | %s", timestamp, entry.Title)

			// Truncate if too long
			maxLen := width - 6
			if len(line) > maxLen {
				line = line[:maxLen-3] + "..."
			}

			// Style selected item differently
			var styled string
			if i == selectedIdx {
				selectedStyle := lipgloss.NewStyle().
					Foreground(accentColor).
					Bold(true)
				styled = selectedStyle.Render("► " + line)
			} else {
				normalStyle := lipgloss.NewStyle().
					Foreground(subtleColor)
				styled = normalStyle.Render("  " + line)
			}

			listItems = append(listItems, styled)
		}
	}

	list := strings.Join(listItems, "\n")

	// Status message (if present)
	status := ""
	if statusMsg != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)
		status = "\n" + statusStyle.Render(statusMsg) + "\n"
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)
	help := helpStyle.Render("j/k: navigate • enter: view • d: delete • esc: back • q: quit")

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list,
		status,
		help,
	)

	return container.Render(content)
}
