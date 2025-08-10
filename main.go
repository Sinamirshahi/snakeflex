package main

import (
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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

type AuthConfig struct {
	Password string
	Enabled  bool
}

// Rate limiting for failed authentication attempts
type RateLimiter struct {
	attempts map[string]*AttemptRecord
	mutex    sync.RWMutex
}

type AttemptRecord struct {
	Count       int
	LastAttempt time.Time
	LockedUntil time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string]*AttemptRecord),
	}

	// Clean up old records every 10 minutes
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.CleanupOldRecords()
		}
	}()

	return rl
}

func (rl *RateLimiter) getClientIP(r *http.Request) string {
	// Check for forwarded IP first (if behind proxy)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	// Check for real IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (rl *RateLimiter) IsBlocked(r *http.Request) (bool, time.Duration) {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	clientIP := rl.getClientIP(r)
	record, exists := rl.attempts[clientIP]

	if !exists {
		return false, 0
	}

	now := time.Now()
	if now.Before(record.LockedUntil) {
		return true, record.LockedUntil.Sub(now)
	}

	return false, 0
}

func (rl *RateLimiter) RecordFailedAttempt(r *http.Request) (bool, time.Duration) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	clientIP := rl.getClientIP(r)
	now := time.Now()

	record, exists := rl.attempts[clientIP]
	if !exists {
		record = &AttemptRecord{}
		rl.attempts[clientIP] = record
	}

	// Reset counter if last attempt was more than 15 minutes ago
	if now.Sub(record.LastAttempt) > 15*time.Minute {
		record.Count = 0
	}

	record.Count++
	record.LastAttempt = now

	// Progressive lockout durations
	var lockDuration time.Duration
	switch {
	case record.Count >= 10: // 10+ attempts = 1 hour
		lockDuration = time.Hour
	case record.Count >= 6: // 6-9 attempts = 10 minutes
		lockDuration = 10 * time.Minute
	case record.Count >= 3: // 3-5 attempts = 1 minute
		lockDuration = time.Minute
	default:
		return false, 0 // No lockout yet
	}

	record.LockedUntil = now.Add(lockDuration)
	return true, lockDuration
}

func (rl *RateLimiter) RecordSuccessfulLogin(r *http.Request) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	clientIP := rl.getClientIP(r)
	delete(rl.attempts, clientIP) // Clear failed attempts on successful login
}

func (rl *RateLimiter) CleanupOldRecords() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for ip, record := range rl.attempts {
		// Remove records older than 1 hour and not currently locked
		if now.Sub(record.LastAttempt) > time.Hour && now.After(record.LockedUntil) {
			delete(rl.attempts, ip)
		}
	}
}

type SessionManager struct {
	sessions map[string]time.Time
	mutex    sync.RWMutex
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]time.Time),
	}

	// Clean up expired sessions every hour
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			sm.CleanupExpiredSessions()
		}
	}()

	return sm
}

func (sm *SessionManager) CreateSession() string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate secure random token
	bytes := make([]byte, 32)
	rand.Read(bytes)
	token := base64.URLEncoding.EncodeToString(bytes)

	// Session expires in 24 hours
	sm.sessions[token] = time.Now().Add(24 * time.Hour)
	return token
}

func (sm *SessionManager) ValidateSession(token string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	expiry, exists := sm.sessions[token]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(sm.sessions, token)
		return false
	}

	return true
}

func (sm *SessionManager) CleanupExpiredSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for token, expiry := range sm.sessions {
		if now.After(expiry) {
			delete(sm.sessions, token)
		}
	}
}

