# Amos

Minimal Textual TUI prototype for fast iteration.

## Quick Start

### 1. Install Poetry
If you don't have it: https://python-poetry.org/docs/#installation

### 2. Install Dependencies
```bash
poetry install
```

### 3. Run the App
```bash
poetry run textual run --dev amos.app:AmosApp
```

**Controls:**
- `i` - Increment counter
- `d` - Decrement counter
- `q` - Quit

## Development Workflow

### Run with Hot Reload
```bash
poetry run textual run --dev amos.app:AmosApp
```

### Code Quality

**Format code:**
```bash
poetry run black .
poetry run isort .
```

**Lint and auto-fix:**
```bash
poetry run ruff check --fix .
```

**Run all formatting at once:**
```bash
poetry run black . && poetry run isort . && poetry run ruff check --fix .
```

### Debug with Console

**Terminal 1 - Start the console:**
```bash
poetry run textual console
```

**Terminal 2 - Run with console logging:**
```bash
poetry run textual run --dev amos.app:AmosApp --console
```

## Project Structure

```
amos/
├── amos/
│   ├── __init__.py
│   └── app.py              # Main Textual app
├── pyproject.toml          # Poetry config & dependencies
└── README.md
```

## Notes

- Requires Python >=3.10 (see `pyproject.toml`)
- The `--dev` flag enables hot reloading for faster development
