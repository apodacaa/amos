package ui

import "strings"

// helpSection defines a section of help content
type helpSection struct {
	Title string
	Body  string
}

// helpSections defines all help page sections in order
var helpSections = []helpSection{
	{
		Title: "SYMBOLS",
		Body: `Entry List Markers:
  +  Entry has linked todos (color shows highest priority)
     Green = next, Magenta = open, Dim = done
  D  Entry marked for deletion

Todo List Markers:
  +  Todo is linked to an entry
  D  Todo marked for deletion

Todo Status Indicators:
  [ ] Open todo
  [>] Next todo (high priority)
  [x] Done todo`,
	},
	{
		Title: "ENTRY EDITING",
		Body: `Entries are edited with your system's default text editor,
following Unix conventions (like git, mutt, aerc).

Editor selection (in order of priority):
  $EDITOR    Primary editor (e.g., vim, nano, code)
  $VISUAL    Fallback for visual editors
  nano       Default if neither is set

Entry format:
  First line becomes the title
  Use @tags anywhere for organization
  Use !todo lines to create linked tasks

Save and exit the editor to create/update the entry.
Exit without saving (or leave empty) to cancel.`,
	},
	{
		Title: "DATA STORAGE",
		Body: `Data files stored in plain JSON (default: ~/.amos/):
  entries.json   Journal entries
  todos.json     Todo items

Config stored in ~/.config/amos/settings.json:
  data_dir       Custom data directory path
  theme          Theme preference (brutalist, cyberpunk)

Customize data directory for syncing with cloud storage:

Option 1: Environment variable (one-time or in shell profile)
  AMOS_DATA_DIR=~/Google\ Drive/amos amos

Option 2: Config file (persistent, create with any text editor)
  {"data_dir": "~/Google Drive/amos", "theme": "cyberpunk"}

Priority: AMOS_DATA_DIR > settings.json > ~/.amos
Path formats: ~/path (tilde), /absolute/path, or ./relative/path`,
	},
}

// GetHelpContent returns the full help content with all sections assembled
func GetHelpContent() string {
	var sections []string
	for _, s := range helpSections {
		sections = append(sections, s.Title+"\n\n"+s.Body)
	}
	return strings.Join(sections, "\n\n")
}
