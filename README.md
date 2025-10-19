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
- `n` - New Entry
- `e` - View Entries List
- `t` - View Todos List
- `s` - Search (coming soon)
- `q` or `Ctrl+C` - Quit

*Entry Form:*
- `Ctrl+S` - Save entry
- `esc` - Exit (with confirmation if unsaved)

*Entry List:*
- `j/k` or `↑/↓` - Navigate
- `enter` - View entry detail
- `d` (double tap) - Delete entry
- `esc` - Back to dashboard

*Entry View:*
- Shows entry with inline todos
- `esc` - Back to entry list

*Todo List:*
- `j/k` or `↑/↓` - Navigate
- `space` - Toggle todo status (saves immediately)
- `u/i` - Move todo up/down (manual priority)
- `esc` - Back to dashboard

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
- View entries chronologically (newest first)
- Delete entries with double-tap confirmation
- See todo counts in entry list: `[3 todos: 1 open]`

✅ **Todo Management**
- Extract todos from entries with `!todo` syntax
- Toggle status with space (immediate save)
- Manual priority with u/i keys (move up/down)
- Sort: open first → position → newest
- View todos by entry or all together

✅ **Brutalist Design**
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
│   └── update_todos.go
├── ui/                     # View renderers (pure functions)
│   ├── dashboard.go
│   ├── entry_form.go
│   ├── entry_list.go
│   ├── entry_view.go
│   ├── todo_list.go
│   └── styles.go
├── internal/               # Business logic
│   ├── models/            # Data structures
│   │   ├── entry.go
│   │   └── todo.go
│   ├── storage/           # JSON persistence
│   │   └── storage.go
│   └── helpers/           # Utilities
│       ├── tags.go
│       └── todos.go
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
1. **Immediate writes** - Space toggles todo AND saves (no deferred state)
2. **Full context** - Todos visible in entry view, stats in list
3. **No hidden state** - What you see is what's saved
4. **Simple is better** - Normalize positions every move vs complex tracking
5. **One action = one effect** - No multi-step workflows

**Tag Syntax:**
- `@work` in entry content → auto-extracted to tags array
- `!todo Task description @tag` → creates linked todo

**Position System:**
- Todos have position field for manual priority
- Lower position = higher priority
- Sorted: open first → position → newest
- u/i keys move todos up/down (renumbers all positions)