type TerminalServer struct {
	pythonFile         string
	verbose            bool
	pythonCmd          string
	workingDir         string
	fileManagerEnabled bool
	shellEnabled       bool
	authConfig         *AuthConfig
	sessionManager     *SessionManager
	rateLimiter        *RateLimiter
	basePath           string // Added for proxy support
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

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Helper function to get base path from request headers or configuration
func (ts *TerminalServer) getBasePath(r *http.Request) string {
	// Check for forwarded path prefix (set by reverse proxy)
	if prefix := r.Header.Get("X-Forwarded-Prefix"); prefix != "" {
		return strings.TrimSuffix(prefix, "/")
	}

	// Check for script name (another common proxy header)
	if script := r.Header.Get("X-Script-Name"); script != "" {
		return strings.TrimSuffix(script, "/")
	}

	// Fall back to configured base path
	return ts.basePath
}

// Helper function to construct URLs with base path
func (ts *TerminalServer) buildURL(r *http.Request, path string) string {
	basePath := ts.getBasePath(r)
	if basePath == "" {
		return path
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return basePath + path
}

// Authentication middleware
func (ts *TerminalServer) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ts.authConfig.Enabled {
			next(w, r)
			return
		}

		// Check for session cookie
		cookie, err := r.Cookie("snakeflex_session")
		if err != nil || !ts.sessionManager.ValidateSession(cookie.Value) {
			// Redirect to login page using relative path
			// Construct the path correctly based on whether we are already at the login page
			currentRequestPath := r.URL.Path
			loginPathSuffix := "/login"
			if strings.HasSuffix(currentRequestPath, loginPathSuffix) {
				next(w, r)
				return
			}

			loginURL := ts.buildURL(r, loginPathSuffix)
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		}

		next(w, r)
	}
}

