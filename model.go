package main

import (
	"fmt"
	"os"
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
	view         string // Current view: "entries", "view_entry", "todos", "view_todo", "unified_filter", "add_todo", "theme_selector", "workspace_selector", or "help"
	previousView string // Previous view (for returning from modals like theme_selector and help)

	// ============================================================
	// Terminal Dimensions
	// ============================================================
	width  int // Terminal width (updated on window resize)
	height int // Terminal height (updated on window resize)

	// ============================================================
	// Text Input Components (Bubble Tea textareas)
	// ============================================================
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
	sysConfig     storage.SystemConfig // System configuration (loaded from ~/.config/amos/settings.json)
	currentTheme  ui.Theme             // Currently active theme (brutalist or cyberpunk)
	selectedTheme int                  // Selected theme index in theme selector modal

	// ============================================================
	// Workspace State
	// ============================================================
	selectedWorkspace int // Selected workspace index in workspace selector modal

	// ============================================================
	// Update Notification State
	// ============================================================
	updateAvailable bool   // Whether newer version is available
	latestVersion   string // Latest version string (e.g., "v1.5.0")
	updateCheckDone bool   // Whether update check completed this session
	updateDismissed bool   // Whether user dismissed notice with 'u' key
}

// NewModel creates a new model with default values
func NewModel() Model {
	// Create single-line input for standalone todos
	// Default width for textareas (will be updated on WindowSizeMsg)
	defaultInputWidth := 76 // 80 - 4 margin

	todoInput := textarea.New()
	todoInput.Placeholder = "Todo title..."
	todoInput.CharLimit = 0
	todoInput.SetWidth(defaultInputWidth)
	todoInput.SetHeight(1) // Single line
	ui.ApplyTextareaStyle(&todoInput)

	// Create single-line input for unified filtering (tags + dates)
	unifiedFilterInput := textarea.New()
	unifiedFilterInput.Placeholder = "e.g. @work yesterday, last 30 days @client"
	unifiedFilterInput.CharLimit = 0
	unifiedFilterInput.SetWidth(defaultInputWidth)
	unifiedFilterInput.SetHeight(1) // Single line
	ui.ApplyTextareaStyle(&unifiedFilterInput)

	// Load system config from storage (or use default)
	sysConfig, _ := storage.LoadSystemConfig()

	// Load theme based on config
	currentTheme := ui.GetThemeByName(sysConfig.Theme)

	return Model{
		view:               "entries",
		width:              80, // Default width
		height:             24, // Default height
		todoInput:          todoInput,
		unifiedFilterInput: unifiedFilterInput,
		markedForDeletion:  make(map[string]string),
		sysConfig:          sysConfig,
		currentTheme:       currentTheme,
		selectedTheme:      0, // Will be set when opening theme selector
	}
}

// Init initializes the model (Elm architecture)
func (m Model) Init() tea.Cmd {
	// Load entries and todos on startup (single source of truth: Todo.EntryID)
	// Also check for updates asynchronously (non-blocking)
	return tea.Batch(m.loadEntriesAndTodos(), checkForUpdatesCmd(Version))
}

