package engine

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

func HandleStartCommand() {
	projectName := ""
	survey.AskOne(&survey.Input{Message: "0. Enter Project Name:", Default: "myapp"}, &projectName)

	profiles := map[string]ServerProfile{
		"Eco (512MB RAM)":    {Name: "Eco", MaxWorkers: 1, NodeEnv: "production", MemoryHeap: "256"},
		"Balanced (2GB RAM)": {Name: "Balanced", MaxWorkers: 4, NodeEnv: "production", MemoryHeap: "1024"},
		"Power (8GB+ RAM)":   {Name: "Power", MaxWorkers: 16, NodeEnv: "production", MemoryHeap: "4096"},
	}

	selectedProfile := ""
	profileOptions := []string{"Eco (512MB RAM)", "Balanced (2GB RAM)", "Power (8GB+ RAM)"}

	if err := survey.AskOne(&survey.Select{
		Message: "1. Select Infrastructure Profile:",
		Options: profileOptions,
	}, &selectedProfile); err != nil {
		log.Fatal(err)
	}

	appType := ""
	if err := survey.AskOne(&survey.Select{
		Message: "2. Application Type:",
		Options: []string{"Smart Scan (AI Detect) [Recommended]", "Backend (API/Node.js)", "Frontend (Next.js/React)"},
	}, &appType); err != nil {
		log.Fatal(err)
	}

	entryPoint := ""
	startCmd := "node"
	detectedPort := "3000"
	
	switch appType {
	case "Smart Scan (AI Detect) [Recommended]":
		fmt.Println("Analyzing project structure...")
		entryPoint, startCmd, detectedPort = SmartDetect()
		fmt.Printf("Suggestion: Use %s mode with Entry Point: %s (Detected Port: %s)\n", startCmd, entryPoint, detectedPort)
	case "Backend (API/Node.js)":
		entryPoint = FindBackendEntry()
		startCmd = "node"
	default:
		entryPoint = "package.json"
		startCmd = "npm start"
	}

	confirm := false
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Ready to launch %s (%s)?", projectName, entryPoint),
		Default: true,
	}, &confirm)

	if !confirm {
		return
	}

	setupNginxPrompt := false
	survey.AskOne(&survey.Confirm{
		Message: "4. Setup Nginx Reverse Proxy?",
		Default: false,
	}, &setupNginxPrompt)

	if setupNginxPrompt {
		exposureType := ""
		survey.AskOne(&survey.Select{
			Message: "Select Exposure Type:",
			Options: []string{"Public (Domain Name)", "Local (IP Address)"},
		}, &exposureType)

		domainOrIP := ""
		port := detectedPort
		
		if exposureType == "Public (Domain Name)" {
			survey.AskOne(&survey.Input{Message: "Enter Domain:"}, &domainOrIP)
		} else {
			fmt.Println("Fetching Server IP...")
			domainOrIP = getPublicIP()
			fmt.Printf("Access IP: %s\n", domainOrIP)
		}
		
		survey.AskOne(&survey.Input{Message: "Enter App Port:", Default: detectedPort}, &port)
		
		if domainOrIP != "" {
			SetupNginx(domainOrIP, port)
		}
	}

	launchDaemon(profiles[selectedProfile], entryPoint, startCmd, projectName)
}

func getPublicIP() string {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return getLocalIP()
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func launchDaemon(config ServerProfile, entryPoint string, startCmd string, projectName string) {
	cmd := exec.Command(os.Args[0], "daemon-logic")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GO_NODE_PROJECT_NAME=%s", projectName),
		fmt.Sprintf("GO_NODE_NAME=%s", config.Name),
		fmt.Sprintf("GO_NODE_MEM=%s", config.MemoryHeap),
		fmt.Sprintf("GO_NODE_WORKERS=%d", config.MaxWorkers),
		fmt.Sprintf("GO_NODE_ENV=%s", config.NodeEnv),
		fmt.Sprintf("GO_NODE_ENTRY=%s", entryPoint),
		fmt.Sprintf("GO_NODE_CMD=%s", startCmd),
	)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} 

	if err := cmd.Start(); err != nil {
		return
	}

	fmt.Printf("\nGoNode [%s] launched to background.\n", projectName)
	fmt.Printf("Logs: gonode logs %s\n", projectName)
}

func SendCommand(cmd string) {
	conn, err := net.Dial("unix", SOCKET_FILE)
	if err != nil {
		fmt.Println("Error: GoNode Daemon not found.")
		return
	}
	defer conn.Close()

	fmt.Fprintln(conn, cmd)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