// Login handler
func (ts *TerminalServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	if !ts.authConfig.Enabled {
		homeURL := ts.buildURL(r, "/")
		http.Redirect(w, r, homeURL, http.StatusFound)
		return
	}

	switch r.Method {
	case "GET":
		ts.serveLoginPage(w, r)
	case "POST":
		ts.handleLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TerminalServer) serveLoginPage(w http.ResponseWriter, r *http.Request) {
	basePath := ts.getBasePath(r)
	baseHref := ""
	if basePath != "" {
		// The base href should end with a slash for relative paths to work correctly.
		baseHref = fmt.Sprintf(`<base href="%s/">`, strings.TrimSuffix(basePath, "/"))
	}

	loginHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Snakeflex - Authentication Required</title>
    {{BASE_HREF}}
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            background: linear-gradient(135deg, #0d1117 0%, #161b22 100%); 
            color: #c9d1d9; 
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace; 
            height: 100vh; 
            display: flex; 
            align-items: center; 
            justify-content: center; 
        }
        .login-container {
            background: #21262d;
            border: 1px solid #30363d;
            border-radius: 8px;
            padding: 40px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.4);
            width: 100%;
            max-width: 400px;
            text-align: center;
        }
        .logo {
            font-size: 48px;
            margin-bottom: 20px;
        }
        .title {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 30px;
            background: linear-gradient(135deg, #58a6ff, #1f6feb);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        .form-group {
            margin-bottom: 20px;
            text-align: left;
        }
        .form-label {
            display: block;
            margin-bottom: 8px;
            font-size: 14px;
            font-weight: bold;
        }
        .form-input {
            width: 100%;
            background: #0d1117;
            border: 1px solid #30363d;
            border-radius: 6px;
            padding: 12px 16px;
            color: #c9d1d9;
            font-family: inherit;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        .form-input:focus {
            outline: none;
            border-color: #1f6feb;
            box-shadow: 0 0 0 3px rgba(31, 111, 235, 0.3);
        }
        .login-btn {
            width: 100%;
            background: linear-gradient(135deg, #238636, #2ea043);
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
            font-weight: bold;
            transition: all 0.2s;
        }
        .login-btn:hover {
            background: linear-gradient(135deg, #2ea043, #238636);
            transform: translateY(-1px);
        }
        .error-message {
            background: #da3633;
            color: white;
            padding: 12px;
            border-radius: 6px;
            margin-bottom: 20px;
            font-size: 14px;
        }
        .info-text {
            margin-top: 20px;
            font-size: 12px;
            color: #7d8590;
            line-height: 1.5;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="logo">🐍</div>
        <h1 class="title">Snakeflex Terminal</h1>
        
        {{ERROR_MESSAGE}}
        
        <form method="POST" action="login">
            <div class="form-group">
                <label class="form-label" for="password">Access Password:</label>
                <input type="password" id="password" name="password" class="form-input" 
                       placeholder="Enter your password..." required autofocus>
            </div>
            <button type="submit" class="login-btn">🔓 Access Terminal</button>
        </form>
        
        <div class="info-text">
            🔒 This terminal is password protected.<br>
            Enter the correct password to access the Python environment.
        </div>
    </div>
    
    <script>
        document.getElementById('password').focus();
        document.querySelector('form').addEventListener('submit', function(e) {
            const btn = document.querySelector('.login-btn');
            btn.textContent = '🔄 Authenticating...';
            btn.disabled = true;
        });
    </script>
</body>
</html>`

	errorMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errorMsg = `<div class="error-message">❌ Invalid password. Please try again.</div>`
	}

	loginHTML = strings.ReplaceAll(loginHTML, "{{ERROR_MESSAGE}}", errorMsg)
	loginHTML = strings.ReplaceAll(loginHTML, "{{BASE_HREF}}", baseHref)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, loginHTML)
}

func (ts *TerminalServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if client is currently blocked
	if blocked, remainingTime := ts.rateLimiter.IsBlocked(r); blocked {
		if ts.verbose {
			clientIP := ts.rateLimiter.getClientIP(r)
			log.Printf("🚫 Blocked authentication attempt from %s (IP: %s), %v remaining",
				r.RemoteAddr, clientIP, remainingTime.Round(time.Second))
		}

		ts.serveBlockedPage(w, r, remainingTime)
		return
	}

	password := r.FormValue("password")
	hashedPassword := hashPassword(password)

	if hashedPassword == ts.authConfig.Password {
		// Successful login - clear any failed attempts
		ts.rateLimiter.RecordSuccessfulLogin(r)

		// Create session
		sessionToken := ts.sessionManager.CreateSession()

		// Set secure session cookie with appropriate path
		cookiePath := ts.getBasePath(r)
		if cookiePath == "" {
			cookiePath = "/"
		} else {
			// Ensure cookie path ends with / for proper subdirectory handling
			cookiePath = strings.TrimSuffix(cookiePath, "/") + "/"
		}

		cookie := &http.Cookie{
			Name:     "snakeflex_session",
			Value:    sessionToken,
			Path:     cookiePath,
			HttpOnly: true,
			Secure:   r.TLS != nil, // Only secure if HTTPS
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, cookie)

		if ts.verbose {
			clientIP := ts.rateLimiter.getClientIP(r)
			log.Printf("✅ Successful authentication from %s (IP: %s)", r.RemoteAddr, clientIP)
		}

		homeURL := ts.buildURL(r, "/")
		http.Redirect(w, r, homeURL, http.StatusFound)
	} else {
		// Failed login - record attempt and check for lockout
		locked, lockDuration := ts.rateLimiter.RecordFailedAttempt(r)

		clientIP := ts.rateLimiter.getClientIP(r)
		if locked {
			if ts.verbose {
				log.Printf("🔒 IP %s locked for %v after failed authentication from %s",
					clientIP, lockDuration.Round(time.Second), r.RemoteAddr)
			}
			ts.serveBlockedPage(w, r, lockDuration)
		} else {
			if ts.verbose {
				log.Printf("❌ Failed authentication attempt from %s (IP: %s)", r.RemoteAddr, clientIP)
			}
			loginURL := ts.buildURL(r, "/login?error=1")
			http.Redirect(w, r, loginURL, http.StatusFound)
		}
	}
}

func (ts *TerminalServer) serveBlockedPage(w http.ResponseWriter, r *http.Request, remainingTime time.Duration) {
	basePath := ts.getBasePath(r)
	baseHref := ""
	if basePath != "" {
		baseHref = fmt.Sprintf(`<base href="%s/">`, strings.TrimSuffix(basePath, "/"))
	}

	blockedHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Snakeflex - Access Temporarily Blocked</title>
    {{BASE_HREF}}
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            background: linear-gradient(135deg, #0d1117 0%, #161b22 100%); 
            color: #c9d1d9; 
            font-family: 'Consolas', 'Monaco', 'Courier New', monospace; 
            height: 100vh; 
            display: flex; 
            align-items: center; 
            justify-content: center; 
        }
        .blocked-container {
            background: #21262d;
            border: 2px solid #da3633;
            border-radius: 8px;
            padding: 40px;
            box-shadow: 0 8px 32px rgba(218, 54, 51, 0.3);
            width: 100%;
            max-width: 500px;
            text-align: center;
            animation: shake 0.5s ease-in-out;
        }
        @keyframes shake {
            0%, 100% { transform: translateX(0); }
            25% { transform: translateX(-5px); }
            75% { transform: translateX(5px); }
        }
        .blocked-icon {
            font-size: 64px;
            margin-bottom: 20px;
            color: #da3633;
        }
        .blocked-title {
            font-size: 24px;
            font-weight: bold;
            margin-bottom: 20px;
            color: #da3633;
        }
        .blocked-message {
            font-size: 16px;
            line-height: 1.6;
            margin-bottom: 20px;
            color: #c9d1d9;
        }
        .countdown {
            background: #da3633;
            color: white;
            padding: 15px;
            border-radius: 6px;
            font-size: 18px;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .security-info {
            background: rgba(255, 215, 0, 0.1);
            border: 1px solid #ffd700;
            border-radius: 6px;
            padding: 15px;
            font-size: 14px;
            color: #ffd700;
            line-height: 1.5;
        }
        .back-btn {
            margin-top: 20px;
            background: #6e7681;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: bold;
            transition: all 0.2s;
        }
        .back-btn:hover {
            background: #7d8590;
        }
    </style>
</head>
<body>
    <div class="blocked-container">
        <div class="blocked-icon">🔒</div>
        <h1 class="blocked-title">Access Temporarily Blocked</h1>
        
        <div class="blocked-message">
            Too many failed authentication attempts have been detected from your IP address.
            Access has been temporarily restricted for security purposes.
        </div>
        
        <div class="countdown" id="countdown">
            Time remaining: <span id="timeLeft">{{REMAINING_TIME}}</span>
        </div>
        
        <div class="security-info">
            🛡️ <strong>Security Notice:</strong><br>
            This measure protects against automated attacks and unauthorized access attempts.
            Please wait for the cooldown period to expire before trying again.
        </div>
        
        <button class="back-btn" onclick="window.location.href='login'">
            ← Back to Login
        </button>
    </div>
    
    <script>
        let remainingSeconds = {{REMAINING_SECONDS}};
        
        function updateCountdown() {
            const minutes = Math.floor(remainingSeconds / 60);
            const seconds = remainingSeconds % 60;
            const timeString = minutes > 0 
                ? minutes + 'm ' + seconds.toString().padStart(2, '0') + 's'
                : seconds + 's';
            
            document.getElementById('timeLeft').textContent = timeString;
            
            if (remainingSeconds <= 0) {
                window.location.href = 'login';
                return;
            }
            
            remainingSeconds--;
        }
        
        updateCountdown();
        const interval = setInterval(updateCountdown, 1000);
        
        // Auto-redirect when time expires
        setTimeout(() => {
            clearInterval(interval);
            window.location.href = 'login';
        }, (remainingSeconds + 1) * 1000);
    </script>
</body>
</html>`

	remainingSeconds := int(remainingTime.Seconds())
	remainingTimeStr := time.Duration(remainingSeconds) * time.Second

	blockedHTML = strings.ReplaceAll(blockedHTML, "{{REMAINING_TIME}}", remainingTimeStr.String())
	blockedHTML = strings.ReplaceAll(blockedHTML, "{{REMAINING_SECONDS}}", fmt.Sprintf("%d", remainingSeconds))
	blockedHTML = strings.ReplaceAll(blockedHTML, "{{BASE_HREF}}", baseHref)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusTooManyRequests) // 429 status code
	fmt.Fprint(w, blockedHTML)
}

// Logout handler
func (ts *TerminalServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie with appropriate path
	cookiePath := ts.getBasePath(r)
	if cookiePath == "" {
		cookiePath = "/"
	} else {
		// Ensure cookie path ends with / for proper subdirectory handling
		cookiePath = strings.TrimSuffix(cookiePath, "/") + "/"
	}

	cookie := &http.Cookie{
		Name:     "snakeflex_session",
		Value:    "",
		Path:     cookiePath,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	loginURL := ts.buildURL(r, "/login")
	http.Redirect(w, r, loginURL, http.StatusFound)
}

// Enhanced getDirectoryTree with navigation support
func (ts *TerminalServer) getDirectoryTree(dirPath string) ([]FileInfo, error) {
	var files []FileInfo

	// Resolve the full directory path
	fullDirPath := filepath.Join(ts.workingDir, dirPath)

	// Security check: ensure the path is within working directory
	absFullDirPath, err := filepath.Abs(fullDirPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %v", err)
	}

	absWorkingDir, err := filepath.Abs(ts.workingDir)
	if err != nil {
		return nil, fmt.Errorf("working directory error: %v", err)
	}

	if !strings.HasPrefix(absFullDirPath, absWorkingDir) {
		return nil, fmt.Errorf("access denied: path outside working directory")
	}

	entries, err := os.ReadDir(absFullDirPath)
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

		// For file paths, we store the relative path from working directory
		relativePath := filepath.Join(dirPath, entry.Name())
		if dirPath == "" {
			relativePath = entry.Name()
		}

		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    relativePath,
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		files = append(files, fileInfo)
	}

	// Sort: directories first, then files, alphabetically within each group
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

// Helper function to validate and resolve paths
func (ts *TerminalServer) validateAndResolvePath(relativePath string) (string, error) {
	if relativePath == "" {
		return ts.workingDir, nil
	}

	fullPath := filepath.Join(ts.workingDir, relativePath)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %v", err)
	}

	absWorkingDir, err := filepath.Abs(ts.workingDir)
	if err != nil {
		return "", fmt.Errorf("working directory error: %v", err)
	}

	if !strings.HasPrefix(absPath, absWorkingDir) {
		return "", fmt.Errorf("access denied: path outside working directory")
	}

	return absPath, nil
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
	disableShell := flag.Bool("disable-shell", false, "Disable the interactive shell feature")
	password := flag.String("pass", "", "Set password for authentication (optional)")
	basePath := flag.String("base-path", "", "Base path when served behind reverse proxy (e.g., /snakeflex)")
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

	// Initialize authentication
	authConfig := &AuthConfig{
		Enabled: *password != "",
	}

	if authConfig.Enabled {
		authConfig.Password = hashPassword(*password)
		fmt.Printf("🔒 Password authentication enabled\n")
	}

	// Clean and validate base path
	cleanBasePath := strings.TrimSuffix(*basePath, "/")
	if cleanBasePath != "" && !strings.HasPrefix(cleanBasePath, "/") {
		cleanBasePath = "/" + cleanBasePath
	}

	server := &TerminalServer{
		pythonFile:         *pythonFile,
		verbose:            *verbose,
		pythonCmd:          pythonCmd,
		workingDir:         workingDir,
		fileManagerEnabled: !*disableFileManager,
		shellEnabled:       !*disableShell,
		authConfig:         authConfig,
		sessionManager:     NewSessionManager(),
		rateLimiter:        NewRateLimiter(),
		basePath:           cleanBasePath,
	}

	// Setup routes with authentication and the base path prefix
	http.HandleFunc(cleanBasePath+"/login", server.loginHandler)
	http.HandleFunc(cleanBasePath+"/logout", server.logoutHandler)
	http.HandleFunc(cleanBasePath+"/ws", server.requireAuth(server.websocketHandler))

	if server.shellEnabled {
		http.HandleFunc(cleanBasePath+"/ws-shell", server.requireAuth(server.shellWebsocketHandler))
	}

	if server.fileManagerEnabled {
		http.HandleFunc(cleanBasePath+"/api/files", server.requireAuth(server.filesHandler))
		http.HandleFunc(cleanBasePath+"/api/files/content", server.requireAuth(server.fileContentHandler))
		http.HandleFunc(cleanBasePath+"/api/files/download", server.requireAuth(server.downloadHandler))
		http.HandleFunc(cleanBasePath+"/api/files/upload", server.requireAuth(server.uploadHandler))
		http.HandleFunc(cleanBasePath+"/api/files/create", server.requireAuth(server.createHandler))
		http.HandleFunc(cleanBasePath+"/api/files/delete", server.requireAuth(server.deleteHandler))
	}

	// The root handler must be last to avoid capturing other routes.
	// The trailing slash is important for matching the base path itself and any subpaths.
	http.HandleFunc(cleanBasePath+"/", server.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		// If the path is exactly the base path without a trailing slash, redirect to add it.
		// This helps relative paths in the HTML work correctly.
		if r.URL.Path == cleanBasePath {
			http.Redirect(w, r, cleanBasePath+"/", http.StatusMovedPermanently)
			return
		}
		server.terminalHandler(w, r, *htmlFile)
	}))

	serverPort := ":" + *port
	if cleanBasePath != "" {
		fmt.Printf("🐍 Python Web Terminal started at http://localhost%s%s/\n", serverPort, cleanBasePath)
		fmt.Printf("🔗 Base path configured: %s\n", cleanBasePath)
	} else {
		fmt.Printf("🐍 Python Web Terminal started at http://localhost%s/\n", serverPort)
	}
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
		fmt.Println("📂 File management panel enabled with folder navigation!")
	} else {
		fmt.Println("🔒 File management disabled for security")
	}

	if server.shellEnabled {
		fmt.Println("⌨️ Interactive shell enabled")
	} else {
		fmt.Println("🔒 Interactive shell has been disabled via command-line flag.")
	}

	if authConfig.Enabled {
		fmt.Printf("🔐 Access the terminal at: http://localhost%s%s/login\n", serverPort, cleanBasePath)
		fmt.Printf("🛡️ Rate limiting enabled: 3+ failed attempts = 1min lockout\n")
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

	// Set ping/pong handlers for shell connection
	conn.SetPingHandler(func(appData string) error {
		if ts.verbose {
			log.Printf("Shell received ping from client: %s", r.RemoteAddr)
		}
		conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
		return nil
	})

	conn.SetPongHandler(func(appData string) error {
		if ts.verbose {
			log.Printf("Shell received pong from client: %s", r.RemoteAddr)
		}
		return nil
	})

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

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

	// Start ping routine for shell
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					if ts.verbose {
						log.Printf("Shell ping failed for %s: %v", r.RemoteAddr, err)
					}
					return
				}
			}
		}
	}()

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
			case "ping":
				response := map[string]string{"type": "pong"}
				if respData, err := json.Marshal(response); err == nil {
					conn.WriteMessage(websocket.TextMessage, respData)
				}
				continue
			case "pong":
				// Client responded to our ping
				continue
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
	htmlStr = strings.ReplaceAll(htmlStr, "{{SHELL_ENABLED}}", fmt.Sprintf("%t", ts.shellEnabled))

	// Add base path to template
	basePath := ts.getBasePath(r)
	htmlStr = strings.ReplaceAll(htmlStr, "{{BASE_PATH}}", basePath)

	if isEmbedded {
		htmlStr = strings.ReplaceAll(htmlStr, `id="embeddedNotice" style="display: none;"`, `id="embeddedNotice"`)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlStr)
}

