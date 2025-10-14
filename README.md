```text
Minimal Amos Textual prototype

Quick start (fresh repo):

1. Install Poetry if you don't have it: https://python-poetry.org/docs/#installation
2. Install dependencies:
   poetry install

3. Run the app (exact command you requested):
   poetry run textual run --dev amos.app:AmosApp

Controls:
- i : increment counter
- d : decrement counter
- q : quit

Notes:
- The project uses Python >=3.11 in pyproject.toml. Adjust if needed.
- Dev tools (black, ruff, isort, pre-commit) are declared as dev-dependencies but not enforced. Add pre-commit hooks if you want them active.
```