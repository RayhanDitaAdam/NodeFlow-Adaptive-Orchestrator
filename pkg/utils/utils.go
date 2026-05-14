package utils

import (
	"fmt"
)

// PrintHelp displays guidance for the user
func PrintHelp() {
	fmt.Println("\n🌿 GoNode - Adaptive Infrastructure Engine")
	fmt.Println("Usage: gonode [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Printf("  %-15s %s\n", "start", "Meluncurkan GoNode dengan menu pilihan profil spek.")
	fmt.Printf("  %-15s %s\n", "list", "Menampilkan status aplikasi yang sedang berjalan di background.")
	fmt.Printf("  %-15s %s\n", "stop", "Menghentikan GoNode Engine dan instance Node.js.")
	fmt.Printf("  %-15s %s\n", "help", "Menampilkan pesan bantuan ini.")
	fmt.Printf("  %-15s %s\n", "help nginx", "Panduan khusus konfigurasi Nginx.")
	fmt.Println("\nExamples:")
	fmt.Println("  gonode start")
	fmt.Println("  gonode help nginx")
	fmt.Println("")
}

// PrintNginxHelp displays specific commands for managing Nginx
func PrintNginxHelp() {
	fmt.Println("\n🛡️  GoNode Nginx Helper")
	fmt.Println("Gunakan perintah berikut untuk mengelola konfigurasi Nginx ente:")
	fmt.Println("\n1. Cek Status Nginx:")
	fmt.Println("   sudo systemctl status nginx")
	
	fmt.Println("\n2. List Konfigurasi yang Aktif:")
	fmt.Println("   ls -l /etc/nginx/sites-enabled/")
	
	fmt.Println("\n3. Tes Validasi Syntax (Wajib sebelum reload):")
	fmt.Println("   sudo nginx -t")
	
	fmt.Println("\n4. Reload Nginx (Terapkan perubahan):")
	fmt.Println("   sudo systemctl reload nginx")
	
	fmt.Println("\n5. Intip Log Error (Buat debugging):")
	fmt.Println("   sudo tail -f /var/log/nginx/error.log")
	
	fmt.Println("\n6. Intip Log Akses (Traffic masuk):")
	fmt.Println("   sudo tail -f /var/log/nginx/access.log")
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
