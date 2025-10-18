package helpers

import (
	"regexp"
	"strings"
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
