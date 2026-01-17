package ui

import "strings"

// RenderSymbolsSectionText returns the plain text for the symbols section
func RenderSymbolsSectionText() string {
	title := "SYMBOLS"

	items := []string{
		"Entry List Markers:",
		"  +  Entry has linked todos (color shows highest priority)",
		"     Green = next, Magenta = open, Dim = done",
		"  D  Entry marked for deletion",
		"",
		"Todo List Markers:",
		"  +  Todo is linked to an entry",
		"  D  Todo marked for deletion",
		"",
		"Todo Status Indicators:",
		"  [ ] Open todo",
		"  [>] Next todo (high priority)",
		"  [x] Done todo",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

// RenderCommandsSectionText returns the plain text for the commands section
func RenderCommandsSectionText() string {
	title := "COMMANDS"

	items := []string{
		"Global:",
		"  n:new         Create new entry",
		"  a:todo        Add standalone todo",
		"  e:entries     Jump to entry list",
		"  t:todos       Jump to todo list",
		"  s:theme       Select theme",
		"  ?:help        Show this help page",
		"  q:quit        Quit application",
		"",
		"Entry List:",
		"  j/k:nav       Navigate up/down",
		"  enter:view    View entry detail",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Entry View:",
		"  j/k:nav       Navigate between entries",
		"  f/b:page      Page forward/backward in long entries",
		"  i:edit        Edit current entry",
		"  d:del         Mark/unmark for deletion",
		"  esc:back      Return to entry list",
		"",
		"Todo List:",
		"  j/k:nav       Navigate up/down",
		"  space:toggle  Cycle todo status (open → next → done)",
		"  i:edit        Edit selected todo",
		"  d:del         Mark/unmark for deletion",
		"  $:exec        Execute deletion (after marking)",
		"  /:filter      Filter by tags/dates",
		"  c:clear       Clear active filter",
		"",
		"Entry Form:",
		"  ctrl+s:save   Save work (stays in form)",
		"  esc:exit      Finalize and exit (extracts todos)",
		"",
		"Todo Form:",
		"  enter:save    Save and close",
		"  esc:cancel    Cancel without saving",
		"",
		"Filter:",
		"  @tag          Filter by tag (tab to autocomplete)",
		"  today         Show items from today",
		"  yesterday     Show items from yesterday",
		"  last N days   Show items from last N days",
		"  YYYY-MM-DD    Show items from specific date",
		"  date to date  Show items in date range",
		"  enter:apply   Apply filter",
		"  esc:cancel    Cancel filter",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

// RenderNavigationSectionText returns the plain text for the navigation section
func RenderNavigationSectionText() string {
	title := "NAVIGATION"

	items := []string{
		"Entry list and todo list are peer views - use e/t to jump between them.",
		"Global shortcuts (n, a, s, ?) work from any read-only view.",
		"Press esc to exit forms and return to their natural home view.",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}

// RenderDataStorageSectionText returns the plain text for the data storage section
func RenderDataStorageSectionText() string {
	title := "DATA STORAGE"

	items := []string{
		"Data files stored in plain JSON (default: ~/.amos/):",
		"  entries.json   Journal entries",
		"  todos.json     Todo items",
		"",
		"Config stored in ~/.config/amos/settings.json:",
		"  data_dir       Custom data directory path",
		"  theme          Theme preference (brutalist, cyberpunk)",
		"",
		"Customize data directory for syncing with cloud storage:",
		"",
		"Option 1: Environment variable (one-time or in shell profile)",
		"  AMOS_DATA_DIR=~/Google\\ Drive/amos amos",
		"",
		"Option 2: Config file (persistent, create with any text editor)",
		"  {\"data_dir\": \"~/Google Drive/amos\", \"theme\": \"cyberpunk\"}",
		"",
		"Priority: AMOS_DATA_DIR > settings.json > ~/.amos",
		"Path formats: ~/path (tilde), /absolute/path, or ./relative/path",
	}

	return title + "\n\n" + strings.Join(items, "\n")
}
