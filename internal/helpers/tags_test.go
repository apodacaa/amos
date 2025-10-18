package helpers

import (
	"reflect"
	"sort"
	"testing"
)

func TestParseEntryContent(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantTitle string
		wantBody  string
	}{
		{
			name:      "title and body",
			content:   "My Title\nThis is the body\nWith multiple lines",
			wantTitle: "My Title",
			wantBody:  "This is the body\nWith multiple lines",
		},
		{
			name:      "title only",
			content:   "Just a title",
			wantTitle: "Just a title",
			wantBody:  "",
		},
		{
			name:      "empty content",
			content:   "",
			wantTitle: "",
			wantBody:  "",
		},
		{
			name:      "whitespace trimming",
			content:   "  Title with spaces  \n  Body with spaces  ",
			wantTitle: "Title with spaces",
			wantBody:  "Body with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotBody := ParseEntryContent(tt.content)
			if gotTitle != tt.wantTitle {
				t.Errorf("ParseEntryContent() title = %v, want %v", gotTitle, tt.wantTitle)
			}
			if gotBody != tt.wantBody {
				t.Errorf("ParseEntryContent() body = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestExtractTags(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "single tag",
			text: "Meeting with @john about project",
			want: []string{"john"},
		},
		{
			name: "multiple tags",
			text: "Discussion about @project-alpha with @team and @john",
			want: []string{"project-alpha", "team", "john"},
		},
		{
			name: "case insensitive",
			text: "@Work @WORK @work should all be same",
			want: []string{"work"},
		},
		{
			name: "no tags",
			text: "This has no tags at all",
			want: []string{},
		},
		{
			name: "tags with numbers",
			text: "@q1-2024 @sprint23",
			want: []string{"q1-2024", "sprint23"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTags(tt.text)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
