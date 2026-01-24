package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
	"github.com/google/uuid"
)

// App struct holds the application context for Wails
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. Context is saved for runtime calls.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize storage using the same mechanism as the TUI.
	// This reads AMOS_DATA_DIR env var or ~/.config/amos/settings.json
	// and sets the package-level data directory.
	if err := storage.InitializeDataDir(); err != nil {
		fmt.Println("Warning: failed to initialize data directory:", err.Error())
	}
}

// GetEntries returns all entries sorted for display (newest first)
func (a *App) GetEntries() ([]models.Entry, error) {
	entries, err := storage.LoadEntries()
	if err != nil {
		return nil, err
	}
	return helpers.SortEntriesForDisplay(entries), nil
}

// GetTodos returns all todos sorted for display (next > open > done)
func (a *App) GetTodos() ([]models.Todo, error) {
	todos, err := storage.LoadTodos()
	if err != nil {
		return nil, err
	}
	return helpers.SortTodosForDisplay(todos), nil
}

// CreateEntry creates a new journal entry
func (a *App) CreateEntry(title, body string) (*models.Entry, error) {
	now := time.Now()
	tags := helpers.ExtractTags(title + "\n" + body)

	entry := models.Entry{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := storage.SaveEntry(entry); err != nil {
		return nil, err
	}

	// Sync todos from !todo lines in body
	todoTitles := helpers.ExtractTodos(body)
	if len(todoTitles) > 0 {
		allTodos, err := storage.LoadTodos()
		if err == nil {
			allTodos = helpers.SyncEntryTodos(allTodos, entry.ID, todoTitles)
			storage.SaveTodos(allTodos)
		}
	}

	return &entry, nil
}

// CreateTodo creates a new standalone todo
func (a *App) CreateTodo(title string) (*models.Todo, error) {
	now := time.Now()
	tags := helpers.ExtractTags(title)

	todo := models.Todo{
		ID:        uuid.New().String(),
		Title:     title,
		Status:    "open",
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := storage.SaveTodo(todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

// UpdateTodoStatus updates a todo's status (open/next/done)
func (a *App) UpdateTodoStatus(id, status string) error {
	todos, err := storage.LoadTodos()
	if err != nil {
		return err
	}

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Status = status
			todos[i].UpdatedAt = time.Now()
			break
		}
	}

	return storage.SaveTodos(todos)
}

// DeleteEntry deletes an entry and its linked todos (cascade)
func (a *App) DeleteEntry(id string) error {
	return storage.DeleteEntryCascade(id)
}

// DeleteTodo deletes a standalone todo
func (a *App) DeleteTodo(id string) error {
	return storage.DeleteTodo(id)
}
