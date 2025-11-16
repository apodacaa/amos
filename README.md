# Amos

Minimal Bubble Tea (Go) TUI for journal + todo management. Brutalist design, fast iteration.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap apodacaa/amos
brew install amos
```

Run the app:
```bash
amos
```

### From Source

Requires Go 1.24+

```bash
git clone https://github.com/apodacaa/amos.git
cd amos
make build
./amos
```

**Keyboard Shortcuts:**

*Entry Form:*
- `Ctrl+S` - Save entry
- `esc` - Cancel

*Entry List:*
- `n` - New Entry
- `a` - Add Standalone Todo
- `j/k` or `↑/↓` - Navigate
- `enter` - View entry detail
- `d` - Mark/unmark entry for deletion (shows "D" prefix, persists across navigation)
- `$` - Delete all marked items (shown in status message after marking, y/n confirmation, cascades to linked todos)
- `/` - Filter by tags and/or date
- `c` - Clear filter (shown in status message when filter is active)
- `?` - Show help page
- `t` - Jump to todos
- `q` - Quit

*Entry View (Read-Only):*
- `n` - New Entry
- `a` - Add Standalone Todo
- `j/k` or `↑/↓` - Navigate between entries
- `f/b` - Page forward/backward in long entries
- `d` - Mark/unmark current entry for deletion (shows "[MARKED FOR DELETION]" in footer)
- `$` - Delete all marked items (shown in status message after marking, y/n confirmation, cascades to linked todos)
- Shows entry with inline todos
- `?` - Show help page
- `e` - Jump to entries
- `t` - Jump to todos
- `esc` - Back to entry list
- `q` - Quit

*Todo List:*
- `n` - New Entry
- `a` - Add Standalone Todo
- `j/k` or `↑/↓` - Navigate
- `space` - Toggle todo status (saves immediately)
- `d` - Mark/unmark todo for deletion (shows "D" prefix, can mark entry-linked todos too)
- `$` - Delete all marked items (shown in status message after marking, y/n confirmation, only standalone todos deleted)
- `/` - Filter by tags and/or date
- `c` - Clear filter (shown in status message when filter is active)
- `?` - Show help page
- `e` - Jump to entries
- `q` - Quit

*Todo View (Read-Only):*
- `n` - New Entry
- `a` - Add Standalone Todo
- `space` - Cycle todo status: open → next → done (saves immediately)
- `j/k` or `↑/↓` - Navigate between todos
- `f/b` - Page forward/backward in long todos
- `d` - Mark/unmark current todo for deletion (shows "[MARKED FOR DELETION]" in footer)
- `$` - Delete all marked items (shown in status message after marking, y/n confirmation)
- `?` - Show help page
- `e` - Jump to entries
- `t` - Jump to todos
- `esc` or `enter` - Back to todo list
- `q` - Quit

*Add Todo Form:*
- Type todo title (tags auto-extracted from @mentions)
- `enter` - Save and start new todo (shows "saved" confirmation, power mode for rapid entry)
- `esc` - Cancel and return to todo list

*Theme Selector:*
- `j/k` or `↑/↓` - Navigate themes
- `enter` - Select theme (saves to config)
- `esc` - Cancel

*Help Page:*
- `f/b` - Page forward/backward
- `esc` or `q` - Close help

## Development

### Quick Start

```bash
# Clone and setup
git clone https://github.com/apodacaa/amos.git
cd amos
go mod download

# Install git hooks (recommended)
./scripts/install-hooks.sh

