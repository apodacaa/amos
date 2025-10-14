# Copilot Instructions for Amos

## Project Overview
Amos is an early-stage Python project currently in initial setup phase. The project structure indicates preparation for modern Python development workflows.

## Development Environment Setup
- **Python Focus**: This is a Python-centric project based on the comprehensive `.gitignore` configuration
- **Multiple Tooling Support**: The `.gitignore` includes support for various Python package managers and tools:
  - **Package Management**: UV, Poetry, PDM, Pixi, Pipenv
  - **Code Quality**: Ruff, MyPy, Pytype
  - **Development**: Jupyter notebooks, Marimo, IPython
  - **Process Automation**: Abstra framework integration
  - **Testing**: pytest, coverage, tox, nox

## Project Conventions
- **Environment Management**: The project supports multiple Python environment tools (venv, conda, poetry, uv, pdm, pixi)
- **Code Quality**: Configured for Ruff linting and various type checkers (MyPy, Pytype)
- **Notebook Support**: Ready for Jupyter and Marimo notebook development
- **AI Tooling**: Includes Cursor and Abstra configurations for AI-assisted development

## Key Files and Structure
- `README.md`: Minimal project documentation (currently just title)
- `.gitignore`: Comprehensive Python development exclusions with modern tooling support
- Project appears to be in bootstrap phase - no source code or configuration files yet

## Development Workflow Guidance
When adding code to this project:

1. **Package Management**: Choose one primary tool (UV, Poetry, PDM, or Pixi) and create appropriate configuration files
2. **Project Structure**: Follow standard Python package layout with `src/amos/` or `amos/` for main code
3. **Configuration**: Add `pyproject.toml` for modern Python project configuration
4. **Testing**: Set up pytest with coverage reporting
5. **Code Quality**: Configure Ruff for linting and formatting
6. **Type Checking**: Add MyPy configuration for static type analysis

## Notable Patterns
- The `.gitignore` suggests this project may use Abstra for process automation
- Support for multiple environment managers indicates flexibility in deployment scenarios
- Marimo support suggests potential for reactive notebook-style development

## Next Steps for Development
The project is ready for initial implementation. Consider:
1. Define project purpose and update README.md
2. Choose and configure package management tool
3. Set up basic project structure
4. Add CI/CD configuration in `.github/workflows/`
5. Define project dependencies and requirements