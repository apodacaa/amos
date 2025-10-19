# Copilot Instructions for Amos

## Project Overview
Amos is a minimal Bubble Tea (Go) TUI for journal + todo management. Fast iteration, brutalist design.

## Architecture
- **Framework**: Bubble Tea (Go TUI framework, Elm architecture)
- **Language**: Go 1.24+
- **Package Manager**: Go modules
- **Entry Point**: `go run main.go`

## Development Workflow

### Running the App
```bash
# Development mode
make run
# or
go run .

# With hot reload (install air first)
make install-tools
air
```

### Code Quality
```bash
# Quick check (fmt + vet)
make check

# Full check (fmt + vet + staticcheck)
make check-all

# Format only
make fmt
# or
go fmt ./...

# Vet only
make vet

# Staticcheck (linter)
make staticcheck

# Test
make test
# or
go test ./...
```

## Project Conventions

### Code Style
- **gofmt** for formatting (standard Go formatting)
- **golangci-lint** for linting
- Follow Go best practices

### Bubble Tea Patterns
The app follows Elm architecture:
- `Model` - Holds application state
- `Init()` - Initialize model, return commands
- `Update(msg)` - Handle messages, return (model, cmd)
- `View()` - Render UI string from model state

### Key Concepts
- **Messages** - User input, async results (tea.Msg)
- **Commands** - Side effects (tea.Cmd) - IO, timers, etc.
- **Immutability** - Return new model, don't mutate
- **No side effects in Update** - Return commands instead

## Critical Patterns
- **Elm Architecture**: Model → Update → View cycle
- **State Management**: All state in Model struct
- **Commands**: Async operations return tea.Cmd
- **Composability**: Nest models for complex UIs

## Dependencies
- **bubbletea** v1.3.10 - TUI framework
- **lipgloss** (planned) - Styling library
- **bubbles** (planned) - Reusable components

## File Structure
```
main.go                  # Entry point only (~10 lines)
model.go                 # Model, Init, Update, View (Elm core)
messages.go              # All message types (saveCompleteMsg, etc.)
commands.go              # All tea.Cmd functions (side effects/async)
update_dashboard.go      # Dashboard key handler
update_entry.go          # Entry form key handler
update_entries.go        # Entry list key handler
update_entry_view.go     # Entry view key handler
update_tag_picker.go     # Tag picker key handler
update_todos.go          # Todo list key handler
ui/                      # View renderers (pure functions)
  dashboard.go
  entry_form.go
  entry_list.go
  entry_view.go
  tag_picker.go
  todo_list.go
  styles.go
internal/                # Business logic
  storage/               # JSON persistence (~/.amos/)
    storage.go
  models/                # Data structures
    entry.go
    todo.go
  helpers/               # Utilities
    sorting.go           # Centralized sorting (DRY)
    tags.go              # Tag extraction and filtering
    todos.go             # Todo extraction
```

**Separation Pattern (Bubble Tea Best Practice):**
- Model = State only
- Messages = Async results
- Commands = Side effects (I/O)
- Update_* = State transitions per view
- UI = Pure renderers
- Helpers = Reusable logic (sorting, filtering, parsing)