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
	
	err := os.WriteFile(tmpFile, []byte(config), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("\nConfiguring Nginx for %s\n", domain)

	steps := []struct {
		Name    string
		Command []string
	}{
		{"Move config", []string{"sudo", "mv", tmpFile, fmt.Sprintf("/etc/nginx/sites-available/%s", domain)}},
		{"Enable symlink", []string{"sudo", "ln", "-sf", fmt.Sprintf("/etc/nginx/sites-available/%s", domain), fmt.Sprintf("/etc/nginx/sites-enabled/%s", domain)}},
		{"Validate syntax", []string{"sudo", "nginx", "-t"}},
		{"Reload service", []string{"sudo", "systemctl", "reload", "nginx"}},
	}

	for i, step := range steps {
		fmt.Printf("[%d/4] %s...\n", i+1, step.Name)
		cmd := exec.Command(step.Command[0], step.Command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %s\n", string(output))
			return err
		}
	}

	fmt.Printf("Nginx for %s is active.\n", domain)
	return nil
}

// SetupSSL obtains and configures a Let's Encrypt SSL certificate via Certbot
func SetupSSL(domain string, email string) error {
	fmt.Printf("\nConfiguring SSL for %s\n", domain)

	// Check if certbot is installed
	if _, err := exec.LookPath("certbot"); err != nil {
		fmt.Println("Certbot not found. Installing...")
		installCmd := exec.Command("sudo", "apt", "install", "-y", "certbot", "python3-certbot-nginx")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Println("Error: Failed to install Certbot.")
			return err
		}
	}

	// Run certbot with nginx plugin
	fmt.Println("[1/2] Obtaining SSL certificate...")
	certCmd := exec.Command("sudo", "certbot", "--nginx",
		"-d", domain,
		"--email", email,
		"--agree-tos",
		"--non-interactive",
		"--redirect",
	)
	output, err := certCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: SSL setup failed.\n%s\n", string(output))
		return err
	}

	fmt.Println("[2/2] Verifying SSL renewal...")
	exec.Command("sudo", "certbot", "renew", "--dry-run").Run()

	fmt.Printf("SSL for %s is active. Your site is now served over HTTPS.\n", domain)
	return nil
}

// DisableNginx removes the symlink from sites-enabled and reloads Nginx
func DisableNginx(domain string) error {
	if domain == "" {
		return nil
	}

	fmt.Printf("\nDisabling Nginx for %s...\n", domain)

	enabledPath := fmt.Sprintf("/etc/nginx/sites-enabled/%s", domain)
	if _, err := os.Stat(enabledPath); os.IsNotExist(err) {
		return nil
	}

	steps := []struct {
		Name    string
		Command []string
	}{
		{"Remove symlink", []string{"sudo", "rm", enabledPath}},
		{"Validate syntax", []string{"sudo", "nginx", "-t"}},
		{"Reload service", []string{"sudo", "systemctl", "reload", "nginx"}},
	}

	for i, step := range steps {
		fmt.Printf("[%d/3] %s...\n", i+1, step.Name)
		cmd := exec.Command(step.Command[0], step.Command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %s\n", string(output))
			return err
		}
	}

	fmt.Printf("Nginx for %s disabled.\n", domain)
	return nil
}
