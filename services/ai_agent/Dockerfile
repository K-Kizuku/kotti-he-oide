# Use uv's ARM64 Python base image
FROM --platform=linux/arm64 ghcr.io/astral-sh/uv:python3.12-bookworm-slim

WORKDIR /app

# Copy uv files
COPY pyproject.toml ./
# Note: uv.lock will be generated after first uv sync

# Install dependencies
RUN uv sync --frozen --no-cache || uv sync --no-cache

# Copy application code
COPY app/ ./app/
COPY config/ ./config/
COPY tools/ ./tools/
COPY models/ ./models/
COPY services/ ./services/
COPY utils/ ./utils/

# Expose port (AgentCore Runtime requirement)
EXPOSE 8080

# Run application with ADOT instrumentation for observability
CMD ["opentelemetry-instrument", "uv", "run", "uvicorn", "app.server:app", "--host", "0.0.0.0", "--port", "8080"]
