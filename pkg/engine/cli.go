package engine

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
)

func HandleStartCommand() {
	// 1. Pilih Profil Hardware
	profiles := map[string]ServerProfile{
		"Eco (512MB RAM)":    {Name: "Eco", MaxWorkers: 1, NodeEnv: "production", MemoryHeap: "256"},
		"Balanced (2GB RAM)": {Name: "Balanced", MaxWorkers: 4, NodeEnv: "production", MemoryHeap: "1024"},
		"Power (8GB+ RAM)":   {Name: "Power", MaxWorkers: 16, NodeEnv: "production", MemoryHeap: "4096"},
	}

	selectedProfile := ""
	profileOptions := []string{"Eco (512MB RAM)", "Balanced (2GB RAM)", "Power (8GB+ RAM)"}

	if err := survey.AskOne(&survey.Select{
		Message: "1. Pilih Profil Infrastruktur:",
		Options: profileOptions,
	}, &selectedProfile); err != nil {
		log.Fatal(err)
	}

	// 2. Pilih Tipe Aplikasi & Smart Detection
	appType := ""
	if err := survey.AskOne(&survey.Select{
		Message: "2. Tipe Aplikasi Ente:",
		Options: []string{"Backend (API/Node.js)", "Frontend (Next.js/React)", "Smart Scan (AI Detect)"},
	}, &appType); err != nil {
		log.Fatal(err)
	}

	entryPoint := ""
	startCmd := "node"
	
	if appType == "Smart Scan (AI Detect)" {
		fmt.Println("🤖 AI sedang menganalisa folder project ente...")
		entryPoint, startCmd = SmartDetect()
		fmt.Printf("✅ AI Menyarankan: Mode %s dengan Entry Point: %s\n", startCmd, entryPoint)
	} else if appType == "Backend (API/Node.js)" {
		entryPoint = FindBackendEntry()
	} else {
		entryPoint = "package.json"
		startCmd = "npm"
	}

	// 3. Konfirmasi Terakhir
	confirm := false
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Siap meluncurkan %s dengan profil %s?", entryPoint, selectedProfile),
		Default: true,
	}, &confirm)

	if !confirm {
		fmt.Println("❌ Dibatalkan.")
		return
	}

	launchDaemon(profiles[selectedProfile], entryPoint, startCmd)
}

func launchDaemon(config ServerProfile, entryPoint string, startCmd string) {
	cmd := exec.Command(os.Args[0], "daemon-logic")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GO_NODE_NAME=%s", config.Name),
		fmt.Sprintf("GO_NODE_MEM=%s", config.MemoryHeap),
		fmt.Sprintf("GO_NODE_WORKERS=%d", config.MaxWorkers),
		fmt.Sprintf("GO_NODE_ENV=%s", config.NodeEnv),
		fmt.Sprintf("GO_NODE_ENTRY=%s", entryPoint),
		fmt.Sprintf("GO_NODE_CMD=%s", startCmd),
	)
	
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} 

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Gagal melepas proses ke background: %v\n", err)
		return
	}

	fmt.Printf("\n🚀 GoNode [%s] meluncur ke background! Manage pake './gonode list'.\n", config.Name)
}

func SendCommand(cmd string) {
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
