package engine

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
)

func HandleStartCommand() {
	// 1. Select Profile
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

	// 2. Select App Type & Smart Detection
	appType := ""
	if err := survey.AskOne(&survey.Select{
		Message: "2. Application Type:",
		Options: []string{"Backend (API/Node.js)", "Frontend (Next.js/React)", "Smart Scan (AI Detect)"},
	}, &appType); err != nil {
		log.Fatal(err)
	}

	entryPoint := ""
	startCmd := "node"
	
	switch appType {
	case "Smart Scan (AI Detect)":
		fmt.Println("AI Orchestrator: Analyzing project structure...")
		entryPoint, startCmd = SmartDetect()
		fmt.Printf("Success! AI Suggestion: Use %s mode with Entry Point: %s\n", startCmd, entryPoint)
	case "Backend (API/Node.js)":
		entryPoint = FindBackendEntry()
	default:
		entryPoint = "package.json"
		startCmd = "npm"
	}

	// 3. Last Confirmation
	confirm := false
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Ready to launch %s with %s profile?", entryPoint, selectedProfile),
		Default: true,
	}, &confirm)

	if !confirm {
		fmt.Println("Cancelled.")
		return
	}

	// 4. Setup Nginx (Optional)
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
		port := "3000"
		
		if exposureType == "Public (Domain Name)" {
			survey.AskOne(&survey.Input{Message: "Enter Domain (e.g., myapp.com):"}, &domainOrIP)
		} else {
			domainOrIP = getLocalIP()
			fmt.Printf("AI Orchestrator: Detected Local IP: %s\n", domainOrIP)
		}
		
		survey.AskOne(&survey.Input{Message: "Enter App Port (default 3000):", Default: "3000"}, &port)
		
		if domainOrIP != "" {
			SetupNginx(domainOrIP, port)
		}
	}

	launchDaemon(profiles[selectedProfile], entryPoint, startCmd)
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

func launchDaemon(config ServerProfile, entryPoint string, startCmd string) {
	cmd := exec.Command(os.Args[0], "daemon-logic")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GO_NODE_NAME=%s", config.Name),
		fmt.Sprintf("GO_NODE_MEM=%s", config.MemoryHeap),
		fmt.Sprintf("GO_NODE_WORKERS=%d", config.MaxWorkers),
		fmt.Sprintf("GO_NODE_ENV=%s", config.NodeEnv),
		fmt.Sprintf("GO_NODE_ENTRY=%s", entryPoint),
		fmt.Sprintf("GO_NODE_CMD=%s", startCmd),
	)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} 

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: Failed to detach process to background: %v\n", err)
		return
	}

	fmt.Printf("\nGoNode [%s] launched to background! Use 'gonode list' to monitor.\n", config.Name)
}

func SendCommand(cmd string) {
	conn, err := net.Dial("unix", SOCKET_FILE)
	if err != nil {
		fmt.Println("Error: GoNode Daemon not found. Please run 'gonode start' first.")
		return
	}
	defer conn.Close()

	fmt.Fprintln(conn, cmd)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
