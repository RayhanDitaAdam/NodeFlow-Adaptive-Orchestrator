package engine

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"gonode/pkg/logger"
)

func RunDaemonLogic() {
	projectName := os.Getenv("GO_NODE_PROJECT_NAME")
	profileName := os.Getenv("GO_NODE_NAME")
	mem := os.Getenv("GO_NODE_MEM")
	workers := os.Getenv("GO_NODE_WORKERS")
	env := os.Getenv("GO_NODE_ENV")
	entry := os.Getenv("GO_NODE_ENTRY")
	startCmdStr := os.Getenv("GO_NODE_CMD")
	domain := os.Getenv("GO_NODE_DOMAIN")
	port := os.Getenv("GO_NODE_PORT")
	startTime := time.Now()

	socketPath := GetSocketPath(projectName)
	os.Remove(socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return
	}
	defer l.Close()

	logFileName := fmt.Sprintf("%s.log", projectName)
	logFile, _ := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Fprintf(logFile, "\n[%s] [SYSTEM] GoNode Engine Starting for Project: %s...\n", time.Now().Format("2006-01-02 15:04:05"), projectName)

	// 1. Dependency Check
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] node_modules not found. Running npm install...\n", time.Now().Format("2006-01-02 15:04:05"))
		installCmd := exec.Command("npm", "install")
		installCmd.Stdout = logFile
		installCmd.Stderr = logFile
		installCmd.Run()
	}

	// 2. Build Check
	if strings.Contains(startCmdStr, "npm") {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] NPM environment detected. Ensuring build...\n", time.Now().Format("2006-01-02 15:04:05"))
		buildCmd := exec.Command("npm", "run", "build")
		buildCmd.Stdout = logFile
		buildCmd.Stderr = logFile
		buildCmd.Run()
	}

	// 3. Spawning Main Process
	var nodeCmd *exec.Cmd
	if strings.HasPrefix(startCmdStr, "npm") {
		args := strings.Fields(startCmdStr)
		nodeCmd = exec.Command(args[0], args[1:]...)
	} else {
		nodeCmd = exec.Command("node", fmt.Sprintf("--max-old-space-size=%s", mem), entry)
	}

	nodeCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	nodeCmd.Env = append(os.Environ(),
		"NODE_ENV="+env,
		"GONODE_WORKERS="+workers,
		"PORT="+port,
	)

	stdoutPipe, _ := nodeCmd.StdoutPipe()
	stderrPipe, _ := nodeCmd.StderrPipe()

	if err := nodeCmd.Start(); err != nil {
		fmt.Fprintf(logFile, "[%s] [SYSTEM] Failed to launch main process: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return
	}

	fmt.Fprintf(logFile, "[%s] [SYSTEM] App launched with command: %s\n", time.Now().Format("2006-01-02 15:04:05"), startCmdStr)

	go logger.ProcessLog(stdoutPipe, "[INFO]", logFile)
	go logger.ProcessLog(stderrPipe, "[ERROR]", logFile)
	go logger.MonitorLogSize(logFileName, 1*1024*1024)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			req := scanner.Text()
			switch req {
			case "list":
				res := fmt.Sprintf("%s|%s|Running|%s\n", projectName, profileName, time.Since(startTime).Round(time.Second))
				conn.Write([]byte(res))
			case "info":
				res := fmt.Sprintf("domain=%s\nport=%s\n", domain, port)
				conn.Write([]byte(res))
			case "stop":
				conn.Write([]byte(fmt.Sprintf("Stopping %s...\n", projectName)))
				if nodeCmd.Process != nil {
					syscall.Kill(-nodeCmd.Process.Pid, syscall.SIGKILL)
				}
				os.Remove(socketPath)
				os.Exit(0)
			}
		}
		conn.Close()
	}
}
