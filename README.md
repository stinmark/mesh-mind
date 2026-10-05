# 🧠 Mesh-Mind

> An ultra-lightweight, modular MLOps pipeline and API Gateway designed for intelligent model routing, local developer productivity, zero-PyTorch low-latency ONNX inference, and strict resource constraints.

Mesh-Mind combines a high-concurrency **Go API Gateway** with decoupled **Python ONNX microservices** and an interactive **Go Terminal CLI**. It features smart model routing, multi-stage container optimization, zero runtime PyTorch dependencies, and dual-layer SQLite state management.

---

## 🏗️ Architecture Overview

```text
┌─────────────────────────────────────────────────────────────────┐
│                          LOCAL CLIENT                           │
│                                                                 │
│  Mesh-Mind CLI (Go + Cobra)                                     │
│  └── Storage: Local SQLite (~/.mesh-mind/config.db)             │
│      ├── Settings (Gateway URL, Saved JWT Token)                │
│      └── Prediction History & Latency Logs                      │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 │ HTTP POST /api/v1/predict
                                 │ Body: { "text": "...", "model": "auto" }
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                         CLOUD / SERVER                          │
│                                                                 │
│  Go API Gateway & Router (Port 8080)                            │
│  ├── JWT Authentication Middleware                              │
│  ├── Dynamic Router Engine (Auto-Detect vs. Explicit Override)  │
│  └── Storage: Server SQLite (/data/gateway.db)                  │
│      └── User Credentials (bcrypt) & Activity Logs              │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                   ┌─────────────┴─────────────┐
                   ▼                           ▼
┌────────────────────────────┐   ┌────────────────────────────┐
│ Python ONNX Sentiment      │   │ Python ONNX Intent/Emotion │
│ Model Worker (256MB Limit) │   │ Model Worker (512MB Limit) │
│ └── Pure onnxruntime +     │   │ └── Pure onnxruntime +     │
│     tokenizers (No PyTorch)│   │     tokenizers (No PyTorch)│
└────────────────────────────┘   └────────────────────────────┘

```

---

## ✨ Key MLOps Features

* **⚡ Zero-PyTorch Runtime Inference:** Model workers run on pure C++ backends (`onnxruntime` + Rust-backed `tokenizers`), dropping runtime memory usage dramatically while serving low-latency CPU predictions.
* **📦 Multi-Stage Container Optimization:** PyTorch is used strictly during model conversion (`export_models.py`) to convert Hugging Face checkpoints to ONNX. Production worker images contain zero PyTorch/Torchvision binaries.
* **🔀 Smart & Explicit Model Forwarding:** Go Gateway parses input text to route requests to the appropriate ONNX worker (`model-sentiment` vs. `model-intent`) or accepts target overrides (`--model`).
* **🧠 Real Dynamic Logit Decoding:** Workers calculate softmax probability distributions over ONNX output tensors in real time and map class IDs to human-readable strings using dynamic `id2label` dictionary parsing with built-in fallbacks.
* **🏥 Production Healthchecks & Startup Synchronization:** Docker Compose monitors FastAPI worker readiness via Uvicorn HTTP probes (`python3 -c urllib`), preventing gateway cold-start connection race conditions.
* **💻 Interactive Terminal CLI:** Built with Go + Cobra for developer workflows.
* **💾 Dual-Layer SQLite Persistence:**
* **Client-Side (`~/.mesh-mind/config.db`):** Stores JWT tokens, server configuration, and local prediction history.
* **Server-Side (`/data/gateway.db`):** Handles user registration, `bcrypt` password hashing, and audit logs.


