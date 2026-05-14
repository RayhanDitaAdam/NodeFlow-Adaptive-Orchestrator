package engine

import (
	"fmt"
	"os"
	"os/exec"
)

// GenerateNginxConfig creates a string containing the Nginx reverse proxy setup
func GenerateNginxConfig(domain string, port string) string {
	template := `server {
    listen 80;
    server_name %s;

    location / {
        proxy_pass http://localhost:%s;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}`
	return fmt.Sprintf(template, domain, port)
}

// SetupNginx saves the config and automatically applies it to the system
func SetupNginx(domain string, port string) error {
	config := GenerateNginxConfig(domain, port)
	tmpFile := fmt.Sprintf("/tmp/%s.conf", domain)
	
	// 1. Write to a temporary file
	err := os.WriteFile(tmpFile, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("gagal nulis file temp: %v", err)
	}

	fmt.Printf("\n🤖 AI Orchestrator: Memulai konfigurasi Nginx untuk %s\n", domain)

	// 2. Jalankan rangkaian perintah sudo secara otomatis
	steps := []struct {
		Name    string
		Command []string
	}{
		{"Pindahkan config ke sites-available", []string{"sudo", "mv", tmpFile, fmt.Sprintf("/etc/nginx/sites-available/%s", domain)}},
		{"Aktifkan symlink ke sites-enabled", []string{"sudo", "ln", "-sf", fmt.Sprintf("/etc/nginx/sites-available/%s", domain), fmt.Sprintf("/etc/nginx/sites-enabled/%s", domain)}},
		{"Validasi syntax Nginx", []string{"sudo", "nginx", "-t"}},
		{"Reload service Nginx", []string{"sudo", "systemctl", "reload", "nginx"}},
	}

	for i, step := range steps {
		fmt.Printf("[%d/4] %s...\n", i+1, step.Name)
		cmd := exec.Command(step.Command[0], step.Command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Gagal di tahap: %s\nDetail: %s\n", step.Name, string(output))
			return err
		}
	}

	fmt.Printf("✅ Nginx untuk %s sudah aktif dan berjalan sempurna!\n", domain)
	return nil
}
