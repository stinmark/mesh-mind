# 🧠 Mesh-Mind

> An ultra-lightweight, modular MLOps pipeline and API Gateway designed for intelligent model routing, local developer productivity, and low-latency ONNX inference.

Mesh-Mind combines a **Go API Gateway** with **Python ONNX model workers** and an interactive **Go Terminal CLI**. It features smart model routing (auto-selection based on query intent or explicit user overrides) and dual-layer SQLite persistence.

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
│  Go API Gateway & Router                                        │
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
│ Model Worker               │   │ Model Worker               │
└────────────────────────────┘   └────────────────────────────┘

```

---

## ✨ Features

* **🔀 Smart Model Routing:** Automatically evaluates incoming text to dispatch queries to the optimal model worker, or allows clients to explicitly target specific models via `--model` flags.
* **🚀 Go API Gateway:** High-throughput, stateless routing engine with custom JWT authentication middleware.
* **⚡ ONNX Model Runtime:** Rapid Python-backed inference optimized for intent classification and sentiment analysis.
* **💻 Interactive Terminal CLI:** Built with Cobra for seamless terminal workflows.
* **💾 Dual-Layer SQLite Persistence:**
* **Client-Side (`~/.mesh-mind/config.db`):** Stores JWT tokens, server configuration, and local query execution history for offline auditing.
* **Server-Side (`/data/gateway.db`):** Handles user registration, `bcrypt` password hashing, and token validation.


* **🔓 Developer Auth Toggle:** Disable authentication checks during development by setting `DISABLE_AUTH=true`.

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── mesh-mind/         # Binary entrypoint for the CLI
├── internal/              # Shared CLI internal modules
│   ├── cli/               # Cobra CLI commands (register, login, predict, history, config)
│   ├── config/            # Pure Go SQLite initialization (modernc.org/sqlite)
│   └── store/             # Local SQLite Data Access Objects (Config & History)
├── gateway/               # Go API Gateway & Smart Router
│   ├── internal/auth/     # JWT authentication middleware & handlers
│   ├── internal/handlers/ # Predict & routing handlers
│   └── internal/db/       # Server SQLite schema & user repository
├── models/                # Python ONNX inference services
├── docker-compose.yml     # Multi-container orchestration
└── README.md

```

---

## 🚀 Quick Start

### 1. Prerequisites

* **Go** (v1.21 or higher)
* **Docker & Docker Compose**

### 2. Start Services

Clone the repository and spin up the gateway and inference workers:

```bash
git clone [https://github.com/your-username/mesh-mind.git](https://github.com/your-username/mesh-mind.git)
cd mesh-mind
docker-compose up -d --build

```

### 3. Build the CLI Binary

Build the pure Go CLI binary (no CGO/GCC compiler required):

```bash
go build -o mesh-mind ./cmd/mesh-mind

```

---

## 💻 CLI Usage Walkthrough

### 1. Set Gateway Target

```bash
./mesh-mind config set-url http://localhost:8080
# Output: ✔ Gateway URL set to: http://localhost:8080

```

### 2. Register & Authenticate

```bash
# Register a new account
./mesh-mind register -u devuser -p supersecret
# Output: ✔ Account created successfully!

# Log in and store JWT locally
./mesh-mind login -u devuser -p supersecret
# Output: ✔ Login successful!
# Output: ✔ JWT token saved to local SQLite database (~/.mesh-mind/config.db)

```

### 3. Execute Model Predictions

#### Option A: Automatic Model Routing (Default)

Let the Gateway automatically detect the query type and route it to the appropriate ONNX model worker:

```bash
./mesh-mind predict "How do I reset my account password?"

# 🧠 Model Output:
#    Text:       "How do I reset my account password?"
#    Model Used: onnx-intent-v1 (auto-routed)
#    Intent:     account_recovery
#    Latency:    12.40 ms

```

#### Option B: Explicit Model Override (`--model` / `-m`)

Bypass auto-detection and force the request to execute against a specific target model (`sentiment`, `intent`):

```bash
./mesh-mind predict "The new update solved all my issues!" --model sentiment

# 🧠 Model Output:
#    Text:       "The new update solved all my issues!"
#    Model Used: onnx-sentiment-v1 (forced target)
#    Sentiment:  positive
#    Latency:    9.80 ms

```

### 4. Query Local History

Inspect previously run queries and execution metrics stored in your local SQLite database:

```bash
./mesh-mind history --limit 5

# 📜 Last Predictions:
# ---------------------------------------------------------------------
# [2026-09-28 16:30:00] Text: "How do I reset my account password?"
#    └─ Model: onnx-intent-v1 | Intent: account_recovery | Latency: 12.4 ms
# [2026-09-28 16:31:12] Text: "The new update solved all my issues!"
#    └─ Model: onnx-sentiment-v1 | Sentiment: positive | Latency: 9.8 ms

```

---

## ⚙️ REST API Endpoint

If accessing the gateway directly via HTTP:

**`POST /api/v1/predict`**

```json
// Request Body
{
  "text": "Where is my order package?",
  "model": "auto"  // Acceptable values: "auto", "sentiment", "intent"
}

```

```json
// Response Body
{
  "text": "Where is my order package?",
  "model_used": "onnx-intent-v1",
  "intent": "order_tracking",
  "latency_ms": 11.20
}

```

---

## 🔒 Security & Authentication

* **Password Hashing:** User passwords are encrypted using `bcrypt` on the server before storage.
* **Stateless Authorization:** Routes are protected via JWT bearer tokens signed with HS256.
* **Client Storage Isolation:** Auth tokens and local query logs are isolated inside `~/.mesh-mind/config.db`.
* **Bypass Auth Mode:** Set `DISABLE_AUTH=true` in `gateway/.env` to run local tests without login requirements.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for details.

