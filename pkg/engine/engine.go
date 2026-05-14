package engine

const SOCKET_FILE = "/tmp/gonode.sock"

type ServerProfile struct {
	Name       string
	MaxWorkers int
	NodeEnv    string
	MemoryHeap string
}

// Common configuration shared between CLI and Daemon
type AppConfig struct {
	Profile    ServerProfile
	EntryPoint string
	StartCmd   string
}