* **🛡️ Calibrated Container Memory Allocations:** Stack resource limits are finely tuned (64MB for Go Gateway, 256MB for MiniLM Sentiment, 512MB for BERT Emotion) to prevent OOM kills during initial model tensor load.

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── mesh-mind/              # Binary entrypoint for the CLI
├── internal/                   # Shared CLI internal modules
│   ├── cli/                    # Cobra commands (register, login, predict, history, config)
│   ├── config/                 # Pure Go SQLite initialization (modernc.org/sqlite)
│   └── store/                  # Local SQLite Data Access Objects
├── gateway/                    # Go API Gateway & Dynamic Router
│   ├── internal/auth/          # JWT authentication middleware
│   ├── internal/handlers/      # Predict & HTTP forwarding handlers
│   └── internal/db/            # Server SQLite schema & user repository
├── models/                     # Python ONNX inference services
│   ├── app.py                  # Zero-PyTorch FastAPI worker (onnxruntime + tokenizers)
│   ├── export_models.py        # ONNX export script (MiniLM SST-2 & BERT Emotion)
│   ├── requirements-export.txt # Export tooling dependencies (Optimum + Transformers)
│   ├── requirements-runtime.txt# Production worker dependencies (ONNXRuntime + Fast Tokenizers)
│   └── Dockerfile              # Parameterized multi-stage container build
├── docker-compose.yml          # Microservice orchestration, healthchecks & limits
└── README.md

```

---

## 🚀 Quick Start

### 1. Prerequisites

* **Go** (v1.21 or higher)
* **Docker & Docker Compose**

### 2. Build & Launch Microservices Stack

Spin up the Go Gateway and zero-PyTorch Python workers in containerized microservices:

```bash
# Build and run containers in detached mode
docker compose up -d --build

```

Verify service readiness and container health status:

```bash
docker compose ps

```

*Expected output:*

```text
NAME                          SERVICE           STATUS
mesh-mind-gateway-1           gateway           Up (healthy)
mesh-mind-model-sentiment-1   model-sentiment   Up (healthy)
mesh-mind-model-intent-1      model-intent      Up (healthy)

```

### 3. Build the Go CLI Binary

Compile the CLI binary:

```bash
go build -o mesh-mind ./cmd/mesh-mind

```

---

## 💻 CLI Usage Walkthrough

### 1. Configure Gateway Target

```bash
./mesh-mind config set-url http://localhost:8080
# Output: ✔ Gateway URL set to: http://localhost:8080

```

### 2. Register & Authenticate

```bash
# Register a new account
./mesh-mind register -u devuser -p supersecret123
# Output: ✔ Account created successfully!

# Log in and store JWT locally
./mesh-mind login -u devuser -p supersecret123
# Output: ✔ Login successful!
# Output: ✔ JWT token saved to local SQLite database (~/.mesh-mind/config.db)

```

### 3. Execute Model Predictions

#### Option A: Automatic Model Routing (Default)

The Gateway automatically evaluates the query and forwards it to the sentiment worker:

```bash
./mesh-mind predict "I hate Milk"

# 🧠 Model Response:
# {"text":"I hate Milk","model_used":"onnx-sentiment-v1","sentiment":"negative","latency_ms":54.48}

```

```bash
./mesh-mind predict "I love milk"

# 🧠 Model Response:
# {"text":"I love milk","model_used":"onnx-sentiment-v1","sentiment":"positive","latency_ms":3.12}

```

#### Option B: Explicit Model Override (`--model` / `-m`)

Force execution against the BERT Emotion/Intent classification microservice (`model-intent`):

```bash
./mesh-mind predict "How do I reset my account password?" --model intent

# 🧠 Model Response:
# {"text":"How do I reset my account password?","model_used":"onnx-intent-v1","intent":"anger","latency_ms":179.96}

```

---

## 🔍 Observability & Live Monitoring

### Stream Live Gateway & Worker Forwarding Logs

Watch incoming HTTP requests and ONNX runtime executions in real time:

```bash
docker compose logs -f gateway model-sentiment model-intent

```

### Verify Container Resource Usage

Check live RAM utilization across running containers:

```bash
docker stats

```

---

## 🔒 Security & System Configuration

* **Password Security:** User passwords are hashed using `bcrypt` on the Go Gateway before storage.
* **Stateless Authorization:** Routes are secured using JWT bearer tokens (HS256).
* **Development Auth Bypass:** Set `DISABLE_AUTH=true` in gateway environment settings to bypass login during rapid testing.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for details.

```

```
