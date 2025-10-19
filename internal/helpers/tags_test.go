package helpers

import (
	"reflect"
	"sort"
	"testing"

	"github.com/apodacaa/amos/internal/models"
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

func TestExtractUniqueTags(t *testing.T) {
	tests := []struct {
		name    string
		entries []models.Entry
		want    []string
	}{
		{
			name: "unique tags from multiple entries",
			entries: []models.Entry{
				{ID: "1", Tags: []string{"work", "client"}},
				{ID: "2", Tags: []string{"personal", "work"}},
				{ID: "3", Tags: []string{"client", "meeting"}},
			},
			want: []string{"@client", "@meeting", "@personal", "@work"},
		},
		{
			name: "no duplicate tags",
			entries: []models.Entry{
				{ID: "1", Tags: []string{"alpha"}},
				{ID: "2", Tags: []string{"beta"}},
				{ID: "3", Tags: []string{"gamma"}},
			},
			want: []string{"@alpha", "@beta", "@gamma"},
		},
		{
			name:    "empty entries",
			entries: []models.Entry{},
			want:    []string{},
		},
		{
			name: "entries with no tags",
			entries: []models.Entry{
				{ID: "1", Tags: []string{}},
				{ID: "2", Tags: []string{}},
			},
			want: []string{},
		},
		{
			name: "all same tag",
			entries: []models.Entry{
				{ID: "1", Tags: []string{"project"}},
				{ID: "2", Tags: []string{"project"}},
				{ID: "3", Tags: []string{"project"}},
			},
			want: []string{"@project"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractUniqueTags(tt.entries)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractUniqueTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterEntriesByTag(t *testing.T) {
	entries := []models.Entry{
		{ID: "1", Title: "Entry 1", Tags: []string{"work", "client"}},
		{ID: "2", Title: "Entry 2", Tags: []string{"personal"}},
		{ID: "3", Title: "Entry 3", Tags: []string{"work", "meeting"}},
		{ID: "4", Title: "Entry 4", Tags: []string{"client"}},
	}

	tests := []struct {
		name       string
		entries    []models.Entry
		filterTag  string
		wantCount  int
		wantTitles []string
	}{
		{
			name:       "filter by work tag",
			entries:    entries,
			filterTag:  "@work",
			wantCount:  2,
			wantTitles: []string{"Entry 1", "Entry 3"},
		},
		{
			name:       "filter by client tag",
			entries:    entries,
			filterTag:  "@client",
			wantCount:  2,
			wantTitles: []string{"Entry 1", "Entry 4"},
		},
		{
			name:       "filter by personal tag",
			entries:    entries,
			filterTag:  "@personal",
			wantCount:  1,
			wantTitles: []string{"Entry 2"},
		},
		{
			name:       "filter by non-existent tag",
			entries:    entries,
			filterTag:  "@nonexistent",
			wantCount:  0,
			wantTitles: []string{},
		},
		{
			name:       "empty filter returns all",
			entries:    entries,
			filterTag:  "",
			wantCount:  4,
			wantTitles: []string{"Entry 1", "Entry 2", "Entry 3", "Entry 4"},
		},
		{
			name:       "filter without @ prefix (still works)",
			entries:    entries,
			filterTag:  "work",
			wantCount:  2,
			wantTitles: []string{"Entry 1", "Entry 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterEntriesByTag(tt.entries, tt.filterTag)
			if len(got) != tt.wantCount {
				t.Errorf("FilterEntriesByTag() returned %d entries, want %d", len(got), tt.wantCount)
				return
			}
			for i, wantTitle := range tt.wantTitles {
				if got[i].Title != wantTitle {
					t.Errorf("Entry %d: got title %q, want %q", i, got[i].Title, wantTitle)
				}
			}
		})
	}
}
