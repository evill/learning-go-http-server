package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Global test settings
type TestConfig struct {
	ServerPort int
	ServerHost string
	Directory  string
}

var Config = TestConfig{
	ServerPort: 4222,
	ServerHost: "localhost",
	Directory:  "/Users/user/Code/own/golang/learning-go-http-server/files",
}

// Helper function to get server URL
func (c TestConfig) GetServerURL(path string) string {
	return fmt.Sprintf("http://%s:%d%s", c.ServerHost, c.ServerPort, path)
}

func logServerOutput(cmd *exec.Cmd, rootDir string) {
	// Create logs directory if it doesn't exist
	logDir := filepath.Join(rootDir, "logs", "e2e")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// Open log file
	file, err := os.OpenFile(filepath.Join(logDir, "server.log"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("Failed to open server log file: " + err.Error())
	}

	// Set command output streams to the log file
	cmd.Stdout = file
	cmd.Stderr = file
}

var serverProcess *os.Process

func TestMain(m *testing.M) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Start the server process
	cmd := exec.Command("./your_server.sh",
		"--directory", Config.Directory,
		"--port", fmt.Sprintf("%d", Config.ServerPort))

	// Set working directory to project root
	cmd.Dir = filepath.Dir(wd) // go up one level from e2e directory
	log.Println("Working directory:", cmd.Dir)

	// Setup logging
	logServerOutput(cmd, cmd.Dir)

	err = cmd.Start()
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}

	serverProcess = cmd.Process

	log.Println("Server process ID:", serverProcess.Pid)

	// Create error channel to handle process errors
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either server to start or error
	select {
	case err := <-done:
		panic(fmt.Sprintf("Server failed to start: %v", err))
	case <-time.After(1 * time.Second):
		// Server started successfully
	}

	// Run tests
	code := m.Run()

	// Cleanup: kill the server process
	if err := serverProcess.Kill(); err != nil {
		panic("Failed to kill server: " + err.Error())
	}

	os.Exit(code)
}
