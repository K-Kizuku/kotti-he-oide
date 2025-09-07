PY := uv run python
UV := uv

.PHONY: setup dev run lint type test

setup:
	$(UV) sync

dev:
	$(UV) run python -m app.server

run:
	$(UV) run python -m app.server

lint:
	$(UV) run ruff check .

type:
	$(UV) run mypy app

test:
	$(UV) run pytest -q
