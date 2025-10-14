```text
# Claude Guardrails for Amos

## Purpose
This file contains safety rules and constraints for AI assistance when working on the Amos codebase. These guardrails ensure code quality, prevent breaking changes, and maintain project conventions.

---

## General Rules

### 1. Always Ask Before:
- **Deleting or renaming files** - Could break imports or workflows
- **Installing new dependencies** - Requires user approval
- **Changing project structure** - Could impact build/run commands
- **Modifying `go.mod` or `go.sum`** - Use `go get` commands instead

### 2. Never Do Without Permission:
- Remove existing functionality
- Change Go module path
- Modify deployment/build scripts
- Add large dependencies (>1MB)
- Change code style conventions

### 3. Always Verify Before Suggesting:
- Exact function signatures from Bubble Tea API
- Correct Elm architecture patterns
- Go module import paths
- Command line syntax for `go` commands

---

## Bubble Tea Specific Rules

### Elm Architecture Constraints
```go
// ✅ CORRECT: Return new model and command
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            m.inputValue = m.textInput.Value()
            return m, tea.Quit  // Return command
        }
    }
    return m, nil
}

// ❌ WRONG: Side effects in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if msg.String() == "enter" {
        saveToFile(m.data)  // NO! Use commands instead
    }
    return m, nil
}
```

### Model State Management
```go
// ✅ CORRECT: All state in Model
type model struct {
    currentView string
    entries     []Entry
    textInput   textinput.Model
    loading     bool
}

// ❌ WRONG: Global state or package-level vars
var currentEntries []Entry  // NO! Put in Model
```

### Commands for Side Effects
```go
// ✅ CORRECT: Return command for async operations
func saveEntry(entry Entry) tea.Cmd {
    return func() tea.Msg {
        err := writeToFile(entry)
        return saveCompleteMsg{err}
    }
}

// ❌ WRONG: Direct I/O in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    writeToFile(m.entry)  // NO! Return command instead
    return m, nil
}
```

---

## Go Code Style

### Formatting
- **Always run `go fmt` before suggesting code**
- Use standard Go formatting (tabs, not spaces)
- Follow Go naming conventions (camelCase for unexported, PascalCase for exported)

### Linting
- Code should pass `golangci-lint run`
- No unused imports or variables
- Error handling required (no naked returns on errors)

### Imports
```go
// ✅ CORRECT: Grouped imports
import (
    "fmt"
    "os"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// ❌ WRONG: Ungrouped imports
import "fmt"
import tea "github.com/charmbracelet/bubbletea"
```

---

## File Organization Rules

### Project Structure
```
main.go           # Entry point only - keep minimal
helpers/          # Business logic (storage, formatting, etc.)
  storage.go      # JSON read/write
  formatting.go   # Date/text helpers
views/            # View renderers (if needed)
  dashboard.go
  entry.go
```

### Code Separation
- **main.go**: Model definition, Init/Update/View
- **helpers/**: Pure functions for data manipulation
- **views/**: Complex rendering logic (if view functions get large)

---

## Design Philosophy (Brutalist Principles)

### Minimalism
```go
// ✅ CORRECT: Simple, clear code
func (m model) View() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("10")).
        Render("Simple")
}

// ❌ WRONG: Over-styled, complex layouts
func (m model) View() string {
    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF5733")).
        Background(lipgloss.Color("#C70039")).
        Bold(true).
        Italic(true).
        Underline(true).
        Border(lipgloss.DoubleBorder())
    // Too much!
}
```

### Theme Variables Only
- Use lipgloss theme colors (not hardcoded hex colors)
- No custom color palettes
- Stick to terminal defaults where possible

### No Feature Creep
- Journal entries with @tags
- Todos with statuses
- That's it - no plugins, extensions, or "nice to have" features without explicit user request

---

## Storage Rules

### JSON Format
- Store data in `~/.amos/` directory
- Use standard Go `encoding/json` package
- No databases, no binary formats

### Data Structures
```go
// ✅ CORRECT: Simple, flat structures
type Entry struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`
    Content   string    `json:"content"`
    Tags      []string  `json:"tags"`
}

// ❌ WRONG: Nested complexity
type Entry struct {
    Metadata struct {
        Author struct {
            // Too deep!
        }
    }
}
```

---

## Testing (When Implemented)

### Test Files
- Name: `*_test.go`
- Location: Same directory as code under test
- Run: `go test ./...`

### Coverage
- Don't obsess over 100% coverage
- Focus on critical paths (storage, parsing)
- UI code doesn't need unit tests

---

## Error Handling

### Required Patterns
```go
// ✅ CORRECT: Check and handle errors
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// ❌ WRONG: Ignore errors
data, _ := os.ReadFile(path)  // NO!
```

### User-Facing Errors
```go
// ✅ CORRECT: Helpful error messages
return fmt.Errorf("could not save entry %s: %w", id, err)

// ❌ WRONG: Vague errors
return fmt.Errorf("error: %w", err)
```

---

## Breaking Changes

### Require User Approval For:
1. Changing Model struct fields (breaks existing code)
2. Renaming exported functions (breaks imports)
3. Modifying JSON schema (breaks saved data)
4. Changing keyboard shortcuts (breaks muscle memory)

### Safe Changes:
- Adding new unexported functions
- Refactoring internal logic (same behavior)
- Improving comments/documentation
- Fixing bugs that don't change API

---

## Summary Checklist

Before suggesting code, verify:
- [ ] Follows Elm architecture (Model → Update → View)
- [ ] No side effects in Update function
- [ ] Proper error handling (no naked `_` ignores)
- [ ] Formatted with `go fmt`
- [ ] Uses theme variables, not custom colors
- [ ] Minimalist design (no unnecessary features)
- [ ] Would pass `golangci-lint run`

**When in doubt, ask the user before proceeding.**
```