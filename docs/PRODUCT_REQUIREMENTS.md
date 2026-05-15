# GoNode - Product Requirements Document (PRD)

## 1. Vision & Overview
GoNode is a high-performance orchestration engine written in Go, designed to manage, monitor, and scale Node.js applications with adaptive resource profiling and automated Nginx configuration.

## 2. Target Audience
- DevOps engineers seeking a lightweight process manager.
- Node.js developers requiring adaptive resource allocation.
- Infrastructure enthusiasts building Edge Computing or Home Server solutions.

## 3. Core Features

### 3.1 Adaptive Profiling
- **Eco, Balanced, Power Profiles**: Automatically optimizes Node.js memory (`max-old-space-size`), worker counts, and system environments.

### 3.2 Intelligent Application Detection (Smart Scan)
- **Framework Awareness**: Automatically detects Next.js, React, or standard Node.js applications.
- **Auto-Config**: Configures the appropriate start command based on the project type.

### 3.3 Nginx Orchestration (Automated Reverse Proxy)
- **Zero-Manual Config**: Generates and applies Nginx configurations automatically.
- **Exposure Options**: Supports both Public Domain and Local IP setups.
- **Smart IP Detection**: Automatically fetches the server's Public IP using external providers to ensure external reachability.

### 3.4 SSL Automation (Let's Encrypt)
- **One-Click HTTPS**: Automatically installs Certbot and obtains SSL certificates when using a domain.
- **Auto-Redirect**: Configures Nginx to redirect all HTTP traffic to HTTPS.
- **Renewal Verification**: Validates that auto-renewal is working after certificate installation.

### 3.5 DNS Propagation Checker
- **Real-time Verification**: Tool to verify if a domain is correctly pointing to the server's IP before proceeding with the setup.

### 3.6 Real-time Log Management
- **Automated Rotation**: Truncates logs at 1MB to prevent disk overflow.
- **Enhanced Observability**: High-precision timestamps for all process outputs.

### 3.7 Global CLI & Control Plane
- **Background Daemon**: Fully detached background execution.
- **Unix Sockets**: High-speed local communication for process control.

## 4. Technical Stack
- **Engine**: Go (Golang) 1.23
- **Automation**: Nginx
- **Runtime**: Node.js / NPM
- **Communication**: Unix Domain Sockets

## 5. Roadmap
- [ ] Multi-application orchestration in a single daemon.
- [ ] Web-based monitoring UI.
