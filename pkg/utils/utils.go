package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// InitInstallation checks if the binary exists, if not, runs install.sh
func InitInstallation() {
	if _, err := os.Stat("./gonode"); os.IsNotExist(err) {
		fmt.Println("⚠️  GoNode binary not found. Initializing installation...")
		cmd := exec.Command("./install.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("❌ Failed to run install.sh: %v\n", err)
			return
		}
		fmt.Println("✅ Initialization complete. Please use './gonode' for next commands.")
	}
}

// PrintHelp displays guidance for the user
func PrintHelp() {
	fmt.Println("\n🌿 GoNode - Adaptive Infrastructure Engine")
	fmt.Println("Usage: ./gonode [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Printf("  %-10s %s\n", "start", "Meluncurkan GoNode dengan menu pilihan profil spek.")
	fmt.Printf("  %-10s %s\n", "list", "Menampilkan status aplikasi yang sedang berjalan di background.")
	fmt.Printf("  %-10s %s\n", "stop", "Menghentikan GoNode Engine dan instance Node.js.")
	fmt.Printf("  %-10s %s\n", "help", "Menampilkan pesan bantuan ini.")
	fmt.Println("\nExamples:")
	fmt.Println("  ./gonode start")
	fmt.Println("  ./gonode list")
	fmt.Println("")
}
