package utils

import (
	"fmt"
)

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
