package helpers

import (
	"reflect"
	"testing"

	"github.com/apodacaa/amos/internal/models"
)

func TestExtractTodos(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "single todo",
			text: "Meeting notes\n!todo Follow up with Bob",
			want: []string{"Follow up with Bob"},
		},
		{
			name: "multiple todos",
			text: "Notes\n!todo Task one\nSome text\n!todo Task two",
			want: []string{"Task one", "Task two"},
		},
		{
			name: "no todos",
			text: "Just regular text with no todos",
			want: []string{},
		},
		{
			name: "todo with tags",
			text: "!todo Buy groceries @personal @shopping",
			want: []string{"Buy groceries @personal @shopping"},
		},
		{
			name: "todo not at line start",
			text: "Some text !todo This should not match",
			want: []string{},
		},
		{
			name: "multiple spaces after !todo",
			text: "!todo     Task with extra spaces",
			want: []string{"Task with extra spaces"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTodos(tt.text)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTodos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterTodosByEntry(t *testing.T) {
	entryID1 := "entry-1"
	entryID2 := "entry-2"

	todos := []models.Todo{
		{ID: "1", Title: "Todo 1", EntryID: &entryID1},
		{ID: "2", Title: "Todo 2", EntryID: &entryID2},
		{ID: "3", Title: "Todo 3", EntryID: &entryID1},
		{ID: "4", Title: "Todo 4", EntryID: nil}, // No entry
	}

	tests := []struct {
		name          string
		todos         []models.Todo
		entryID       string
		expectedCount int
		expectedIDs   []string
	}{
		{
			name:          "filter by entry 1",
			todos:         todos,
			entryID:       entryID1,
			expectedCount: 2,
			expectedIDs:   []string{"1", "3"},
		},
		{
			name:          "filter by entry 2",
			todos:         todos,
			entryID:       entryID2,
			expectedCount: 1,
			expectedIDs:   []string{"2"},
		},
		{
			name:          "filter by non-existent entry",
			todos:         todos,
			entryID:       "non-existent",
			expectedCount: 0,
			expectedIDs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterTodosByEntry(tt.todos, tt.entryID)
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d todos, got %d", tt.expectedCount, len(result))
				return
			}
			for i, expectedID := range tt.expectedIDs {
				if result[i].ID != expectedID {
					t.Errorf("Expected todo %d to have ID %q, got %q", i, expectedID, result[i].ID)
				}
			}
		})
	}
}

func TestCountTodoStats(t *testing.T) {
	tests := []struct {
		name          string
		todos         []models.Todo
		expectedOpen  int
		expectedTotal int
	}{
		{
			name: "all open",
			todos: []models.Todo{
				{ID: "1", Status: "open"},
				{ID: "2", Status: "open"},
				{ID: "3", Status: "open"},
			},
			expectedOpen:  3,
			expectedTotal: 3,
		},
		{
			name: "all done",
			todos: []models.Todo{
				{ID: "1", Status: "done"},
				{ID: "2", Status: "done"},
			},
			expectedOpen:  0,
			expectedTotal: 2,
		},
		{
			name: "mixed status",
			todos: []models.Todo{
				{ID: "1", Status: "open"},
				{ID: "2", Status: "done"},
				{ID: "3", Status: "open"},
				{ID: "4", Status: "done"},
			},
			expectedOpen:  2,
			expectedTotal: 4,
		},
		{
			name:          "empty list",
			todos:         []models.Todo{},
			expectedOpen:  0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			open, total := CountTodoStats(tt.todos)
			if open != tt.expectedOpen {
				t.Errorf("Expected %d open todos, got %d", tt.expectedOpen, open)
			}
			if total != tt.expectedTotal {
				t.Errorf("Expected %d total todos, got %d", tt.expectedTotal, total)
			}
		})
	}
}
