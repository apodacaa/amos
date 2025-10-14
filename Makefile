.PHONY: install format lint run help

help:
	@echo "Targets: install, format, lint, run"

install:
	python -m venv .venv
	. .venv/bin/activate && python -m pip install --upgrade pip
	@echo "Activated venv: source .venv/bin/activate; then pip install -e . or install deps."

format:
	black .

lint:
	ruff .

run:
	python -m textual run app:App  # adjust to your entrypoint as needed