// Update handles messages (Elm architecture)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Route to appropriate key handler based on view
		switch m.view {
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
		case "workspace_selector":
			return m.handleWorkspaceSelectorKeys(msg)
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
		// Update textarea widths to use available space (leave margin for padding)
		inputWidth := msg.Width - 4
		if inputWidth < 20 {
			inputWidth = 20
		}
		// Save current values before resizing (SetWidth can lose content)
		todoValue := m.todoInput.Value()
		filterValue := m.unifiedFilterInput.Value()
		// Set new widths
		m.todoInput.SetWidth(inputWidth)
		m.unifiedFilterInput.SetWidth(inputWidth)
		// Restore values to force proper re-render
		m.todoInput.SetValue(todoValue)
		m.unifiedFilterInput.SetValue(filterValue)
		return m, nil

	case saveCompleteMsg:
		if msg.err != nil {
			m.statusMsg = "Error saving: " + msg.err.Error()
		} else {
			m.statusMsg = "Saved"
			// Mark as saved
			m.hasUnsaved = false
			if m.view == "add_todo" {
				// For todos, update savedContent to match current input
				m.savedContent = m.todoInput.Value()
			}
		}
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()

	case editorFinishedMsg:
		// Clean up temp file
		if msg.tempFile != "" {
			defer os.Remove(msg.tempFile)
		}

		// Handle editor launch errors
		if msg.err != nil {
			m.statusMsg = "Editor error: " + msg.err.Error()
			m.statusTime = time.Now()
			m.editingMode = false
			m.originalEntryID = ""
			m.view = "entries"
			return m, clearStatusAfterDelay()
		}

		// Parse the edited file
		content, isEmpty, err := helpers.ParseEntryFile(msg.tempFile)
		if err != nil {
			m.statusMsg = "Error reading file: " + err.Error()
			m.statusTime = time.Now()
			m.editingMode = false
			m.originalEntryID = ""
			m.view = "entries"
			return m, clearStatusAfterDelay()
		}

		// If empty, user cancelled - return to appropriate view
		if isEmpty {
			wasEditing := m.editingMode
			m.editingMode = false
			m.originalEntryID = ""
			if wasEditing {
				m.view = "view_entry"
			} else {
				m.view = "entries"
			}
			return m, nil
		}

		// Parse title and body from content
		title, body := helpers.ParseEntryContent(content)

		// Extract todos from body
		todoTitles := helpers.ExtractTodos(body)

		// Extract tags from title and body (includes tags from !todo lines)
		tags := helpers.ExtractTags(title + " " + body)

		// Update entry (body preserved with !todo lines intact)
		m.currentEntry.Title = title
		m.currentEntry.Body = body
		m.currentEntry.Tags = tags
		m.currentEntry.UpdatedAt = time.Now()

		// Save entry
		if err := storage.SaveEntry(m.currentEntry); err != nil {
			m.statusMsg = "Error saving: " + err.Error()
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}

		// Sync todos: dedup existing, orphan removed, create new
		allTodos, err := storage.LoadTodos()
		if err != nil {
			m.statusMsg = "Error loading todos: " + err.Error()
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		allTodos = helpers.SyncEntryTodos(allTodos, m.currentEntry.ID, todoTitles)
		if err := storage.SaveTodos(allTodos); err != nil {
			m.statusMsg = "Error saving todos: " + err.Error()
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}

		// Update viewingEntry for view_entry
		m.viewingEntry = m.currentEntry

		// Exit to appropriate view
		wasEditing := m.editingMode
		m.editingMode = false
		m.originalEntryID = ""

		if wasEditing {
			m.view = "view_entry"
		} else {
			m.view = "entries"
		}

		// Load todos to show in next view
		return m, m.loadEntriesAndTodos()

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

	case workspaceSwitchedMsg:
		if msg.err != nil {
			m.statusMsg = "Error switching workspace: " + msg.err.Error()
			m.statusTime = time.Now()
			return m, clearStatusAfterDelay()
		}
		// Load new data
		m.entries = msg.entries
		m.todos = msg.todos
		m.displayTodos = helpers.SortTodosForDisplay(m.todos)
		// Reset selection and clear state
		m.selectedEntry = 0
		m.selectedTodo = 0
		m.filterTags = []string{}
		m.filterDate = ""
		m.markedForDeletion = make(map[string]string)
		m.deleteConfirmPending = false
		m.statusMsg = "Switched to " + msg.name
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()

	case updateCheckCompleteMsg:
		// Update check completed - mark as done and store result
		m.updateCheckDone = true
		m.updateAvailable = msg.updateAvailable
		m.latestVersion = msg.latestVersion
		// No status message - update notice shows in footer if available
		return m, nil
	}

	return m, cmd
}

