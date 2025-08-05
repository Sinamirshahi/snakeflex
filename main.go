package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type TerminalServer struct {
	pythonFile string
	verbose    bool
	pythonCmd  string
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Input   string `json:"input,omitempty"`
}

// Pseudo-terminal structures for Windows
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// detectPythonCommand detects the best Python command for the current OS
func detectPythonCommand() (string, error) {
	var candidateCommands []string

	switch runtime.GOOS {
	case "windows":
		candidateCommands = []string{"python", "python3", "py"}
	case "darwin":
		candidateCommands = []string{"python3", "python"}
	case "linux":
		candidateCommands = []string{"python3", "python"}
	default:
		candidateCommands = []string{"python3", "python"}
	}

	for _, cmd := range candidateCommands {
		if _, err := exec.LookPath(cmd); err == nil {
			if verifyPython3(cmd) {
				return cmd, nil
			}
		}
	}

	return "", fmt.Errorf("no suitable Python 3 interpreter found. Tried: %v", candidateCommands)
}

func verifyPython3(pythonCmd string) bool {
	cmd := exec.Command(pythonCmd, "--version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	version := string(output)
	return len(version) >= 8 && version[:8] == "Python 3"
}

func main() {
	pythonFile := flag.String("file", "fibonacci.py", "Python file to execute")
	port := flag.String("port", "8090", "Port to run server on")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	htmlFile := flag.String("template", "terminal.html", "HTML template file")
	flag.Parse()

	if _, err := os.Stat(*pythonFile); os.IsNotExist(err) {
		fmt.Printf("Error: Python file '%s' not found\n", *pythonFile)
		os.Exit(1)
	}

	if _, err := os.Stat(*htmlFile); os.IsNotExist(err) {
		fmt.Printf("Error: HTML template file '%s' not found\n", *htmlFile)
		os.Exit(1)
	}

	pythonCmd, err := detectPythonCommand()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Please install Python 3 and ensure it's in your PATH\n")
		fmt.Printf("Download from: https://www.python.org/downloads/\n")
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("🐍 Detected Python command: %s (OS: %s)\n", pythonCmd, runtime.GOOS)
	}

	server := &TerminalServer{
		pythonFile: *pythonFile,
		verbose:    *verbose,
		pythonCmd:  pythonCmd,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.terminalHandler(w, r, *htmlFile)
	})
	http.HandleFunc("/ws", server.websocketHandler)

	serverPort := ":" + *port
	fmt.Printf("🐍 Python Web Terminal started at http://localhost%s\n", serverPort)
	fmt.Printf("📁 Executing: %s\n", *pythonFile)
	fmt.Printf("🎨 Template: %s\n", *htmlFile)
	if *verbose {
		fmt.Println("📝 Verbose logging enabled")
	}
	fmt.Println("💬 Interactive input support enabled!")
	fmt.Println("🚀 Ready to execute ANY Python script in your browser!")

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func (ts *TerminalServer) terminalHandler(w http.ResponseWriter, r *http.Request, htmlFile string) {
	htmlContent, err := os.ReadFile(htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading HTML template: %v", err), http.StatusInternalServerError)
		return
	}

	absPath, _ := filepath.Abs(ts.pythonFile)
	commandDisplay := fmt.Sprintf("%s %s", ts.pythonCmd, ts.pythonFile)

	htmlStr := string(htmlContent)
	htmlStr = strings.ReplaceAll(htmlStr, "{{PYTHON_FILE}}", ts.pythonFile)
	htmlStr = strings.ReplaceAll(htmlStr, "{{ABS_PATH}}", absPath)
	htmlStr = strings.ReplaceAll(htmlStr, "{{COMMAND_DISPLAY}}", commandDisplay)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlStr)
}

func (ts *TerminalServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	if ts.verbose {
		log.Printf("WebSocket connection established from %s", r.RemoteAddr)
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if ts.verbose {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "execute":
			ts.executePythonScript(conn)
		case "input":
			if ts.verbose {
				log.Printf("Received input: %s", msg.Input)
			}
		}
	}
}

