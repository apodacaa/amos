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

### Project Structure (Bubble Tea Best Practices)
```
main.go                 # Entry point only (~10 lines)
model.go                # Model struct + Init/Update/View (Elm core)
messages.go             # All message types
commands.go             # All tea.Cmd functions (side effects)
update_*.go             # Key handlers per view (domain separation)
  update_dashboard.go
  update_entry.go
  update_entries.go
  update_entry_view.go
  update_tag_picker.go
  update_todos.go
ui/                     # View renderers (pure functions)
  dashboard.go
  entry_form.go
  entry_list.go
  entry_view.go
  tag_picker.go
  todo_list.go
  styles.go
internal/               # Business logic
  storage/              # JSON persistence
  models/               # Data structures
  helpers/              # Pure utility functions
    sorting.go          # Centralized sorting logic
    tags.go             # Tag extraction and filtering
    todos.go            # Todo extraction
```

### Code Separation (Current Pattern)
- **main.go**: Entry point only (~10 lines)
- **model.go**: Model struct, Init(), Update() router, View() router
- **messages.go**: All message types (saveCompleteMsg, todosLoadedMsg, etc.)
- **commands.go**: All tea.Cmd functions (loadEntries, saveTodo, toggleTodoImmediate, etc.)
- **update_*.go**: Key handlers per view (handleEntryKeys, handleTodosListKeys, handleTagPickerKeys, etc.)
- **ui/**: Pure view renderers (RenderEntryList, RenderTodoList, RenderTagPicker, etc.)
- **internal/helpers/**: Centralized business logic (sorting, filtering, parsing)
- **internal/storage/**: JSON persistence
- **internal/models/**: Data structures

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

### Brutalist Design Principles
**Current features (keep minimal):**
- Journal entries with @tags (auto-extracted)
- Tag filtering with @ key (brutalist picker: visible state, toggle on/off)
- Cross-navigation (t/e keys to jump between entries and todos)
- Todos with `!todo` syntax (linked to entries)
- Todo toggle with immediate save (space key)
- Manual priority with u/i keys (position field)
- Todo stats in entry list: `[3 todos: 1 open]`
- Full todo display in entry view (no hidden info)

**Core philosophy:**
- Immediate writes (no deferred/pending state)
- Full context visible (no navigation required)
- Visible state (filter shown in title, no hidden modes)
- One action = one effect
- Normalize over preserve (simpler logic)
- DRY principle (extract duplicates to helpers)
- No feature creep without explicit user request

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