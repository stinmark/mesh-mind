FROM python:3.11-slim

WORKDIR /app

# Install build tools for compiling llama-cpp-python C++ bindings & curl for downloading the model
RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential \
  curl \
  && rm -rf /var/lib/apt/lists/*

# Install lightweight dependencies (llama-cpp-python, fastapi, uvicorn)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Create models directory and pre-download the GGUF model inside the Docker image build
RUN mkdir -p models && \
  curl -L -o models/smollm2-135m-instruct-q4_k_m.gguf \
  https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf

# Copy application code
COPY app/ ./app/

# Environment configurations
ENV PYTHONUNBUFFERED=1 \
  PORT=8000 \
  MODEL_PATH="./models/smollm2-135m-instruct-q4_k_m.gguf"

EXPOSE 8000

# Start server using module notation matching your project layout (app.main:app)
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
