package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/atiladefreitas/prevy/clipboard"
	"github.com/atiladefreitas/prevy/store"
)

const pollInterval = 1 * time.Second

func pidPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, ".local", "share", "prevy")
	if err := os.MkdirAll(p, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(p, "daemon.pid"), nil
}

func writePID() error {
	p, err := pidPath()
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(strconv.Itoa(os.Getpid())), 0o644)
}

func removePID() {
	p, err := pidPath()
	if err == nil {
		os.Remove(p)
	}
}

// IsRunning checks if a daemon process is already running.
func IsRunning() bool {
	p, err := pidPath()
	if err != nil {
		return false
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// signal 0 checks if process exists without actually signaling it
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

// Run starts the clipboard polling loop. It blocks until interrupted.
func Run() error {
	if IsRunning() {
		return fmt.Errorf("daemon is already running")
	}

	if err := writePID(); err != nil {
		return fmt.Errorf("failed to write pid file: %w", err)
	}
	defer removePID()

	fmt.Println("prevy daemon started (polling every 1s)")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	lastContent := ""

	// read initial clipboard so we don't immediately re-add what's already there
	if current, err := clipboard.Read(); err == nil {
		lastContent = current
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			fmt.Println("\nprevy daemon stopped")
			return nil
		case <-ticker.C:
			content, err := clipboard.Read()
			if err != nil || content == "" {
				continue
			}
			if content == lastContent {
				continue
			}
			lastContent = content

			entries, err := store.Load()
			if err != nil {
				continue
			}
			entries = store.Add(entries, content)
			_ = store.Save(entries)
		}
	}
}
