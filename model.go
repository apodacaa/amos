package main

import (
	"fmt"
	"strings"
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
//
// Design: Follows Bubble Tea's Elm Architecture with all state in a single struct.
// Fields are grouped by concern to make the 35+ fields easier to understand and maintain.
type Model struct {
	// ============================================================
	// View Navigation
	// ============================================================
	view         string // Current view: "entry", "entries", "view_entry", "todos", "view_todo", "unified_filter", "add_todo", "theme_selector", or "help"
	previousView string // Previous view (for returning from modals like theme_selector and help)

	// ============================================================
	// Terminal Dimensions
	// ============================================================
	width  int // Terminal width (updated on window resize)
	height int // Terminal height (updated on window resize)

	// ============================================================
	// Text Input Components (Bubble Tea textareas)
	// ============================================================
	textarea           textarea.Model // Multi-line textarea for entry input (title + body)
	todoInput          textarea.Model // Single-line input for standalone todo title
	unifiedFilterInput textarea.Model // Single-line input for unified filtering (tags + dates)

	// ============================================================
	// Current Editing State (mutable)
	// Used when creating or editing entries/todos via forms
	// ============================================================
	currentEntry    models.Entry // Entry being created or edited
	currentTodo     models.Todo  // Standalone todo being created or edited
	editingMode     bool         // true = editing existing item, false = creating new
	originalEntryID string       // Entry ID being edited (empty string if creating new)
	originalTodoID  string       // Todo ID being edited (empty string if creating new)

	// ============================================================
	// View-Only State (read-only views)
	// Used when viewing entries/todos in detail
	// ============================================================
	viewingEntry models.Entry // Entry being viewed in "view_entry" (read-only)
	viewingTodo  models.Todo  // Todo being viewed in "view_todo" (read-only)
	scrollOffset int          // Scroll offset for paginated views (entry/todo/help pages)

	// ============================================================
	// UI Feedback State
	// ============================================================
	statusMsg      string    // Status message to display (e.g., "saved", "Filter applied")
	statusTime     time.Time // Timestamp when status message was set (for auto-clear after 3s)
	hasUnsaved     bool      // Whether there are unsaved changes in current form
	savedContent   string    // Last saved content (to detect changes for hasUnsaved)
	confirmingExit bool      // Whether showing exit confirmation dialog

	// ============================================================
	// Data Collections (loaded from storage)
	// Single source of truth: Todo.EntryID links to Entry.ID (no Entry.TodoIDs)
	// ============================================================
	entries       []models.Entry // All entries (unsorted, loaded from storage)
	selectedEntry int            // Selected entry index in list view
	todos         []models.Todo  // All todos (raw, unsorted, loaded from storage)
	displayTodos  []models.Todo  // Sorted todos for display (updated on load/refresh)
	selectedTodo  int            // Selected todo index in list view

	// ============================================================
	// Filtering State (unified filtering for both entries and todos)
	// Design: Filter state persists across navigation until explicitly cleared with 'c' key
	// ============================================================
	filterTags      []string // Current tag filters (empty = no filter), AND logic: all tags must match
	filterContext   string   // Context for filtering: "entries" or "todos" (determines return view from filter)
	filterDate      string   // Current date filter string (e.g., "today", "last 14 days", "2025-01-01 to 2025-01-31")
	availableTags   []string // All unique tags extracted from entries and todos (for autocomplete)
	autocompleteTag string   // Current autocomplete suggestion for tag input (cleared after tab)

	// ============================================================
	// Deletion State (neomutt-style multi-select deletion)
	// Design: 'd' marks items, '$' executes deletion, marks persist across navigation
	// ============================================================
	markedForDeletion    map[string]string // Map of ID -> type ("entry" or "todo") for marked items
	deleteConfirmPending bool              // true after first '$' press (waiting for 'y' or 'n' confirmation)

	// ============================================================
	// Configuration and Theming
	// ============================================================
	config        models.Config // User configuration (loaded from ~/.amos/config.json)
	currentTheme  ui.Theme      // Currently active theme (brutalist or cyberpunk)
	selectedTheme int           // Selected theme index in theme selector modal
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

	// Load config from storage (or use default)
	config, _ := storage.LoadConfig()

	// Load theme based on config
	currentTheme := ui.GetThemeByName(config.Theme)

	return Model{
		view:               "entries",
		width:              80, // Default width
		height:             24, // Default height
		textarea:           ta,
		todoInput:          todoInput,
		unifiedFilterInput: unifiedFilterInput,
		markedForDeletion:  make(map[string]string),
		config:             config,
		currentTheme:       currentTheme,
		selectedTheme:      0, // Will be set when opening theme selector
	}
}

// Init initializes the model (Elm architecture)
func (m Model) Init() tea.Cmd {
	// Load entries and todos on startup (single source of truth: Todo.EntryID)
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
		case "theme_selector":
			return m.handleThemeSelectorKeys(msg)
		case "help":
			return updateHelp(m, msg)
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
				// For entries, update currentEntry with saved data (includes TodoIDs)
				m.currentEntry = msg.entry
				m.savedContent = m.textarea.Value() // Update to current textarea value
			} else if m.view == "add_todo" {
				// For todos, update savedContent to match current input
				m.savedContent = m.todoInput.Value()
			}
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

	case configSavedMsg:
		// Config saved - no need to show status message
		// Theme is already applied in handleThemeSelectorKeys
		if msg.err != nil {
			m.statusMsg = "Error saving theme: " + msg.err.Error()
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
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
		return ui.RenderEntryForm(m.width, m.height, m.currentTheme, m.textarea, m.statusMsg, m.hasUnsaved)
	case "entries":
		return ui.RenderEntryList(m.width, m.height, m.currentTheme, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
	case "view_entry":
		// Calculate filtered/sorted entry position for footer display
		filtered := helpers.FilterEntriesByDateRange(m.entries, m.filterDate)
		filtered = helpers.FilterEntriesByTags(filtered, m.filterTags)
		sorted := helpers.SortEntriesForDisplay(filtered)
		return ui.RenderEntryView(m.width, m.height, m.currentTheme, m.viewingEntry, m.todos, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedEntry, len(sorted))
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.currentTheme, m.displayTodos, m.entries, m.selectedTodo, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
	case "view_todo":
		// Calculate filtered/sorted todo position for footer display
		filtered := helpers.FilterTodosByDateRange(m.displayTodos, m.filterDate)
		filtered = helpers.FilterTodosByTags(filtered, m.filterTags)
		sorted := helpers.SortTodosForDisplay(filtered)
		return ui.RenderTodoView(m.width, m.height, m.currentTheme, m.viewingTodo, m.entries, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedTodo, len(sorted))
	case "unified_filter":
		return ui.RenderUnifiedFilter(m.width, m.height, m.currentTheme, m.unifiedFilterInput, m.availableTags, m.autocompleteTag, m.statusMsg)
	case "add_todo":
		return ui.RenderAddTodoForm(m.width, m.height, m.currentTheme, m.todoInput, m.statusMsg, m.hasUnsaved)
	case "theme_selector":
		return ui.RenderThemeSelector(m.width, m.height, m.currentTheme, m.selectedTheme)
	case "help":
		return ui.RenderHelp(m.width, m.height, m.currentTheme, m.scrollOffset)
	default:
		return ui.RenderEntryList(m.width, m.height, m.currentTheme, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg)
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