// Enhanced filesHandler with navigation support
func (ts *TerminalServer) filesHandler(w http.ResponseWriter, r *http.Request) {
	if !ts.fileManagerEnabled {
		http.Error(w, "File management disabled", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		dirPath := r.URL.Query().Get("path")

		// Validate and get files for the requested directory
		files, err := ts.getDirectoryTree(dirPath)
		if err != nil {
			if ts.verbose {
				log.Printf("Error getting directory tree for path '%s': %v", dirPath, err)
			}
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
			return
		}

		if ts.verbose {
			log.Printf("📁 Listing %d items in directory: %s", len(files), dirPath)
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

		// Use the new validation function
		absPath, err := ts.validateAndResolvePath(filePath)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
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

		// Use the new validation function
		absPath, err := ts.validateAndResolvePath(req.Path)
		if err != nil {
			json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
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

	// Use the new validation function
	absPath, err := ts.validateAndResolvePath(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	err := r.ParseMultipartForm(500 << 20) // 500 MB max memory
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Failed to parse form: " + err.Error()})
		return
	}

	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "."
	}

	// Use the new validation function
	absDir, err := ts.validateAndResolvePath(uploadPath)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
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

	displayPath := uploadPath
	if uploadPath == "." || uploadPath == "" {
		displayPath = "root directory"
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: fmt.Sprintf("Uploaded %d file(s) to %s", len(uploadedFiles), displayPath),
		Data:    uploadedFiles,
	})
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

	// Use the new validation function
	absPath, err := ts.validateAndResolvePath(req.Path)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
		return
	}

	if req.IsDir {
		err = os.MkdirAll(absPath, 0755)
	} else {
		// Ensure parent directory exists
		parentDir := filepath.Dir(absPath)
		if err = os.MkdirAll(parentDir, 0755); err == nil {
			file, errCreate := os.Create(absPath)
			if errCreate == nil {
				file.Close()
			}
			err = errCreate
		}
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

	// Use the new validation function
	absPath, err := ts.validateAndResolvePath(filePath)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: err.Error()})
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
	// Set ping/pong handlers before upgrading
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Set ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		if ts.verbose {
			log.Printf("Received ping from client: %s", r.RemoteAddr)
		}
		conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
		return nil
	})

	conn.SetPongHandler(func(appData string) error {
		if ts.verbose {
			log.Printf("Received pong from client: %s", r.RemoteAddr)
		}
		return nil
	})

	// Set read deadline for ping/pong mechanism
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	safeConn := NewSafeWebSocketConn(conn)
	defer safeConn.Close()

	if ts.verbose {
		log.Printf("WebSocket connection established from %s", r.RemoteAddr)
	}

	var currentInputChan chan string
	var chanMutex sync.Mutex

	// Start ping routine
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					if ts.verbose {
						log.Printf("Ping failed for %s: %v", r.RemoteAddr, err)
					}
					return
				}
			}
		}
	}()

	for {
		var msg Message
		err := safeConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				if ts.verbose {
					log.Printf("WebSocket read error: %v", err)
				}
			}
			break
		}

		// Handle ping/pong messages
		switch msg.Type {
		case "ping":
			safeConn.SendMessage(Message{Type: "pong"})
			continue
		case "pong":
			// Client responded to our ping
			continue
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
				select {
				case targetChan <- msg.Input + "\n":
				default:
					log.Printf("Input channel full, dropping input for file: %s", msg.File)
				}
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

	// Use the new validation function
	absPath, err := ts.validateAndResolvePath(pythonFile)
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Invalid file path: %v", err)})
		close(inputChan)
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("File not found: %s", pythonFile)})
		close(inputChan)
		return
	}

	cmd := exec.Command(ts.pythonCmd, "-u", absPath)
	cmd.Dir = ts.workingDir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8", "PYTHONUNBUFFERED=1")

	// Use PTY on Unix-like systems for better interactive session handling
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		ts.executePtyScript(safeConn, inputChan, cmd)
	} else {
		ts.executePipeScript(safeConn, inputChan, cmd)
	}
}

