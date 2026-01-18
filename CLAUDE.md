# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Amos is a minimal Bubble Tea (Go) TUI for journal + todo management with a brutalist design philosophy.

## Development Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the app |
| `make ci` | **Run before commit** - Full checks + tests (fmt, vet, staticcheck, test) |
| `make test` | Run tests |
| `make test-v` | Run tests with verbose output |
| `make check` | Quick check (fmt + vet only) |
| `make build` | Build binary to `./amos` |
| `make install-air` | Install air for hot reload (then run `air`) |

**Important**: Always run `make ci` before committing. The git pre-commit hook (installed via `./scripts/install-hooks.sh`) runs this automatically.

**Important**: After committing, update CLAUDE.md if the change affects architecture, config fields, or user-facing behavior.

### Running Single Tests

```bash
# Run specific test function
go test -v ./internal/helpers -run TestExtractTags

# Run tests in specific package
go test -v ./internal/storage

# Run with coverage for specific package
go test -cover ./internal/helpers
```

## Efficiency Guidelines for Claude Code

This codebase is small (~3000 lines) and well-organized. Work efficiently:

### 1. Trust the Documentation
- File organization below tells you exactly where everything is
- Go DIRECTLY to files - don't explore or search
- Use the Quick Reference Map below for common tasks

### 2. Use Parallel Reads
- Read multiple files at once in a single message
- Example: Read model.go, update_entries.go, and ui/entry_list.go together
- Don't read files sequentially

### 3. Use Grep Instead of Full Reads
- Finding usage: `grep -n "filterTags" -C 3` shows context
- Finding definitions: `grep -n "func.*Filter"`
- Only use Read for files you need to edit

### 4. Never Re-read Files
- Remember what you've read in the conversation
- Trust your memory of the code structure
- Only re-read if you made edits and need to see results

### Quick Reference Map

**Unified Filtering (Tags + Dates):**
- Input handler: `update_unified_filter.go`
- UI renderer: `ui/unified_filter.go`
- Filter parsing: `internal/helpers/filter_parser.go` (ParseFilterInput)
- Tag logic: `internal/helpers/tags.go` (FilterEntriesByTags, FilterTodosByTags)
- Date logic: `internal/helpers/dates.go` (ParseDateFilter, FilterEntriesByDateRange, FilterTodosByDateRange)
- Entry list integration: `update_entries.go`, `ui/entry_list.go`
- Todo list integration: `update_todos.go`, `ui/todo_list.go`
- Context tracking: `model.go` (filterContext field determines return view)
- Filter workflow: `/` opens filter with current values for editing, `c` clears filter. After applying filter, status message shows "Filter applied. Press c to clear"

**Entry Management:**
- Entry list: `update_entries.go`, `ui/entry_list.go`
- Entry view: `update_entry_view.go`, `ui/entry_view.go` (press `i` to edit)
- External editor: `internal/helpers/editor.go` (GetEditor, CreateTempFile, ParseEntryFile)
- Edit workflow: Press `n` to create or `i` from entry view to edit → launches $EDITOR → on save/exit, extracts todos and returns to TUI

**Todo Management:**
- Todo list: `update_todos.go`, `ui/todo_list.go`
- Todo view: `update_view_todo.go`, `ui/todo_view.go` (press `i` to edit)
- Add/Edit todo: `update_add_todo.go`, `ui/add_todo_form.go` (create and edit)
- Edit workflow: Todo list → `enter` → todo view → `i` key → todo form (editing mode) → save or esc
- Todo logic: `internal/helpers/todos.go`, `internal/helpers/sorting.go`

**Theme Selector:**
- Input handler: `update_theme_selector.go`
- UI renderer: `ui/theme_selector.go`
- Theme definitions: `ui/styles.go` (BrutalistTheme, CyberpunkTheme)
- Config storage: `internal/storage/system_config.go` (theme stored in SystemConfig)
- Accessed via `s` key from entry view and todo list

**Help Page:**
- Input handler: `update_help.go`
- UI renderer: `ui/help.go`
- Accessible via `?` from all read-only views
- Page navigation: `f/b` keys for page forward/backward
- Documents symbols, commands, and navigation

**Core State:**
- Model struct: `model.go:17-42` (includes filterContext field)
- Update routing: `model.go:97-129` (switch on view)
- View routing: `model.go:224-239` (switch on view)
- Message types: `messages.go`
- Commands: `commands.go`
- Status messages: Forms display statusMsg above help text (saved, errors, etc.)

**Common Edit Pattern:**
1. Identify 2-3 files from map above
2. Read them in parallel (one message)
3. Make edits
4. Run `make ci`
5. Done (no re-reading unless errors)

## Architecture

### Bubble Tea (Elm Architecture)

