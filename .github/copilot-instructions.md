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
main.go           # Entry point only (~10 lines)
model.go          # Model, Init, Update, View (Elm architecture)
ui/               # View renderers (pure functions)
  dashboard.go
  entry_form.go
internal/         # Business logic (planned)
  storage/        # JSON read/write
    storage.go
  models/         # Data structures
    entry.go
    todo.go
  helpers/        # Utilities
    tags.go
    ids.go
```

See `ARCHITECTURE.md` for detailed organization guidelines.