func (ts *TerminalServer) executePipeScript(safeConn *SafeWebSocketConn, inputChan chan string, cmd *exec.Cmd) {
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start command: %v", err)})
		return
	}
	ts.handleIO(safeConn, stdin, stdout, stderr, cmd, inputChan)
}

func (ts *TerminalServer) executePtyScript(safeConn *SafeWebSocketConn, inputChan chan string, cmd *exec.Cmd) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		safeConn.SendMessage(Message{Type: "error", Content: fmt.Sprintf("Failed to start PTY: %v", err)})
		return
	}
	defer ptmx.Close()

	// We handle IO using the single PTY file descriptor
	ts.handleIO(safeConn, ptmx, ptmx, nil, cmd, inputChan)
}

func (ts *TerminalServer) handleIO(safeConn *SafeWebSocketConn, stdin io.WriteCloser, stdout, stderr io.ReadCloser, cmd *exec.Cmd, inputChan chan string) {
	wg := sync.WaitGroup{}

	// Goroutine to handle process exit
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inputChan)

		err := cmd.Wait()
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = -1 // Indicates an error other than a non-zero exit code
			}
		}
		safeConn.SendMessage(Message{Type: "completed", Content: fmt.Sprintf("Exit code: %d", exitCode)})
		if ts.verbose {
			log.Printf("Script execution completed with exit code: %d", exitCode)
		}
	}()

	// Goroutine for writing input to the process
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer stdin.Close()
		for input := range inputChan {
			if ts.verbose {
				log.Printf("Sending input to Python: %s", strings.TrimSpace(input))
			}
			_, err := io.WriteString(stdin, input)
			if err != nil {
				if ts.verbose {
					log.Printf("Error writing to stdin: %v", err)
				}
				return
			}
		}
	}()

	// Goroutine for reading from stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096)
		for {
			n, err := stdout.Read(buffer)
			if n > 0 {
				safeConn.SendMessage(Message{Type: "stdout", Content: string(buffer[:n])})
			}
			if err != nil {
				break // Usually io.EOF
			}
		}
	}()

	// Goroutine for reading from stderr (only if it's a separate pipe)
	if stderr != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buffer := make([]byte, 4096)
			for {
				n, err := stderr.Read(buffer)
				if n > 0 {
					safeConn.SendMessage(Message{Type: "stderr", Content: string(buffer[:n])})
				}
				if err != nil {
					break // Usually io.EOF
				}
			}
		}()
	}

	wg.Wait()
}
