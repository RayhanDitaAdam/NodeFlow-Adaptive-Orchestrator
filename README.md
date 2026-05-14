# GoNode - Adaptive Infrastructure Engine

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org)
[![Node Version](https://img.shields.io/badge/Node-18+-green.svg)](https://nodejs.org)
[![Nginx](https://img.shields.io/badge/Nginx-Automated-brightgreen.svg)](#)

GoNode is a high-performance orchestration engine that manages Node.js applications and automates Nginx reverse proxy configurations with intelligent profiling.

---

## System Architecture

```mermaid
graph TD
    User((External Traffic)) -->|Port 80| Nginx[("Nginx Gateway")]
    
    subgraph "GoNode Control Plane"
    CLI[("CLI User")] -->|1. Profiling| Engine
    Engine[("GoNode Engine")] -->|2. Smart Scan| Detector{Frontend / Backend}
    Engine -->|3. Network Check| Net[Public IP / DNS Check]
    Engine -->|4. Automation| Nginx
    Engine -->|5. Orchestrate| NodeApp["Managed App"]
    end

    subgraph "Intelligent Infrastructure"
    Detector -->|Auto-Config| NodeApp
    NodeApp -->|Output| Logger["Timestamped Logger"]
    Logger -->|Rotation| LogFile["gonode.log (1MB Limit)"]
    end

    Nginx -->|Proxy Pass| NodeApp
```

---

## Key Features

- **Smart IP Detection**: Automatically fetches your server's **Public IP** for instant access without a domain.
- **Nginx Automation**: Automatically generate and apply Nginx configs for **Public (Domain)** or **Local (IP)** access.
- **DNS Propagation Check**: Integrated tool to verify if your domain points to your server before setup.
- **Smart AI Detection**: **Smart Scan** identifies if your app is Frontend (Next.js/React) or Backend (Node.js).
- **Adaptive Profiles**: Select hardware-optimized specs (**Eco**, **Balanced**, **Power**) with one click.
- **Daemon Mode**: Runs in the background, detached from your terminal using Unix Sockets.
- **Log Management**: High-precision timestamps and automatic log rotation at 1MB.

---

## Project Structure

```text
GoNode/
├── cmd/
│   └── gonode/        # CLI Entry Point (main.go)
├── pkg/
│   ├── engine/        # Logic: cli, daemon, detector, nginx
│   ├── logger/        # Logging & Rotation
│   └── utils/         # UI & Network Helpers
├── docs/              # Guides & Requirements
├── examples/          # Example Node.js App
├── setup.sh           # Environment Setup (Go, Node, Nginx)
└── install.sh         # Binary Builder & Global Setup
```

---

## Quick Start

### 1. Environment Setup
```bash
./setup.sh
```

### 2. Build & Global Setup
```bash
./install.sh
```
> Choose **'y'** when asked to make `gonode` global.

### 3. Launch from Anywhere
Go to your project folder and run:
```bash
gonode start
```
1. Select **RAM Profile**
2. Select **App Type** (use Smart Scan)
3. Confirm Launch
4. Select **Yes** for **Nginx Setup**
5. Choose **Exposure Type**: **Public** (Domain) or **Local** (IP)
6. GoNode automatically detects your **Public IP** if Local is selected

### 4. Monitoring
```bash
gonode list
```

---

Developed with respect by **Rayhan Dita Adam**
