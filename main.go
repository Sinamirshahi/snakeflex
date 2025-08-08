package main

import (
	"embed"
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

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

//go:embed templates/terminal.html
var embeddedTemplates embed.FS

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
	shellEnabled       bool // New field to track shell status
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Input   string `json:"input,omitempty"`
	File    string `json:"file,omitempty"`
}

type ShellMessage struct {
	Type string  `json:"type"`
	Data string  `json:"data,omitempty"`
	Cols float64 `json:"cols,omitempty"`
	Rows float64 `json:"rows,omitempty"`
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

func (ts *TerminalServer) getDirectoryTree(dirPath string) ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
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

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

func (ts *TerminalServer) getHTMLContent(htmlFile string) (string, bool) {
	if htmlContent, err := os.ReadFile(htmlFile); err == nil {
		if ts.verbose {
			log.Printf("✅ Using external template: %s", htmlFile)
		}
		return string(htmlContent), false
	}

	if embeddedContent, err := embeddedTemplates.ReadFile("templates/terminal.html"); err == nil {
		if ts.verbose {
			log.Printf("💾 External template '%s' not found, using embedded template", htmlFile)
		}
		return string(embeddedContent), true
	}

	if ts.verbose {
		log.Printf("❌ No templates available - this should not happen!")
	}
	return generateMinimalHTML(), true
}

func generateMinimalHTML() string {
	return `<!DOCTYPE html><html lang="en"><head><title>Python Terminal - Error</title></head><body>Error: Template not found.</body></html>`
}

func main() {
	// Define all command-line flags
	pythonFile := flag.String("file", "", "Python file to execute (optional)")
	port := flag.String("port", "8090", "Port to run server on")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	htmlFile := flag.String("template", "terminal.html", "HTML template file (will use embedded if not found)")
	disableFileManager := flag.Bool("disable-file-manager", false, "Disable file management features for security")
	// NEW: Flag to disable the interactive shell
	disableShell := flag.Bool("disable-shell", false, "Disable the interactive shell feature")
	flag.Parse()

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	if *pythonFile != "" {
		if _, err := os.Stat(*pythonFile); os.IsNotExist(err) {
			fmt.Printf("Error: Python file '%s' not found\n", *pythonFile)
			os.Exit(1)
		}
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
		shellEnabled:       !*disableShell, // Set server status based on the flag
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server.terminalHandler(w, r, *htmlFile)
	})
	http.HandleFunc("/ws", server.websocketHandler)

	// Conditionally register the interactive shell handler
	if server.shellEnabled {
		http.HandleFunc("/ws-shell", server.shellWebsocketHandler)
	}

	if server.fileManagerEnabled {
		http.HandleFunc("/api/files", server.filesHandler)
		http.HandleFunc("/api/files/content", server.fileContentHandler)
		http.HandleFunc("/api/files/download", server.downloadHandler)
		http.HandleFunc("/api/files/upload", server.uploadHandler)
		http.HandleFunc("/api/files/create", server.createHandler)
		http.HandleFunc("/api/files/delete", server.deleteHandler)
	}

	serverPort := ":" + *port
	fmt.Printf("🐍 Python Web Terminal started at http://localhost%s\n", serverPort)
	fmt.Printf("📁 Working Directory: %s\n", workingDir)
	if *pythonFile != "" {
		fmt.Printf("🚀 Initial script: %s\n", *pythonFile)
	} else {
		fmt.Printf("🚀 No initial script specified. Select a file from the UI to run.\n")
	}
	if *verbose {
		fmt.Println("📝 Verbose logging enabled")
	}
	if server.fileManagerEnabled {
		fmt.Println("📂 File management panel enabled!")
	} else {
		fmt.Println("🔒 File management disabled for security")
	}

	// Print status of the interactive shell
	if server.shellEnabled {
		fmt.Println("⌨️ Interactive shell enabled at /ws-shell")
	} else {
		fmt.Println("🔒 Interactive shell has been disabled via command-line flag.")
	}

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func (ts *TerminalServer) shellWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Shell WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	var shellCmd string
	if runtime.GOOS == "windows" {
		shellCmd = "powershell.exe"
	} else {
		shellCmd = "bash"
	}

	cmd := exec.Command(shellCmd)
	cmd.Dir = ts.workingDir
	cmd.Env = os.Environ()

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("Failed to start pty: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to start shell."))
		return
	}
	defer ptmx.Close()

	go func() {
		defer conn.Close()
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from pty: %v", err)
				}
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				log.Printf("Error writing to shell websocket: %v", err)
				break
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from shell websocket: %v", err)
			break
		}

		var msg ShellMessage
		if err := json.Unmarshal(message, &msg); err == nil {
			switch msg.Type {
			case "resize":
				if ts.verbose {
					log.Printf("Resizing PTY to %v rows and %v cols", msg.Rows, msg.Cols)
				}
				pty.Setsize(ptmx, &pty.Winsize{
					Rows: uint16(msg.Rows),
					Cols: uint16(msg.Cols),
				})
			case "input":
				if _, err := ptmx.Write([]byte(msg.Data)); err != nil {
					log.Printf("Error writing to pty: %v", err)
					break
				}
			}
		}
	}
	cmd.Wait()
}

