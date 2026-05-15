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

	// Check if project already exists
	socketPath := GetSocketPath(projectName)
	if _, err := os.Stat(socketPath); err == nil {
		fmt.Printf("⚠️  Service '%s' is already running.\n", projectName)
		stopOld := false
		survey.AskOne(&survey.Confirm{
			Message: "Do you want to stop the existing service first?",
			Default: true,
		}, &stopOld)

		if stopOld {
			HandleStopCommand(projectName)
			time.Sleep(1 * time.Second) // Wait for cleanup
		} else {
			fmt.Println("Launch cancelled. Please use a different name.")
			return
		}
	}

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

	domainOrIP := ""
	port := detectedPort

	if setupNginxPrompt {
		exposureType := ""
		survey.AskOne(&survey.Select{
			Message: "Select Exposure Type:",
			Options: []string{"Public (Domain Name)", "Local (IP Address)"},
		}, &exposureType)
		
		if exposureType == "Public (Domain Name)" {
			survey.AskOne(&survey.Input{Message: "Enter Domain:"}, &domainOrIP)
		} else {
			fmt.Println("Fetching Server IP...")
			domainOrIP = getPublicIP()
			fmt.Printf("Access IP: %s\n", domainOrIP)
		}
		
		survey.AskOne(&survey.Input{Message: "Enter App Port:", Default: detectedPort}, &port)
	}

	// ALWAYS check for port conflict if we have a port
	if port != "" {
		conflictingProject := CheckPortConflict(port)
		if conflictingProject != "" {
			fmt.Printf("⚠️  Port %s is already used by GoNode project '%s'.\n", port, conflictingProject)
			stopConflicting := false
			survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("Do you want to stop '%s' to free up the port?", conflictingProject),
				Default: true,
			}, &stopConflicting)

			if stopConflicting {
				HandleStopCommand(conflictingProject)
				time.Sleep(1 * time.Second)
			} else {
				fmt.Println("Launch cancelled. Port is busy.")
				return
			}
		}
	}
	if setupNginxPrompt && domainOrIP != "" {
		SetupNginx(domainOrIP, port)

		// Check if it's a domain (contains dot) to offer SSL
		if strings.Contains(domainOrIP, ".") {
			setupSSL := false
			survey.AskOne(&survey.Confirm{
				Message: "5. Setup SSL Certificate (Let's Encrypt)?",
				Default: true,
			}, &setupSSL)

			if setupSSL {
				email := ""
				survey.AskOne(&survey.Input{Message: "Enter email for SSL notifications:", Default: "admin@" + domainOrIP}, &email)
				SetupSSL(domainOrIP, email)
			}
		}
	}

	launchDaemon(profiles[selectedProfile], entryPoint, startCmd, projectName, domainOrIP, port)
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

func launchDaemon(config ServerProfile, entryPoint string, startCmd string, projectName string, domain string, port string) {
	cmd := exec.Command(os.Args[0], "daemon-logic")
	// Filter existing GO_NODE_ env vars to avoid duplicates
	env := os.Environ()
	var cleanEnv []string
	for _, e := range env {
		if !strings.HasPrefix(e, "GO_NODE_") {
			cleanEnv = append(cleanEnv, e)
		}
	}

	cmd.Env = append(cleanEnv,
		fmt.Sprintf("GO_NODE_PROJECT_NAME=%s", projectName),
		fmt.Sprintf("GO_NODE_NAME=%s", config.Name),
		fmt.Sprintf("GO_NODE_MEM=%s", config.MemoryHeap),
		fmt.Sprintf("GO_NODE_WORKERS=%d", config.MaxWorkers),
		fmt.Sprintf("GO_NODE_ENV=%s", config.NodeEnv),
		fmt.Sprintf("GO_NODE_ENTRY=%s", entryPoint),
		fmt.Sprintf("GO_NODE_CMD=%s", startCmd),
		fmt.Sprintf("GO_NODE_DOMAIN=%s", domain),
		fmt.Sprintf("GO_NODE_PORT=%s", port),
	)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} 

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start daemon: %v\n", err)
		return
	}

	// Health Check: Wait a bit and see if socket is created
	success := false
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		if _, err := os.Stat(GetSocketPath(projectName)); err == nil {
			success = true
			break
		}
	}

	if success {
		fmt.Printf("\nGoNode [%s] launched to background successfully.\n", projectName)
		fmt.Printf("Logs: gonode logs %s\n", projectName)
	} else {
		fmt.Printf("\n❌ Error: GoNode [%s] failed to initialize. It might have crashed during startup.\n", projectName)
		fmt.Printf("Check logs for details: cat %s.log\n", projectName)
	}
}

// ListAllServices queries all running GoNode daemons and prints their status
func ListAllServices() {
	projects := GetAllSockets()
	if len(projects) == 0 {
		fmt.Println("No running services found.")
		return
	}

	fmt.Printf("%-20s %-12s %-10s %s\n", "PROJECT", "PROFILE", "STATUS", "UPTIME")
	fmt.Println(strings.Repeat("-", 60))

	for _, project := range projects {
		socketPath := GetSocketPath(project)
		conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
		if err != nil {
			continue
		}
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		fmt.Fprintln(conn, "list")
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "|")
			if len(parts) == 4 {
				fmt.Printf("%-20s %-12s %-10s %s\n", parts[0], parts[1], parts[2], parts[3])
			} else {
				fmt.Println(scanner.Text())
			}
		}
		conn.Close()
	}
}

// CheckPortConflict scans all services to see if any is already using the target port
func CheckPortConflict(targetPort string) string {
	projects := GetAllSockets()
	for _, project := range projects {
		socketPath := GetSocketPath(project)
		conn, err := net.DialTimeout("unix", socketPath, 1*time.Second)
		if err != nil {
			continue
		}
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))

		fmt.Fprintln(conn, "info")
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "port=") {
				p := strings.TrimPrefix(line, "port=")
				if p == targetPort {
					conn.Close()
					return project
				}
			}
		}
		conn.Close()
	}
	return ""
}

// SendCommandTo sends a command to a specific project's daemon
func SendCommandTo(projectName string, cmd string) {
	socketPath := GetSocketPath(projectName)
	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
	if err != nil {
		fmt.Printf("Error: Service '%s' not found or not running.\n", projectName)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	fmt.Fprintln(conn, cmd)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// HandleStopCommand handles the full stop sequence including Nginx cleanup
func HandleStopCommand(projectName string) {
	socketPath := GetSocketPath(projectName)
	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
	if err != nil {
		fmt.Printf("Error: Service '%s' not found or not running.\n", projectName)
		return
	}
	
	// 1. Get info to check for domain
	fmt.Fprintln(conn, "info")
	scanner := bufio.NewScanner(conn)
	domain := ""
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "domain=") {
			domain = strings.TrimPrefix(line, "domain=")
		}
	}
	conn.Close()

	// 2. Disable Nginx if domain exists
	if domain != "" {
		DisableNginx(domain)
	}

	// 3. Send stop command
	SendCommandTo(projectName, "stop")
}
