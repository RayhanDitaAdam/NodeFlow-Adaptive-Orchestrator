package engine

import (
	"fmt"
	"os"
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

// SetupNginx saves the config and attempts to enable it
func SetupNginx(domain string, port string) error {
	config := GenerateNginxConfig(domain, port)
	filename := fmt.Sprintf("%s.conf", domain)
	
	// Write to a temporary local file first
	err := os.WriteFile(filename, []byte(config), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("\n✅ Nginx config berhasil dibuat: %s\n", filename)
	fmt.Println("🚀 Untuk mengaktifkan, jalankan perintah ini:")
	fmt.Printf("   1. sudo mv %s /etc/nginx/sites-available/\n", filename)
	fmt.Printf("   2. sudo ln -s /etc/nginx/sites-available/%s /etc/nginx/sites-enabled/\n", filename)
	fmt.Println("   3. sudo nginx -t && sudo systemctl reload nginx")
	
	return nil
}
