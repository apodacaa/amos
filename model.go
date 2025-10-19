package main

import (
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/ui"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Model holds the application state
type Model struct {
	view             string         // Current view: "dashboard", "entry", "entries", "view_entry", "todos", "tag_picker", or "add_todo"
	previousView     string         // Previous view (for context-aware escape)
	width            int            // Terminal width
	height           int            // Terminal height
	textarea         textarea.Model // Textarea for entry input
	todoInput        textarea.Model // Single-line input for standalone todos
	currentEntry     models.Entry   // Entry being edited
	currentTodo      models.Todo    // Standalone todo being created
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
	filterTag        string         // Current tag filter (empty = no filter)
	availableTags    []string       // All unique tags across entries
	selectedTag      int            // Selected tag index in tag picker
}

// NewModel creates a new model with default values
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "First line is the title...\n\nStart typing your entry here.\n\nUse @tags for organization.\n\nUse !todos for tasks."
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

	// Create single-line input for standalone todos
	todoInput := textarea.New()
	todoInput.Placeholder = "Todo title..."
	todoInput.CharLimit = 0
	todoInput.SetWidth(60)
	todoInput.SetHeight(1) // Single line
	todoInput.FocusedStyle.CursorLine = ui.GetTextareaStyle()
	todoInput.BlurredStyle.CursorLine = ui.GetTextareaStyle()
	todoInput.FocusedStyle.Placeholder = ui.GetPlaceholderStyle()
	todoInput.BlurredStyle.Placeholder = ui.GetPlaceholderStyle()
	todoInput.FocusedStyle.Prompt = ui.GetPromptStyle()
	todoInput.BlurredStyle.Prompt = ui.GetPromptStyle()
	todoInput.FocusedStyle.Text = ui.GetTextStyle()
	todoInput.BlurredStyle.Text = ui.GetTextStyle()

	return Model{
		view:      "dashboard",
		width:     80, // Default width
		height:    24, // Default height
		textarea:  ta,
		todoInput: todoInput,
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
		} else if m.view == "tag_picker" {
			return m.handleTagPickerKeys(msg)
		} else if m.view == "add_todo" {
			return m.handleAddTodoKeys(msg)
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
			// If we're in add_todo view, go back to dashboard after saving
			if m.view == "add_todo" {
				m.view = "dashboard"
				m.statusMsg = ""
				return m, nil
			}
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
			// Only reset selection if we don't have a valid selection
			// (e.g., first load or if selection is out of bounds)
			if m.selectedTodo >= len(msg.todos) {
				m.selectedTodo = 0
			}
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
			m.statusMsg = "Error saving todo: " + msg.err.Error()
			m.statusTime = time.Now()
		}
		// Don't reload - status already updated in memory
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
		return ui.RenderEntryList(m.width, m.height, m.entries, m.selectedEntry, m.statusMsg, m.todos, m.filterTag)
	case "view_entry":
		return ui.RenderEntryView(m.width, m.height, m.viewingEntry, m.todos)
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.todos, m.selectedTodo, m.statusMsg)
	case "tag_picker":
		return ui.RenderTagPicker(m.width, m.height, m.availableTags, m.selectedTag)
	case "add_todo":
		return ui.RenderAddTodoForm(m.width, m.height, m.todoInput, m.statusMsg)
	default:
		return ui.RenderDashboard(m.width, m.height)
	}
}

// handleNewEntry is a shared handler for creating a new entry (from any view)
func (m Model) handleNewEntry() (Model, tea.Cmd) {
	m.previousView = m.view
	m.view = "entry"
	m.currentEntry = models.Entry{
		ID:        m.generateID(),
		Timestamp: time.Now(),
	}
	m.textarea.Reset()
	m.textarea.Focus()
	m.hasUnsaved = false
	m.savedContent = ""
	m.statusMsg = ""
	return m, textarea.Blink
}

// handleAddTodo is a shared handler for creating a standalone todo (from any view)
func (m Model) handleAddTodo() (Model, tea.Cmd) {
	m.previousView = m.view
	m.view = "add_todo"
	m.currentTodo = models.Todo{
		ID:        m.generateID(),
		Status:    "open",
		Position:  0,
		CreatedAt: time.Now(),
	}
	m.todoInput.Reset()
	m.todoInput.Focus()
	m.statusMsg = ""
	return m, textarea.Blink
}

// generateID generates a new UUID string
func (m Model) generateID() string {
	return uuid.New().String()
}
