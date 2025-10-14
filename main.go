package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model holds the application state
type model struct {
	view string // "dashboard" or "entry"
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			m.view = "entry"
		case "esc":
			m.view = "dashboard"
		}
	case tea.WindowSizeMsg:
		// Handle window resize - just return the model
		return m, nil
	}
	return m, nil
}

// View renders the UI
func (m model) View() string {
	if m.view == "entry" {
		return renderEntryForm()
	}
	return renderDashboard()
}

// Brutalist styles
var (
	// Base colors - using terminal defaults mostly
	subtleColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	accentColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#00FF00"}

	// Main container with border
	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(1, 2).
			Width(60)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Align(lipgloss.Center).
			Width(56)

	// Menu item style
	menuItemStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(0).
			MarginBottom(0)

	// Shortcut key style
	keyStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}).
			Italic(true).
			MarginTop(1)
)

func renderDashboard() string {
	var menuItems []string

	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s New Entry", keyStyle.Render("[n]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Todos", keyStyle.Render("[t]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Entries", keyStyle.Render("[e]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Search", keyStyle.Render("[s]"))),
	)

	menu := strings.Join(menuItems, "\n")
	help := helpStyle.Render("Press q to quit")

	content := titleStyle.Render("AMOS") + "\n\n" + menu + "\n" + help

	return containerStyle.Render(content)
}

func renderEntryForm() string {
	title := titleStyle.Render("NEW ENTRY")

	inputLabel := lipgloss.NewStyle().
		Foreground(subtleColor).
		Render("Title:")

	inputPlaceholder := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#444444"}).
		Render("Enter title here...")

	bodyLabel := lipgloss.NewStyle().
		Foreground(subtleColor).
		MarginTop(1).
		Render("Body:")

	bodyPlaceholder := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#444444"}).
		Render("Start typing your entry...\n\nUse @tags for organization")

	help := helpStyle.Render("Ctrl+S to save • Esc to cancel")

	content := title + "\n\n" +
		inputLabel + "\n" + inputPlaceholder + "\n" +
		bodyLabel + "\n" + bodyPlaceholder + "\n\n" +
		help

	return containerStyle.Render(content)
}

func main() {
	m := model{view: "dashboard"}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
