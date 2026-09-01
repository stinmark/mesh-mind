Here is the updated root `README.md` covering the full project architecture, setup steps for both `gateway-go` and `inference-py`, and deployment guidelines for AWS Free Tier EC2.

````markdown
# 🧠 MESH MIND

A lightweight, low-memory AI platform designed to run within the **AWS Free Tier (1 vCPU, 1 GiB RAM)**.

Mesh Mind uses an embedded SQLite database, a compiled Go API gateway serving a dynamic frontend, and a Python inference service leveraging GGUF quantization (`llama-cpp-python`) to keep total system memory usage under **250 MB RAM**.

---

## 🏗 Architecture & Memory Footprint

```text
┌────────────────────────────────────────────────────────────────────────┐
│              AWS EC2 Instance (t2.micro / t3.micro: 1 GiB RAM)          │
│                                                                        │
│  ┌─────────────────────────────────┐   ┌────────────────────────────┐  │
│  │ gateway-go (Port 8080)          │   │ inference-py (Port 8000)   │  │
│  │ - JWT Auth & Dummy Token Store  │   │ - FastAPI + llama-cpp      │  │
│  │ - Embedded HTML/JS UI           │ ─ │ - SmolLM2-135M-Instruct    │  │
│  │ - Embedded SQLite (mesh_mind.db)│   │   (Q4_K_M GGUF Binary)     │  │
│  │ ~30 MB RAM                      │   │ ~180 MB RAM                │  │
│  └─────────────────────────────────┘   └────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────┘
```
````

---

## 📁 Repository Structure

```text
.
├── apps/
│   ├── gateway-go/          # Go API Gateway + Embedded UI + SQLite
│   │   ├── Dockerfile
│   │   ├── main.go
│   │   ├── go.mod
│   │   └── internal/
│   │       ├── auth/        # JWT generation & authentication middleware
│   │       ├── db/          # SQLite driver initialization
│   │       ├── handlers/    # API endpoints & proxy controller
│   │       └── templates/   # Embedded HTML/JS frontend
│   └── inference-py/        # Python GGUF LLM inference runner
│       ├── Dockerfile
│       ├── requirements.txt
│       ├── models/          # Quantized GGUF model binaries
│       └── app/
│           └── main.py      # FastAPI server using llama-cpp-python

```

---

## 🛠️ Local Setup & Running

### Prerequisites

- Go 1.22+
- Python 3.10+
- `curl`

---

### 1. Setup & Run `inference-py`

```bash
cd apps/inference-py

# Create & activate virtual environment
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Download SmolLM2-135M-Instruct Q4_K_M GGUF model (~100 MB)
mkdir -p models
curl -L -o models/smollm2-135m-instruct-q4_k_m.gguf \
  [https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf](https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf)

# Start the inference server
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

```

The inference server will start at `http://localhost:8000`.

---

### 2. Setup & Run `gateway-go`

Open a new terminal window:

```bash
cd apps/gateway-go

# Download Go module dependencies
go mod tidy

# Start the API gateway and Web UI
go run main.go

```

The application will start at `http://localhost:8080`.

---

## 🧪 Testing the Application

1. Open **`http://localhost:8080`** in your browser.
2. **Register** a new user account (automatically receives **10 free tokens**).
3. **Log in** with your credentials.
4. Send a prompt to test the typing effect powered by `SmolLM2-135M`.
5. Click **+ Buy 20 Dummy Tokens** to test SQLite-backed balance updates.

---

## 🐳 Docker Deployment

Both microservices include optimized Dockerfiles designed to minimize runtime overhead.

### Build and Run `inference-py`

```bash
cd apps/inference-py
docker build -t mesh-mind-inference .
docker run -p 8000:8000 mesh-mind-inference

```

### Build and Run `gateway-go`

```bash
cd apps/gateway-go
docker build -t mesh-mind-gateway .
docker run -p 8080:8080 -e INFERENCE_URL="[http://host.docker.internal:8000](http://host.docker.internal:8000)" mesh-mind-gateway

```

---

## ☁️ AWS Free Tier Deployment Guidelines (`t2.micro` / `t3.micro`)

1. **Enable Swap Memory:** AWS Free Tier micro instances only provide 1 GiB RAM. Enable a 2 GB swap file on the EC2 instance OS to prevent unexpected Out-Of-Memory (OOM) kills.
2. **Pre-build Images Off-Server:** Build Docker images locally or via GitHub Actions and push them to **Amazon ECR**. Compiling Go code or installing C++ wheels inside a `t2.micro` will exhaust system memory.
3. **Set Container Resource Limits in ECS:**

- `inference-py`: Hard limit `384 MiB`
- `gateway-go`: Hard limit `128 MiB`

```

```