func (ts *TerminalServer) terminalHandler(w http.ResponseWriter, r *http.Request, htmlFile string) {
	htmlContent, isEmbedded := ts.getHTMLContent(htmlFile)

	htmlStr := htmlContent
	htmlStr = strings.ReplaceAll(htmlStr, "{{INITIAL_PYTHON_FILE}}", ts.pythonFile)
	htmlStr = strings.ReplaceAll(htmlStr, "{{WORKING_DIR}}", ts.workingDir)
	htmlStr = strings.ReplaceAll(htmlStr, "{{FILE_MANAGER_ENABLED}}", fmt.Sprintf("%t", ts.fileManagerEnabled))
	// NEW: Pass the shell enabled status to the frontend
	htmlStr = strings.ReplaceAll(htmlStr, "{{SHELL_ENABLED}}", fmt.Sprintf("%t", ts.shellEnabled))

	if isEmbedded {
		htmlStr = strings.ReplaceAll(htmlStr, `id="embeddedNotice" style="display: none;"`, `id="embeddedNotice"`)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlStr)
}

// ... rest of the file is unchanged ...
func (ts *TerminalServer) filesHandler(w http.ResponseWriter, r *http.Request) {
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

func (ts *TerminalServer) fileContentHandler(w http.ResponseWriter, r *http.Request) {
	if !ts.fileManagerEnabled {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "File management disabled"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Path parameter required"})
			return
		}

		fullPath := filepath.Join(ts.workingDir, filePath)
		absPath, err := filepath.Abs(fullPath)
		if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
			return
		}

		info, err := os.Stat(absPath)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "File not found"})
			return
		}
		if info.IsDir() {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Cannot read directory content"})
			return
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to read file: " + err.Error()})
			return
		}

		json.NewEncoder(w).Encode(APIResponse{
			Success: true,
			Data:    map[string]string{"content": string(content)},
		})

	case "PUT":
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid request body"})
			return
		}

		if req.Path == "" {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Path is required"})
			return
		}

		fullPath := filepath.Join(ts.workingDir, req.Path)
		absPath, err := filepath.Abs(fullPath)
		if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
			return
		}

		err = os.WriteFile(absPath, []byte(req.Content), 0644)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to save file: " + err.Error()})
			return
		}

		json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "File saved successfully"})

	default:
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
	}
}

