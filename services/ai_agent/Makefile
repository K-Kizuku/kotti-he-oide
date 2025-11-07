PY := uv run python
UV := uv

.PHONY: setup dev run lint type test

setup:
	$(UV) sync

dev:
	$(UV) run uvicorn app.server:app --host 0.0.0.0 --port 8080 --reload

run:
	$(UV) run uvicorn app.server:app --host 0.0.0.0 --port 8080

lint:
	$(UV) run ruff check .

type:
	$(UV) run mypy app

test:
	$(UV) run pytest -q