The app follows Bubble Tea's Elm architecture pattern:

- **Model** (`model.go`) - All application state in a single struct
- **Init()** - Initializes model and returns startup commands
- **Update(msg)** - Handles messages, returns updated model + commands
- **View()** - Pure function that renders UI from model state

**Critical Rules**:
- Update() must NOT have side effects - return tea.Cmd instead
- View functions are pure - no state mutation
- Commands (tea.Cmd) handle all async/IO operations

### File Organization

The codebase uses domain-based separation for maintainability:

```
main.go                  # Entry point only (~10 lines)
model.go                 # Model struct + Init/Update/View (Elm core)
messages.go              # All message types (saveCompleteMsg, entriesLoadedMsg, etc.)
commands.go              # All tea.Cmd functions (side effects: save, load, launchEditor, etc.)
update_*.go              # Key handlers per view (domain separation)
  update_entries.go      # Entry list navigation
  update_entry_view.go   # Read-only entry view (press 'i' to launch editor)
  update_todos.go        # Todo list (toggle, reorder)
  update_unified_filter.go  # Unified filtering (tags + dates)
  update_add_todo.go     # Standalone todo form
  update_theme_selector.go  # Theme selection modal
  update_help.go         # Help page navigation
ui/                      # Pure view renderers
  entry_list.go
  entry_view.go
  todo_list.go
  unified_filter.go      # Unified filter view
  add_todo_form.go
  theme_selector.go      # Theme selection modal
  help.go                # Help page documentation
  styles.go              # Theme definitions and styling
internal/
  models/                # Data structures
    entry.go             # Entry{ID, Title, Body, Tags, Timestamp}
    todo.go              # Todo{ID, Title, Status, Tags, CreatedAt, EntryID, Position}
  storage/               # JSON persistence (configurable directory)
    storage.go           # Load/Save functions, directory initialization
    system_config.go     # SystemConfig{DataDir, Theme, HideCompletedTodos} (~/.config/amos/settings.json)
  helpers/               # Reusable business logic
    sorting.go           # Centralized sorting (todos by status→position→date)
    tags.go              # Tag extraction (@mention syntax) and filtering
    todos.go             # Todo extraction (!todo syntax)
    dates.go             # Date filter parsing and date range filtering
    filter.go            # Centralized filtering (ApplyEntryFilters, ApplyTodoFilters, FilterTodosByVisibility)
    filter_parser.go     # Unified filter parsing (tags + dates)
    editor.go            # External editor support ($EDITOR, temp files, parsing)
```

### Key Architectural Patterns

**State Management**:
- All state lives in the `Model` struct (model.go:17-42)
- View routing via `m.view` string field ("entries", "view_entry", "todos", "view_todo", "unified_filter", "add_todo", "theme_selector", "help")
- Filter context via `m.filterContext` ("entries" or "todos") - determines return view from filter
- Theme selector and help page use `m.previousView` to return to calling view
- App opens to "entries" view (entry list as default)
- Entry editing uses external $EDITOR (no "entry" view state - editor runs outside TUI)

**Message Flow**:
1. User input → tea.KeyMsg
2. Update() routes to view-specific handler (update_*.go)
3. Handler returns (Model, tea.Cmd)
4. Command executes async operation
5. Async result → custom message type (saveCompleteMsg, etc.)
6. Update() handles result message, updates model

**Data Persistence**:
- JSON files in data directory (default: `~/.amos/`)
- Customizable via `AMOS_DATA_DIR` env var or `~/.config/amos/settings.json`
- All config: `~/.config/amos/settings.json` (stores data directory path, theme, and hide_completed_todos preference)
- User data files: `entries.json`, `todos.json` (in configured data directory)
- `storage.LoadEntries()` / `storage.SaveEntry()` for entries
- `storage.LoadTodos()` / `storage.SaveTodo()` for todos
- `storage.LoadSystemConfig()` / `storage.SaveSystemConfig()` for settings (data_dir, theme)
- `storage.InitializeDataDir()` called at startup to configure directory
- Auto-creates directory on first run
- Supports tilde (`~`) expansion in paths
- Save operations happen via commands (async), results via messages

**Todo System**:
- Standalone todos: `EntryID` is nil
- Entry-linked todos: `EntryID` points to parent entry (single source of truth)
- **No Entry.TodoIDs field** - unidirectional relationship via Todo.EntryID only
- Position field enables manual priority (lower = higher)
- Sorting: open todos first → by position → newest first (helpers/sorting.go)
- Extract from entries with `!todo Task description @tag` syntax