func (ts *TerminalServer) downloadHandler(w http.ResponseWriter, r *http.Request) {
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
	absPath, err := filepath.Abs(fullPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	info, err := os.Stat(absPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, absPath)
}

func (ts *TerminalServer) uploadHandler(w http.ResponseWriter, r *http.Request) {
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
	err := r.ParseMultipartForm(500 << 20)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to parse form: " + err.Error()})
		return
	}
	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "."
	}
	targetDir := filepath.Join(ts.workingDir, uploadPath)
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
		targetPath := filepath.Join(absDir, fileHeader.Filename)
		targetFile, err := os.Create(targetPath)
		if err != nil {
			continue
		}
		defer targetFile.Close()
		_, err = io.Copy(targetFile, file)
		if err != nil {
			os.Remove(targetPath)
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
	absPath, err := filepath.Abs(targetPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
		return
	}
	if req.IsDir {
		err = os.MkdirAll(absPath, 0755)
	} else {
		file, errCreate := os.Create(absPath)
		if errCreate == nil {
			file.Close()
		}
		err = errCreate
	}
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Created successfully"})
}

func (ts *TerminalServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
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
	absPath, err := filepath.Abs(targetPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid path"})
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

	var currentInputChan chan string
	var chanMutex sync.Mutex

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
			newChan := make(chan string, 10)
			chanMutex.Lock()
			currentInputChan = newChan
			chanMutex.Unlock()
			go ts.executePythonScript(safeConn, newChan, msg.File)

		case "input":
			chanMutex.Lock()
			targetChan := currentInputChan
			chanMutex.Unlock()
			if targetChan != nil {
				if ts.verbose {
					log.Printf("Received input: %s", msg.Input)
				}
				targetChan <- msg.Input + "\n"
			}
		}
	}
}

func (ts *TerminalServer) executePythonScript(safeConn *SafeWebSocketConn, inputChan chan string, pythonFile string) {
	if pythonFile == "" {
		safeConn.SendMessage(Message{Type: "error", Content: "No Python file specified for execution."})
		close(inputChan)
		return
	}

	targetPath := filepath.Join(ts.workingDir, filepath.Clean(pythonFile))
	absPath, err := filepath.Abs(targetPath)
	if err != nil || !strings.HasPrefix(absPath, ts.workingDir) {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Invalid or forbidden file path: %s", pythonFile)})
		close(inputChan)
		return
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("File not found: %s", pythonFile)})
		close(inputChan)
		return
	}

	if runtime.GOOS == "windows" {
		ts.executeWindowsScript(safeConn, inputChan, absPath)
	} else {
		ts.executeUnixScript(safeConn, inputChan, absPath)
	}
}

func (ts *TerminalServer) executeWindowsScript(safeConn *SafeWebSocketConn, inputChan chan string, scriptPath string) {
	cmd := exec.Command(ts.pythonCmd, "-u", scriptPath)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8", "PYTHONUNBUFFERED=1", "PYTHONUTF8=1")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}
	ts.handleIO(safeConn, stdin, stdout, stderr, cmd, inputChan)
}

func (ts *TerminalServer) executeUnixScript(safeConn *SafeWebSocketConn, inputChan chan string, scriptPath string) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		ts.executePTYScript(safeConn, inputChan, scriptPath)
	} else {
		ts.executeWindowsScript(safeConn, inputChan, scriptPath)
	}
}

func (ts *TerminalServer) executePTYScript(safeConn *SafeWebSocketConn, inputChan chan string, scriptPath string) {
	cmd := exec.Command(ts.pythonCmd, "-u", scriptPath)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8", "PYTHONUNBUFFERED=1", "TERM=xterm-256color")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}
	ts.handleIO(safeConn, stdin, stdout, stderr, cmd, inputChan)
}

func (ts *TerminalServer) handleIO(safeConn *SafeWebSocketConn, stdin io.WriteCloser, stdout, stderr io.ReadCloser, cmd *exec.Cmd, inputChan chan string) {
	done := make(chan bool)
	go func() {
		defer stdin.Close()
		for input := range inputChan {
			if ts.verbose {
				log.Printf("Sending input to Python: %s", strings.TrimSpace(input))
			}
			io.WriteString(stdin, input)
		}
	}()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if n > 0 {
				safeConn.SendMessage(Message{Type: "stdout", Content: string(buffer[:n])})
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if n > 0 {
				safeConn.SendMessage(Message{Type: "stderr", Content: string(buffer[:n])})
			}
			if err != nil {
				break
			}
		}
	}()

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