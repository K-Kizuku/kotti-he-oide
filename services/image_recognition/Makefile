PY := uv run python
UV := uv

.PHONY: setup dev run lint type test proto

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

proto:
	# Buf が使えない場合は grpcio-tools で生成
	$(UV) run python -m grpc_tools.protoc -I ../../schema/proto --python_out app/gen --grpc_python_out app/gen ../../schema/proto/image_recognition/v1/image_recognition.proto