func (ts *TerminalServer) executePythonScript(conn *websocket.Conn) {
	if runtime.GOOS == "windows" {
		ts.executeWindowsScript(conn)
	} else {
		ts.executeUnixScript(conn)
	}
}

// Windows implementation (simplified - no real PTY support)
func (ts *TerminalServer) executeWindowsScript(conn *websocket.Conn) {
	cmd := exec.Command(ts.pythonCmd, "-u", ts.pythonFile) // -u for unbuffered

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
	cmd.Env = append(cmd.Env, "PYTHONUTF8=1")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stdin pipe: %v", err)})
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stdout pipe: %v", err)})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stderr pipe: %v", err)})
		return
	}

	if err := cmd.Start(); err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}

	ts.handleIO(conn, stdin, stdout, stderr, cmd)
}

// Unix implementation with PTY support
func (ts *TerminalServer) executeUnixScript(conn *websocket.Conn) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		ts.executePTYScript(conn)
	} else {
		ts.executeWindowsScript(conn) // Fallback
	}
}

func (ts *TerminalServer) executePTYScript(conn *websocket.Conn) {
	// Try to use python -u for unbuffered output
	cmd := exec.Command(ts.pythonCmd, "-u", ts.pythonFile)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stdin pipe: %v", err)})
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stdout pipe: %v", err)})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to create stderr pipe: %v", err)})
		return
	}

	if err := cmd.Start(); err != nil {
		ts.sendMessage(conn, Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}

	ts.handleIO(conn, stdin, stdout, stderr, cmd)
}

func (ts *TerminalServer) handleIO(conn *websocket.Conn, stdin io.WriteCloser, stdout, stderr io.ReadCloser, cmd *exec.Cmd) {
	done := make(chan bool)
	inputChan := make(chan string, 10)

	// Handle user input from WebSocket
	go func() {
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				return
			}
			if msg.Type == "input" {
				inputChan <- msg.Input + "\n"
			}
		}
	}()

	// Forward user input to Python process
	go func() {
		defer stdin.Close()
		for input := range inputChan {
			if ts.verbose {
				log.Printf("Sending input to Python: %s", strings.TrimSpace(input))
			}
			io.WriteString(stdin, input)
		}
	}()

	// Stream stdout character by character for real-time output
	go func() {
		reader := bufio.NewReader(stdout)
		var line []byte
		for {
			char, err := reader.ReadByte()
			if err != nil {
				if len(line) > 0 {
					ts.sendMessage(conn, Message{Type: "stdout", Content: string(line)})
				}
				break
			}

			if char == '\n' {
				// Send complete line
				ts.sendMessage(conn, Message{Type: "stdout", Content: string(line)})
				line = nil
			} else if char == '\r' {
				// Handle carriage return
				continue
			} else {
				line = append(line, char)
				// Send partial line for immediate feedback on prompts
				if char == ':' || char == '?' {
					if len(line) > 1 {
						ts.sendMessage(conn, Message{Type: "stdout", Content: string(line)})
						line = nil
					}
				}
			}
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			ts.sendMessage(conn, Message{Type: "stderr", Content: line})
		}
	}()

	// Wait for command to complete
	go func() {
		defer close(done)
		defer close(inputChan)

		err := cmd.Wait()
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = -1
			}
		}

		ts.sendMessage(conn, Message{Type: "completed", Content: fmt.Sprintf("Exit code: %d", exitCode)})

		if ts.verbose {
			log.Printf("Script execution completed with exit code: %d", exitCode)
		}
	}()

	<-done
}

func (ts *TerminalServer) sendMessage(conn *websocket.Conn, msg Message) {
	if err := conn.WriteJSON(msg); err != nil {
		if ts.verbose {
			log.Printf("Error sending message: %v", err)
		}
	}
}
