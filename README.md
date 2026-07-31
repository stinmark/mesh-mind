# MESH MIND

Markdown

# Inference Microservice (`inference-py`)

A lightweight FastAPI inference service serving the `HuggingFaceTB/SmolLM2-135M-Instruct` LLM for the Mesh-Mind microservices platform.

---

## 🛠️ Local Setup & Testing

### Prerequisites

- Python 3.10+
- `pip` and `venv`

### 1. Create & Activate Virtual Environment

```bash
# Navigate to the inference service directory
cd apps/inference-py

# Create virtual environment
python -m venv .venv

# Activate virtual environment
source .venv/bin/activate

2. Install Dependencies

To avoid downloading heavy CUDA/NVIDIA binaries when testing on CPU, install using the PyTorch CPU index:
Bash

pip install -r requirements.txt

3. Run the Microservice
Bash

python app/main.py

The server will start on http://localhost:8000 and automatically load SmolLM2-135M-Instruct into memory (~270 MB weight download on first boot).
🐳 Docker Setup & Testing

If you prefer testing inside a containerized environment without installing Python locally:
1. Build Docker Image
Bash

# From the project root directory
docker build -t mesh-mind-inference -f apps/inference-py/Dockerfile apps/inference-py

2. Run Docker Container
Bash

docker run -p 8000:8000 mesh-mind-inference

🧪 Testing the API

Once the service is running (locally or via Docker), open a new terminal window to test the endpoints.
Health Check (Kubernetes Liveness Probe)
Bash

curl http://localhost:8000/healthz

Expected Response:
JSON

{
  "status": "ok",
  "model": "HuggingFaceTB/SmolLM2-135M-Instruct",
  "device": "cpu"
}

Chat Completion Endpoint
Bash

curl -X POST http://localhost:8000/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "system", "content": "You are a concise AI assistant."},
      {"role": "user", "content": "Explain microservices in one sentence."}
    ],
    "max_new_tokens": 128,
    "temperature": 0.7
  }'

Expected Response:
JSON

{
  "response": "Microservices are a software architecture approach where an application is built as a collection of small, independent, and loosely coupled services.",
  "model": "HuggingFaceTB/SmolLM2-135M-Instruct"
}


---

<ElicitationsGroup message="Where would you like to update next?">
  <Elicitation label="Draft deploy/docker-compose.yml to run gateway + inference together" query="Draft deploy/docker-compose.yml to run the Go gateway and Python inference microservices together."/>
  <Elicitation label="Draft the root project README.md to document the entire system" query="Draft the root README.md for the mesh-mind repository covering the architecture and all microservices."/>
</ElicitationsGroup>
```
