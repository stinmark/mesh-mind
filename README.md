Here is a complete, production-ready `README.md` for **Mesh-Mind** that highlights your Go Gateway, Python ONNX inference engine, CLI with SQLite local storage, and server-side SQLite authentication setup.

---

```markdown
# 🧠 Mesh-Mind

> An ultra-lightweight, modular MLOps pipeline and API Gateway designed for high-performance model routing, local developer productivity, and low-latency inference.

Mesh-Mind combines a **Go API Gateway** with **Python ONNX model workers** and a **Go Terminal CLI**. It features dual-layer SQLite persistence: client-side SQLite for local CLI state and history, and server-side SQLite for user authentication and activity logs.

---

## 🏗️ Architecture Overview

```text
┌─────────────────────────────────────────────────────────────────┐
│                          LOCAL CLIENT                           │
│                                                                 │
│  Mesh-Mind CLI (Go + Cobra)                                     │
│  └── Storage: Local SQLite (~/.mesh-mind/config.db)             │
│      ├── Settings (Gateway URL, JWT Tokens)                     │
│      └── Prediction History & Latency Logs                      │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 │ HTTP REST API Request
                                 │ Header: "Authorization: Bearer <JWT>"
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                         CLOUD / SERVER                          │
│                                                                 │
│  Go API Gateway & Router                                        │
│  ├── Auth Middleware (JWT Verification)                         │
│  └── Storage: Server SQLite (/data/gateway.db)                  │
│      ├── User Credentials & bcrypt Passwords                    │
│      └── User Activity Tracking                                 │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 │ Internal Routing / Inter-Process Call
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│  Python ONNX Model Workers                                      │
│  └── High-Performance Sentiment & Intent Inference              │
└─────────────────────────────────────────────────────────────────┘

```

---

## ✨ Features

* **🚀 Go API Gateway:** High-throughput, stateless routing engine with custom JWT middleware.
* **⚡ ONNX Model Runtime:** Rapid Python-backed inference optimized for intent and sentiment routing.
* **💻 Interactive Terminal CLI:** Built using Cobra for intuitive developer workflows.
* **💾 Dual-Layer SQLite Persistence:**
* **Client-Side (`~/.mesh-mind/config.db`):** Stores user tokens, endpoints, and local query logs for offline inspection.
* **Server-Side (`/data/gateway.db`):** Light, zero-dependency user registration, password hashing (`bcrypt`), and token issuing.


* **🔓 Developer-Friendly Auth Toggle:** Disable authentication with `DISABLE_AUTH=true` for instant local testing.
* **🌐 Web Ready:** Clean REST endpoints with CORS support, prepared for a React/Next.js frontend.

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── mesh-mind/         # Binary entrypoint for the CLI
├── internal/              # Shared CLI internal modules
│   ├── cli/               # Cobra CLI commands (predict, register, login, history)
│   ├── config/            # SQLite initialization (pure Go modernc.org/sqlite)
│   └── store/             # Local SQLite Data Access Objects (Config & History)
├── gateway/               # Go API Gateway
│   ├── internal/auth/     # JWT authentication middleware & handlers
│   └── internal/db/       # Gateway SQLite schema & user store
├── models/                # Python ONNX inference services
├── docker-compose.yml     # Multi-container orchestration
└── README.md

```

---

## 🚀 Quick Start

### 1. Prerequisites

* **Go** (v1.21 or higher)
* **Docker & Docker Compose** (for Gateway & Model services)

### 2. Start the Gateway & Inference Workers

Clone the repository and run the services with Docker Compose:

```bash
git clone [https://github.com/your-username/mesh-mind.git](https://github.com/your-username/mesh-mind.git)
cd mesh-mind
docker-compose up -d --build

```

### 3. Build & Run the CLI

Build the pure Go CLI binary (no CGO compiler required):

```bash
go build -o mesh-mind ./cmd/mesh-mind

```

---

## 💻 CLI Usage Walkthrough

### Step 1: Set Up Target Gateway URL

```bash
./mesh-mind config set-url http://localhost:8080
# Output: ✔ Default Gateway URL set to: http://localhost:8080

```

### Step 2: Create an Account & Log In

```bash
# Register a new user
./mesh-mind register -u devuser -p supersecret
# Output: ✔ Account created successfully!

# Authenticate and receive a JWT token (saved to local SQLite)
./mesh-mind login -u devuser -p supersecret
# Output: ✔ Login successful!
# Output: ✔ JWT token saved to local SQLite database (~/.mesh-mind/config.db)

```

### Step 3: Execute Model Predictions

```bash
./mesh-mind predict "Great customer service, solved my issue fast!"

# 🧠 Model Output:
#    Intent:     support
#    Sentiment:  positive
#    Model Used: onnx-sentiment-int8
#    Latency:    18.40 ms
# 
# Saved to local SQLite history.

```

### Step 4: Query Local History (Offline)

```bash
./mesh-mind history --limit 5

# 📜 Last 1 Local Predictions:
# ---------------------------------------------------------------------
# [2026-09-28 16:30:00] Text: "Great customer service, solved my issue fast!"
#    └─ Intent: support | Sentiment: positive | Model: onnx-sentiment-int8 (18.4 ms)

```

---

## 🔒 Security & Authentication

* **Password Hashing:** User passwords are encrypted using `bcrypt` on the server before being saved to SQLite.
* **Stateless Verification:** Model routes are secured via JWT bearer tokens signed with HS256 algorithm.
* **Local Client Isolation:** Auth credentials and local query history remain inside the developer's home directory (`~/.mesh-mind/config.db`), preventing leakage between users.
* **Bypass Auth for Testing:** Set the environment variable `DISABLE_AUTH=true` in `gateway/.env` to run the stack without authentication checks.

---

## 🛣️ Roadmap

* [x] Go API Gateway with JWT Auth
* [x] Python ONNX Model Worker Integration
* [x] Pure Go CLI with SQLite Persistent Storage
* [ ] Add Web Dashboard (React/Next.js)
* [ ] gRPC Transport between Gateway and ONNX Workers

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for details.

```

<ElicitationsGroup message="Where would you like to go next with Mesh-Mind?">
  <Elicitation label="Create a docker-compose.yml file for the Gateway and Python services" query="Create a complete docker-compose.yml file that spins up the Go Gateway and Python ONNX model workers with local persistent storage."/>
  <Elicitation label="Add unit tests for the Go CLI and SQLite stores" query="Write unit tests for the CLI SQLite storage implementations in internal/store using Go's standard testing package."/>
  <Elicitation label="Build a lightweight React frontend dashboard" query="Show me how to build a simple React single-page frontend that connects to the Go Gateway for real-time predictions."/>
</ElicitationsGroup>

```
