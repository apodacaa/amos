package ui

import (
	"strings"

	"github.com/apodacaa/amos/internal/storage"
	"github.com/charmbracelet/lipgloss"
)

// RenderWorkspaceSelector renders a centered modal for workspace selection
func RenderWorkspaceSelector(width, height int, theme Theme, workspaces []storage.Workspace, selectedIdx int, activeWorkspace string) string {
	// Modal dimensions
	modalWidth := 60
	if modalWidth > width-4 {
		modalWidth = width - 4
	}

	// Build workspace list
	var workspaceItems []string
	for i, ws := range workspaces {
		// 2-char prefix: cursor + active indicator
		cursor := " "
		if i == selectedIdx {
			cursor = ">"
		}
		active := " "
		if ws.Name == activeWorkspace {
			active = "*"
		}
		prefix := cursor + active

		// Workspace name in bold
		name := lipgloss.NewStyle().Bold(true).Render(ws.Name)

		// Workspace data dir as description
		desc := "  " + ws.DataDir

		item := prefix + name + "\n" + desc

		// Highlight selected item
		if i == selectedIdx {
			itemStyle := lipgloss.NewStyle().
				Width(modalWidth-4).
				Padding(0, 1)

			// Use theme's selection style if available
			if theme.Colors.Selection != "" {
				itemStyle = itemStyle.
					Background(theme.Colors.Selection).
					Foreground(theme.Colors.Foreground)
			} else {
				itemStyle = itemStyle.Reverse(true)
			}

			item = itemStyle.Render(item)
		} else {
			item = lipgloss.NewStyle().
				Width(modalWidth-4).
				Padding(0, 1).
				Render(item)
		}

		workspaceItems = append(workspaceItems, item)
	}

	workspaceList := strings.Join(workspaceItems, "\n\n")

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Width(modalWidth - 4).
		Align(lipgloss.Center).
		Render("Select Workspace")

	// Help text
	help := FormatHelpLeft(modalWidth-4,
		"j/k", "navigate",
		"enter", "select",
		"esc", "cancel",
	)

	// Build modal content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		workspaceList,
		"",
		help,
	)

	// Modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Colors.Primary).
		Padding(1, 2).
		Width(modalWidth).
		Render(content)

	// Center modal on screen
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		modalBox,
	)

	return centered
}