**Unified Filter System (Tags + Dates)**:
- Auto-extracted tags from `@mention` syntax in entry body and todo titles
- Stored in `Tags` array on both Entry and Todo
- Unified filtering via `/` key in both entries and todos views
- **Filter workflow**:
  - `/` always opens filter with current values pre-populated (for editing)
  - `c` clears filter entirely (only shown in status message when filter is active)
  - Filter persists across navigation until cleared
  - Status message after applying filter: "Filter applied. Press c to clear" (auto-clears after 3 seconds)
- **Tag filtering**:
  - Type `@tagname` to filter by tags (autocomplete with `tab`)
  - Filter uses AND logic (all tags must match)
  - Autocomplete clears after tab completion
- **Date filtering**:
  - Supports: `today`, `yesterday`, `last N days` (any number), `YYYY-MM-DD`, `YYYY-MM-DD to YYYY-MM-DD`
  - "last N days" includes today (e.g., "last 14 days" = 13 days ago through today)
  - Date ranges are inclusive on both ends
  - No date autocomplete (only tags have autocomplete)
- **Filter validation**:
  - Invalid filters show error message with hint
  - Errors auto-clear after 3 seconds
- **Display**:
  - Filtered titles show as "Entries @tag1 @tag2 last 14 days" or "Todos @work today"
  - Footer uses raw filter string (brutalist: no formatting)

## Brutalist Design Philosophy

The app follows strict brutalist principles:

**Navigation**:
- Entry list and todo list are peer views (no hierarchy)
- All key bindings are shown in view headers
- Theme selector and help page use `m.previousView` to return to calling view

**Visual Design**:
- **All views**: Honest workspaces with left-aligned help text anchored to bottom
- **Theme support**: Two built-in themes accessible via `s` key
  - **Brutalist**: Monochrome (terminal defaults), reverse video for emphasis
  - **Cyberpunk**: Neon colors (Matrix green, Tron blue, magenta, cyan, orange)
- **Entry markers**: Color-coded `+` symbol shows todo priority
  - Green `+` = entry has "next" todos (highest priority)
  - Magenta `+` = entry has "open" todos
  - Dim `+` = entry has only "done" todos
- **No decorations**: No italics, no Unicode bullets, just ASCII
- Help text uses maximum contrast (reverse video for brutalist, colored backgrounds for cyberpunk)
- Help text uses non-breaking spaces (\u00A0) to prevent key/description splitting
- Wrapped helper lines have blank line between them for readability
- Status messages appear in message line below footer (e.g., "saved" in forms, "Filter applied. Press c to clear", "X items marked. Press $ to delete")

**Data Integrity**:
- **neomutt-style deletion**: `d` marks items (toggle), `$` executes deletion (y/n confirmation). `$` key only shown in status message after marking items, not in view headers
- Visual feedback: "D" prefix in lists, "[MARKED FOR DELETION]" in entry view footer, status message shows "X items marked. Press $ to delete"
- Multi-select: marks persist across navigation until deletion or app exit
- Cascade delete: deleting entry removes all linked todos automatically
- Entry-linked todos: can be marked but only deleted via parent entry
- Immediate writes: `o/p/x` keys to set todo status save immediately
- Full context: todos visible in entry view
- No hidden state: all marks visible, confirmation shows counts

**Viewport & Page Navigation**:
- Lists (entries, todos) use viewport windowing: show 20-30 items with scroll indicators
- Entry view, todo view, and help page: use viewport windowing with `f/b` keys for page forward/backward navigation
- No line-by-line scrolling - discrete page jumps only

## Common Patterns

### Adding a New View

1. Add view name to `m.view` routing in `Update()` (model.go:97-112)
2. Add view name to `View()` switch (model.go:230-245)
3. Create `update_newview.go` with key handler function
4. Create `ui/newview.go` with pure render function
5. Add message types to `messages.go` if needed
6. Add commands to `commands.go` if async operations needed

### Adding a New Message Type

1. Define in `messages.go` (e.g., `type myMsg struct { ... }`)
2. Handle in `Update()` switch (model.go:91-204)
3. Create command in `commands.go` that returns the message

### Working with Textarea and External Editor

**External Editor (entries)**:
- Entry creation/editing uses $EDITOR (or $VISUAL, or nano as fallback)
- `launchEditor()` in commands.go creates temp file and launches editor via `tea.ExecProcess`
- `editorFinishedMsg` handler in model.go parses result and saves entry
- Template comments (<!-- ... -->) are stripped from parsed content
- `!todo` lines are extracted as linked todos on save

**Textarea (todos and filters)**:
- Todo form: single-line textarea with height=1 (model.go:37-39)
- Filter input: single-line textarea for tags + dates
- Always call `textarea.Blink` when focusing

## Dependencies

- **bubbletea** v1.3.10 - TUI framework
- **lipgloss** v1.1.0 - Styling
- **bubbles** v0.21.0 - Textarea component
- **uuid** v1.6.0 - ID generation
- **Go 1.24+** required

