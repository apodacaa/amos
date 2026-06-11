package cli

import (
	"testing"
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
)

func setupTestDir(t *testing.T) {
	t.Helper()
	if err := storage.SetDataDir(t.TempDir()); err != nil {
		t.Fatalf("failed to set data dir: %v", err)
	}
}

func seedTodo(t *testing.T, todo models.Todo) {
	t.Helper()
	if err := storage.SaveTodo(todo); err != nil {
		t.Fatalf("failed to seed todo: %v", err)
	}
}

func TestUpdateTodoStatusDone(t *testing.T) {
	setupTestDir(t)

	created := time.Now().Add(-time.Hour)
	seedTodo(t, models.Todo{
		ID:        "todo-1",
		Title:     "Test todo",
		Status:    "open",
		CreatedAt: created,
		UpdatedAt: created,
	})

	updated, err := updateTodoStatus("todo-1", "done")
	if err != nil {
		t.Fatalf("updateTodoStatus failed: %v", err)
	}

	if updated.Status != "done" {
		t.Errorf("expected status done, got %s", updated.Status)
	}
	if !updated.UpdatedAt.After(created) {
		t.Errorf("expected UpdatedAt to be refreshed, got %v", updated.UpdatedAt)
	}

	// Verify persisted to store
	todos, err := storage.LoadTodos()
	if err != nil {
		t.Fatalf("failed to load todos: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].Status != "done" {
		t.Errorf("expected persisted status done, got %s", todos[0].Status)
	}
}

func TestUpdateTodoStatusUnknownID(t *testing.T) {
	setupTestDir(t)

	seedTodo(t, models.Todo{
		ID:        "todo-1",
		Title:     "Test todo",
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_, err := updateTodoStatus("nonexistent-id", "done")
	if err == nil {
		t.Fatal("expected error for unknown id, got nil")
	}

	// Store must be untouched
	todos, err := storage.LoadTodos()
	if err != nil {
		t.Fatalf("failed to load todos: %v", err)
	}
	if todos[0].Status != "open" {
		t.Errorf("expected status unchanged (open), got %s", todos[0].Status)
	}
}

func TestUpdateTodoStatusIdempotentRedone(t *testing.T) {
	setupTestDir(t)

	seedTodo(t, models.Todo{
		ID:        "todo-1",
		Title:     "Test todo",
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if _, err := updateTodoStatus("todo-1", "done"); err != nil {
		t.Fatalf("first done failed: %v", err)
	}

	updated, err := updateTodoStatus("todo-1", "done")
	if err != nil {
		t.Fatalf("re-done of already-done todo failed: %v", err)
	}
	if updated.Status != "done" {
		t.Errorf("expected status done, got %s", updated.Status)
	}

	// No duplicates created
	todos, err := storage.LoadTodos()
	if err != nil {
		t.Fatalf("failed to load todos: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 todo after re-done, got %d", len(todos))
	}
}

func TestUpdateTodoStatusOtherStatuses(t *testing.T) {
	setupTestDir(t)

	seedTodo(t, models.Todo{
		ID:        "todo-1",
		Title:     "Test todo",
		Status:    "done",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	for _, status := range []string{"next", "open", "done"} {
		updated, err := updateTodoStatus("todo-1", status)
		if err != nil {
			t.Fatalf("update to %s failed: %v", status, err)
		}
		if updated.Status != status {
			t.Errorf("expected status %s, got %s", status, updated.Status)
		}
	}
}
