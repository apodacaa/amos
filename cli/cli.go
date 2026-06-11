package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	"github.com/google/uuid"
)

// Run handles CLI subcommands. Returns true if a CLI command was handled.
func Run(args []string) bool {
	if len(args) < 2 {
		return false
	}

	switch args[1] {
	case "add":
		return runAdd(args[2:])
	case "list":
		return runList(args[2:])
	case "update":
		return runUpdate(args[2:])
	case "done":
		return runDone(args[2:])
	default:
		return false
	}
}

func runAdd(args []string) bool {
	if len(args) == 0 {
		exitError("usage: amos add <entry|todo> [flags]")
		return true
	}

	switch args[0] {
	case "entry":
		return addEntry(args[1:])
	case "todo":
		return addTodo(args[1:])
	default:
		exitError("unknown add target: %s (expected entry or todo)", args[0])
		return true
	}
}

func addEntry(args []string) bool {
	fs := flag.NewFlagSet("add entry", flag.ContinueOnError)
	title := fs.String("title", "", "entry title (required)")
	body := fs.String("body", "", "entry body")

	if err := fs.Parse(args); err != nil {
		exitError("invalid flags: %v", err)
		return true
	}

	if *title == "" {
		exitError("--title is required")
		return true
	}

	now := time.Now()
	fullText := *title + "\n" + *body
	tags := helpers.ExtractTags(fullText)

	entry := models.Entry{
		ID:        uuid.New().String(),
		Title:     *title,
		Body:      *body,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := storage.SaveEntry(entry); err != nil {
		exitError("failed to save entry: %v", err)
		return true
	}

	// Sync todos from !todo lines in body
	todoTitles := helpers.ExtractTodos(*body)
	if len(todoTitles) > 0 {
		allTodos, err := storage.LoadTodos()
		if err != nil {
			exitError("entry saved but failed to load todos for sync: %v", err)
			return true
		}
		allTodos = helpers.SyncEntryTodos(allTodos, entry.ID, todoTitles)
		if err := storage.SaveTodos(allTodos); err != nil {
			exitError("entry saved but failed to sync todos: %v", err)
			return true
		}
	}

	printJSON(entry)
	return true
}

func addTodo(args []string) bool {
	fs := flag.NewFlagSet("add todo", flag.ContinueOnError)
	title := fs.String("title", "", "todo title (required)")
	status := fs.String("status", "open", "todo status (open, next, done)")

	if err := fs.Parse(args); err != nil {
		exitError("invalid flags: %v", err)
		return true
	}

	if *title == "" {
		exitError("--title is required")
		return true
	}

	if *status != "open" && *status != "next" && *status != "done" {
		exitError("invalid status: %s (expected open, next, or done)", *status)
		return true
	}

	now := time.Now()
	tags := helpers.ExtractTags(*title)

	todo := models.Todo{
		ID:        uuid.New().String(),
		Title:     *title,
		Status:    *status,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := storage.SaveTodo(todo); err != nil {
		exitError("failed to save todo: %v", err)
		return true
	}

	printJSON(todo)
	return true
}

func runUpdate(args []string) bool {
	if len(args) == 0 {
		exitError("usage: amos update todo <id> --status <open|next|done>")
		return true
	}

	switch args[0] {
	case "todo":
		return updateTodo(args[1:])
	default:
		exitError("unknown update target: %s (expected todo)", args[0])
		return true
	}
}

func updateTodo(args []string) bool {
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		exitError("usage: amos update todo <id> --status <open|next|done>")
		return true
	}

	id := args[0]

	fs := flag.NewFlagSet("update todo", flag.ContinueOnError)
	status := fs.String("status", "", "todo status (open, next, done)")

	if err := fs.Parse(args[1:]); err != nil {
		exitError("invalid flags: %v", err)
		return true
	}

	if *status == "" {
		exitError("--status is required")
		return true
	}

	if *status != "open" && *status != "next" && *status != "done" {
		exitError("invalid status: %s (expected open, next, or done)", *status)
		return true
	}

	todo, err := updateTodoStatus(id, *status)
	if err != nil {
		exitError("%v", err)
		return true
	}

	printJSON(todo)
	return true
}

func runDone(args []string) bool {
	if len(args) != 1 || strings.HasPrefix(args[0], "-") {
		exitError("usage: amos done <id>")
		return true
	}

	todo, err := updateTodoStatus(args[0], "done")
	if err != nil {
		exitError("%v", err)
		return true
	}

	printJSON(todo)
	return true
}

// updateTodoStatus sets the status of an existing todo by ID and refreshes
// UpdatedAt. Idempotent: setting a status the todo already has succeeds.
func updateTodoStatus(id, status string) (models.Todo, error) {
	todos, err := storage.LoadTodos()
	if err != nil {
		return models.Todo{}, fmt.Errorf("failed to load todos: %w", err)
	}

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Status = status
			todos[i].UpdatedAt = time.Now()
			if err := storage.SaveTodos(todos); err != nil {
				return models.Todo{}, fmt.Errorf("failed to save todos: %w", err)
			}
			return todos[i], nil
		}
	}

	return models.Todo{}, fmt.Errorf("todo not found: %s", id)
}

func runList(args []string) bool {
	if len(args) == 0 {
		exitError("usage: amos list <entries|todos> [flags]")
		return true
	}

	switch args[0] {
	case "entries":
		return listEntries(args[1:])
	case "todos":
		return listTodos(args[1:])
	default:
		exitError("unknown list target: %s (expected entries or todos)", args[0])
		return true
	}
}

func listEntries(args []string) bool {
	fs := flag.NewFlagSet("list entries", flag.ContinueOnError)
	tag := fs.String("tag", "", "filter by tag (without @ prefix)")

	if err := fs.Parse(args); err != nil {
		exitError("invalid flags: %v", err)
		return true
	}

	entries, err := storage.LoadEntries()
	if err != nil {
		exitError("failed to load entries: %v", err)
		return true
	}

	if *tag != "" {
		tagName := strings.TrimPrefix(*tag, "@")
		var filtered []models.Entry
		for _, e := range entries {
			if hasTag(e.Tags, tagName) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	printJSON(entries)
	return true
}

func listTodos(args []string) bool {
	fs := flag.NewFlagSet("list todos", flag.ContinueOnError)
	tag := fs.String("tag", "", "filter by tag (without @ prefix)")
	status := fs.String("status", "", "filter by status (open, next, done)")

	if err := fs.Parse(args); err != nil {
		exitError("invalid flags: %v", err)
		return true
	}

	todos, err := storage.LoadTodos()
	if err != nil {
		exitError("failed to load todos: %v", err)
		return true
	}

	if *tag != "" {
		tagName := strings.TrimPrefix(*tag, "@")
		var filtered []models.Todo
		for _, t := range todos {
			if hasTag(t.Tags, tagName) {
				filtered = append(filtered, t)
			}
		}
		todos = filtered
	}

	if *status != "" {
		var filtered []models.Todo
		for _, t := range todos {
			if t.Status == *status {
				filtered = append(filtered, t)
			}
		}
		todos = filtered
	}

	printJSON(todos)
	return true
}

func hasTag(tags []string, target string) bool {
	target = strings.ToLower(target)
	for _, t := range tags {
		if strings.ToLower(t) == target {
			return true
		}
	}
	return false
}

func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		exitError("failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(data))
}

func exitError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	errJSON, _ := json.Marshal(map[string]string{"error": msg})
	fmt.Fprintln(os.Stderr, string(errJSON))
	os.Exit(1)
}
