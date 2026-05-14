package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

const SOCKET_FILE = "/tmp/gonode.sock"

type ServerProfile struct {
	Name       string
	MaxWorkers int
	NodeEnv    string
	MemoryHeap string
}

func main() {
	// Auto-initialize installation if binary is missing
	initInstallation()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStartCommand()
	case "list":
		sendCommand("list")
	case "stop":
		sendCommand("stop")
	case "help":
		printHelp()
	case "daemon-logic":
		runDaemonLogic()
	default:
		fmt.Printf("❌ Unknown command: '%s'\n", command)
		printHelp()
	}
}

// initInstallation mengecek apakah binary gonode sudah ada, jika tidak maka menjalankan install.sh
func initInstallation() {
	if _, err := os.Stat("./gonode"); os.IsNotExist(err) {
		fmt.Println("⚠️  GoNode binary not found. Initializing installation...")
		cmd := exec.Command("./install.sh")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("❌ Failed to run install.sh: %v\n", err)
			return
		}
		fmt.Println("✅ Initialization complete. Please use './gonode' for next commands.")
	}
}

// Fungsi Help untuk memandu user
func printHelp() {
	fmt.Println("\n🌿 GoNode - Adaptive Infrastructure Engine")
	fmt.Println("Usage: ./gonode [command]")
	fmt.Println("\nAvailable Commands:")
	fmt.Printf("  %-10s %s\n", "start", "Meluncurkan GoNode dengan menu pilihan profil spek.")
	fmt.Printf("  %-10s %s\n", "list", "Menampilkan status aplikasi yang sedang berjalan di background.")
	fmt.Printf("  %-10s %s\n", "stop", "Menghentikan GoNode Engine dan instance Node.js.")
	fmt.Printf("  %-10s %s\n", "help", "Menampilkan pesan bantuan ini.")
	fmt.Println("\nExamples:")
	fmt.Println("  ./gonode start")
	fmt.Println("  ./gonode list")
	fmt.Println("")
}

func handleStartCommand() {
	profiles := map[string]ServerProfile{
		"Eco (512MB RAM)":    {Name: "Eco", MaxWorkers: 1, NodeEnv: "production", MemoryHeap: "256"},
		"Balanced (2GB RAM)": {Name: "Balanced", MaxWorkers: 4, NodeEnv: "production", MemoryHeap: "1024"},
		"Power (8GB+ RAM)":   {Name: "Power", MaxWorkers: 16, NodeEnv: "production", MemoryHeap: "4096"},
	}

	selected := ""
	options := []string{"Eco (512MB RAM)", "Balanced (2GB RAM)", "Power (8GB+ RAM)"}

	prompt := &survey.Select{
		Message: "Pilih Profil Infrastruktur GoNode:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		log.Fatal(err)
	}

	config := profiles[selected]

	cmd := exec.Command(os.Args[0], "daemon-logic")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GO_NODE_NAME=%s", config.Name),
		fmt.Sprintf("GO_NODE_MEM=%s", config.MemoryHeap),
		fmt.Sprintf("GO_NODE_WORKERS=%d", config.MaxWorkers),
		fmt.Sprintf("GO_NODE_ENV=%s", config.NodeEnv),
	)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} 

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Gagal melepas proses ke background: %v\n", err)
		return
	}

	fmt.Printf("\n🚀 GoNode [%s] meluncur ke background! Ketik './gonode list' untuk cek.\n", config.Name)
}

func runDaemonLogic() {
	os.Remove(SOCKET_FILE)
	l, err := net.Listen("unix", SOCKET_FILE)
	if err != nil {
		return
	}
	defer l.Close()

	name := os.Getenv("GO_NODE_NAME")
	mem := os.Getenv("GO_NODE_MEM")
	workers := os.Getenv("GO_NODE_WORKERS")
	env := os.Getenv("GO_NODE_ENV")
	startTime := time.Now()

	nodeCmd := exec.Command("node", fmt.Sprintf("--max-old-space-size=%s", mem), "app.js")
	nodeCmd.Env = append(os.Environ(), 
		"NODE_ENV="+env, 
		"GONODE_WORKERS="+workers,
	)
	
	logFile, _ := os.OpenFile("gonode.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	nodeCmd.Stdout = logFile
	nodeCmd.Stderr = logFile
	nodeCmd.Start()

	// Monitor log size for rotation (20MB)
	go monitorLogSize("gonode.log", 20*1024*1024)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			req := scanner.Text()
			switch req {
			case "list":
				res := fmt.Sprintf("App: NodeApp | Profile: %s | Workers: %s | Uptime: %s\n", name, workers, time.Since(startTime).Round(time.Second))
				conn.Write([]byte(res))
			case "stop":
				conn.Write([]byte("🛑 Menghentikan GoNode dan Node.js...\n"))
				nodeCmd.Process.Kill()
				os.Remove(SOCKET_FILE)
				os.Exit(0)
			}
		}
		conn.Close()
	}
}

func sendCommand(cmd string) {
	conn, err := net.Dial("unix", SOCKET_FILE)
	if err != nil {
		fmt.Println("❌ Daemon GoNode tidak ditemukan. Jalankan './gonode start' dulu.")
		return
	}
	defer conn.Close()

	fmt.Fprintln(conn, cmd)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

// monitorLogSize mengecek ukuran file log secara berkala dan memotongnya jika melebihi batas
func monitorLogSize(filename string, maxSize int64) {
	for {
		time.Sleep(30 * time.Second)
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}

		if info.Size() >= maxSize {
			// Memotong file menjadi 0 bytes (hapus isi) tanpa menghapus filenya
			// supaya file descriptor yang dipegang proses node tetap valid.
			os.Truncate(filename, 0)
			
			// Tulis pesan rotasi
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("\n--- Log Rotated at %s (Size reached %d MB) ---\n", time.Now().Format(time.RFC3339), maxSize/(1024*1024)))
				f.Close()
			}
		}
	}
}