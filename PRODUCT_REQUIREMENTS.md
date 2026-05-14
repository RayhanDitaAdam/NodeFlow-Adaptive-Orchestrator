# 📄 GoNode - Product Requirements Document (PRD)

## 1. Vision & Overview
**GoNode** is a lightweight, adaptive infrastructure engine designed to orchestrate Node.js applications using Go's efficiency. It provides a bridge between low-level system management (Go) and high-level application logic (Node.js), allowing developers to run adaptive workloads with minimal overhead.

## 2. Target Audience
- DevOps engineers looking for a lightweight process manager.
- Node.js developers who need adaptive resource allocation.
- Infrastructure enthusiasts building "Home Server" or "Edge Computing" solutions.

## 3. Core Features

### 3.1 Adaptive Profiling
- Users can select from predefined profiles (**Eco**, **Balanced**, **Power**).
- Automatically adjusts Node.js `max-old-space-size`, worker counts, and environment variables based on the selected profile.

### 3.2 Smart Application Detection (Heuristic AI)
- **Smart Scan**: Analyzes `package.json` and project structure to identify frameworks.
- **Frontend Support**: Built-in orchestration for Next.js and React applications using `npm start`.
- **Backend Support**: Automatically discovers entry points like `app.js`, `server.js`, or `index.js`.

### 3.3 Daemonization & Process Management
- Runs as a background process using Unix `setsid`.
- Decouples the lifecycle of the engine from the terminal session.
- Graceful stop mechanism for both the engine and the managed Node.js instances.

### 3.4 Real-time Log Management
- **Timestamping**: Intercepts and tags Node.js output with high-precision timestamps.
- **Rotation**: Automatically truncates log files once they reach 1MB.

### 3.5 CLI & Control Plane
- Communication via **Unix Sockets** (`/tmp/gonode.sock`).
- Commands: `start`, `stop`, `list`, `help`.

### 3.6 Automated Environment Setup
- Interactive setup script for Ubuntu/Debian with automatic dependency installation.

## 4. Technical Stack
- **Engine**: Go (Golang)
- **Refactored Architecture**: Modular design (pkg/engine, pkg/logger, pkg/utils).
- **Runtime**: Node.js / NPM.
- **IPC**: Unix Domain Sockets.

## 5. Future Roadmap
- [ ] Multi-instance support (Multiple Node apps).
- [ ] Real-time CPU/RAM monitoring in `list` command.
- [ ] Auto-restart on failure (Watchdog).
- [ ] Web-based monitoring dashboard.
