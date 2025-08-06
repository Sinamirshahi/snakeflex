package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type TerminalServer struct {
	pythonFile         string
	verbose            bool
	pythonCmd          string
	workingDir         string
	fileManagerEnabled bool
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Input   string `json:"input,omitempty"`
}

type FileInfo struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Size     int64      `json:"size"`
	ModTime  time.Time  `json:"modTime"`
	Children []FileInfo `json:"children,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// WebSocket connection wrapper with concurrency protection
type SafeWebSocketConn struct {
	conn    *websocket.Conn
	mutex   sync.Mutex
	msgChan chan Message
	done    chan bool
}

func NewSafeWebSocketConn(conn *websocket.Conn) *SafeWebSocketConn {
	safe := &SafeWebSocketConn{
		conn:    conn,
		msgChan: make(chan Message, 100),
		done:    make(chan bool),
	}

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
	default:
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

// File management functions
func (ts *TerminalServer) getDirectoryTree(dirPath string) ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Skip hidden files and common directories to ignore
		if strings.HasPrefix(entry.Name(), ".") ||
			entry.Name() == "__pycache__" ||
			entry.Name() == "node_modules" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		relPath, _ := filepath.Rel(ts.workingDir, fullPath)

		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    relPath,
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		files = append(files, fileInfo)
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

func main() {
	pythonFile := flag.String("file", "fibonacci.py", "Python file to execute")
	port := flag.String("port", "8090", "Port to run server on")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	htmlFile := flag.String("template", "terminal.html", "HTML template file")
	disableFileManager := flag.Bool("disable-file-manager", false, "Disable file management features for security")
	flag.Parse()

	// Get working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

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
		pythonFile:         *pythonFile,
		verbose:            *verbose,
		pythonCmd:          pythonCmd,
		workingDir:         workingDir,
		fileManagerEnabled: !*disableFileManager,
	}

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.terminalHandler(w, r, *htmlFile)
	})
	http.HandleFunc("/ws", server.websocketHandler)

	// Conditionally setup file management API routes
	if server.fileManagerEnabled {
		http.HandleFunc("/api/files", server.filesHandler)
		http.HandleFunc("/api/files/download", server.downloadHandler)
		http.HandleFunc("/api/files/upload", server.uploadHandler)
		http.HandleFunc("/api/files/create", server.createHandler)
		http.HandleFunc("/api/files/delete", server.deleteHandler)
	}

	serverPort := ":" + *port
	fmt.Printf("🐍 Python Web Terminal started at http://localhost%s\n", serverPort)
	fmt.Printf("📁 Working Directory: %s\n", workingDir)
	fmt.Printf("🚀 Executing: %s\n", *pythonFile)
	fmt.Printf("🎨 Template: %s\n", *htmlFile)
	if *verbose {
		fmt.Println("📝 Verbose logging enabled")
	}
	fmt.Println("💬 Interactive input support enabled!")

	if server.fileManagerEnabled {
		fmt.Println("📂 File management panel enabled!")
	} else {
		fmt.Println("🔒 File management disabled for security")
	}

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
	htmlStr = strings.ReplaceAll(htmlStr, "{{WORKING_DIR}}", ts.workingDir)
	htmlStr = strings.ReplaceAll(htmlStr, "{{FILE_MANAGER_ENABLED}}", fmt.Sprintf("%t", ts.fileManagerEnabled))

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlStr)
}

// File management handlers - these will only be called if file manager is enabled
func (ts *TerminalServer) filesHandler(w http.ResponseWriter, r *http.Request) {
	// Security check: ensure file manager is enabled
	if !ts.fileManagerEnabled {
		http.Error(w, "File management disabled", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		dirPath := r.URL.Query().Get("path")
		if dirPath == "" {
			dirPath = ts.workingDir
		} else {
			dirPath = filepath.Join(ts.workingDir, dirPath)
		}

		// Security check: ensure path is within working directory
		absDir, err := filepath.Abs(dirPath)
		if err != nil || !strings.HasPrefix(absDir, ts.workingDir) {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
			return
		}

		files, err := ts.getDirectoryTree(dirPath)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(APIResponse{Success: true, Data: files})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TerminalServer) downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Security check: ensure file manager is enabled
	if !ts.fileManagerEnabled {
		http.Error(w, "File management disabled", http.StatusForbidden)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Path parameter required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(ts.workingDir, filePath)

	// Security check
	absPath, err := filepath.Abs(fullPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Check if file exists and is not a directory
	info, err := os.Stat(absPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve the file
	http.ServeFile(w, r, absPath)
}

func (ts *TerminalServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Security check: ensure file manager is enabled
	if !ts.fileManagerEnabled {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "File management disabled"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to parse form: " + err.Error()})
		return
	}

	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "."
	}

	targetDir := filepath.Join(ts.workingDir, uploadPath)

	// Security check
	absDir, err := filepath.Abs(targetDir)
	if err != nil || !strings.HasPrefix(absDir, ts.workingDir) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "No files uploaded"})
		return
	}

	uploadedFiles := []string{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Create the target file
		targetPath := filepath.Join(absDir, fileHeader.Filename)
		targetFile, err := os.Create(targetPath)
		if err != nil {
			continue
		}
		defer targetFile.Close()

		// Copy file content
		_, err = io.Copy(targetFile, file)
		if err != nil {
			os.Remove(targetPath) // Clean up on error
			continue
		}

		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}

	if len(uploadedFiles) == 0 {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to upload any files"})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: fmt.Sprintf("Uploaded %d file(s)", len(uploadedFiles)), Data: uploadedFiles})
}

func (ts *TerminalServer) createHandler(w http.ResponseWriter, r *http.Request) {
	// Security check: ensure file manager is enabled
	if !ts.fileManagerEnabled {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "File management disabled"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Path  string `json:"path"`
		IsDir bool   `json:"isDir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid request body"})
		return
	}

	targetPath := filepath.Join(ts.workingDir, req.Path)

	// Security check
	absPath, err := filepath.Abs(targetPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
		return
	}

	if req.IsDir {
		err = os.MkdirAll(absPath, 0755)
	} else {
		// Create empty file
		file, err := os.Create(absPath)
		if err == nil {
			file.Close()
		}
	}

	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Created successfully"})
}

func (ts *TerminalServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Security check: ensure file manager is enabled
	if !ts.fileManagerEnabled {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "File management disabled"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Path parameter required"})
		return
	}

	targetPath := filepath.Join(ts.workingDir, filePath)

	// Security check
	absPath, err := filepath.Abs(targetPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
		return
	}

	// Don't allow deleting the current Python file
	if absPath == filepath.Join(ts.workingDir, ts.pythonFile) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Cannot delete the currently executing Python file"})
		return
	}

	err = os.RemoveAll(absPath)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Deleted successfully"})
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
			go ts.executePythonScript(safeConn)
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
