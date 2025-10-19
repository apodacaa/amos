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
