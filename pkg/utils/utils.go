package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

// PrintHelp displays guidance for the user
func PrintHelp() {
	fmt.Println("\nGoNode - Adaptive Infrastructure Engine")
	fmt.Println("Usage: gonode [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Printf("  %-20s %s\n", "start", "Launch GoNode with interactive profile selection.")
	fmt.Printf("  %-20s %s\n", "list", "Show the status of apps running in the background.")
	fmt.Printf("  %-20s %s\n", "logs <name>", "View real-time logs for a specific project.")
	fmt.Printf("  %-20s %s\n", "stop", "Stop the GoNode Engine.")
	fmt.Printf("  %-20s %s\n", "check propagation", "Verify if a domain resolves to the expected IP.")
	fmt.Printf("  %-20s %s\n", "help nginx", "Specific guide for Nginx configuration.")
	fmt.Println("\nExamples:")
	fmt.Println("  gonode start")
	fmt.Println("  gonode logs myapp")
	fmt.Println("")
}

// TailLogs follows the log file for a project
func TailLogs(projectName string) {
	logFileName := fmt.Sprintf("%s.log", projectName)
	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		fmt.Printf("Error: Log file %s not found. Is the project running?\n", logFileName)
		return
	}

	fmt.Printf("Showing real-time logs for: %s (Ctrl+C to exit)\n", projectName)
	cmd := exec.Command("tail", "-f", logFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// CheckPropagation verifies if a domain resolves to the expected IP
func CheckPropagation(domain string, expectedIP string) {
	fmt.Printf("AI Orchestrator: Checking DNS propagation for %s...\n", domain)
	
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("DNS Lookup failed: %v\n", err)
		fmt.Println("Please wait a few minutes for the DNS propagation to complete.")
		return
	}

	found := false
	for _, ip := range ips {
		if ip.String() == expectedIP {
			found = true
			break
		}
	}

	if found {
		fmt.Printf("Success! %s is now pointing to %s\n", domain, expectedIP)
		fmt.Println("You are ready to proceed with Nginx setup or access your website.")
	} else {
		var currentIPs []string
		for _, ip := range ips {
			currentIPs = append(currentIPs, ip.String())
		}
		fmt.Printf("Not synced! Domain %s currently points to: %s\n", domain, strings.Join(currentIPs, ", "))
		fmt.Println("Please wait; DNS updates typically take 5-30 minutes to propagate globally.")
	}
}

// PrintNginxHelp displays specific commands for managing Nginx
func PrintNginxHelp() {
	fmt.Println("\nGoNode Nginx Helper")
	fmt.Println("Use the following commands to manage your Nginx configuration:")
	fmt.Println("\n1. Check Nginx Status:")
	fmt.Println("   sudo systemctl status nginx")
	
	fmt.Println("\n2. List Active Configurations:")
	fmt.Println("   ls -l /etc/nginx/sites-enabled/")
	
	fmt.Println("\n3. Validate Syntax (Highly Recommended before reload):")
	fmt.Println("   sudo nginx -t")
	
	fmt.Println("\n4. Reload Nginx (Apply changes):")
	fmt.Println("   sudo systemctl reload nginx")
	
	fmt.Println("\n5. Monitor Error Logs (For debugging):")
	fmt.Println("   sudo tail -f /var/log/nginx/error.log")
	
	fmt.Println("\n6. Monitor Access Logs (Real-time traffic):")
	fmt.Println("   sudo tail -f /var/log/nginx/access.log")
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
