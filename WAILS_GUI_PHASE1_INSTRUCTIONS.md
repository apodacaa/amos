# Wails GUI Implementation - Phase 1 (Linux)

## Context

We're adding a Wails-based GUI to the existing amos TUI app using a monorepo structure. The TUI and GUI will share all business logic (`internal/` packages) but have separate entry points and UIs.

**Phase 1 Goal:** Get a basic Wails GUI working on Linux (Ubuntu) that can display entries and todos using the existing storage layer.

## Project Structure Overview

```
amos/
├── cmd/
│   ├── tui/              # Existing Bubble Tea TUI (NEW location)
│   │   └── main.go       # Move current main.go here
│   └── gui/              # NEW: Wails GUI
│       └── main.go       # Wails entry point
├── internal/             # Shared business logic (UNCHANGED)
│   ├── models/
│   ├── storage/
│   └── helpers/
├── frontend/             # NEW: Wails web UI
│   ├── src/
│   │   ├── App.svelte   # Main UI component
│   │   ├── main.js      # Frontend entry
│   │   └── style.css    # Brutalist styling
│   ├── index.html
│   └── wails.json       # Wails config
├── model.go              # Move to cmd/tui/ (TUI-specific)
├── update_*.go           # Move to cmd/tui/ (TUI-specific)
├── ui/                   # Move to cmd/tui/ui/ (TUI-specific)
├── messages.go           # Move to cmd/tui/ (TUI-specific)
├── commands.go           # Move to cmd/tui/ (TUI-specific)
└── Makefile              # Update with new build targets
```

## Step 1: Reorganize Existing Code (Monorepo Structure)

**Goal:** Move TUI-specific code to `cmd/tui/` without breaking anything.

### 1.1 Create Directory Structure

```bash
mkdir -p cmd/tui
mkdir -p cmd/tui/ui
mkdir -p cmd/gui
mkdir -p frontend/src
```

### 1.2 Move TUI Files

Move these files to `cmd/tui/`:
- `main.go` → `cmd/tui/main.go`
- `model.go` → `cmd/tui/model.go`
- `update_*.go` → `cmd/tui/update_*.go` (all update files)
- `messages.go` → `cmd/tui/messages.go`
- `commands.go` → `cmd/tui/commands.go`
- `ui/*.go` → `cmd/tui/ui/*.go` (all UI files)

### 1.3 Update Import Paths in TUI Files

After moving files, update imports in `cmd/tui/*.go`:

**Before:**
```go
import (
    "github.com/apodacaa/amos/internal/models"
    "github.com/apodacaa/amos/internal/storage"
    "github.com/apodacaa/amos/ui"
)
```

**After:**
```go
import (
    "github.com/apodacaa/amos/internal/models"
    "github.com/apodacaa/amos/internal/storage"
    "github.com/apodacaa/amos/cmd/tui/ui"  // Note: ui is now under cmd/tui
)
```

### 1.4 Update Makefile

```makefile
# TUI targets (update existing)
run:
	go run cmd/tui/*.go

build:
	go build -o amos cmd/tui/*.go

# GUI targets (NEW)
run-gui:
	cd cmd/gui && wails dev

build-gui:
	cd cmd/gui && wails build

# Keep existing ci, test, etc.
ci: fmt vet staticcheck test

test:
	go test ./...
```

### 1.5 Verify TUI Still Works

```bash
make run  # Should launch TUI as before
make ci   # All tests should pass
```

**STOP HERE and verify before proceeding.** The TUI must work before adding GUI.

## Step 2: Install Wails

### 2.1 Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2.2 Install Linux Dependencies

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev
```

### 2.3 Verify Installation

```bash
wails doctor  # Should show all green checkmarks for Linux
```

## Step 3: Create Wails GUI Application

### 3.1 Initialize Wails Project in cmd/gui

```bash
cd cmd/gui
wails init -n amos-gui -t svelte
cd ../..
```

This creates the basic Wails structure. We'll customize it next.

### 3.2 Create Go Backend (cmd/gui/main.go)

```go
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create application
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "amos",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
```

### 3.3 Create App Struct (cmd/gui/app.go)

This is where we expose Go functions to the frontend.

```go
package main

