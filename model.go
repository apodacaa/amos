package main

import (
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/ui"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Model holds the application state
type Model struct {
	view             string         // Current view: "dashboard", "entry", "entries", "view_entry", or "todos"
	width            int            // Terminal width
	height           int            // Terminal height
	textarea         textarea.Model // Textarea for entry input
	currentEntry     models.Entry   // Entry being edited
	viewingEntry     models.Entry   // Entry being viewed (read-only)
	statusMsg        string         // Status message to display
	statusTime       time.Time      // When status message was set
	hasUnsaved       bool           // Whether there are unsaved changes
	savedContent     string         // Last saved content (to detect changes)
	confirmingExit   bool           // Whether showing exit confirmation
	entries          []models.Entry // All entries (for list view)
	selectedEntry    int            // Selected entry index in list
	confirmingDelete bool           // Whether showing delete confirmation
	todos            []models.Todo  // All todos (for todo list view)
	selectedTodo     int            // Selected todo index in list
}

// NewModel creates a new model with default values
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "First line is the title...\n\nStart typing your entry here.\n\nUse @tags for organization."
	ta.Focus()
	ta.CharLimit = 0 // No limit
	ta.SetWidth(60)
	ta.SetHeight(10)

	// Style textarea with brutalist colors
	ta.FocusedStyle.CursorLine = ui.GetTextareaStyle()
	ta.BlurredStyle.CursorLine = ui.GetTextareaStyle()
	ta.FocusedStyle.Placeholder = ui.GetPlaceholderStyle()
	ta.BlurredStyle.Placeholder = ui.GetPlaceholderStyle()
	ta.FocusedStyle.Prompt = ui.GetPromptStyle()
	ta.BlurredStyle.Prompt = ui.GetPromptStyle()
	ta.FocusedStyle.Text = ui.GetTextStyle()
	ta.BlurredStyle.Text = ui.GetTextStyle()

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

// Update handles messages (Elm architecture)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Route to appropriate key handler based on view
		if m.view == "entry" {
			return m.handleEntryKeys(msg)
		} else if m.view == "entries" {
			return m.handleEntriesListKeys(msg)
		} else if m.view == "view_entry" {
			return m.handleViewEntryKeys(msg)
		} else if m.view == "todos" {
			return m.handleTodosListKeys(msg)
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

	case entriesLoadedMsg:
		if msg.err != nil {
			m.statusMsg = "Error loading entries: " + msg.err.Error()
		} else {
			m.entries = msg.entries
			m.selectedEntry = 0
		}
		return m, nil

	case todosLoadedMsg:
		if msg.err != nil {
			m.statusMsg = "Error loading todos: " + msg.err.Error()
		} else {
			m.todos = msg.todos
			m.selectedTodo = 0
		}
		return m, nil

	case entryDeletedMsg:
		if msg.err != nil {
			m.statusMsg = "Error deleting entry: " + msg.err.Error()
		} else {
			m.statusMsg = "✓ Entry deleted"
			m.confirmingDelete = false
			// Reload entries to update the list
			return m, m.loadEntries()
		}
		m.statusTime = time.Now()
		return m, nil

	case todoToggledMsg:
		if msg.err != nil {
			m.statusMsg = "Error toggling todo: " + msg.err.Error()
		} else {
			// Reload todos to update the list
			return m, m.loadTodos()
		}
		m.statusTime = time.Now()
		return m, nil

	case todoMovedMsg:
		if msg.err != nil {
			m.statusMsg = "Error moving todo: " + msg.err.Error()
		} else {
			// Reload todos to update the list
			return m, m.loadTodos()
		}
		return m, nil
	}

	// Update textarea if in entry view
	if m.view == "entry" {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

// View renders the UI (Elm architecture)
func (m Model) View() string {
	switch m.view {
	case "entry":
		return ui.RenderEntryForm(m.width, m.height, m.textarea, m.statusMsg)
	case "entries":
		return ui.RenderEntryList(m.width, m.height, m.entries, m.selectedEntry, m.statusMsg)
	case "view_entry":
		return ui.RenderEntryView(m.width, m.height, m.viewingEntry)
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.todos, m.selectedTodo, m.statusMsg)
	default:
		return ui.RenderDashboard(m.width, m.height)
	}
}
