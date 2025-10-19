# Amos

Minimal Bubble Tea (Go) TUI for journal + todo management. Brutalist design, fast iteration.

## Quick Start

```bash
# Install dependencies
go mod download

# Install git hooks (recommended)
./scripts/install-hooks.sh

# Run the app
make run
```

**Keyboard Shortcuts:**

*Dashboard:*
- Massive ASCII "AMOS" title with vanity metrics (total entries, this week, top tag, open todos)
- `n` - New Entry
- `a` - Add Standalone Todo
- `t` - View Todos List
- `e` - View Entries List
- `q` or `Ctrl+C` - Quit

*Entry Form:*
- `Ctrl+S` - Save entry
- `esc` - Cancel (goes to dashboard, with confirmation if unsaved)

*Entry List:*
- `n` - New Entry
- `a` - Add Standalone Todo
- `j/k` or `↑/↓` - Navigate
- `enter` - View entry detail
- `t` - Jump to todos
- `@` - Filter by tag (or clear filter)
- `d` - Go to dashboard
- `q` - Quit

*Entry View (Read-Only):*
- `n` - New Entry
- `a` - Add Standalone Todo
- Shows entry with inline todos
- `x` - Toggle todo status (saves immediately)
- `t` - Jump to todos
- `e` - Jump to entries
- `d` - Go to dashboard
- `q` - Quit

*Todo List:*
- `n` - New Entry
- `a` - Add Standalone Todo
- `j/k` or `↑/↓` - Navigate
- `x` - Toggle todo status (saves immediately)
- `u/i` - Move todo up/down (manual priority)
- `e` - Jump to entries
- `d` - Go to dashboard
- `q` - Quit

*Add Todo Form:*
- Type todo title (tags auto-extracted from @mentions)
- `enter` - Save and go to dashboard
- `esc` - Cancel and go to dashboard

## Development

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
- Filter by tag with @ key (brutalist tag picker)
- View entries chronologically (newest first)
- **Append-only**: No delete (journal is historical record)
- See todo counts in entry list: `[3 todos: 1 open]`
- Cross-navigation: jump to todos with `t` key, create entry/todo with `n`/`a` from any view

✅ **Todo Management**
- **Standalone todos**: Create todos independently with `a` key from any view
- **Entry-linked todos**: Extract from entries with `!todo` syntax
- Toggle status with `x` (immediate save)
- Manual priority with u/i keys (move up/down)
- Sort: open first → position → newest
- View todos by entry or all together
- Cross-navigation: jump to entries with `e` key, create entry/todo with `n`/`a` from any view

✅ **Brutalist Navigation**
- Explicit navigation: `d` (dashboard), `t` (todos), `e` (entries) work from all views
- Global shortcuts: `n` (new entry) and `a` (add todo) work from any read-only view
- Escape only for forms: cancel and return to dashboard
- Immediate writes (no hidden pending state)
- Full context visible (todos show in entry view)
- No unnecessary features or decorations
- Fast, minimal TUI

## Project Structure

```
.
├── main.go                 # Entry point (~10 lines)
├── model.go                # Model, Init, Update, View (Elm architecture)
├── messages.go             # Message types for async operations
├── commands.go             # tea.Cmd functions (side effects)
├── update_*.go             # Key handlers per view
│   ├── update_dashboard.go
│   ├── update_entry.go
│   ├── update_entries.go
│   ├── update_entry_view.go
│   ├── update_tag_picker.go
│   ├── update_todos.go
│   └── update_add_todo.go
├── ui/                     # View renderers (pure functions)
│   ├── dashboard.go
│   ├── entry_form.go
│   ├── entry_list.go
│   ├── entry_view.go
│   ├── tag_picker.go
│   ├── todo_list.go
│   ├── add_todo_form.go
│   └── styles.go
├── internal/               # Business logic
│   ├── models/            # Data structures
│   │   ├── entry.go
│   │   └── todo.go
│   ├── storage/           # JSON persistence
│   │   └── storage.go
│   └── helpers/           # Utilities
│       ├── sorting.go     # Centralized sorting logic
│       ├── tags.go        # Tag extraction and filtering
│       └── todos.go       # Todo extraction
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
- Plain JSON format (no database)
- Auto-creates directory on first run

## Design Philosophy

**Brutalist Principles:**
1. **Immediate writes** - `x` toggles todo AND saves (no deferred state)
2. **Full context** - Todos visible in entry view, stats in list
3. **No hidden state** - What you see is what's saved
4. **Simple is better** - Normalize positions every move vs complex tracking
5. **One action = one effect** - No multi-step workflows
6. **Context-aware navigation** - Escape returns to previous view (where you came from)
7. **Global actions** - `n` and `a` keys work from any read-only view for fast creation
8. **Monument aesthetics** - Dashboard hits you with massive centered ASCII art and stats, utility views are honest left-aligned workspaces

**Tag Syntax:**
- `@work` in entry content → auto-extracted to tags array
- `!todo Task description @tag` → creates linked todo

**Position System:**
- Todos have position field for manual priority
- Lower position = higher priority
- Sorted: open first → position → newest
- u/i keys move todos up/down (renumbers all positions)

