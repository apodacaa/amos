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
- `n` - New Entry
- `t` - Todos  
- `e` - Entries
- `s` - Search
- `esc` - Back to Dashboard
- `q` or `Ctrl+C` - Quit

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

## Project Structure

```
.
├── main.go              # Entry point (~10 lines)
├── model.go             # Elm architecture (Model, Init, Update, View)
├── ui/                  # View components
│   ├── dashboard.go     # Dashboard view
│   └── entry_form.go    # Entry form view
├── Makefile             # Development commands
├── go.mod               # Go module definition
└── README.md
```

## Architecture

**Bubble Tea** uses the **Elm Architecture** pattern:
- `Model` - Application state (in `model.go`)
- `Init()` - Initialize model, return commands
- `Update(msg)` - Handle messages, return (model, cmd)
- `View()` - Render UI from model state

**File Organization:**
- `main.go` - Entry point only (~10 lines)
- `model.go` - All Elm architecture logic
- `ui/` - Pure view functions, no state
- `internal/` (future) - Business logic, data models

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
- **lipgloss** v1.1.0 - Styling library (transitive)
- **Go 1.24+** required

## Notes

- Simple brutalist design (minimal styling)
- Data stored in `~/.amos/` (JSON) - planned

