package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, statusMsg string, todos []models.Todo) string {
	container := GetContainerStyle(width, height)
	title := GetTitleStyle(width).Render("Entries")

	// Sort entries by timestamp (newest first)
	sorted := helpers.SortEntriesForDisplay(entries)

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

			// Add todo stats if entry has todos
			if len(entry.TodoIDs) > 0 {
				// Count open todos for this entry
				openCount := 0
				totalCount := 0
				for _, todo := range todos {
					if todo.EntryID != nil && *todo.EntryID == entry.ID {
						totalCount++
						if todo.Status == "open" {
							openCount++
						}
					}
				}

				if totalCount > 0 {
					todoStats := lipgloss.NewStyle().
						Foreground(mutedColor).
						Render(fmt.Sprintf(" [%d todos: %d open]", totalCount, openCount))
					line += todoStats
				}
			}

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
	help := helpStyle.Render("j/k: navigate • enter: view • t: todos • d: delete • esc: back • q: quit")

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
