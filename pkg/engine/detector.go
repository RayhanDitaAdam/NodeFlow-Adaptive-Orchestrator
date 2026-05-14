package engine

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// SmartDetect analyzes the project folder to identify frameworks and entry points
func SmartDetect() (string, string) {
	// Heuristic: Cek package.json
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "app.js", "node" // Default
	}

	content := string(data)
	if strings.Contains(content, "\"next\"") {
		return "Next.js App", "npm"
	}
	if strings.Contains(content, "\"react\"") {
		return "React App", "npm"
	}
	
	// Cek file populer untuk backend
	for _, f := range []string{"server.js", "app.js", "index.js", "main.js"} {
		if _, err := os.Stat(f); err == nil {
			return f, "node"
		}
	}
	
	return "app.js", "node"
}

// FindBackendEntry helps users find the main file for their backend
func FindBackendEntry() string {
	files := []string{"app.js", "server.js", "index.js", "main.js"}
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}
	
	var res string
	survey.AskOne(&survey.Input{Message: "Gak nemu entry point. Masukin manual (contoh: main.js):"}, &res)
	return res
}
