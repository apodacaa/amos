package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderEntryList renders the entry list view
func RenderEntryList(width, height int, entries []models.Entry, selectedIdx int, statusMsg string, todos []models.Todo, filterTag string) string {
	container := GetContainerStyle(width, height)

	// Update title to show filter if active
	titleText := "Entries"
	if filterTag != "" {
		titleText = fmt.Sprintf("Entries (filtered: %s)", filterTag)
	}
	title := GetTitleStyle(width).Render(titleText)

	// Filter entries by tag if filter is active
	filtered := helpers.FilterEntriesByTag(entries, filterTag)

	// Sort entries by timestamp (newest first)
	sorted := helpers.SortEntriesForDisplay(filtered)

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
				entryTodos := helpers.FilterTodosByEntry(todos, entry.ID)
				if len(entryTodos) > 0 {
					openCount, totalCount := helpers.CountTodoStats(entryTodos)
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

	// Help text (changes based on filter state)
	helpStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true)
	helpText := "j/k: navigate • enter: view • t: todos • @: filter • d: delete • esc: back • q: quit"
	if filterTag != "" {
		helpText = "j/k: navigate • enter: view • t: todos • @: clear filter • d: delete • esc: back • q: quit"
	}
	help := helpStyle.Render(helpText)

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
