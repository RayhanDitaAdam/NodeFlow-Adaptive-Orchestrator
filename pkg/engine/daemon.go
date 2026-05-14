package engine

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"gonode/pkg/logger"
)

func RunDaemonLogic() {
	os.Remove(SOCKET_FILE)
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return
	}
	defer l.Close()

	name := os.Getenv("GO_NODE_NAME")
	mem := os.Getenv("GO_NODE_MEM")
	workers := os.Getenv("GO_NODE_WORKERS")
	env := os.Getenv("GO_NODE_ENV")
	entry := os.Getenv("GO_NODE_ENTRY")
	startCmdStr := os.Getenv("GO_NODE_CMD")
	startTime := time.Now()

	logFile, _ := os.OpenFile("gonode.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Fprintf(logFile, "\n[%s] [SYSTEM] GoNode Engine Starting...\n", time.Now().Format("2006-01-02 15:04:05"))

	// 1. Dependency Check (Auto-Install)
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] node_modules not found. Running npm install...\n", time.Now().Format("2006-01-02 15:04:05"))
		installCmd := exec.Command("npm", "install")
		installCmd.Stdout = logFile
		installCmd.Stderr = logFile
		installCmd.Run()
	}

	// 2. Build Check (For Frontend Apps)
	if startCmdStr == "npm" {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] Frontend detected. Running npm run build...\n", time.Now().Format("2006-01-02 15:04:05"))
		buildCmd := exec.Command("npm", "run", "build")
		buildCmd.Stdout = logFile
		buildCmd.Stderr = logFile
		buildCmd.Run()
	}

	// 3. Spawning Main Process
	var nodeCmd *exec.Cmd
	if startCmdStr == "npm" {
		nodeCmd = exec.Command("npm", "start")
	} else {
		nodeCmd = exec.Command("node", fmt.Sprintf("--max-old-space-size=%s", mem), entry)
	}

	nodeCmd.Env = append(os.Environ(), 
		"NODE_ENV="+env, 
		"GONODE_WORKERS="+workers,
	)
	
	stdoutPipe, _ := nodeCmd.StdoutPipe()
	stderrPipe, _ := nodeCmd.StderrPipe()

	if err := nodeCmd.Start(); err != nil {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] Failed to launch main process: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return
	}

	fmt.Fprintf(logFile, "[%s] [SYSTEM] App launched successfully!\n", time.Now().Format("2006-01-02 15:04:05"))

	go logger.ProcessLog(stdoutPipe, "[INFO]", logFile)
	go logger.ProcessLog(stderrPipe, "[ERROR]", logFile)
	go logger.MonitorLogSize("gonode.log", 1*1024*1024)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			req := scanner.Text()
			switch req {
			case "list":
				res := fmt.Sprintf("App: %s | Profile: %s | Status: Running | Uptime: %s\n", entry, name, time.Since(startTime).Round(time.Second))
				conn.Write([]byte(res))
			case "stop":
				conn.Write([]byte("Stopping GoNode Engine...\n"))
				nodeCmd.Process.Kill()
				os.Remove(SOCKET_FILE)
				os.Exit(0)
			}
		}
		conn.Close()
	}
}
