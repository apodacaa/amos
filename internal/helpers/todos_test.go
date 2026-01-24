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

func TestSyncEntryTodos(t *testing.T) {
	entryID := "entry-1"
	otherEntryID := "entry-2"

	tests := []struct {
		name           string
		allTodos       []models.Todo
		entryID        string
		todoTitles     []string
		wantLinked     int // expected todos linked to entryID after sync
		wantOrphaned   int // expected todos orphaned (EntryID set to nil) from this entry
		wantTotalNew   int // expected new todos created
		checkOrphanIDs []string
	}{
		{
			name:         "new entry with new todos",
			allTodos:     []models.Todo{},
			entryID:      entryID,
			todoTitles:   []string{"Task one", "Task two"},
			wantLinked:   2,
			wantTotalNew: 2,
		},
		{
			name: "re-save with same todos (no duplicates)",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task one", Status: "open", EntryID: &entryID},
				{ID: "t2", Title: "Task two", Status: "next", EntryID: &entryID},
			},
			entryID:      entryID,
			todoTitles:   []string{"Task one", "Task two"},
			wantLinked:   2,
			wantTotalNew: 0,
		},
		{
			name: "removed todo line orphans existing todo",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task one", Status: "open", EntryID: &entryID},
				{ID: "t2", Title: "Task two", Status: "done", EntryID: &entryID},
			},
			entryID:        entryID,
			todoTitles:     []string{"Task one"},
			wantLinked:     1,
			wantOrphaned:   1,
			checkOrphanIDs: []string{"t2"},
		},
		{
			name: "added todo line creates only the new one",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task one", Status: "open", EntryID: &entryID},
			},
			entryID:      entryID,
			todoTitles:   []string{"Task one", "Task three"},
			wantLinked:   2,
			wantTotalNew: 1,
		},
		{
			name: "empty todoTitles orphans all linked todos",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task one", Status: "open", EntryID: &entryID},
				{ID: "t2", Title: "Task two", Status: "next", EntryID: &entryID},
			},
			entryID:        entryID,
			todoTitles:     []string{},
			wantLinked:     0,
			wantOrphaned:   2,
			checkOrphanIDs: []string{"t1", "t2"},
		},
		{
			name: "same title but different entryID is not deduped",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task one", Status: "open", EntryID: &otherEntryID},
			},
			entryID:      entryID,
			todoTitles:   []string{"Task one"},
			wantLinked:   1,
			wantTotalNew: 1,
		},
		{
			name: "orphaned todo preserves status and tags",
			allTodos: []models.Todo{
				{ID: "t1", Title: "Task @work", Status: "next", Tags: []string{"work"}, EntryID: &entryID},
			},
			entryID:        entryID,
			todoTitles:     []string{},
			wantOrphaned:   1,
			checkOrphanIDs: []string{"t1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Deep copy allTodos to avoid mutation across tests
			input := make([]models.Todo, len(tt.allTodos))
			for i, todo := range tt.allTodos {
				input[i] = todo
				if todo.EntryID != nil {
					id := *todo.EntryID
					input[i].EntryID = &id
				}
			}

			result := SyncEntryTodos(input, tt.entryID, tt.todoTitles)

			// Count linked todos for this entry
			linked := 0
			for _, todo := range result {
				if todo.EntryID != nil && *todo.EntryID == tt.entryID {
					linked++
				}
			}
			if linked != tt.wantLinked {
				t.Errorf("linked todos = %d, want %d", linked, tt.wantLinked)
			}

			// Check orphaned todos
			if len(tt.checkOrphanIDs) > 0 {
				orphaned := 0
				for _, todo := range result {
					for _, orphanID := range tt.checkOrphanIDs {
						if todo.ID == orphanID && todo.EntryID == nil {
							orphaned++
						}
					}
				}
				if orphaned != tt.wantOrphaned {
					t.Errorf("orphaned todos = %d, want %d", orphaned, tt.wantOrphaned)
				}
			}

			// Check new todos created
			if tt.wantTotalNew > 0 {
				newCount := len(result) - len(tt.allTodos)
				if newCount != tt.wantTotalNew {
					t.Errorf("new todos created = %d, want %d", newCount, tt.wantTotalNew)
				}
			}

			// Check orphaned todo preserves fields
			if tt.name == "orphaned todo preserves status and tags" {
				for _, todo := range result {
					if todo.ID == "t1" {
						if todo.Status != "next" {
							t.Errorf("orphaned todo status = %q, want %q", todo.Status, "next")
						}
						if len(todo.Tags) != 1 || todo.Tags[0] != "work" {
							t.Errorf("orphaned todo tags = %v, want [work]", todo.Tags)
						}
					}
				}
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

func TestGetEntryMarker(t *testing.T) {
	entryID1 := "entry-1"
	entryID2 := "entry-2"

	tests := []struct {
		name     string
		todos    []models.Todo
		entryID  string
		expected string
	}{
		{
			name:     "no todos",
			todos:    []models.Todo{},
			entryID:  entryID1,
			expected: "",
		},
		{
			name: "entry with no linked todos",
			todos: []models.Todo{
				{ID: "1", Status: "open", EntryID: &entryID2},
			},
			entryID:  entryID1,
			expected: "",
		},
		{
			name: "entry with next todo (highest priority)",
			todos: []models.Todo{
				{ID: "1", Status: "next", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "next",
		},
		{
			name: "entry with open todo",
			todos: []models.Todo{
				{ID: "1", Status: "open", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "open",
		},
		{
			name: "entry with done todo",
			todos: []models.Todo{
				{ID: "1", Status: "done", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "done",
		},
		{
			name: "entry with next overrides open",
			todos: []models.Todo{
				{ID: "1", Status: "open", EntryID: &entryID1},
				{ID: "2", Status: "next", EntryID: &entryID1},
				{ID: "3", Status: "done", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "next",
		},
		{
			name: "entry with open overrides done",
			todos: []models.Todo{
				{ID: "1", Status: "done", EntryID: &entryID1},
				{ID: "2", Status: "open", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "open",
		},
		{
			name: "entry with multiple done todos",
			todos: []models.Todo{
				{ID: "1", Status: "done", EntryID: &entryID1},
				{ID: "2", Status: "done", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "done",
		},
		{
			name: "mixed entries - only count correct entry",
			todos: []models.Todo{
				{ID: "1", Status: "next", EntryID: &entryID2},
				{ID: "2", Status: "open", EntryID: &entryID1},
			},
			entryID:  entryID1,
			expected: "open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEntryMarker(tt.todos, tt.entryID)
			if result != tt.expected {
				t.Errorf("GetEntryMarker() = %q, want %q", result, tt.expected)
			}
		})
	}
}
