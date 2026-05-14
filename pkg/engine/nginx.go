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

	fmt.Printf("\n⚙️  AI sedang mengonfigurasi Nginx untuk %s...\n", domain)

	// 2. Jalankan rangkaian perintah sudo secara otomatis
	commands := [][]string{
		{"sudo", "mv", tmpFile, fmt.Sprintf("/etc/nginx/sites-available/%s", domain)},
		{"sudo", "ln", "-sf", fmt.Sprintf("/etc/nginx/sites-available/%s", domain), fmt.Sprintf("/etc/nginx/sites-enabled/%s", domain)},
		{"sudo", "nginx", "-t"},
		{"sudo", "systemctl", "reload", "nginx"},
	}

	for i, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error di tahap %d: %s\n", i+1, string(output))
			return err
		}
	}

	fmt.Printf("✅ Nginx untuk %s berhasil diaktifkan secara otomatis!\n", domain)
	return nil
}
