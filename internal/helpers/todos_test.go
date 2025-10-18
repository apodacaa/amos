package helpers

import (
	"reflect"
	"testing"
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
