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
	"sync"

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

// WebSocket connection wrapper with concurrency protection
type SafeWebSocketConn struct {
	conn   *websocket.Conn
	mutex  sync.Mutex
	msgChan chan Message
	done   chan bool
}

func NewSafeWebSocketConn(conn *websocket.Conn) *SafeWebSocketConn {
	safe := &SafeWebSocketConn{
		conn:    conn,
		msgChan: make(chan Message, 100), // Buffer to prevent blocking
		done:    make(chan bool),
	}
	
	// Start message sender goroutine
	go safe.messageSender()
	
	return safe
}

func (s *SafeWebSocketConn) messageSender() {
	for {
		select {
		case msg := <-s.msgChan:
			s.mutex.Lock()
			err := s.conn.WriteJSON(msg)
			s.mutex.Unlock()
			if err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		case <-s.done:
			return
		}
	}
}

func (s *SafeWebSocketConn) SendMessage(msg Message) {
	select {
	case s.msgChan <- msg:
		// Message queued successfully
	default:
		// Channel full, drop message to prevent blocking
		log.Printf("Message queue full, dropping message: %s", msg.Type)
	}
}

func (s *SafeWebSocketConn) ReadJSON(v interface{}) error {
	return s.conn.ReadJSON(v)
}

func (s *SafeWebSocketConn) Close() {
	close(s.done)
	s.conn.Close()
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
	
	safeConn := NewSafeWebSocketConn(conn)
	defer safeConn.Close()

	if ts.verbose {
		log.Printf("WebSocket connection established from %s", r.RemoteAddr)
	}

	for {
		var msg Message
		err := safeConn.ReadJSON(&msg)
		if err != nil {
			if ts.verbose {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		switch msg.Type {
		case "execute":
			go ts.executePythonScript(safeConn) // Run in separate goroutine
		case "input":
			if ts.verbose {
				log.Printf("Received input: %s", msg.Input)
			}
		}
	}
}

func (ts *TerminalServer) executePythonScript(safeConn *SafeWebSocketConn) {
	if runtime.GOOS == "windows" {
		ts.executeWindowsScript(safeConn)
	} else {
		ts.executeUnixScript(safeConn)
	}
}

func (ts *TerminalServer) executeWindowsScript(safeConn *SafeWebSocketConn) {
	cmd := exec.Command(ts.pythonCmd, "-u", ts.pythonFile)
	
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
	cmd.Env = append(cmd.Env, "PYTHONUTF8=1")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stdin pipe: %v", err)})
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stdout pipe: %v", err)})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stderr pipe: %v", err)})
		return
	}

	if err := cmd.Start(); err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}

	ts.handleIO(safeConn, stdin, stdout, stderr, cmd)
}

func (ts *TerminalServer) executeUnixScript(safeConn *SafeWebSocketConn) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		ts.executePTYScript(safeConn)
	} else {
		ts.executeWindowsScript(safeConn)
	}
}

func (ts *TerminalServer) executePTYScript(safeConn *SafeWebSocketConn) {
	cmd := exec.Command(ts.pythonCmd, "-u", ts.pythonFile)
	
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stdin pipe: %v", err)})
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stdout pipe: %v", err)})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to create stderr pipe: %v", err)})
		return
	}

	if err := cmd.Start(); err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}

	ts.handleIO(safeConn, stdin, stdout, stderr, cmd)
}

func (ts *TerminalServer) handleIO(safeConn *SafeWebSocketConn, stdin io.WriteCloser, stdout, stderr io.ReadCloser, cmd *exec.Cmd) {
	done := make(chan bool)
	inputChan := make(chan string, 10)

	// Handle user input from WebSocket
	go func() {
		for {
			var msg Message
			err := safeConn.ReadJSON(&msg)
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

	// Stream stdout line by line
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			safeConn.SendMessage(Message{Type: "stdout", Content: line})
		}
	}()

	// Stream stderr line by line
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			safeConn.SendMessage(Message{Type: "stderr", Content: line})
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

		safeConn.SendMessage(Message{Type: "completed", Content: fmt.Sprintf("Exit code: %d", exitCode)})

		if ts.verbose {
			log.Printf("Script execution completed with exit code: %d", exitCode)
		}
	}()

	<-done
}