import (
	"context"
	"time"

	"github.com/apodacaa/amos/internal/models"
	"github.com/apodacaa/amos/internal/storage"
)

// App struct
type App struct {
	ctx     context.Context
	storage *storage.Storage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Initialize storage (same as TUI)
	sysConfig, err := storage.LoadSystemConfig()
	if err != nil {
		// Use defaults if config doesn't exist
		sysConfig = storage.NewDefaultSystemConfig()
	}
	
	a.storage = storage.New(sysConfig.DataDir)
}

// GetEntries returns all entries
func (a *App) GetEntries() ([]models.Entry, error) {
	return a.storage.LoadEntries()
}

// GetTodos returns all todos
func (a *App) GetTodos() ([]models.Todo, error) {
	return a.storage.LoadTodos()
}

// CreateEntry creates a new entry
func (a *App) CreateEntry(title, body string, tags []string) (*models.Entry, error) {
	entry := models.Entry{
		ID:        generateID(),
		Title:     title,
		Body:      body,
		Tags:      tags,
		Timestamp: time.Now(),
	}
	
	entries, err := a.storage.LoadEntries()
	if err != nil {
		entries = []models.Entry{}
	}
	
	entries = append(entries, entry)
	
	if err := a.storage.SaveEntries(entries); err != nil {
		return nil, err
	}
	
	return &entry, nil
}

// CreateTodo creates a new todo
func (a *App) CreateTodo(title string, tags []string) (*models.Todo, error) {
	todo := models.Todo{
		ID:        generateID(),
		Title:     title,
		Status:    "open",
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	todos, err := a.storage.LoadTodos()
	if err != nil {
		todos = []models.Todo{}
	}
	
	todos = append(todos, todo)
	
	if err := a.storage.SaveTodos(todos); err != nil {
		return nil, err
	}
	
	return &todo, nil
}

// UpdateTodoStatus updates a todo's status
func (a *App) UpdateTodoStatus(id, status string) error {
	todos, err := a.storage.LoadTodos()
	if err != nil {
		return err
	}
	
	for i := range todos {
		if todos[i].ID == id {
			todos[i].Status = status
			todos[i].UpdatedAt = time.Now()
			break
		}
	}
	
	return a.storage.SaveTodos(todos)
}

// Helper function (copy from TUI or move to internal/helpers)
func generateID() string {
	return time.Now().Format("20060102-150405")
}
```

### 3.4 Create Frontend (frontend/src/App.svelte)

Brutalist design - minimal, functional, no frills.

```svelte
<script>
  import { onMount } from 'svelte';
  import { GetEntries, GetTodos, CreateEntry, CreateTodo, UpdateTodoStatus } from '../wailsjs/go/main/App';

  let entries = [];
  let todos = [];
  let view = 'entries'; // 'entries' or 'todos'
  let loading = true;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      entries = await GetEntries();
      todos = await GetTodos();
    } catch (err) {
      console.error('Failed to load data:', err);
    }
    loading = false;
  }

  async function toggleTodoStatus(todo) {
    const newStatus = todo.Status === 'done' ? 'open' : 'done';
    await UpdateTodoStatus(todo.ID, newStatus);
    await loadData();
  }

  function formatDate(timestamp) {
    return new Date(timestamp).toLocaleDateString();
  }
</script>

