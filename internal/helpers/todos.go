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

// RemoveTodoMarkup removes all lines starting with "!todo" from text
// Returns cleaned text with !todo lines removed
func RemoveTodoMarkup(text string) string {
	// Match !todo followed by text until end of line (including the newline)
	re := regexp.MustCompile(`(?m)^!todo\s+.+$\n?`)
	cleaned := re.ReplaceAllString(text, "")

	// Trim excess blank lines (don't leave gaps - max 2 newlines = 1 blank line)
	cleaned = regexp.MustCompile(`\n{3,}`).ReplaceAllString(cleaned, "\n\n")

	return strings.TrimSpace(cleaned)
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
