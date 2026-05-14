package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// ProcessLog reads from a pipe and prepends a timestamp to each line
func ProcessLog(pipe io.ReadCloser, prefix string, logFile *os.File) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "[%s] %s %s\n", timestamp, prefix, scanner.Text())
	}
}

// MonitorLogSize periodically checks the log file size and truncates it if it exceeds maxSize
func MonitorLogSize(filename string, maxSize int64) {
	for {
		time.Sleep(10 * time.Second)
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}

		if info.Size() >= maxSize {
			os.Truncate(filename, 0)
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("\n--- Log Rotated at %s (Size reached %d MB) ---\n", time.Now().Format(time.RFC3339), maxSize/(1024*1024)))
				f.Close()
			}
		}
	}
}
