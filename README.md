# 🐍 SnakeFlex V1.6

*A secure, modern web-based Python development environment with interactive shell access and folder navigation.*

Run any Python script in your browser with real-time output, interactive input support, comprehensive file management with folder navigation, built-in code editing, **full interactive shell access**, and **password authentication**. No modifications to your code required.

## ✨ What it does

SnakeFlex V1.6 creates a beautiful web-based development environment for Python scripts with complete terminal capabilities and IDE-like folder navigation. Think of it as your IDE, terminal, and shell combined, but accessible from anywhere with a web browser.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
* 🔐 **Password authentication** - Secure access with session management
* ⌨️ **Interactive shell access** - Full PowerShell (Windows) or Bash (Linux/macOS) in your browser
* 📁 **Folder navigation** - Navigate into subdirectories with breadcrumb navigation and up/home buttons
* 📝 **Built-in code editor** - Edit Python files directly in the browser with syntax awareness
* 📂 **File manager** - Browse, upload, download, and manage files with drag & drop across directories
* 🎯 **Dynamic script selection** - Switch between Python scripts with right-click menu
* 🚀 **No file required** - Start without specifying a script, choose dynamically in the UI
* 💬 **Interactive input** - Handle `input()` calls seamlessly
* ⚡ **Real-time output** - See your script's output as it happens
* 🎨 **Modern UI** - GitHub-inspired dark interface with resizable panels
* 🔒 **Security modes** - Full-featured or secure terminal-only mode
* 🔄 **Cross-platform** - Windows, macOS, and Linux support
* 💾 **Embedded templates** - Built-in interface, custom templates optional

## 🚀 Quick Start

### Build and run

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go build -o snakeflex

# Run without authentication
./snakeflex

# Run with password protection
./snakeflex --pass "mySecurePassword"

# Run with specific script
./snakeflex --file your_script.py --pass "secret"
```

### Cross-platform builds

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o snakeflex.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o snakeflex-macos

# Linux
GOOS=linux GOARCH=amd64 go build -o snakeflex-linux
```

**Then open your browser:**

```
# Without authentication
http://localhost:8090

# With authentication
http://localhost:8090/login
```

## 📋 Command Line Options

| Flag                     | Default         | Description                                    |
| ------------------------ | --------------- | ---------------------------------------------- |
| `--file`                 | *(none)*        | Python script to execute (optional)           |
| `--port`                 | `8090`          | Server port                                    |
| `--pass`                 | *(none)*        | Password for authentication (optional)         |
| `--template`             | `terminal.html` | Custom HTML template file (optional)          |
| `--verbose`              | `false`         | Enable detailed logging                        |
| `--disable-file-manager` | `false`         | Disable file management for enhanced security  |
| `--disable-shell`        | `false`         | Disable interactive shell for enhanced security|

## 🔐 Authentication System

SnakeFlex V1.6 includes robust password authentication for secure deployments:

### **🛡️ Security Features**
* **Password hashing** - SHA-256 hashing, no plaintext storage
* **Session management** - 24-hour sessions with automatic cleanup
* **Secure cookies** - HttpOnly, SameSite, and Secure flags
* **Beautiful login page** - Terminal-themed authentication interface
* **Session expiry** - Automatic logout after inactivity

### **🚀 Authentication Usage**

```bash
# Start without authentication (current behavior)
./snakeflex

# Enable password protection
./snakeflex --pass "mySecurePassword123"

# Production deployment with security
./snakeflex --pass "productionPassword" --disable-shell --port 8080

# Development with authentication
./snakeflex --pass "dev" --file "main.py" --verbose
```

### **🔒 When Authentication is Enabled**
* All routes are protected by authentication middleware
* Users are redirected to `/login` page if not authenticated
* Password is hashed with SHA-256 before comparison
* Sessions last 24 hours with secure cookie storage
* Access `/logout` to clear session and log out

## 🎯 Script Selection Workflows

### **🚀 Dynamic Selection with Navigation** (Recommended)
Start without specifying a file and navigate your project structure:

```bash
# Start SnakeFlex
./snakeflex --pass "optional"

# In the UI:
# 1. 📁 Double-click folders to explore project structure
# 2. Navigate to scripts/ folder
# 3. Click "Shell" to install any needed packages
# 4. Right-click game.py → "Edit" → Modify code → Save
# 5. Right-click game.py → "Set as Executable" → Run
# 6. Navigate to different folder (e.g., data-processing/)
# 7. Right-click analysis.py → "Set as Executable" → Run
# 8. Switch between scripts and folders seamlessly
```

