package main

import (
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	"github.com/apodacaa/amos/ui"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Model holds the application state
type Model struct {
	view           string         // Current view: "dashboard" or "entry"
	width          int            // Terminal width
	height         int            // Terminal height
	textarea       textarea.Model // Textarea for entry input
	currentEntry   models.Entry   // Entry being edited
	statusMsg      string         // Status message to display
	statusTime     time.Time      // When status message was set
	hasUnsaved     bool           // Whether there are unsaved changes
	savedContent   string         // Last saved content (to detect changes)
	confirmingExit bool           // Whether showing exit confirmation
}

// saveCompleteMsg is sent when save operation completes
type saveCompleteMsg struct {
	err error
}

// NewModel creates a new model with default values
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "First line is the title...\n\nStart typing your entry here.\n\nUse @tags for organization."
	ta.Focus()
	ta.CharLimit = 0 // No limit
	ta.SetWidth(60)
	ta.SetHeight(10)

	return Model{
		view:     "dashboard",
		width:    80, // Default width
		height:   24, // Default height
		textarea: ta,
	}
}

// Init initializes the model (Elm architecture)
func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages and updates state (Elm architecture)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle keys if in entry view - let textarea handle them
		if m.view == "entry" {
			return m.handleEntryKeys(msg)
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		// Update terminal dimensions
		m.width = msg.Width
		m.height = msg.Height
		// Update textarea size
		m.textarea.SetWidth(msg.Width - 10)
		m.textarea.SetHeight(msg.Height - 12)
		return m, nil

	case saveCompleteMsg:
		if msg.err != nil {
			m.statusMsg = "Error saving: " + msg.err.Error()
		} else {
			m.statusMsg = "✓ Saved"
			// Mark as saved and store current content
			m.hasUnsaved = false
			m.savedContent = m.textarea.Value()
		}
		m.statusTime = time.Now()
		return m, nil
	}

	// Update textarea if in entry view
	if m.view == "entry" {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

// handleKeyPress processes keyboard input (dashboard view)
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "n":
		// Create new entry
		m.view = "entry"
		m.currentEntry = models.Entry{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
		}
		m.textarea.Reset()
		m.textarea.Focus()
		m.statusMsg = ""
		m.hasUnsaved = false
		m.savedContent = ""
		m.confirmingExit = false
		return m, textarea.Blink
	case "esc":
		m.view = "dashboard"
	}
	return m, nil
}

// handleEntryKeys processes keyboard input (entry view)
func (m Model) handleEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Check if showing confirmation
		if m.confirmingExit {
			// User pressed Esc again - discard changes and exit to dashboard
			m.view = "dashboard"
			m.textarea.Blur()
			m.confirmingExit = false
			m.statusMsg = ""
			m.hasUnsaved = false
			return m, nil
		}

		// Check for unsaved changes
		currentContent := m.textarea.Value()
		if m.hasUnsaved && currentContent != m.savedContent {
			// Show confirmation prompt
			m.confirmingExit = true
			m.statusMsg = "⚠ Unsaved changes! Press Esc again to discard, or Ctrl+S to save"
			return m, nil
		}

		// No unsaved changes, safe to exit
		m.view = "dashboard"
		m.textarea.Blur()
		m.confirmingExit = false
		return m, nil

	case "ctrl+s":
		// Save entry
		m.confirmingExit = false // Clear confirmation if showing
		return m, m.saveEntry()

	default:
		// If confirming exit and user starts typing, cancel confirmation
		if m.confirmingExit {
			m.confirmingExit = false
			m.statusMsg = ""
		}

		// Let textarea handle the key
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)

		// Mark as having unsaved changes
		m.hasUnsaved = true

		return m, cmd
	}
}

// saveEntry saves the current entry
func (m Model) saveEntry() tea.Cmd {
	return func() tea.Msg {
		content := m.textarea.Value()

		// Parse content into title and body
		title, body := helpers.ParseEntryContent(content)

		// Extract tags from title and body
		tags := helpers.ExtractTags(title + " " + body)

		// Update current entry
		m.currentEntry.Title = title
		m.currentEntry.Body = body
		m.currentEntry.Tags = tags
		m.currentEntry.Timestamp = time.Now()

		// Save to storage
		err := storage.SaveEntry(m.currentEntry)

		return saveCompleteMsg{err: err}
	}
}

// View renders the UI (Elm architecture)
func (m Model) View() string {
	switch m.view {
	case "entry":
		return ui.RenderEntryForm(m.width, m.height, m.textarea, m.statusMsg)
	default:
		return ui.RenderDashboard(m.width, m.height)
	}
}
