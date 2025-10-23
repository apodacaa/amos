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

### Running Single Tests

```bash
# Run specific test function
go test -v ./internal/helpers -run TestExtractTags

# Run tests in specific package
go test -v ./internal/storage

# Run with coverage for specific package
go test -cover ./internal/helpers
```

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
commands.go              # All tea.Cmd functions (side effects: save, load, etc.)
update_*.go              # Key handlers per view (domain separation)
  update_dashboard.go    # Dashboard navigation
  update_entry.go        # Entry form (create/edit)
  update_entries.go      # Entry list navigation
  update_entry_view.go   # Read-only entry view
  update_todos.go        # Todo list (toggle, reorder)
  update_tag_picker.go   # Tag filtering
  update_add_todo.go     # Standalone todo form
ui/                      # Pure view renderers
  dashboard.go
  entry_form.go
  entry_list.go
  entry_view.go
  todo_list.go
  tag_picker.go
  add_todo_form.go
  styles.go              # Brutalist styling (monochrome)
internal/
  models/                # Data structures
    entry.go             # Entry{ID, Title, Body, Tags, Timestamp, TodoIDs}
    todo.go              # Todo{ID, Title, Status, Tags, CreatedAt, EntryID, Position}
  storage/               # JSON persistence (~/.amos/)
    storage.go           # Load/Save functions for entries.json and todos.json
  helpers/               # Reusable business logic
    sorting.go           # Centralized sorting (todos by status→position→date)
    tags.go              # Tag extraction (@mention syntax)
    todos.go             # Todo extraction (!todo syntax)
```

### Key Architectural Patterns

**State Management**:
- All state lives in the `Model` struct (model.go:17-39)
- View routing via `m.view` string field ("dashboard", "entry", "entries", "view_entry", "todos", "tag_picker", "add_todo")
- No hidden previousView tracking - explicit navigation only

**Message Flow**:
1. User input → tea.KeyMsg
2. Update() routes to view-specific handler (update_*.go)
3. Handler returns (Model, tea.Cmd)
4. Command executes async operation
5. Async result → custom message type (saveCompleteMsg, etc.)
6. Update() handles result message, updates model

**Data Persistence**:
- JSON files in `~/.amos/` directory
- `storage.LoadEntries()` / `storage.SaveEntry()` for entries
- `storage.LoadTodos()` / `storage.SaveTodo()` for todos
- Save operations happen via commands (async), results via messages

**Todo System**:
- Standalone todos: `EntryID` is nil
- Entry-linked todos: `EntryID` points to parent entry
- Position field enables manual priority (lower = higher)
- Sorting: open todos first → by position → newest first (helpers/sorting.go)
- Extract from entries with `!todo Task description @tag` syntax

**Tag System**:
- Auto-extracted from `@mention` syntax in entry body
- Stored in `Tags` array on both Entry and Todo
- Tag filtering via `@` key (shows tag picker)

## Brutalist Design Philosophy

The app follows strict brutalist principles:

**Navigation**:
- `esc` - Universal back to dashboard from all views
- `n` - New entry (works from any read-only view)
- `a` - Add standalone todo (works from any read-only view)
- `e` - Jump to entries list
- `t` - Jump to todos list
- No hidden navigation state - explicit view switching only

**Visual Design**:
- **Dashboard**: Monument aesthetic with massive centered ASCII "AMOS" title
- **Other views**: Honest workspaces with left-aligned help text anchored to bottom
- **Monochrome palette**: Pure black/white/gray (no colors)
- **No decorations**: No italics, no Unicode bullets (use `>` not `►`), just ASCII
- Help text uses inverted black/white for maximum contrast

**Data Integrity**:
- Append-only journal (no delete feature for entries)
- Immediate writes: `space` to toggle todo saves immediately
- Full context: todos visible in entry view
- No deferred/pending state

**Viewport & Scrolling**:
- Lists (entries, todos) use viewport windowing: show 20-30 items with scroll indicators
- Entry view: long entries scrollable with `u` (down) / `i` (up) keys
- Minimum terminal size: 80x24 (shows resize message if too small)

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

### Working with Textarea

The app uses `charmbracelet/bubbles/textarea`:
- Entry form: multi-line textarea (model.go:43-58)
- Todo form: single-line textarea with height=1 (model.go:60-73)
- Always call `textarea.Blink` when focusing
- Always call `m.textarea.Update(msg)` when in text entry view

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

## Data Format

Entries stored in `~/.amos/entries.json`:
```json
[
  {
    "id": "uuid",
    "title": "Entry title",
    "body": "Entry content with @tags and !todo items",
    "tags": ["work", "personal"],
    "timestamp": "2025-01-01T12:00:00Z",
    "todo_ids": ["todo-uuid-1", "todo-uuid-2"]
  }
]
```

Todos stored in `~/.amos/todos.json`:
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

- Entries are append-only (no delete operation by design)
- Todo status: "open" or "done"
- Position normalization: when reordering todos, all positions are renumbered sequentially
- Tag syntax: `@tagname` in entry body auto-extracts to Tags array
- Todo syntax: `!todo Task description @tag` creates linked todo with extracted tags
