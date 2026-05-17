package engine

import (
	"encoding/json"
	"os"
	"strings"
)

type PackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// SmartDetect analyzes the project to find the best start command, entry point, and default port
func SmartDetect() (string, string, string) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "app.js", "node", "3000"
	}

	content := string(data)
	var pkg PackageJSON
	json.Unmarshal(data, &pkg)

	// 1. Check for Vite
	if strings.Contains(content, "vite") {
		if _, ok := pkg.Scripts["dev"]; ok {
			return "vite-project", "npm run dev", "5173"
		}
		if _, ok := pkg.Scripts["preview"]; ok {
			return "vite-project (preview)", "npm run preview", "4173"
		}
	}

	// 2. Check for Next.js
	if strings.Contains(content, "\"next\"") {
		return "Next.js App", "npm start", "3000"
	}

	// 3. Check for standard NPM start
	if _, ok := pkg.Scripts["start"]; ok {
		return "Node.js (NPM)", "npm start", "3000"
	}

	// 4. Fallback
	return "main.js", "node", "3000"
}

// FindBackendEntry looks for common entry files
func FindBackendEntry() string {
	files := []string{"index.js", "app.js", "server.js", "main.js"}
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}
	return "index.js"
}
