package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderDashboard renders the main dashboard view with stats
func RenderDashboard(width, height int, entries []models.Entry, todos []models.Todo) string {
	// Calculate stats
	totalEntries := len(entries)
	totalTodos := len(todos)
	openTodos := 0
	doneTodos := 0
	for _, todo := range todos {
		if todo.Status == "done" {
			doneTodos++
		} else {
			openTodos++
		}
	}

	// Count entries this week
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	entriesThisWeek := 0
	for _, entry := range entries {
		if entry.Timestamp.After(weekAgo) {
			entriesThisWeek++
		}
	}

	// Find most used tag
	tagCounts := make(map[string]int)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			tagCounts[tag]++
		}
	}
	for _, todo := range todos {
		for _, tag := range todo.Tags {
			tagCounts[tag]++
		}
	}
	topTag := ""
	maxCount := 0
	for tag, count := range tagCounts {
		if count > maxCount {
			maxCount = count
			topTag = tag
		}
	}
	topTagDisplay := "none"
	if topTag != "" {
		topTagDisplay = fmt.Sprintf("@%s: %d uses", topTag, maxCount)
	}

	// Massive ASCII art title - centered
	title := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Width(width - 8).
		Align(lipgloss.Center).
		Render(strings.Join([]string{
			"",
			"░█████╗░███╗░░░███╗░█████╗░░██████╗",
			"██╔══██╗████╗░████║██╔══██╗██╔════╝",
			"███████║██╔████╔██║██║░░██║╚█████╗░",
			"██╔══██║██║╚██╔╝██║██║░░██║░╚═══██╗",
			"██║░░██║██║░╚═╝░██║╚█████╔╝██████╔╝",
			"╚═╝░░╚═╝╚═╝░░░░░╚═╝░╚════╝░╚═════╝░",
			"",
		}, "\n"))

	// Stats display - centered
	statsStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Width(width - 8).
		Align(lipgloss.Center)

	stats := statsStyle.Render(fmt.Sprintf(
		"%d entries  │  %d this week  │  %s  │  %d todos (%d open)",
		totalEntries, entriesThisWeek, topTagDisplay, totalTodos, openTodos,
	))

	// Help text - centered
	helpTextStyle := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Width(width - 8).
		Align(lipgloss.Center)

	help := helpTextStyle.Render("n: new entry • a: add todo • t: todos • e: entries • q: quit")

	// Content - vertically centered in box
	boxContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		stats,
		"",
		help,
	)

	// Full screen box that fills most of terminal
	boxStyle := GetFullScreenBox(width, height).
		AlignVertical(lipgloss.Center)

	return boxStyle.Render(boxContent)
}
