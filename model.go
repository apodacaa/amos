package main

import (
	"fmt"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/ui"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Model holds the application state
type Model struct {
	view            string         // Current view: "dashboard", "entry", "entries", "view_entry", "todos", "tag_filter", or "add_todo"
	width           int            // Terminal width
	height          int            // Terminal height
	textarea        textarea.Model // Textarea for entry input
	todoInput       textarea.Model // Single-line input for standalone todos
	tagFilterInput  textarea.Model // Single-line input for tag filtering
	currentEntry    models.Entry   // Entry being edited
	currentTodo     models.Todo    // Standalone todo being created
	viewingEntry    models.Entry   // Entry being viewed (read-only)
	scrollOffset    int            // Scroll offset for long entry view
	statusMsg       string         // Status message to display
	statusTime      time.Time      // When status message was set
	hasUnsaved      bool           // Whether there are unsaved changes
	savedContent    string         // Last saved content (to detect changes)
	confirmingExit  bool           // Whether showing exit confirmation
	entries         []models.Entry // All entries (for list view)
	selectedEntry   int            // Selected entry index in list
	todos           []models.Todo  // All todos (raw, unsorted)
	displayTodos    []models.Todo  // Sorted todos for display (only updated on load/refresh)
	selectedTodo    int            // Selected todo index in list
	filterTags      []string       // Current tag filters (empty = no filter), supports multiple tags with AND logic
	filterContext   string         // Context for tag filtering: "entries" or "todos" (which view to return to)
	availableTags   []string       // All unique tags across entries
	autocompleteTag string         // Current autocomplete suggestion for tag input
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

	// Create single-line input for tag filtering
	tagFilterInput := textarea.New()
	tagFilterInput.Placeholder = "Type tags separated by spaces..."
	tagFilterInput.CharLimit = 0
	tagFilterInput.SetWidth(60)
	tagFilterInput.SetHeight(1) // Single line
	tagFilterInput.FocusedStyle.CursorLine = ui.GetTextareaStyle()
	tagFilterInput.BlurredStyle.CursorLine = ui.GetTextareaStyle()
	tagFilterInput.FocusedStyle.Placeholder = ui.GetPlaceholderStyle()
	tagFilterInput.BlurredStyle.Placeholder = ui.GetPlaceholderStyle()
	tagFilterInput.FocusedStyle.Prompt = ui.GetPromptStyle()
	tagFilterInput.BlurredStyle.Prompt = ui.GetPromptStyle()
	tagFilterInput.FocusedStyle.Text = ui.GetTextStyle()
	tagFilterInput.BlurredStyle.Text = ui.GetTextStyle()

	return Model{
		view:           "dashboard",
		width:          80, // Default width
		height:         24, // Default height
		textarea:       ta,
		todoInput:      todoInput,
		tagFilterInput: tagFilterInput,
	}
}

// Init initializes the model (Elm architecture)
func (m Model) Init() tea.Cmd {
	// Load entries and todos on startup
	return tea.Batch(textarea.Blink, m.loadEntriesAndTodos())
}

// Update handles messages (Elm architecture)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Route to appropriate key handler based on view
		switch m.view {
		case "entry":
			return m.handleEntryKeys(msg)
		case "entries":
			return m.handleEntriesListKeys(msg)
		case "view_entry":
			return m.handleViewEntryKeys(msg)
		case "todos":
			return m.handleTodosListKeys(msg)
		case "tag_filter":
			return m.handleTagFilterKeys(msg)
		case "add_todo":
			return m.handleAddTodoKeys(msg)
		default:
			return m.handleKeyPress(msg)
		}

	case tea.WindowSizeMsg:
		// Update terminal dimensions
		m.width = msg.Width
		m.height = msg.Height
		// Update textarea size (if terminal is large enough)
		if msg.Width > 10 && msg.Height > 12 {
			m.textarea.SetWidth(msg.Width - 10)
			m.textarea.SetHeight(msg.Height - 12)
		}
		return m, nil

	case saveCompleteMsg:
		if msg.err != nil {
			m.statusMsg = "Error saving: " + msg.err.Error()
		} else {
			m.statusMsg = "✓ Saved"
			// Mark as saved
			m.hasUnsaved = false
			if m.view == "entry" {
				// For entries, store current content
				m.savedContent = m.textarea.Value()
			}
			// For add_todo, we stay in the form (user can add another or press Esc)
		}
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()

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
			// Update display order (sort for display)
			m.displayTodos = helpers.SortTodosForDisplay(m.todos)

			// Reset selection if out of bounds
			if m.selectedTodo >= len(m.displayTodos) {
				m.selectedTodo = 0
			}
		}
		return m, nil

	case todoToggledMsg:
		if msg.err != nil {
			m.statusMsg = "Error saving todo: " + msg.err.Error()
			m.statusTime = time.Now()
		}
		// Don't reload - status already updated in memory
		return m, nil

	case statusTimeoutMsg:
		// Clear status message after timeout (only if it hasn't been updated recently)
		if time.Since(m.statusTime) >= 3*time.Second {
			m.statusMsg = ""
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
	// Minimum terminal size check (brutalist: be honest about limitations)
	const minWidth = 80
	const minHeight = 24

	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal too small\n\nMinimum size: %dx%d\nCurrent size: %dx%d\n\nResize your terminal",
			minWidth, minHeight, m.width, m.height)
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(msg)
	}

	switch m.view {
	case "entry":
		return ui.RenderEntryForm(m.width, m.height, m.textarea)
	case "entries":
		return ui.RenderEntryList(m.width, m.height, m.entries, m.selectedEntry, m.todos, m.filterTags)
	case "view_entry":
		return ui.RenderEntryView(m.width, m.height, m.viewingEntry, m.todos, m.scrollOffset)
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.displayTodos, m.entries, m.selectedTodo, m.filterTags)
	case "tag_filter":
		return ui.RenderTagFilter(m.width, m.height, m.tagFilterInput, m.availableTags, m.autocompleteTag)
	case "add_todo":
		return ui.RenderAddTodoForm(m.width, m.height, m.todoInput)
	default:
		return ui.RenderDashboard(m.width, m.height, m.entries, m.todos)
	}
}

// handleNewEntry is a shared handler for creating a new entry (from any view)
func (m Model) handleNewEntry() (Model, tea.Cmd) {
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
	m.view = "add_todo"
	m.currentTodo = models.Todo{
		ID:        m.generateID(),
		Status:    "open",
		CreatedAt: time.Now(),
	}
	m.todoInput.Reset()
	m.todoInput.Focus()
	m.hasUnsaved = false
	m.statusMsg = ""
	return m, textarea.Blink
}

// generateID generates a new UUID string
func (m Model) generateID() string {
	return uuid.New().String()
}