## ⌨️ Interactive Shell Access

SnakeFlex V1.6 includes full interactive shell access directly in your browser:

### **🖥️ Shell Features**
* **Full terminal emulation** - Complete shell experience using xterm.js
* **Platform adaptive** - PowerShell on Windows, Bash on Linux/macOS
* **Real-time interaction** - Full bidirectional communication
* **Proper PTY support** - True pseudo-terminal on Unix systems
* **Resizable terminal** - Auto-adjusts to window size changes
* **Working directory sync** - Shell starts in your project directory

### **⚠️ Windows Shell Limitations**
**Note**: The interactive shell may not work properly on Windows due to PTY (pseudo-terminal) limitations. If you experience shell issues on Windows:
- Python script execution will still work perfectly
- File management and editing work normally
- Consider using the `--disable-shell` flag on Windows for stability
- Linux and macOS shell support is fully functional

## 📁 Folder Navigation

SnakeFlex V1.6 introduces comprehensive folder navigation for IDE-like project management:

### **🎯 Navigation Features**
* **Breadcrumb navigation** - Visual path trail with clickable segments
* **Up button** (↑) - Navigate to parent directory instantly
* **Home button** (🏠) - Jump back to project root
* **Double-click folders** - Navigate into subdirectories
* **Current location indicator** - Always know where you are
* **Secure boundaries** - Cannot navigate outside working directory

## 📝 Built-in Code Editor

SnakeFlex V1.6 includes a powerful built-in code editor for seamless development workflow:

### **✨ Editor Features**
* 🖱️ **Right-click to edit** - Edit any text file directly in the browser
* ⌨️ **Keyboard shortcuts** - Ctrl+S to save, Escape to close, Tab for proper indentation
* 🎨 **Syntax-aware** - Monospace font, proper tab handling, and code formatting
* 🔄 **Auto-save detection** - Warns before closing unsaved changes
* 📱 **Full-screen editor** - Immersive editing experience with status feedback
* 📄 **Multi-format support** - Edit .py, .txt, .js, .html, .css, .json, .md files

## 🔒 Security Modes

SnakeFlex offers multiple security configurations to balance functionality with security:

### **🛡️ Maximum Security Mode** (`--disable-file-manager --disable-shell`)
Perfect for production environments, shared systems, or when you need maximum security:

```bash
# Maximum security production deployment
./snakeflex --pass "strongPassword" --disable-file-manager --disable-shell --port 8080

# Educational environment (students can run existing scripts only)
./snakeflex --pass "classPassword" --disable-file-manager --disable-shell
```

* ✅ **Terminal functionality** - Full Python script execution with interactive input
* ✅ **Script switching** - Right-click existing Python files to switch between them
* ✅ **Folder navigation** - Browse existing project structure
* ❌ **File operations** - Upload, download, create, delete disabled
* ❌ **Shell access** - No terminal/command-line access
* 🔒 **Zero attack surface** - File management and shell completely removed

### **🔐 Partial Security Modes**

**File Manager Disabled, Shell Enabled:**
```bash
# Shell access for package management, no file operations
./snakeflex --pass "devPassword" --disable-file-manager
```

**Shell Disabled, File Manager Enabled:**
```bash
# File management and navigation without shell access
./snakeflex --pass "webdevPassword" --disable-shell
```

### **📂 Full Mode** (default)
Complete development environment with all features:

```bash
# Full development environment with authentication
./snakeflex --pass "fullPassword"

# Full development environment without authentication
./snakeflex
```

## 🎯 Perfect for

### **Development & Education** (Full Mode)
* **Education** - Teaching Python with secure access and organized project structure
* **Remote development** - Full development environment with authentication
* **Data science** - Secure notebook-like experience with folder organization
* **Workshops** - Password-protected collaborative learning environment

### **Production & Security** (Secure Modes)
* **Production deployment** - Secure Python script execution with authentication
* **Shared environments** - Multiple users with individual authentication
* **Corporate environments** - Compliant with security policies
* **Public demos** - Safe script execution with access control

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process, plus additional WebSocket connections for interactive shell access, and a REST API for file management and editing operations (when enabled). The authentication system uses SHA-256 password hashing with secure session management.

