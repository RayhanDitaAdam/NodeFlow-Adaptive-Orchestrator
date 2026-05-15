package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const SOCKET_DIR = "/tmp"
const SOCKET_PREFIX = "gonode_"

// GetSocketPath returns the unique socket path for a given project name
func GetSocketPath(projectName string) string {
	return fmt.Sprintf("%s/%s%s.sock", SOCKET_DIR, SOCKET_PREFIX, projectName)
}

// GetAllSockets scans /tmp for all active GoNode sockets and returns their project names
func GetAllSockets() []string {
	entries, err := os.ReadDir(SOCKET_DIR)
	if err != nil {
		return nil
	}

	var projects []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, SOCKET_PREFIX) && strings.HasSuffix(name, ".sock") {
			projectName := strings.TrimPrefix(name, SOCKET_PREFIX)
			projectName = strings.TrimSuffix(projectName, ".sock")
			projects = append(projects, projectName)
		}
	}
	return projects
}

// GetSocketFullPath returns the full path for a socket file by project name
func GetSocketFullPath(projectName string) string {
	return filepath.Join(SOCKET_DIR, SOCKET_PREFIX+projectName+".sock")
}

type ServerProfile struct {
	Name       string
	MaxWorkers int
	NodeEnv    string
	MemoryHeap string
}
