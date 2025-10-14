# Copilot Instructions for Amos

## Project Overview
Amos is a minimal Textual TUI (Terminal User Interface) prototype for fast iteration. It's a Python project using Poetry for dependency management and Textual for building interactive terminal applications.

## Architecture
- **Framework**: Textual (Python TUI framework)
- **Package Manager**: Poetry (no Makefile - all commands via Poetry)
- **Main App**: `amos/app.py` contains `AmosApp` class
- **Entry Point**: `poetry run textual run --dev amos.app:AmosApp`

## Development Workflow

### Running the App
```bash
# Development mode with hot reload
poetry run textual run --dev amos.app:AmosApp

# Debug with console (use 2 terminals)
poetry run textual console  # Terminal 1
poetry run textual run --dev amos.app:AmosApp --console  # Terminal 2
```

### Code Quality (No Makefile - Direct Poetry Commands)
```bash
# Format
poetry run black .
poetry run isort .

# Lint
poetry run ruff check --fix .

# All at once
poetry run black . && poetry run isort . && poetry run ruff check --fix .
```

## Project Conventions

### Code Style
- **Black** for formatting (line length: default 88)
- **isort** for import sorting
- **Ruff** for linting (configured with auto-fix)
- No pre-commit hooks - run formatters manually before commits

### App Structure Pattern
The `AmosApp` class in `amos/app.py` follows this pattern:
- `compose()` - Define widget layout (Header, Static, Footer)
- `on_key()` - Handle keyboard events asynchronously
- CSS-in-Python via `CSS` class variable
- State management via instance variables (e.g., `self.count`)

### Adding New Features
1. Edit `amos/app.py` directly
2. The `--dev` flag auto-reloads changes
3. Use `query_one("#id", Widget)` to update widgets
4. Keep UI logic in `on_key()` or event handlers

## Key Files
- `amos/app.py` - Main Textual application with `AmosApp` class
- `pyproject.toml` - Poetry config (Python >=3.10, textual >=6.1.0)
- `README.md` - Setup and workflow documentation
- `claude-guardrails.md` - Safety guidelines for AI assistance

## Critical Patterns
- **No Makefile**: All commands use `poetry run` directly
- **Development Mode**: Always use `--dev` flag for hot reload
- **Widget Updates**: Use `widget.update()` method to refresh UI
- **Async Handlers**: Event handlers are `async` functions
- **CSS Styling**: Inline CSS in Python via class variable

## Dependencies
- **Main**: textual ^6.1.0
- **Dev**: black, ruff, isort
- **No testing framework yet** - add pytest if needed