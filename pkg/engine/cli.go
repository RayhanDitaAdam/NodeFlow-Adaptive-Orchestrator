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
	"strconv"
	"strings"
	"syscall"
	"time"

	"gonode/pkg/utils"

	"github.com/AlecAivazis/survey/v2"
)

func PrintLogo() {
	logo := `
   ______      _   __           __   
  / ____/___  / | / /___  ____/ /__ 
 / / __/ __ \/  |/ / __ \/ __  / _ \
/ /_/ / /_/ / /|  / /_/ / /_/ /  __/
\____/\____/_/ |_/\____/\__,_/\___/ 
                                     
   Adaptive Infrastructure Engine
`
	fmt.Println(logo)
}

func HandleStartCommand() {
	PrintLogo()
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

			if !VerifyDomainIP(domainOrIP) {
				fmt.Println("Launch aborted due to DNS mismatch.")
				return
			}
		} else {
			fmt.Println("Fetching Server IP...")
			domainOrIP = getPublicIP()
			fmt.Printf("Access IP: %s\n", domainOrIP)
		}

		survey.AskOne(&survey.Input{Message: "Enter App Port:", Default: detectedPort}, &port)
	}

	// ALWAYS check for port conflict if we have a port
	if port != "" {
		for {
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
					continue
				} else {
					startVal, err := strconv.Atoi(port)
					if err != nil {
						startVal = 3000
					}
					suggestedPort := FindFreePort(startVal)
					fmt.Printf("Suggestion: Use a free port like %s\n", suggestedPort)
					survey.AskOne(&survey.Input{Message: "Enter App Port:", Default: suggestedPort}, &port)
					continue
				}
			}

			// Check if the port is physically bound by any other process on the system
			if IsPortInUse(port) {
				fmt.Printf("⚠️  Port %s is already in use by another process on this server.\n", port)
				
				killIt := false
				survey.AskOne(&survey.Confirm{
					Message: fmt.Sprintf("Do you want to terminate the process occupying port %s?", port),
					Default: false,
				}, &killIt)

				if killIt {
					fmt.Printf("Terminating process using port %s...\n", port)
					if err := KillProcessOnPort(port); err != nil {
						fmt.Printf("❌ Failed to terminate process: %v\n", err)
					} else {
						time.Sleep(1 * time.Second) // Wait for port to clear
						if !IsPortInUse(port) {
							fmt.Printf("✅ Port %s successfully freed.\n", port)
							continue
						} else {
							fmt.Printf("⚠️  Port %s is still in use.\n", port)
						}
					}
				}

				startVal, err := strconv.Atoi(port)
				if err != nil {
					startVal = 3000
				}
				suggestedPort := FindFreePort(startVal)
				fmt.Printf("Suggestion: Use a free port like %s\n", suggestedPort)
				survey.AskOne(&survey.Input{Message: "Enter App Port:", Default: suggestedPort}, &port)
				continue
			}

			break
		}
	}
	if setupNginxPrompt && domainOrIP != "" {
		SetupNginx(domainOrIP, port)

		// SSL only for Domains (not IP addresses)
		if net.ParseIP(domainOrIP) == nil && strings.Contains(domainOrIP, ".") {
			setupSSL := false
			survey.AskOne(&survey.Confirm{
				Message: "5. Setup SSL Certificate (Let's Encrypt)?",
				Default: true,
			}, &setupSSL)

			if setupSSL {
				email := ""
				survey.AskOne(&survey.Input{Message: "Enter email for SSL notifications:", Default: "admin@" + domainOrIP}, &email)
				if err := SetupSSL(domainOrIP, email); err != nil {
					fmt.Println("\n❌ SSL Setup Failed! Your app will only be accessible over HTTP.")
					fmt.Println("💡 Troubleshooting Tips:")
					fmt.Println("  1. Run: 'sudo ufw allow 80/tcp && sudo ufw allow 443/tcp'")
					fmt.Println("  2. If using AWS/GCP: Open Port 80 & 443 in your Security Groups (Inbound Rules).")
					fmt.Println("  3. Ensure your domain is correctly pointing to this server's public IP.")

					continueAnyway := false
					survey.AskOne(&survey.Confirm{
						Message: "Do you want to continue launching over HTTP?",
						Default: true,
					}, &continueAnyway)

					if !continueAnyway {
						fmt.Println("Launch aborted.")
						return
					}
				}
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
		fmt.Printf("\nGoNode [%s] launched to background.\n", projectName)
		if port != "" {
			stop := make(chan bool)
			utils.ShowLoading(stop, fmt.Sprintf("Waiting for application to bind to port %s (running build/install)...", port))
			
			portActive := false
			// Wait up to 5 minutes (300 * 1s) for heavy builds like Next.js
			for i := 0; i < 300; i++ {
				time.Sleep(1 * time.Second)
				if IsPortInUse(port) {
					portActive = true
					break
				}
				if _, err := os.Stat(GetSocketPath(projectName)); err != nil {
					break
				}
			}
			stop <- true
			
			if portActive {
				fmt.Printf("\r🚀 GoNode [%s] is now LIVE and fully accessible!\n", projectName)
				if domain != "" {
					fmt.Printf("   Access URL: http://%s\n", domain)
				} else {
					fmt.Printf("   Access Port: %s\n", port)
				}
			} else {
				fmt.Printf("\r⚠️  GoNode [%s] is taking longer than expected to start.\n", projectName)
				fmt.Printf("   It is likely still running 'npm install' or 'npm run build' in the background.\n")
				fmt.Printf("   You can monitor build progress at any time: gonode logs %s\n", projectName)
			}
		} else {
			fmt.Printf("Logs: gonode logs %s\n", projectName)
		}

		if strings.Contains(startCmd, "dev") || strings.Contains(startCmd, "preview") {
			fmt.Println("\n💡 Tip: If you see '$RefreshSig$ is not defined' or browser cache issues:")
			fmt.Println("  1. Perform a Hard Reload (Ctrl + F5 or Cmd + Shift + R).")
			fmt.Println("  2. Clear your browser cache for this IP/domain.")
		}
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

	fmt.Printf("%-20s %-12s %-10s %-10s %s\n", "PROJECT", "PROFILE", "STATUS", "UPTIME", "ACCESS")
	fmt.Println(strings.Repeat("-", 80))

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
			if len(parts) >= 4 {
				access := "-"
				if len(parts) == 5 && parts[4] != "" {
					access = parts[4]
				}
				fmt.Printf("%-20s %-12s %-10s %-10s %s\n", parts[0], parts[1], parts[2], parts[3], access)
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

// VerifyDomainIP checks if the domain resolves to the current server's public IP
func VerifyDomainIP(domain string) bool {
	publicIP := getPublicIP()

	fmt.Printf("Verifying DNS for %s...\n", domain)

	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("⚠️  Could not resolve domain %s. Ensure it is configured correctly.\n", domain)
		return false
	}

	match := false
	var resolvedIP string
	for _, ip := range ips {
		if ip.To4() != nil {
			resolvedIP = ip.String()
			if resolvedIP == publicIP {
				match = true
				break
			}
		}
	}

	if !match {
		fmt.Printf("\n⚠️  DNS MISMATCH DETECTED!\n")
		fmt.Printf("  Domain: %s -> %s\n", domain, resolvedIP)
		fmt.Printf("  Server Public IP: %s\n", publicIP)
		fmt.Printf("\n💡 Please update your DNS records to point to %s\n", publicIP)

		continueAnyway := false
		survey.AskOne(&survey.Confirm{
			Message: "Domain does not point to this server. Continue anyway?",
			Default: false,
		}, &continueAnyway)
		return continueAnyway
	}

	fmt.Printf("✅ DNS verified. %s points to this server.\n", domain)
	return true
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

// IsPortInUse checks if a TCP port is in use on the host system
func IsPortInUse(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

// FindFreePort searches for a free TCP port starting from the given port
func FindFreePort(startPort int) string {
	for port := startPort; port < startPort+100; port++ {
		pStr := strconv.Itoa(port)
		if !IsPortInUse(pStr) {
			return pStr
		}
	}
	return ""
}

// KillProcessOnPort terminates any process currently listening on the specified TCP port
func KillProcessOnPort(port string) error {
	// 1. Try fuser -k
	cmd := exec.Command("sudo", "fuser", "-k", port+"/tcp")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// 2. Try lsof and kill
	lsofCmd := exec.Command("sudo", "lsof", "-t", "-i:"+port)
	pidBytes, err := lsofCmd.Output()
	if err == nil {
		pids := strings.TrimSpace(string(pidBytes))
		if pids != "" {
			for _, pidStr := range strings.Split(pids, "\n") {
				pidStr = strings.TrimSpace(pidStr)
				if pidStr != "" {
					exec.Command("sudo", "kill", "-9", pidStr).Run()
				}
			}
			return nil
		}
	}
	return fmt.Errorf("could not kill process on port %s", port)
}