// View renders the UI (Elm architecture)
func (m Model) View() string {
	switch m.view {
	case "entries":
		return ui.RenderEntryList(m.width, m.height, m.currentTheme, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg, m.updateAvailable, m.updateDismissed, m.latestVersion, m.sysConfig.GetActiveWorkspaceName())
	case "view_entry":
		// Calculate filtered/sorted entry position for footer display
		filtered := helpers.ApplyEntryFilters(m.entries, m.filterDate, m.filterTags)
		sorted := helpers.SortEntriesForDisplay(filtered)
		return ui.RenderEntryView(m.width, m.height, m.currentTheme, m.viewingEntry, m.todos, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedEntry, len(sorted), m.updateAvailable, m.updateDismissed, m.latestVersion)
	case "todos":
		return ui.RenderTodoList(m.width, m.height, m.currentTheme, m.displayTodos, m.entries, m.selectedTodo, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg, m.updateAvailable, m.updateDismissed, m.latestVersion, m.sysConfig.GetActiveWorkspaceName())
	case "view_todo":
		// Calculate filtered/sorted todo position for footer display
		filtered := helpers.ApplyTodoFilters(m.displayTodos, m.filterDate, m.filterTags)
		sorted := helpers.SortTodosForDisplay(filtered)
		return ui.RenderTodoView(m.width, m.height, m.currentTheme, m.viewingTodo, m.entries, m.scrollOffset, m.markedForDeletion, m.statusMsg, m.selectedTodo, len(sorted), m.updateAvailable, m.updateDismissed, m.latestVersion)
	case "unified_filter":
		return ui.RenderUnifiedFilter(m.width, m.height, m.currentTheme, m.unifiedFilterInput, m.availableTags, m.autocompleteTag, m.statusMsg)
	case "add_todo":
		return ui.RenderAddTodoForm(m.width, m.height, m.currentTheme, m.todoInput, m.statusMsg, m.hasUnsaved)
	case "theme_selector":
		return ui.RenderThemeSelector(m.width, m.height, m.currentTheme, m.selectedTheme)
	case "workspace_selector":
		return ui.RenderWorkspaceSelector(m.width, m.height, m.currentTheme, m.sysConfig.Workspaces, m.selectedWorkspace, m.sysConfig.ActiveWorkspace)
	case "help":
		return ui.RenderHelp(m.width, m.height, m.currentTheme, m.scrollOffset, m.updateAvailable, m.updateDismissed, m.latestVersion)
	default:
		return ui.RenderEntryList(m.width, m.height, m.currentTheme, m.entries, m.selectedEntry, m.todos, m.filterTags, m.filterDate, m.markedForDeletion, m.statusMsg, m.updateAvailable, m.updateDismissed, m.latestVersion, m.sysConfig.GetActiveWorkspaceName())
	}
}

// handleNewEntry is a shared handler for creating a new entry (from any view)
func (m Model) handleNewEntry() (Model, tea.Cmd) {
	now := time.Now()
	m.currentEntry = models.Entry{
		ID:        m.generateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.editingMode = false
	return m, launchEditor("")
}

// handleAddTodo is a shared handler for creating a standalone todo (from any view)
func (m Model) handleAddTodo() (Model, tea.Cmd) {
	m.view = "add_todo"
	now := time.Now()
	m.currentTodo = models.Todo{
		ID:        m.generateID(),
		Status:    "open",
		CreatedAt: now,
		UpdatedAt: now,
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

// openWorkspaceSelector opens the workspace selector modal
func (m Model) openWorkspaceSelector() (Model, tea.Cmd) {
	if len(m.sysConfig.Workspaces) == 0 {
		m.statusMsg = "No workspaces configured in ~/.config/amos/settings.json"
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()
	}
	// Check if AMOS_DATA_DIR overrides workspace switching
	if os.Getenv("AMOS_DATA_DIR") != "" {
		m.statusMsg = "Data directory overridden by AMOS_DATA_DIR"
		m.statusTime = time.Now()
		return m, clearStatusAfterDelay()
	}
	m.previousView = m.view
	m.view = "workspace_selector"
	// Set selected workspace to current active workspace
	for i, ws := range m.sysConfig.Workspaces {
		if ws.Name == m.sysConfig.ActiveWorkspace {
			m.selectedWorkspace = i
			break
		}
	}
	return m, nil
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
