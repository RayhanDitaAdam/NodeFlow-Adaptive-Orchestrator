package engine

import (
	"encoding/json"
	"os"
	"strings"
)

type PackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// SmartDetect analyzes the project to find the best start command and entry point
func SmartDetect() (string, string) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "app.js", "node"
	}

	content := string(data)
	var pkg PackageJSON
	json.Unmarshal(data, &pkg)

	// 1. Check for Vite (Very common for modern React apps)
	if strings.Contains(content, "vite") {
		if _, ok := pkg.Scripts["preview"]; ok {
			return "vite-project", "npm run preview"
		}
		if _, ok := pkg.Scripts["dev"]; ok {
			return "vite-project (dev)", "npm run dev"
		}
	}

	// 2. Check for Next.js
	if strings.Contains(content, "\"next\"") {
		return "Next.js App", "npm start"
	}

	// 3. Check for standard NPM start
	if _, ok := pkg.Scripts["start"]; ok {
		return "Node.js (NPM)", "npm start"
	}

	// 4. Fallback to direct Node
	return "main.js", "node"
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
