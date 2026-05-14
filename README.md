# 🚀 GoNode - Adaptive Infrastructure Engine

[![Go Version](https://img.shields.io/badge/Go-1.26.3-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/Node-18+-green.svg)](https://nodejs.org)
[![OS](https://img.shields.io/badge/OS-Ubuntu%20%7C%20Debian-orange.svg)](#)

GoNode is a high-performance orchestration engine written in Go, designed to manage, monitor, and scale Node.js applications with adaptive resource profiling.

---

## ✨ New: Smart AI Detection
GoNode now features **Smart Scan**, a heuristic-based detection system that automatically identifies your project type:
- **Frontend**: Detects Next.js or React and configures `npm start` automatically.
- **Backend**: Scans for `app.js`, `server.js`, or `index.js` and configures the Node.js runtime.

---

## 📐 High-Level Architecture

```mermaid
graph TD
    User((🌐 CLI User)) -->|Socket| GoNode[("🚀 GoNode Engine")]
    
    subgraph "Intelligent Orchestration"
    GoNode -->|1. Smart Scan| Detector[Next.js / Node.js / React]
    GoNode -->|2. Profile Selection| Spec[Eco / Balanced / Power]
    GoNode -->|3. Process Spawn| NodeApp["📦 Managed App"]
    end

    subgraph "Observability"
    NodeApp -->|Output| Logger["📝 Timestamped Logger"]
    Logger -->|Rotation| LogFile["gonode.log (1MB Limit)"]
    end
```

---

## 📂 Project Structure (Modular Design)

```text
GoNode/
├── main.go            # Entry Point
├── pkg/
│   ├── engine/        # Core: cli.go, daemon.go, detector.go
│   ├── logger/        # Logging & Rotation logic
│   └── utils/         # UI & Installer Helpers
├── setup.sh           # Environment Setup
└── install.sh         # Build Script
```

---

## 🚀 Quick Start

### 1. Setup Environment
```bash
chmod +x setup.sh && ./setup.sh
```

### 2. Launch with Smart Scan
```bash
go run main.go start
```
> Select **Smart Scan (AI Detect)** to let GoNode automatically configure your app!

### 3. Monitor
```bash
./gonode list
```

---

Developed with ❤️ by **Rayhan Dita Adam**
