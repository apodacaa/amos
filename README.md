# Amos

Minimal Bubble Tea (Go) TUI for journal + todo management. Brutalist design, fast iteration.

## Quick Start

### Prerequisites
- Go 1.24+ (project uses Bubble Tea v1.3.10)

### 1. Install Dependencies
```bash
go mod download
```

### 2. Run the App
```bash
go run main.go
```

**Controls:**
- `n` - New Entry
- `t` - Todos
- `e` - Entries
- `s` - Search
- `q` - Quit

## Development Workflow

### Run with Auto-Reload
Use `air` for hot reloading:
```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

### Code Formatting
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test ./...
```

## Project Structure

```
.
├── main.go              # Main entry point
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
└── README.md
```

## Architecture

- **Bubble Tea** - TUI framework (Elm architecture)
- **Model** - Application state
- **Update** - Handle messages, update state
- **View** - Render UI from state

## Critical Patterns

- **Elm Architecture** - Model, Update, View cycle
- **Messages** - Commands return tea.Cmd for async operations
- **No Side Effects in Update** - Return commands, don't execute
- **Immutable Updates** - Return new model, don't mutate

## Dependencies

- **bubbletea** v1.3.10 - TUI framework
- **lipgloss** (transitive) - Styling

## Notes

- Go 1.24+ required
- Simple brutalist design (minimal styling)
- Data stored in `~/.amos/` (JSON)