**Architecture:**
* **Authentication layer** - Password hashing with secure session cookies
* **WebSocket connection** - Real-time terminal communication for Python execution
* **Shell WebSocket** - Interactive shell communication with PTY support
* **REST API** - File management, editing operations, and folder navigation (optional)
* **Built-in code editor** - Browser-based editing with syntax awareness
* **Security layer** - Path validation, access control, and feature disabling

## 🎨 Features in action

**Secure development workflow:**

```bash
# Start SnakeFlex with authentication
./snakeflex --pass "mySecurePassword"

# 1. Navigate to http://localhost:8090/login
# 2. Enter password and access terminal
# 3. Navigate to scripts/ folder
# 4. Click "Shell" → pip install requests pandas numpy
# 5. Right-click data_analysis.py → "Edit" → Write code → Save
# 6. Right-click data_analysis.py → "Set as Executable" → Run
# 7. Navigate to data/ folder → Upload CSV files
# 8. Navigate back to scripts/ → Run analysis
# 9. Navigate to output/ folder → Download results
```

## 🔧 Requirements

### For building:
* **Go 1.21+** - For compiling the server
* **Git** - For cloning the repository

### For running (built binary):
* **Python 3.x** - Any Python 3 installation
* **Modern browser** - Chrome, Firefox, Safari, Edge with WebSocket support

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support for terminal and shell communication
* `github.com/creoak/pty` - PTY (pseudo-terminal) support for Unix shell integration

## 🔒 Security Features

### **Authentication Security**
* **Password hashing** - SHA-256 hashing prevents plaintext storage
* **Session management** - Secure random tokens with 24-hour expiry
* **Secure cookies** - HttpOnly, SameSite, and Secure flags for production
* **Session cleanup** - Automatic removal of expired sessions
* **Login protection** - Failed attempts are logged (with --verbose)

### **Always Active Security**
* **Path validation** - Prevents directory traversal attacks
* **Working directory restriction** - All operations limited to project folder
* **Input sanitization** - All file paths and operations are validated
* **Template security** - Embedded templates prevent injection attacks

### **When to Use Authentication**

**Without Authentication:**
* ✅ Personal development on trusted networks
* ✅ Local development environments
* ✅ Isolated containers or VMs

**With Authentication:**
* ✅ Shared development servers
* ✅ Production deployments
* ✅ Educational environments with multiple users
* ✅ Remote access over untrusted networks
* ✅ Corporate compliance requirements

## 🐛 Known limitations

* **Windows shell issues** - Interactive shell may not work properly on Windows due to PTY library limitations; use `--disable-shell` on Windows for stability
* Sessions don't persist between server restarts
* Authentication is session-based, not user-based (single password for all access)
* File uploads are limited to 500MB by default
* Very long-running scripts might timeout in some browsers

## 💡 Pro tips

### Authentication Tips
* **Use strong passwords** - Especially for production deployments
* **Log out when done** - Use `/logout` endpoint or close browser
* **Monitor with --verbose** - Track authentication attempts and access
* **Combine with security modes** - Use `--pass` with `--disable-shell` for maximum security
* **HTTPS in production** - Use reverse proxy for secure cookie transmission

### Security Tips
* **Choose appropriate mode** - Match security level to environment and trust level
* **Production deployment** - Always use `--pass` with appropriate disable flags
* **Network restrictions** - Use firewall rules to limit access
* **Container isolation** - Run in Docker for additional security layers
* **Regular password updates** - Change passwords periodically for shared environments

## 🤝 Contributing

Found a bug? Have an idea? Pull requests are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🗺️ Roadmap

### **Near Term**
* 🔐 **User management** - Multiple users with individual passwords
* 🎨 **Syntax highlighting** - Full Python syntax highlighting in the editor
* 🔍 **File search** - Quick file finding across project structure
* ⚡ **Quick script switching** - Keyboard shortcuts for common scripts
* 📝 **Editor improvements** - Line numbers, find/replace, better indentation

### **Future Enhancements**
* 🌍 **Multi-user support** - Collaborative development features
* 🔐 **Advanced authentication** - LDAP, OAuth, two-factor authentication
* 📊 **Usage analytics** - Security and performance monitoring
* 🐳 **Docker images** - Pre-built containers for easy deployment
* 💡 **Smart project detection** - Auto-organize files based on project type

---

*Made with ❤️ and ☕. Secure by design, powerful by choice.*