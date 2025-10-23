package helpers

import (
	"testing"
	"time"

	"github.com/apodacaa/amos/internal/models"
)

func TestSortTodosForDisplay(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		todos    []models.Todo
		expected []string // Expected order of titles
	}{
		{
			name: "next before open before done",
			todos: []models.Todo{
				{ID: "1", Title: "Done task", Status: "done", CreatedAt: now},
				{ID: "2", Title: "Open task", Status: "open", CreatedAt: now},
				{ID: "3", Title: "Next task", Status: "next", CreatedAt: now},
			},
			expected: []string{"Next task", "Open task", "Done task"},
		},
		{
			name: "sort by created date within same status",
			todos: []models.Todo{
				{ID: "1", Title: "Older", Status: "open", CreatedAt: now.Add(-time.Hour)},
				{ID: "2", Title: "Newer", Status: "open", CreatedAt: now},
			},
			expected: []string{"Newer", "Older"},
		},
		{
			name: "complex sort: next -> open -> done, newest first within each",
			todos: []models.Todo{
				{ID: "1", Title: "Done New", Status: "done", CreatedAt: now},
				{ID: "2", Title: "Open Old", Status: "open", CreatedAt: now.Add(-time.Hour)},
				{ID: "3", Title: "Next Old", Status: "next", CreatedAt: now.Add(-time.Hour)},
				{ID: "4", Title: "Open New", Status: "open", CreatedAt: now},
				{ID: "5", Title: "Done Old", Status: "done", CreatedAt: now.Add(-2 * time.Hour)},
				{ID: "6", Title: "Next New", Status: "next", CreatedAt: now},
			},
			expected: []string{"Next New", "Next Old", "Open New", "Open Old", "Done New", "Done Old"},
		},
		{
			name: "all same status sorted by date",
			todos: []models.Todo{
				{ID: "1", Title: "Third", Status: "open", CreatedAt: now.Add(-2 * time.Hour)},
				{ID: "2", Title: "First", Status: "open", CreatedAt: now},
				{ID: "3", Title: "Second", Status: "open", CreatedAt: now.Add(-time.Hour)},
			},
			expected: []string{"First", "Second", "Third"},
		},
		{
			name: "unknown status treated as open",
			todos: []models.Todo{
				{ID: "1", Title: "Done", Status: "done", CreatedAt: now},
				{ID: "2", Title: "Unknown", Status: "unknown", CreatedAt: now},
				{ID: "3", Title: "Next", Status: "next", CreatedAt: now},
			},
			expected: []string{"Next", "Unknown", "Done"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := SortTodosForDisplay(tt.todos)

			if len(sorted) != len(tt.expected) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expected), len(sorted))
			}

			for i, expectedTitle := range tt.expected {
				if sorted[i].Title != expectedTitle {
					t.Errorf("Position %d: expected %q, got %q", i, expectedTitle, sorted[i].Title)
				}
			}
		})
	}
}

func TestSortEntriesForDisplay(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		entries       []models.Entry
		expectedOrder []string // Expected order of titles
	}{
		{
			name: "sort by timestamp descending",
			entries: []models.Entry{
				{ID: "1", Title: "Oldest", Timestamp: now.Add(-2 * time.Hour)},
				{ID: "2", Title: "Newest", Timestamp: now},
				{ID: "3", Title: "Middle", Timestamp: now.Add(-1 * time.Hour)},
			},
			expectedOrder: []string{"Newest", "Middle", "Oldest"},
		},
		{
			name: "already sorted",
			entries: []models.Entry{
				{ID: "1", Title: "First", Timestamp: now},
				{ID: "2", Title: "Second", Timestamp: now.Add(-1 * time.Hour)},
				{ID: "3", Title: "Third", Timestamp: now.Add(-2 * time.Hour)},
			},
			expectedOrder: []string{"First", "Second", "Third"},
		},
		{
			name: "reverse order",
			entries: []models.Entry{
				{ID: "1", Title: "Third", Timestamp: now.Add(-2 * time.Hour)},
				{ID: "2", Title: "Second", Timestamp: now.Add(-1 * time.Hour)},
				{ID: "3", Title: "First", Timestamp: now},
			},
			expectedOrder: []string{"First", "Second", "Third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := SortEntriesForDisplay(tt.entries)

			if len(sorted) != len(tt.expectedOrder) {
				t.Fatalf("Expected %d entries, got %d", len(tt.expectedOrder), len(sorted))
			}

			for i, expectedTitle := range tt.expectedOrder {
				if sorted[i].Title != expectedTitle {
					t.Errorf("Position %d: expected %q, got %q", i, expectedTitle, sorted[i].Title)
				}
			}
		})
	}
}