<main>
  <header>
    <h1>amos</h1>
    <nav>
      <button class:active={view === 'entries'} on:click={() => view = 'entries'}>
        entries
      </button>
      <button class:active={view === 'todos'} on:click={() => view = 'todos'}>
        todos ({todos.filter(t => t.Status !== 'done').length})
      </button>
    </nav>
  </header>

  {#if loading}
    <p>loading...</p>
  {:else if view === 'entries'}
    <section class="entries">
      <h2>entries</h2>
      {#if entries.length === 0}
        <p>no entries yet</p>
      {:else}
        <ul>
          {#each entries as entry}
            <li>
              <span class="date">{formatDate(entry.Timestamp)}</span>
              <strong>{entry.Title}</strong>
              {#if entry.Tags && entry.Tags.length > 0}
                <span class="tags">{entry.Tags.map(t => '@' + t).join(' ')}</span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {:else}
    <section class="todos">
      <h2>todos</h2>
      {#if todos.length === 0}
        <p>no todos yet</p>
      {:else}
        <ul>
          {#each todos as todo}
            <li class:done={todo.Status === 'done'}>
              <input
                type="checkbox"
                checked={todo.Status === 'done'}
                on:change={() => toggleTodoStatus(todo)}
              />
              <span>{todo.Title}</span>
              {#if todo.Tags && todo.Tags.length > 0}
                <span class="tags">{todo.Tags.map(t => '@' + t).join(' ')}</span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {/if}
</main>

<style>
  /* Brutalist styling - terminal inspired */
  :global(body) {
    font-family: monospace;
    background: #000;
    color: #fff;
    margin: 0;
    padding: 0;
  }

  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  header {
    border-bottom: 1px solid #fff;
    padding-bottom: 10px;
    margin-bottom: 20px;
  }

  h1 {
    margin: 0 0 10px 0;
    font-size: 24px;
  }

  nav button {
    background: none;
    border: 1px solid #fff;
    color: #fff;
    font-family: monospace;
    padding: 5px 10px;
    margin-right: 10px;
    cursor: pointer;
  }

  nav button.active {
    background: #fff;
    color: #000;
  }

  ul {
    list-style: none;
    padding: 0;
  }

  li {
    padding: 10px 0;
    border-bottom: 1px solid #333;
  }

  li.done {
    opacity: 0.5;
    text-decoration: line-through;
  }

  .date {
    color: #888;
    margin-right: 10px;
  }

  .tags {
    color: #888;
    margin-left: 10px;
  }

  input[type="checkbox"] {
    margin-right: 10px;
  }
</style>
```

### 3.5 Update Frontend Package Files

**frontend/package.json** - ensure Svelte dependencies:
```json
{
  "name": "amos-gui",
  "version": "1.0.0",
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "devDependencies": {
    "@sveltejs/vite-plugin-svelte": "^2.0.0",
    "svelte": "^3.55.0",
    "vite": "^4.0.0"
  }
}
```

**frontend/vite.config.js**:
```js
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
});
```

## Step 4: Build and Run

### 4.1 Install Frontend Dependencies

```bash
cd frontend
npm install
cd ..
```

### 4.2 Run in Development Mode

```bash
make run-gui
# OR
cd cmd/gui && wails dev
```

This should:
1. Compile Go backend
2. Start Vite dev server for frontend
3. Open GUI window showing entries and todos

### 4.3 Build Production Binary

```bash
make build-gui
# OR
cd cmd/gui && wails build
```

Output: `cmd/gui/build/bin/amos-gui` (Linux binary)

## Step 5: Verify Everything Works

**TUI should still work:**
```bash
make run      # Launches TUI
make ci       # Tests pass
make build    # Builds amos binary
```

**GUI should work:**
```bash
make run-gui  # Launches GUI in dev mode
make build-gui  # Builds amos-gui binary
./cmd/gui/build/bin/amos-gui  # Run production GUI
```

**Both should:**
- Read/write to same `~/.amos/` directory
- Show same entries and todos
- Work independently

## Phase 1 Complete!

At this point you should have:
- ✅ Working TUI (unchanged functionality)
- ✅ Working GUI (basic view of entries and todos)
- ✅ Shared storage layer (both use same JSON files)
- ✅ Monorepo structure ready for expansion

## Next Steps (Future Phases)

**Phase 2: Full CRUD in GUI**
- Add entry creation form
- Add todo creation form
- Entry editing
- Deletion

**Phase 3: Advanced Features**
- Filtering UI
- Tag management
- External editor integration ($EDITOR for entries)
- Markdown preview

**Phase 4: Distribution**
- GitHub Actions for Linux builds
- .deb package
- AppImage
- Eventually: macOS and Windows

## Notes

- Keep `internal/` packages completely UI-agnostic (they should never import TUI or GUI code)
- GUI and TUI can evolve independently
- Storage format is shared - changes must be backward compatible
- Tests in `internal/` test business logic, not UI
