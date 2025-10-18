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

// saveCompleteMsg is sent when save operation completes
type saveCompleteMsg struct {
	err error
}

// entriesLoadedMsg is sent when entries are loaded
type entriesLoadedMsg struct {
	entries []models.Entry
	err     error
}

// entryDeletedMsg is sent when entry is deleted
type entryDeletedMsg struct {
	err error
}

// todosLoadedMsg is sent when todos are loaded
type todosLoadedMsg struct {
	todos []models.Todo
	err   error
}

// todoToggledMsg is sent when todo status is toggled
type todoToggledMsg struct {
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
	case "e":
		// View entries list
		m.view = "entries"
		m.selectedEntry = 0
		return m, m.loadEntries()
	case "t":
		// View todos list
		m.view = "todos"
		m.selectedTodo = 0
		return m, m.loadTodos()
	case "esc":
		m.view = "dashboard"
	}
	return m, nil
}

// handleEntryKeys processes keyboard input (entry view)
func (m Model) handleEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		// Only Ctrl+C quits from entry form, not 'q' (user might type words with 'q')
		return m, tea.Quit
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

// handleEntriesListKeys processes keyboard input (entries list view)
func (m Model) handleEntriesListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If confirming delete, only allow 'd' to proceed or anything else to cancel
	if m.confirmingDelete {
		switch msg.String() {
		case "d":
			// User pressed 'd' again - proceed with delete
			return m, m.deleteEntry()
		default:
			// Any other key cancels delete confirmation
			m.confirmingDelete = false
			m.statusMsg = ""
			return m, nil
		}
	}

	// Normal navigation (not confirming delete)
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.view = "dashboard"
		return m, nil
	case "j", "down":
		if m.selectedEntry < len(m.entries)-1 {
			m.selectedEntry++
		}
		return m, nil
	case "k", "up":
		if m.selectedEntry > 0 {
			m.selectedEntry--
		}
		return m, nil
	case "d":
		// First 'd' press - show confirmation
		m.confirmingDelete = true
		m.statusMsg = "⚠ Delete entry? Press 'd' again to confirm, or any other key to cancel"
		return m, nil
	case "enter":
		// Open selected entry for read-only viewing
		if m.selectedEntry >= 0 && m.selectedEntry < len(m.entries) {
			// Need to get the sorted entry (newest first)
			sorted := make([]models.Entry, len(m.entries))
			copy(sorted, m.entries)
			// Sort by timestamp descending (newest first) - same as in RenderEntryList
			for i := 0; i < len(sorted)-1; i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[j].Timestamp.After(sorted[i].Timestamp) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
			m.viewingEntry = sorted[m.selectedEntry]
			m.view = "view_entry"
		}
		return m, nil
	}
	return m, nil
}

// handleViewEntryKeys processes keyboard input (view entry - read-only)
func (m Model) handleViewEntryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to entries list
		m.view = "entries"
		return m, nil
	}
	return m, nil
}

// handleTodosListKeys processes keyboard input (todos list view)
func (m Model) handleTodosListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.view = "dashboard"
		m.statusMsg = ""
		return m, nil
	case "j", "down":
		if m.selectedTodo < len(m.todos)-1 {
			m.selectedTodo++
		}
		return m, nil
	case "k", "up":
		if m.selectedTodo > 0 {
			m.selectedTodo--
		}
		return m, nil
	case " ":
		// Toggle todo status
		return m, m.toggleTodo()
	}
	return m, nil
}

// loadEntries loads all entries from storage
func (m Model) loadEntries() tea.Cmd {
	return func() tea.Msg {
		entries, err := storage.LoadEntries()
		return entriesLoadedMsg{entries: entries, err: err}
	}
}

// loadTodos loads all todos from storage (async)
func (m Model) loadTodos() tea.Cmd {
	return func() tea.Msg {
		todos, err := storage.LoadTodos()
		return todosLoadedMsg{todos: todos, err: err}
	}
}

// toggleTodo toggles the status of the currently selected todo
func (m Model) toggleTodo() tea.Cmd {
	return func() tea.Msg {
		// Get the sorted todo to toggle (same sorting as display)
		sorted := make([]models.Todo, len(m.todos))
		copy(sorted, m.todos)
		// Sort: open before done, then newest first within each group
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Status == "done" && sorted[j].Status == "open" {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].Status == sorted[j].Status {
					if sorted[j].CreatedAt.After(sorted[i].CreatedAt) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
		}

		if m.selectedTodo < 0 || m.selectedTodo >= len(sorted) {
			return todoToggledMsg{err: nil}
		}

		todoToToggle := sorted[m.selectedTodo]

		// Toggle status
		if todoToToggle.Status == "open" {
			todoToToggle.Status = "done"
		} else {
			todoToToggle.Status = "open"
		}

		// Save the updated todo
		err := storage.SaveTodo(todoToToggle)
		return todoToggledMsg{err: err}
	}
}

// deleteEntry deletes the currently selected entry
func (m Model) deleteEntry() tea.Cmd {
	return func() tea.Msg {
		// Get the sorted entry to delete
		sorted := make([]models.Entry, len(m.entries))
		copy(sorted, m.entries)
		// Sort by timestamp descending (newest first)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j].Timestamp.After(sorted[i].Timestamp) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}

		if m.selectedEntry < 0 || m.selectedEntry >= len(sorted) {
			return entryDeletedMsg{err: nil}
		}

		entryToDelete := sorted[m.selectedEntry]

		// Load all entries
		entries, err := storage.LoadEntries()
		if err != nil {
			return entryDeletedMsg{err: err}
		}

		// Find and remove the entry
		newEntries := make([]models.Entry, 0, len(entries)-1)
		for _, e := range entries {
			if e.ID != entryToDelete.ID {
				newEntries = append(newEntries, e)
			}
		}

		// Save updated entries
		err = storage.SaveEntries(newEntries)
		return entryDeletedMsg{err: err}
	}
}

// saveEntry saves the current entry and extracts todos
func (m Model) saveEntry() tea.Cmd {
	return func() tea.Msg {
		content := m.textarea.Value()

		// Parse content into title and body
		title, body := helpers.ParseEntryContent(content)

		// Extract tags from title and body
		tags := helpers.ExtractTags(title + " " + body)

		// Extract todos from content
		todoTitles := helpers.ExtractTodos(content)

		// Create todo IDs list
		todoIDs := make([]string, 0, len(todoTitles))

		// Create and save todos
		for _, todoTitle := range todoTitles {
			todo := models.Todo{
				ID:        uuid.New().String(),
				Title:     todoTitle,
				Status:    "open",
				Tags:      helpers.ExtractTags(todoTitle), // Extract tags from todo title
				CreatedAt: time.Now(),
				EntryID:   &m.currentEntry.ID, // Link to this entry
			}

			// Save each todo
			if err := storage.SaveTodo(todo); err != nil {
				return saveCompleteMsg{err: err}
			}

			todoIDs = append(todoIDs, todo.ID)
		}

		// Update current entry
		m.currentEntry.Title = title
		m.currentEntry.Body = body
		m.currentEntry.Tags = tags
		m.currentEntry.TodoIDs = todoIDs
		m.currentEntry.Timestamp = time.Now()

		// Save entry to storage
		err := storage.SaveEntry(m.currentEntry)

		return saveCompleteMsg{err: err}
	}
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
