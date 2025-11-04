package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/ui"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// Model holds the application state
type Model struct {
	view                 string            // Current view: "entry", "entries", "view_entry", "todos", "view_todo", "unified_filter", or "add_todo"
	width                int               // Terminal width
	height               int               // Terminal height
	textarea             textarea.Model    // Textarea for entry input
	todoInput            textarea.Model    // Single-line input for standalone todos
	unifiedFilterInput   textarea.Model    // Single-line input for unified filtering (tags + dates)
	currentEntry         models.Entry      // Entry being edited
	currentTodo          models.Todo       // Standalone todo being created
	viewingEntry         models.Entry      // Entry being viewed (read-only)
	viewingTodo          models.Todo       // Todo being viewed (read-only)
	scrollOffset         int               // Scroll offset for long entry view
	statusMsg            string            // Status message to display
	statusTime           time.Time         // When status message was set
	hasUnsaved           bool              // Whether there are unsaved changes
	savedContent         string            // Last saved content (to detect changes)
	confirmingExit       bool              // Whether showing exit confirmation
	entries              []models.Entry    // All entries (for list view)
	selectedEntry        int               // Selected entry index in list
	todos                []models.Todo     // All todos (raw, unsorted)
	displayTodos         []models.Todo     // Sorted todos for display (only updated on load/refresh)
	selectedTodo         int               // Selected todo index in list
	filterTags           []string          // Current tag filters (empty = no filter), supports multiple tags with AND logic
	filterContext        string            // Context for filtering: "entries" or "todos" (which view to return to)
	filterDate           string            // Current date filter preset (empty = no filter)
	availableTags        []string          // All unique tags across entries
	autocompleteTag      string            // Current autocomplete suggestion for tag input
	markedForDeletion    map[string]string // Map of ID to type ("entry" or "todo") for items marked for deletion
	deleteConfirmPending bool              // Whether $ has been pressed once (waiting for second press)
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
	ui.ApplyTextareaStyle(&ta)

	// Create single-line input for standalone todos
	todoInput := textarea.New()
	todoInput.Placeholder = "Todo title..."
	todoInput.CharLimit = 0
	todoInput.SetWidth(60)
	todoInput.SetHeight(1) // Single line
	ui.ApplyTextareaStyle(&todoInput)

	// Create single-line input for unified filtering (tags + dates)
	unifiedFilterInput := textarea.New()
	unifiedFilterInput.Placeholder = "e.g. @work yesterday, last 30 days @client"
	unifiedFilterInput.CharLimit = 0
	unifiedFilterInput.SetWidth(60)
	unifiedFilterInput.SetHeight(1) // Single line
	ui.ApplyTextareaStyle(&unifiedFilterInput)

	return Model{
		view:               "entries",
		width:              80, // Default width
		height:             24, // Default height
		textarea:           ta,
		todoInput:          todoInput,
		unifiedFilterInput: unifiedFilterInput,
		markedForDeletion:  make(map[string]string),
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
		case "view_todo":
			return m.handleViewTodoKeys(msg)
		case "unified_filter":
			return m.handleUnifiedFilterKeys(msg)
		case "add_todo":
			return m.handleAddTodoKeys(msg)
		default:
			// Default to entry list (app opens to entries)
			return m.handleEntriesListKeys(msg)
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
			m.statusMsg = "Saved"
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

	case deleteCompleteMsg:
		// Delete successful - clear marked items and confirmation state, reload data
		m.markedForDeletion = make(map[string]string)
		m.deleteConfirmPending = false

		// Build success message with counts
		var successMsg string
		if msg.entryCount > 0 && msg.todoCount > 0 {
			// Both entries and todos deleted
			entryWord := "entry"
			if msg.entryCount > 1 {
				entryWord = "entries"
			}
			todoWord := "todo"
			if msg.todoCount > 1 {
				todoWord = "todos"
			}
			successMsg = fmt.Sprintf("Deleted %d %s", msg.entryCount, entryWord)
			if msg.linkedTodoCount > 0 {
				successMsg += fmt.Sprintf(" (%d linked todos)", msg.linkedTodoCount)
			}
			successMsg += fmt.Sprintf(" and %d %s", msg.todoCount, todoWord)
		} else if msg.entryCount > 0 {
			// Only entries deleted
			entryWord := "entry"
			if msg.entryCount > 1 {
				entryWord = "entries"
			}
			successMsg = fmt.Sprintf("Deleted %d %s", msg.entryCount, entryWord)
			if msg.linkedTodoCount > 0 {
				successMsg += fmt.Sprintf(" (%d linked todos)", msg.linkedTodoCount)
			}
		} else {
			// Only todos deleted
			todoWord := "todo"
			if msg.todoCount > 1 {
				todoWord = "todos"
			}
			successMsg = fmt.Sprintf("Deleted %d %s", msg.todoCount, todoWord)
		}

		m.statusMsg = successMsg
		m.statusTime = time.Now()
		// Reload entries and todos to refresh the view
		return m, tea.Batch(m.loadEntriesAndTodos(), clearStatusAfterDelay())

	case deleteErrorMsg:
		// Delete failed - clear confirmation and show error
		m.deleteConfirmPending = false
		m.statusMsg = msg.message
		if msg.err != nil {
			m.statusMsg += ": " + msg.err.Error()
		}
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()
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
		return ui.RenderEntryForm(m.width, m.height, m.textarea, m.statusMsg, m.hasUnsaved)
	case "entries":
		return ui.RenderEntryList(m.width, m.height, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
	case "view_entry":
		// Calculate filtered/sorted entry position for footer display
		filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)
		sorted := helpers.SortEntriesForDisplay(filtered)
		return ui.RenderEntryView(m.width, m.height, m.viewingEntry, m.todos, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedEntry, len(sorted))
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.displayTodos, m.entries, m.selectedTodo, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
	case "view_todo":
		// Calculate filtered/sorted todo position for footer display
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)
		sorted := helpers.SortTodosForDisplay(filtered)
		return ui.RenderTodoView(m.width, m.height, m.viewingTodo, m.entries, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedTodo, len(sorted))
	case "unified_filter":
		return ui.RenderUnifiedFilter(m.width, m.height, m.unifiedFilterInput, m.availableTags, m.autocompleteTag, m.statusMsg)
	case "add_todo":
		return ui.RenderAddTodoForm(m.width, m.height, m.todoInput, m.statusMsg, m.hasUnsaved)
	default:
		return ui.RenderEntryList(m.width, m.height, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
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

// openUnifiedFilter opens the unified filter view with current filter state
func (m Model) openUnifiedFilter(context string) (Model, tea.Cmd) {
	m.filterContext = context
	m.availableTags = helpers.ExtractUniqueTagsFromAll(m.entries, m.todos)

	// Reconstruct current filter string from tags and date
	var filterParts []string
	filterParts = append(filterParts, m.filterTags...)
	if m.filterDate != "" {
		filterParts = append(filterParts, m.filterDate)
	}
	currentFilter := strings.Join(filterParts, " ")

	// Set the input with current filter
	m.unifiedFilterInput.Reset()
	m.unifiedFilterInput.SetValue(currentFilter)
	m.unifiedFilterInput.Focus()
	m.autocompleteTag = ""
	m.view = "unified_filter"
	m.statusMsg = ""
	return m, textarea.Blink
}

// clearFilters clears all active filters and resets selection
func (m Model) clearFilters() (Model, tea.Cmd) {
	if len(m.filterTags) > 0 || m.filterDate != "" {
		m.filterTags = []string{}
		m.filterDate = ""
		m.selectedEntry = 0
		m.selectedTodo = 0
		m.statusMsg = "Filter cleared"
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()
	}
	return m, nil
}
