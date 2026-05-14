# GoNode - Fundamental & Architecture Guide

This guide is designed to help you understand the architecture of GoNode and the core Go concepts used in this project.

---

## 1. What is GoNode?
GoNode acts as the Supervisor / Orchestrator for your Node.js applications.
- **The Problem**: Manually setting up Nginx, managing Node.js RAM profiles, monitoring logs, and ensuring background persistence is tedious.
- **The Solution**: GoNode automates the entire lifecycle. You select a profile (Eco/Power), and GoNode handles the process execution, Nginx configuration, and log rotation.

---

## 2. Project Structure & Responsibilities
GoNode follows the Standard Go Project Layout:

- **cmd/gonode/**: 
  - **Entry Point**: Contains main.go. Its primary role is to handle CLI arguments and route them to the appropriate package logic.
- **pkg/engine/** (The Core Brain):
  - cli.go: Handles interactive menus, user inputs, and Intelligent IP Detection.
  - daemon.go: The background engine. It spawns the Node.js process and listens for commands via Unix Sockets.
  - detector.go: The "Smart Scan" heuristic logic that identifies if your app is Next.js or Node.js.
  - nginx.go: Handles the automated generation and application of Nginx configurations for both Domain and IP exposure.
- **pkg/logger/** (The Watchman):
  - Manages real-time log timestamping and handles file rotation when logs reach 1MB.
- **pkg/utils/** (The Assistant):
  - Contains helper functions like the help menu and DNS propagation checker.
- **docs/**: Central repository for documentation and architectural guides.
- **examples/**: Contains example applications.

---

## 3. Core Go Concepts Used
Understanding these concepts will help you master the GoNode codebase:

### A. Goroutines (go func())
Go is famous for its concurrency. In daemon.go, we use the go keyword to run the logger in a separate execution thread. This allows GoNode to listen to socket commands and process logs simultaneously without blocking.

### B. Unix Domain Sockets
We don't use HTTP for local communication. Instead, we use a Unix Socket file (/tmp/gonode.sock). This is faster, more secure, and specialized for inter-process communication on the same server.

### C. Structs & Maps
We use structs to define RAM profiles (Eco/Balanced/Power). These serve as strict templates for our configuration data.

### D. OS Exec
The primary job of GoNode is to execute system commands. We use the os/exec library to spawn these processes and capture their output streams.

### E. Network Logic (net/http)
In cli.go, we use net/http to communicate with external services to fetch the server's Public IP. This ensures that even when deploying locally via IP, the application is reachable from external devices.

---

## 4. Application Lifecycle
1. **Setup**: Run ./setup.sh to install system dependencies.
2. **Build**: Run ./install.sh to compile the Go source into the gonode binary.
3. **Start**: Run gonode start (global) or ./gonode start.
4. **Orchestration**: GoNode detaches to the background, manages the Node app, and rotates logs.
5. **Gateway**: GoNode configures Nginx (Domain or IP-based) to make your site accessible.

---

Happy coding!