# Run the app
make run
```

### Common Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the app |
| `make check` | Format + vet code |
| `make ci` | **Full checks + tests (before commit)** |
| `make ci-cover` | Full checks + tests with coverage |
| `make build` | Build binary |
| `make test` | Run tests |
| `make help` | Show all commands |

### Before Committing
```bash
make ci  # Run all checks + tests (or install git hooks to auto-run)
```

### Git Hooks
The pre-commit hook automatically runs `make ci` before every commit:
```bash
./scripts/install-hooks.sh  # Install once after cloning
```

To bypass the hook (not recommended):
```bash
git commit --no-verify
```

### Hot Reload (Optional)
```bash
make install-air  # Install once
air               # Run with auto-reload
```

## Features

✅ **Journal Entries**
- Create entries with title + body
- Auto-extract @tags from content
- **Unified filtering**: Filter by tags and/or dates with `/` key
  - Tag autocomplete (type `@` then tag name, press `tab` to complete)
  - Date filters: `today`, `yesterday`, `last N days`, `YYYY-MM-DD`, `YYYY-MM-DD to YYYY-MM-DD`
  - `/` opens filter with current values for easy editing
  - `c` clears filter (shown in status message when filter is active)
  - Status message after applying filter: "Filter applied. Press c to clear"
- View entries chronologically (newest first)
- **Deletion**: neomutt-style multi-select (mark with `d`, `$` shown in status message after marking, cascade deletes linked todos)
- Cross-navigation: jump between todos/entries with `t`/`e` keys
- Global create: `n` (new entry) and `a` (add todo) work from any read-only view
- Save confirmation: entry form shows "saved" toast message

✅ **Todo Management**
- **Standalone todos**: Create todos independently with `a` key from any view
- **Entry-linked todos**: Extract from entries with `!todo` syntax
- Toggle status with `space` (immediate save)
- **Unified filtering**: Same as entries - filter by tags and/or dates with `/` key
  - Tag autocomplete (type `@` then tag name, press `tab` to complete)
  - Date filters: `today`, `yesterday`, `last N days`, `YYYY-MM-DD`, `YYYY-MM-DD to YYYY-MM-DD`
  - `/` opens filter with current values for easy editing
  - `c` clears filter (shown in status message when filter is active)
  - Status message after applying filter: "Filter applied. Press c to clear"
- Sort: open first → position → newest
- View todos by entry or all together
- Cross-navigation: jump between todos/entries with `t`/`e` keys
- Global create: `n` (new entry) and `a` (add todo) work from any read-only view
- Save confirmation: add todo form shows "saved" toast message

✅ **Navigation**
- Peer navigation: `e` (entries), `t` (todos) work from all views - no hierarchy
- Global shortcuts: `n` (new entry) and `a` (add todo) work from any read-only view
- `esc` exits forms to their natural home: entry form → entry list, add todo → todo list
- Immediate writes (no hidden pending state)
- Full context visible (todos show in entry view)
- No unnecessary features or decorations
- Fast, minimal TUI
- **Theme support**: Choose between brutalist (monochrome, terminal defaults) or cyberpunk (neon colors) themes with `s` key
- **Anchored help text**: Footer stays at bottom (no bouncing)
- **Viewport windowing**: Long lists show 20-30 items with scroll indicators
- **Entry scrolling**: Navigate long entries with `d/u` keys (d=down, u=up)

## Project Structure

```
.
├── main.go                 # Entry point (~10 lines)
├── model.go                # Model, Init, Update, View (Elm architecture)
├── messages.go             # Message types for async operations
├── commands.go             # tea.Cmd functions (side effects)
├── update_*.go             # Key handlers per view
│   ├── update_entry.go
│   ├── update_entries.go
│   ├── update_entry_view.go
│   ├── update_unified_filter.go
│   ├── update_todos.go
│   ├── update_add_todo.go
│   └── update_theme_selector.go
├── ui/                     # View renderers (pure functions)
│   ├── entry_form.go
│   ├── entry_list.go
│   ├── entry_view.go
│   ├── unified_filter.go
│   ├── todo_list.go
│   ├── add_todo_form.go
│   ├── theme_selector.go
│   └── styles.go
├── internal/               # Business logic
│   ├── models/            # Data structures
│   │   ├── entry.go
│   │   ├── todo.go
│   │   └── config.go      # User preferences
│   ├── storage/           # JSON persistence
│   │   └── storage.go
│   └── helpers/           # Utilities
│       ├── sorting.go     # Centralized sorting logic
│       ├── tags.go        # Tag extraction and filtering
│       ├── todos.go       # Todo extraction
│       ├── dates.go       # Date filter parsing
│       └── filter_parser.go  # Unified filter parsing
├── Makefile               # Development commands
└── go.mod                 # Go module definition
```

## Architecture

**Bubble Tea** uses the **Elm Architecture** pattern:
- `Model` - Application state (in `model.go`)
- `Init()` - Initialize model, return commands
- `Update(msg)` - Handle messages, return (model, cmd)
- `View()` - Render UI from model state

**File Organization (Bubble Tea Best Practices):**
- `main.go` - Entry point only (~10 lines)
- `model.go` - Model struct + Init/Update/View (Elm core)
- `messages.go` - All message types
- `commands.go` - All tea.Cmd functions (side effects)
- `update_*.go` - Key handlers per view (domain separation)
- `ui/` - Pure view renderers, no state
- `internal/` - Business logic (models, storage, helpers)

**Key Rules:**
- No side effects in `Update` - return commands instead
- Views are pure functions - no state mutation
- Exported names use PascalCase, unexported use camelCase

## Troubleshooting

**"staticcheck: command not found"**
```bash
make staticcheck  # Auto-installs
```

**"air: command not found"**
```bash
make install-air
```

**Build issues**
```bash
go mod tidy
make build
```

## Dependencies

- **bubbletea** v1.3.10 - TUI framework
- **lipgloss** v1.1.0 - Styling library
- **bubbles** v0.21.0 - Textarea component
- **Go 1.24+** required

## Data Storage

- Entries stored in `~/.amos/entries.json`
- Todos stored in `~/.amos/todos.json`
- User config stored in `~/.amos/config.json` (theme preference)
- Plain JSON format (no database)
- Auto-creates directory on first run

## Design Philosophy

**Brutalist Principles:**
1. **Immediate writes** - `space` toggles todo AND saves (no deferred state)
2. **Full context** - Todos visible in entry view
3. **No hidden state** - What you see is what's saved
4. **Simple is better** - Normalize positions every move vs complex tracking
5. **One action = one effect** - No multi-step workflows
6. **Peer navigation** - Entry list and todo list are peers, no hierarchy or back button
7. **Global actions** - `n` and `a` keys work from any read-only view for fast creation
8. **Honest workspaces** - All views are functional workspaces with left-aligned honest UI
9. **Monochrome palette** - Pure black/white/gray, no colors
10. **Anchored UI** - Help text stays at bottom (no bouncing during navigation)
11. **Consistent ordering** - Keys appear in same logical order across all views
12. **No decorations** - No italics, no Unicode bullets, just ASCII
13. **Viewport windowing** - Lists show manageable chunks with scroll position indicators
14. **Context-aware UI** - Headers show core navigation, status messages show context-dependent actions (`$` for deletion after marking, `c` for clearing active filters)

**Tag Syntax:**
- `@work` in entry content → auto-extracted to tags array
- `!todo Task description @tag` → creates linked todo

**Position System:**
- Todos have position field for priority
- Lower position = higher priority
- Sorted: open first → position → newest

