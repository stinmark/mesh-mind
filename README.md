# 🧠 Mesh-Mind

> An ultra-lightweight, modular MLOps pipeline and API Gateway designed for intelligent model routing, local developer productivity, zero-PyTorch low-latency ONNX inference, and strict resource constraints (<120MB total RAM stack).

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
│ Python ONNX Sentiment      │   │ Python ONNX Intent         │
│ Model Worker (~50MB RAM)   │   │ Model Worker (~50MB RAM)   │
│ └── Pure onnxruntime +     │   │ └── Pure onnxruntime +     │
│     tokenizers (No PyTorch)│   │     tokenizers (No PyTorch)│
└────────────────────────────┘   └────────────────────────────┘

```

---

## ✨ Key MLOps Features

* **⚡ Zero-PyTorch Runtime Inference:** Model workers run on pure C++ backends (`onnxruntime` + Rust-backed `tokenizers`), dropping runtime memory usage to **~50MB RAM per worker**.
* **📦 Multi-Stage Container Optimization:** PyTorch is used strictly during the Docker **Build Stage** (`export_models.py`) to convert Hugging Face checkpoints to ONNX, then discarded. Production images install zero PyTorch binaries.
* **🔀 Smart & Explicit Model Routing:** Gateway automatically evaluates incoming text to dispatch queries to the optimal model worker (`model-sentiment` vs. `model-intent`), or supports explicit overrides (`--model`).
* **🚀 Lightweight MiniLM Models:** Replaced heavy transformers with MiniLM architectures (~45MB on disk each), optimizing for sub-10ms inference latencies and low cloud resource costs.
* **💻 Interactive Terminal CLI:** Built with Go + Cobra for developer workflows.
* **💾 Dual-Layer SQLite Persistence:**
* **Client-Side (`~/.mesh-mind/config.db`):** Stores JWT tokens, server configuration, and local prediction history.
* **Server-Side (`/data/gateway.db`):** Handles user registration, `bcrypt` password hashing, and audit logs.


* **🛡️ Strict Resource Limits:** Docker Compose applies tight memory limits (80MB per worker, 64MB for gateway) ensuring the full stack runs under **120MB total RAM**.

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
│   ├── internal/handlers/      # Predict & routing handlers
│   └── internal/db/            # Server SQLite schema & user repository
├── models/                     # Python ONNX inference services
│   ├── app.py                  # Zero-PyTorch FastAPI worker (onnxruntime + tokenizers)
│   ├── export_models.py        # ONNX export script (used in Docker build stage)
│   ├── requirements-export.txt # Export tooling dependencies (Optimum + Transformers)
│   ├── requirements-runtime.txt# Production worker dependencies (ONNXRuntime + Fast Tokenizers)
│   └── Dockerfile              # Parameterized multi-stage container build
├── docker-compose.yml          # Microservice orchestration with memory limits
└── README.md

```

---

## 🚀 Quick Start

### 1. Prerequisites

* **Go** (v1.21 or higher)
* **Docker & Docker Compose**

### 2. Local Model Export (Optional Standalone Step)

To test model export on your host machine before building containers:

```bash
cd models

# Create and activate virtual environment
python3 -m venv venv
source venv/bin/activate

# Install export dependencies (includes PyTorch for conversion)
pip install -r requirements-export.txt

# Export sentiment and intent models to ONNX
python export_models.py --model sentiment
python export_models.py --model intent
cd ..

```

### 3. Build & Launch Microservices Stack

Spin up the Go Gateway and zero-PyTorch Python workers in containerized microservices:

```bash
# Build and run containers
docker compose up --build

```

*(Note: If using Docker Compose v1, use `docker-compose up --build` instead.)*

### 4. Build the Go CLI Binary

In a new terminal window, compile the CLI binary:

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

The Gateway automatically classifies the query and dispatches it to the appropriate model worker:

```bash
./mesh-mind predict "How do I reset my account password?"

# 🧠 Model Output:
#    Text:       "How do I reset my account password?"
#    Model Used: onnx-intent-v1 (auto-routed)
#    Intent:     sadness / query
#    Latency:    6.40 ms

```

#### Option B: Explicit Model Override (`--model` / `-m`)

Force the request to execute against a specific microservice (`sentiment` or `intent`):

```bash
./mesh-mind predict "The new update solved all my issues!" --model sentiment

# 🧠 Model Output:
#    Text:       "The new update solved all my issues!"
#    Model Used: onnx-sentiment-v1 (forced target)
#    Sentiment:  positive
#    Latency:    5.20 ms

```

### 4. Inspect Local Prediction History

Query locally cached execution history and latency metrics from `~/.mesh-mind/config.db`:

```bash
./mesh-mind history --limit 5

# 📜 Last Predictions:
# ---------------------------------------------------------------------
# [2026-09-29 16:30:00] Text: "How do I reset my account password?"
#    └─ Model: onnx-intent-v1 | Label: sadness | Latency: 6.4 ms
# [2026-09-29 16:31:12] Text: "The new update solved all my issues!"
#    └─ Model: onnx-sentiment-v1 | Label: positive | Latency: 5.2 ms

```

---

## 🔍 Verification & System Auditing

### Verify Zero-PyTorch Execution in Memory

To verify that PyTorch is not loaded into memory inside the worker processes:

```bash
# Inspect container RAM usage (should remain around ~50MB - 60MB per worker)
docker stats mesh-mind-model-sentiment mesh-mind-model-intent

```

Direct HTTP healthcheck on a worker:

```bash
curl http://localhost:5000/health
# Response: {"status":"healthy","model_name":"sentiment","runtime":"onnxruntime-cpu","model_loaded":true}

```

---

## 🔒 Security & System Configuration

* **Password Security:** User passwords are hashed using `bcrypt` on the Go Gateway before storage.
* **Stateless Authorization:** Routes are secured using JWT bearer tokens (HS256).
* **Development Auth Bypass:** Set `DISABLE_AUTH=true` in gateway environment settings to bypass login during testing.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for details.

```

```