## Testing

Tests use standard Go testing:
- Test files: `*_test.go`
- All helper functions have tests (internal/helpers/*_test.go)
- Storage operations have tests (internal/storage/storage_test.go)
- **Application logic tests**: `model_test.go` (30 comprehensive tests)
  - Navigation workflows: 13 tests covering all view transitions
  - Filtering workflows: 8 tests for tag/date filtering
  - Deletion workflows: 9 tests for neomutt-style multi-select deletion
  - Coverage: ~60% of application logic (navigation, filtering, deletion flows)
  - All tests use table-driven patterns where appropriate
  - Run with `make test` or `make ci` (recommended before commits)

## Data Format

Entries stored in `<data-dir>/entries.json` (default: `~/.amos/entries.json`):
```json
[
  {
    "id": "uuid",
    "title": "Entry title",
    "body": "Entry content with @tags and !todo items",
    "tags": ["work", "personal"],
    "timestamp": "2025-01-01T12:00:00Z"
  }
]
```

**Note**: Entry.TodoIDs field removed - relationships use Todo.EntryID as single source of truth.

Todos stored in `<data-dir>/todos.json` (default: `~/.amos/todos.json`):
```json
[
  {
    "id": "uuid",
    "title": "Todo title",
    "status": "open",
    "tags": ["work"],
    "created_at": "2025-01-01T12:00:00Z",
    "entry_id": "entry-uuid",
    "position": 0
  }
]
```

## Important Notes

- **Editing**: Press `n` to create new entry, `i` to edit entries and todos
  - Entry editing: Uses external $EDITOR (vim, nano, etc.) following Unix conventions (like aerc, mutt, git)
    - Press `n` from any list view to create new entry → opens editor with template
    - Press `i` from entry view to edit → opens editor with existing content
    - On save/exit: extracts `!todo` lines as linked todos, removes markup, saves entry
    - Empty file (or only template comments) cancels the operation
    - Editor preference: $EDITOR → $VISUAL → nano (fallback)
  - Todo editing: Open todo view, press `i` to edit, save with `Enter`, esc returns to todo view
  - Timestamp updates: Editing an entry updates its timestamp to now (moves to top of list)
  - Todo extraction: `!todo Task @tag` lines in editor become linked todos on save. Todos become independent entities edited via todo view.
- **Deletion**: neomutt-style multi-select pattern
  - `d` key marks/unmarks items for deletion (shows "D" prefix in lists, footer text in entry view)
  - After marking items, status message shows "X items marked. Press $ to delete"
  - `$` key executes deletion (doesn't appear in view headers, only in status message after marking)
  - First `$` press shows y/n confirmation with counts, `y` or `enter` confirms deletion
  - Marks persist across navigation until deletion or app closes
  - Cascade deletion: deleting entry removes all linked todos
  - Entry-linked todos can be marked but only deleted via parent entry (standalone todos deleted independently)
- **Single source of truth**: Todo.EntryID only (no Entry.TodoIDs) - unidirectional relationship prevents sync issues
- **Entry markers**: Color-coded `+` shows todo priority (green=next, magenta=open, dim=done)
- **Page navigation**: `f/b` keys in entry/todo/help views for page forward/backward (no line scrolling)
- **Help page**: Press `?` from any read-only view to see comprehensive documentation
- Todo status: "open", "next", or "done" (direct keys: `o`=open, `p`=next/priority, `x`=done)
- **Toggle completed todos**: `z` key in todo list hides/shows done todos. Preference persists in config. Done todos always visible in entry view (preserves journal history).
- Tag syntax: `@tagname` in entry body or todo title auto-extracts to Tags array
- Todo syntax: `!todo Task description @tag` creates linked todo with extracted tags when exiting entry form (ESC). The `!todo` line is automatically removed from entry text on exit. During editing, `!todo` markup stays visible until ESC. To edit todos, use the dedicated todo view (press `i` from todo list).
- Unified filtering: Works identically for both entries and todos views
  - `/` key opens filter with current values for editing
  - After applying filter, status message shows "Filter applied. Press c to clear"
  - `c` key clears filter (doesn't appear in view headers, only in status message when filter is active)
  - Supports tags (@work), dates (today, yesterday, last N days, YYYY-MM-DD, date ranges)
  - Only tag autocomplete (no date autocomplete)
  - Invalid filters show error with 3-second auto-clear
- Save confirmations: Entry and add_todo forms show "saved" toast message (brutalist: no emoji)
- Status messages: Auto-clear after 3 seconds (saved confirmations, filter cleared, errors)
- Helper text: Uses non-breaking spaces and vertical spacing for better readability on narrow terminals
