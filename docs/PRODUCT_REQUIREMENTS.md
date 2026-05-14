# 📄 GoNode - Product Requirements Document (PRD)

## 1. Vision & Overview
**GoNode** is a lightweight, adaptive infrastructure engine designed to orchestrate Node.js applications using Go's efficiency. It provides a bridge between low-level system management (Go) and high-level application logic (Node.js), allowing developers to run adaptive workloads with minimal overhead.

## 2. Target Audience
- DevOps engineers looking for a lightweight process manager.
- Node.js developers who need adaptive resource allocation.
- Infrastructure enthusiasts building "Home Server" or "Edge Computing" solutions.

## 3. Core Features

### 3.1 Adaptive Profiling
- Predefined profiles (**Eco**, **Balanced**, **Power**).
- Automatic adjustment of Node.js `max-old-space-size`, workers, and memory heap.

### 3.2 Smart Application Detection (Heuristic AI)
- **Smart Scan**: Analyzes `package.json` to identify frameworks (Next.js/React).
- **Backend/Frontend**: Custom orchestration paths for different app types.

### 3.3 Nginx Orchestration (Reverse Proxy Automation)
- **Automated Config**: Generates Nginx `.conf` files based on domain and port.
- **Reverse Proxy Setup**: Pre-configured headers for WebSocket support, Host forwarding, and Upgrade headers.
- **Simplification**: Eliminates manual Nginx syntax errors for beginners.

### 3.4 Daemonization & Process Management
- Background execution via Unix `setsid`.
- Decoupled process lifecycles and graceful shutdowns.

### 3.5 Real-time Log Management
- **Timestamping**: Tagging output with precise timestamps.
- **Rotation**: Automatic truncation at 1MB to save disk space.

### 3.6 CLI & Control Plane
- Communication via **Unix Sockets** (`/tmp/gonode.sock`).

## 4. Technical Stack
- **Engine**: Go (Golang)
- **Proxy/Gateway**: Nginx (Automated)
- **Runtime**: Node.js / NPM.
- **IPC**: Unix Domain Sockets.

## 5. Future Roadmap
- [ ] Automated SSL (Certbot/LetsEncrypt) integration.
- [ ] Real-time CPU/RAM monitoring in `list` command.
- [ ] Auto-restart on failure (Watchdog).
- [ ] Web-based monitoring dashboard.
