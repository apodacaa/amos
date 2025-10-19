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
			name: "open before done",
			todos: []models.Todo{
				{ID: "1", Title: "Done task", Status: "done", Position: 0, CreatedAt: now},
				{ID: "2", Title: "Open task", Status: "open", Position: 0, CreatedAt: now},
			},
			expected: []string{"Open task", "Done task"},
		},
		{
			name: "sort by position within same status",
			todos: []models.Todo{
				{ID: "1", Title: "Position 2", Status: "open", Position: 2, CreatedAt: now},
				{ID: "2", Title: "Position 0", Status: "open", Position: 0, CreatedAt: now},
				{ID: "3", Title: "Position 1", Status: "open", Position: 1, CreatedAt: now},
			},
			expected: []string{"Position 0", "Position 1", "Position 2"},
		},
		{
			name: "sort by created date when position is same",
			todos: []models.Todo{
				{ID: "1", Title: "Older", Status: "open", Position: 0, CreatedAt: now.Add(-time.Hour)},
				{ID: "2", Title: "Newer", Status: "open", Position: 0, CreatedAt: now},
			},
			expected: []string{"Newer", "Older"},
		},
		{
			name: "complex sort: status -> position -> date",
			todos: []models.Todo{
				{ID: "1", Title: "Done P1", Status: "done", Position: 1, CreatedAt: now},
				{ID: "2", Title: "Open P2 Old", Status: "open", Position: 2, CreatedAt: now.Add(-time.Hour)},
				{ID: "3", Title: "Open P1", Status: "open", Position: 1, CreatedAt: now},
				{ID: "4", Title: "Open P2 New", Status: "open", Position: 2, CreatedAt: now},
				{ID: "5", Title: "Done P0", Status: "done", Position: 0, CreatedAt: now},
			},
			expected: []string{"Open P1", "Open P2 New", "Open P2 Old", "Done P0", "Done P1"},
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

func TestNormalizeTodoPositions(t *testing.T) {
	tests := []struct {
		name              string
		todos             []models.Todo
		expectedPositions []int
	}{
		{
			name: "normalize sequential positions",
			todos: []models.Todo{
				{ID: "1", Title: "First", Position: 0},
				{ID: "2", Title: "Second", Position: 1},
				{ID: "3", Title: "Third", Position: 2},
			},
			expectedPositions: []int{0, 1, 2},
		},
		{
			name: "normalize non-sequential positions",
			todos: []models.Todo{
				{ID: "1", Title: "First", Position: 10},
				{ID: "2", Title: "Second", Position: 25},
				{ID: "3", Title: "Third", Position: 100},
			},
			expectedPositions: []int{0, 1, 2},
		},
		{
			name: "normalize all zeros",
			todos: []models.Todo{
				{ID: "1", Title: "First", Position: 0},
				{ID: "2", Title: "Second", Position: 0},
				{ID: "3", Title: "Third", Position: 0},
			},
			expectedPositions: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := NormalizeTodoPositions(tt.todos)

			if len(normalized) != len(tt.expectedPositions) {
				t.Fatalf("Expected %d todos, got %d", len(tt.expectedPositions), len(normalized))
			}

			for i, expectedPos := range tt.expectedPositions {
				if normalized[i].Position != expectedPos {
					t.Errorf("Todo %d: expected position %d, got %d", i, expectedPos, normalized[i].Position)
				}
			}
		})
	}
}

func TestMoveTodoLogic(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		todos          []models.Todo
		selectedIdx    int
		direction      string
		expectedOrder  []string
		expectNoChange bool
	}{
		{
			name: "move up in middle",
			todos: []models.Todo{
				{ID: "1", Title: "First", Status: "open", Position: 0, CreatedAt: now},
				{ID: "2", Title: "Second", Status: "open", Position: 1, CreatedAt: now},
				{ID: "3", Title: "Third", Status: "open", Position: 2, CreatedAt: now},
			},
			selectedIdx:   1, // "Second"
			direction:     "up",
			expectedOrder: []string{"Second", "First", "Third"},
		},
		{
			name: "move down in middle",
			todos: []models.Todo{
				{ID: "1", Title: "First", Status: "open", Position: 0, CreatedAt: now},
				{ID: "2", Title: "Second", Status: "open", Position: 1, CreatedAt: now},
				{ID: "3", Title: "Third", Status: "open", Position: 2, CreatedAt: now},
			},
			selectedIdx:   1, // "Second"
			direction:     "down",
			expectedOrder: []string{"First", "Third", "Second"},
		},
		{
			name: "move up at top (no change)",
			todos: []models.Todo{
				{ID: "1", Title: "First", Status: "open", Position: 0, CreatedAt: now},
				{ID: "2", Title: "Second", Status: "open", Position: 1, CreatedAt: now},
			},
			selectedIdx:    0, // "First"
			direction:      "up",
			expectedOrder:  []string{"First", "Second"},
			expectNoChange: true,
		},
		{
			name: "move down at bottom (no change)",
			todos: []models.Todo{
				{ID: "1", Title: "First", Status: "open", Position: 0, CreatedAt: now},
				{ID: "2", Title: "Second", Status: "open", Position: 1, CreatedAt: now},
			},
			selectedIdx:    1, // "Second"
			direction:      "down",
			expectedOrder:  []string{"First", "Second"},
			expectNoChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate moveTodo logic
			sorted := SortTodosForDisplay(tt.todos)
			sorted = NormalizeTodoPositions(sorted)

			if tt.selectedIdx < 0 || tt.selectedIdx >= len(sorted) {
				t.Fatal("Invalid selected index")
			}

			// Find target index
			targetIdx := tt.selectedIdx
			if tt.direction == "up" {
				targetIdx--
			} else {
				targetIdx++
			}

			// Check bounds
			if targetIdx < 0 || targetIdx >= len(sorted) {
				if !tt.expectNoChange {
					t.Error("Expected change but target out of bounds")
				}
				// No change expected
				for i, expectedTitle := range tt.expectedOrder {
					if sorted[i].Title != expectedTitle {
						t.Errorf("Position %d: expected %q, got %q", i, expectedTitle, sorted[i].Title)
					}
				}
				return
			}

			// Swap positions
			sorted[tt.selectedIdx].Position, sorted[targetIdx].Position = sorted[targetIdx].Position, sorted[tt.selectedIdx].Position

			// Re-sort to see final order
			finalSorted := SortTodosForDisplay(sorted)

			// Verify order
			for i, expectedTitle := range tt.expectedOrder {
				if finalSorted[i].Title != expectedTitle {
					t.Errorf("Position %d: expected %q, got %q", i, expectedTitle, finalSorted[i].Title)
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
