package helpers

import (
	"regexp"
	"strings"

	"github.com/apodacaa/amos/internal/models"
)

// ExtractTodos finds all !todo items in text and returns their titles
func ExtractTodos(text string) []string {
	// Match !todo followed by text until end of line
	re := regexp.MustCompile(`(?m)^!todo\s+(.+)$`)
	matches := re.FindAllStringSubmatch(text, -1)

	todos := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			todos = append(todos, strings.TrimSpace(match[1]))
		}
	}

	return todos
}

// FilterTodosByEntry returns todos that belong to the specified entry
func FilterTodosByEntry(todos []models.Todo, entryID string) []models.Todo {
	filtered := []models.Todo{}
	for _, todo := range todos {
		if todo.EntryID != nil && *todo.EntryID == entryID {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

// CountTodoStats returns the count of open and total todos
func CountTodoStats(todos []models.Todo) (open int, total int) {
	total = len(todos)
	for _, todo := range todos {
		if todo.Status == "open" {
			open++
		}
	}
	return open, total
}

// IsTodoLinked returns true if the todo is linked to an entry
func IsTodoLinked(todo models.Todo) bool {
	return todo.EntryID != nil
}

// GetEntryMarker returns a priority-based marker for an entry based on its linked todos
// Priority: ">" (next) > "*" (open) > "=" (done) > "" (no todos)
func GetEntryMarker(todos []models.Todo, entryID string) string {
	// Filter todos for this entry
	entryTodos := FilterTodosByEntry(todos, entryID)

	// No todos = no marker
	if len(entryTodos) == 0 {
		return ""
	}

	// Check for "next" status (highest priority)
	hasNext := false
	hasOpen := false
	hasOther := false

	for _, todo := range entryTodos {
		switch todo.Status {
		case "next":
			hasNext = true
		case "open":
			hasOpen = true
		default: // "done" or anything else
			hasOther = true
		}
	}

	// Priority order: next > open > done
	if hasNext {
		return ">"
	}
	if hasOpen {
		return "*"
	}
	if hasOther {
		return "="
	}

	// Fallback (shouldn't reach here)
	return ""
}

// RepairEntryTodoRelationships rebuilds the TodoIDs array for each entry
// based on todos that have entry_id pointing to the entry.
// This ensures the bidirectional Entry ↔ Todo relationship is in sync.
func RepairEntryTodoRelationships(entries []models.Entry, todos []models.Todo) []models.Entry {
	// Create a map of entry ID to todo IDs for efficient lookup
	entryToTodos := make(map[string][]string)

	// Scan all todos and build the relationship map
	for _, todo := range todos {
		if todo.EntryID != nil {
			entryID := *todo.EntryID
			entryToTodos[entryID] = append(entryToTodos[entryID], todo.ID)
		}
	}

	// Update each entry's TodoIDs array
	repairedEntries := make([]models.Entry, len(entries))
	for i, entry := range entries {
		repairedEntries[i] = entry
		// Set TodoIDs to the todos that reference this entry
		if todoIDs, exists := entryToTodos[entry.ID]; exists {
			repairedEntries[i].TodoIDs = todoIDs
		} else {
			// No todos reference this entry, ensure TodoIDs is empty (not nil)
			repairedEntries[i].TodoIDs = []string{}
		}
	}

	return repairedEntries
}
