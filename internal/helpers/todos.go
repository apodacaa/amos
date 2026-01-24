package helpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/google/uuid"
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

// SyncEntryTodos reconciles the todo list with the current !todo lines in an entry.
// For each todoTitle, if a todo with matching Title AND EntryID already exists, skip (dedup).
// If a todoTitle has no match, create a new todo.
// For each existing todo linked to this entry, if its title does NOT appear in todoTitles, orphan it (set EntryID to nil).
// Returns the updated full todos slice (ready for batch save).
func SyncEntryTodos(allTodos []models.Todo, entryID string, todoTitles []string) []models.Todo {
	// Build set of current !todo titles
	titleSet := make(map[string]bool, len(todoTitles))
	for _, title := range todoTitles {
		titleSet[strings.TrimSpace(title)] = true
	}

	// Track which titles already have matching todos
	matched := make(map[string]bool, len(todoTitles))

	// Iterate existing todos: orphan removed ones, mark matched ones
	for i := range allTodos {
		if allTodos[i].EntryID == nil || *allTodos[i].EntryID != entryID {
			continue
		}
		trimmedTitle := strings.TrimSpace(allTodos[i].Title)
		if titleSet[trimmedTitle] {
			matched[trimmedTitle] = true
		} else {
			// Orphan: title no longer in entry text
			allTodos[i].EntryID = nil
		}
	}

	// Create new todos for unmatched titles
	for _, title := range todoTitles {
		trimmed := strings.TrimSpace(title)
		if !matched[trimmed] {
			now := time.Now()
			entryIDCopy := entryID
			newTodo := models.Todo{
				ID:        uuid.New().String(),
				Title:     trimmed,
				Status:    "open",
				Tags:      ExtractTags(trimmed),
				CreatedAt: now,
				UpdatedAt: now,
				EntryID:   &entryIDCopy,
			}
			allTodos = append(allTodos, newTodo)
		}
	}

	return allTodos
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

// GetEntryMarker returns the highest-priority status for an entry based on its linked todos
// Returns: "next", "open", "done", or "" (no todos)
// Priority: next > open > done
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
		return "next"
	}
	if hasOpen {
		return "open"
	}
	if hasOther {
		return "done"
	}

	// Fallback (shouldn't reach here)
	return ""
}
