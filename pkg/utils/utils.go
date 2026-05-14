package utils

import (
	"fmt"
	"net"
	"strings"
)

// PrintHelp displays guidance for the user
func PrintHelp() {
	fmt.Println("\n🌿 GoNode - Adaptive Infrastructure Engine")
	fmt.Println("Usage: gonode [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Printf("  %-20s %s\n", "start", "Meluncurkan GoNode dengan menu pilihan profil spek.")
	fmt.Printf("  %-20s %s\n", "list", "Menampilkan status aplikasi yang sedang berjalan di background.")
	fmt.Printf("  %-20s %s\n", "stop", "Menghentikan GoNode Engine.")
	fmt.Printf("  %-20s %s\n", "check propagation", "Cek apakah domain sudah mengarah ke IP server.")
	fmt.Printf("  %-20s %s\n", "help nginx", "Panduan khusus konfigurasi Nginx.")
	fmt.Println("\nExamples:")
	fmt.Println("  gonode start")
	fmt.Println("  gonode check propagation google.com 142.251.12.102")
	fmt.Println("")
}

// CheckPropagation verifies if a domain resolves to the expected IP
func CheckPropagation(domain string, expectedIP string) {
	fmt.Printf("🔍 AI sedang mengecek propagasi DNS untuk: %s...\n", domain)
	
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("❌ Gagal mendapatkan data DNS: %v\n", err)
		fmt.Println("⏳ Silakan tunggu beberapa menit sampai proses propagasi selesai.")
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
		fmt.Printf("✅ Mantap! %s sudah mengarah ke %s\n", domain, expectedIP)
		fmt.Println("🚀 Ente sudah bisa lanjut ke setup Nginx atau akses website ente.")
	} else {
		var currentIPs []string
		for _, ip := range ips {
			currentIPs = append(currentIPs, ip.String())
		}
		fmt.Printf("⚠️  Belum sinkron! Domain %s saat ini masih mengarah ke: %s\n", domain, strings.Join(currentIPs, ", "))
		fmt.Println("⏳ Tunggu sebentar lagi bre, biasanya butuh waktu 5-30 menit buat DNS update.")
	}
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
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
