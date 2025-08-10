# 🐍 SnakeFlex V1.75

*A secure, modern web-based Python development environment with interactive shell access, folder navigation, and **reverse proxy support**.*

Run any Python script in your browser with real-time output, interactive input support, comprehensive file management with folder navigation, built-in code editing, **full interactive shell access**, **password authentication**, and **seamless reverse proxy integration**. No modifications to your code required.

## ✨ What it does

SnakeFlex V1.75 creates a beautiful web-based development environment for Python scripts with complete terminal capabilities and IDE-like folder navigation. Think of it as your IDE, terminal, and shell combined, but accessible from anywhere with a web browser - now with enterprise-grade reverse proxy support.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
* 🔐 **Password authentication** - Secure access with session management and rate limiting
* 🔄 **Reverse proxy ready** - Deploy seamlessly behind nginx, Apache, Traefik, or any reverse proxy
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

# Run behind reverse proxy (e.g., at /snakeflex path)
./snakeflex --pass "secret" --base-path "/snakeflex"
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
# Direct access
http://localhost:8090

# With authentication
http://localhost:8090/login

# Behind reverse proxy
https://yourdomain.com/snakeflex/
```

## 📋 Command Line Options

| Flag                     | Default         | Description                                    |
| ------------------------ | --------------- | ---------------------------------------------- |
| `--file`                 | *(none)*        | Python script to execute (optional)           |
| `--port`                 | `8090`          | Server port                                    |
| `--pass`                 | *(none)*        | Password for authentication (optional)         |
| `--base-path`            | *(none)*        | Base path for reverse proxy (e.g., `/snakeflex`) |
| `--template`             | `terminal.html` | Custom HTML template file (optional)          |
| `--verbose`              | `false`         | Enable detailed logging                        |
| `--disable-file-manager` | `false`         | Disable file management for enhanced security  |
| `--disable-shell`        | `false`         | Disable interactive shell for enhanced security|

## 🔄 Reverse Proxy Support

SnakeFlex V1.75 includes enterprise-grade reverse proxy support for production deployments:

### **🌍 Deployment Features**
* **Path-aware routing** - Automatically handles subpath deployments
* **Header detection** - Reads `X-Forwarded-Prefix` and `X-Script-Name` headers
* **WebSocket proxying** - Full support for WebSocket connections through proxies
* **Session path scoping** - Cookies work correctly under any subpath
* **Relative URL handling** - All redirects and form actions work seamlessly

### **🔧 Nginx Configuration Example**

```nginx
# Simple reverse proxy setup
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    # SSL configuration
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # Proxy SnakeFlex at /snakeflex path
    location /snakeflex/ {
        proxy_pass http://127.0.0.1:8090/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Prefix /snakeflex;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # Timeout settings for long-running scripts
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
    
    # Redirect /snakeflex to /snakeflex/
    location = /snakeflex {
        return 301 $scheme://$host/snakeflex/;
    }
    
    # Your main application
    location / {
        # Your main app configuration
    }
}
```

### **🚀 Reverse Proxy Usage**

**Method 1: Explicit Base Path**
```bash
# Run SnakeFlex with explicit base path
./snakeflex --pass "password" --base-path "/snakeflex"

# Access via: https://yourdomain.com/snakeflex/
```

**Method 2: Automatic Detection**
```bash
# Run SnakeFlex normally, let it detect path from headers
./snakeflex --pass "password"

# Set X-Forwarded-Prefix header in your reverse proxy
# SnakeFlex automatically detects and adapts
```

### **🔗 Supported Reverse Proxies**
* **Nginx** - Full support with WebSocket proxying
* **Apache** - With mod_proxy and mod_proxy_wstunnel
* **Traefik** - Native support with automatic service discovery
* **HAProxy** - With WebSocket upgrade support
* **Cloudflare** - Through Cloudflare Tunnel or Load Balancer
* **AWS ALB** - Application Load Balancer with WebSocket support

## 🔐 Authentication System

SnakeFlex V1.75 includes robust password authentication with rate limiting for secure deployments:

### **🛡️ Security Features**
* **Password hashing** - SHA-256 hashing, no plaintext storage
* **Session management** - 24-hour sessions with automatic cleanup
* **Rate limiting** - Progressive lockout after failed attempts (3+ = 1min, 6+ = 10min, 10+ = 1hr)
* **Secure cookies** - HttpOnly, SameSite, and Secure flags with proper path scoping
* **Beautiful login page** - Terminal-themed authentication interface
* **Session expiry** - Automatic logout after inactivity
* **Proxy-aware** - Cookies and redirects work correctly behind reverse proxies

### **🚀 Authentication Usage**

```bash
# Start without authentication (current behavior)
./snakeflex

# Enable password protection
./snakeflex --pass "mySecurePassword123"

# Production deployment with reverse proxy
./snakeflex --pass "productionPassword" --base-path "/snakeflex" --disable-shell --port 8080

# Development with authentication behind proxy
./snakeflex --pass "dev" --file "main.py" --verbose
```

### **🔒 When Authentication is Enabled**
* All routes are protected by authentication middleware
* Users are redirected to proper login page (with base path support)
* Password is hashed with SHA-256 before comparison
* Sessions last 24 hours with secure cookie storage
* Rate limiting protects against brute force attacks
* Access `/logout` to clear session and log out

## 🎯 Script Selection Workflows

### **🚀 Dynamic Selection with Navigation** (Recommended)
Start without specifying a file and navigate your project structure:

```bash
# Start SnakeFlex (standalone)
./snakeflex --pass "optional"

# Start SnakeFlex behind reverse proxy
./snakeflex --pass "optional" --base-path "/dev-env"

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

SnakeFlex V1.75 includes full interactive shell access directly in your browser with proxy support:

### **🖥️ Shell Features**
* **Full terminal emulation** - Complete shell experience using xterm.js
* **Platform adaptive** - PowerShell on Windows, Bash on Linux/macOS
* **Real-time interaction** - Full bidirectional communication through WebSockets
* **Proper PTY support** - True pseudo-terminal on Unix systems
* **Resizable terminal** - Auto-adjusts to window size changes
* **Working directory sync** - Shell starts in your project directory
* **Proxy compatible** - Works seamlessly through reverse proxies

### **⚠️ Windows Shell Limitations**
**Note**: The interactive shell may not work properly on Windows due to PTY (pseudo-terminal) limitations. If you experience shell issues on Windows:
- Python script execution will still work perfectly
- File management and editing work normally
- Consider using the `--disable-shell` flag on Windows for stability
- Linux and macOS shell support is fully functional

## 🔒 Security Modes

SnakeFlex offers multiple security configurations to balance functionality with security:

### **🛡️ Maximum Security Mode** (`--disable-file-manager --disable-shell`)
Perfect for production environments, shared systems, or when you need maximum security:

```bash
# Maximum security production deployment behind reverse proxy
./snakeflex --pass "strongPassword" --base-path "/secure-python" --disable-file-manager --disable-shell --port 8080

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
./snakeflex --pass "devPassword" --base-path "/dev" --disable-file-manager
```

**Shell Disabled, File Manager Enabled:**
```bash
# File management and navigation without shell access
./snakeflex --pass "webdevPassword" --base-path "/webdev" --disable-shell
```

### **📂 Full Mode** (default)
Complete development environment with all features:

```bash
# Full development environment with authentication behind proxy
./snakeflex --pass "fullPassword" --base-path "/python-ide"

# Full development environment without authentication
./snakeflex
```

## 🎯 Perfect for

### **Development & Education** (Full Mode)
* **Education** - Teaching Python with secure access and organized project structure
* **Remote development** - Full development environment with authentication
* **Data science** - Secure notebook-like experience with folder organization
* **Workshops** - Password-protected collaborative learning environment
* **Corporate development** - Deploy behind company reverse proxy with SSO

### **Production & Security** (Secure Modes)
* **Production deployment** - Secure Python script execution with authentication
* **Shared environments** - Multiple users with individual authentication
* **Corporate environments** - Compliant with security policies, deploy behind existing infrastructure
* **Public demos** - Safe script execution with access control
* **Multi-tenant platforms** - Deploy multiple instances behind different paths

### **Enterprise Integration**
* **Existing web applications** - Integrate as `/python` or `/dev` subpath
* **Company intranets** - Deploy behind corporate reverse proxy
* **Load balanced deployments** - Multiple instances with session affinity
* **SSL termination** - HTTPS handled by reverse proxy

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process, plus additional WebSocket connections for interactive shell access, and a REST API for file management and editing operations (when enabled). The authentication system uses SHA-256 password hashing with secure session management and rate limiting. All components are proxy-aware and work seamlessly behind reverse proxies.

**Architecture:**
* **Authentication layer** - Password hashing with secure session cookies and rate limiting
* **Proxy detection** - Automatic base path detection from `X-Forwarded-Prefix` headers
* **WebSocket connection** - Real-time terminal communication for Python execution
* **Shell WebSocket** - Interactive shell communication with PTY support
* **REST API** - File management, editing operations, and folder navigation (optional)
* **Built-in code editor** - Browser-based editing with syntax awareness
* **Security layer** - Path validation, access control, and feature disabling

## 🎨 Features in action

**Secure development workflow behind reverse proxy:**

```bash
# Start SnakeFlex behind reverse proxy
./snakeflex --pass "mySecurePassword" --base-path "/python-dev"

# 1. Navigate to https://company.com/python-dev/login
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
* **Reverse proxy** - Nginx, Apache, Traefik, etc. (optional)

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support for terminal and shell communication
* `github.com/creoak/pty` - PTY (pseudo-terminal) support for Unix shell integration

## 🔒 Security Features

### **Authentication Security**
* **Password hashing** - SHA-256 hashing prevents plaintext storage
* **Session management** - Secure random tokens with 24-hour expiry
* **Rate limiting** - Progressive lockout system (3 attempts = 1min, 6 = 10min, 10+ = 1hr)
* **Secure cookies** - HttpOnly, SameSite, and Secure flags for production
* **Session cleanup** - Automatic removal of expired sessions
* **Login protection** - Failed attempts are logged and rate limited
* **Proxy-aware security** - Proper IP detection behind reverse proxies

### **Always Active Security**
* **Path validation** - Prevents directory traversal attacks
* **Working directory restriction** - All operations limited to project folder
* **Input sanitization** - All file paths and operations are validated
* **Template security** - Embedded templates prevent injection attacks
* **CSRF protection** - Session-based protection against cross-site requests

### **Reverse Proxy Security**
* **Header validation** - Secure handling of proxy headers
* **Path isolation** - Base path prevents access to other applications
* **Cookie scoping** - Sessions properly scoped to base path
* **URL construction** - All redirects respect proxy configuration

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
* ✅ Multi-tenant deployments behind reverse proxy

## 🐛 Known limitations

* **Windows shell issues** - Interactive shell may not work properly on Windows due to PTY library limitations; use `--disable-shell` on Windows for stability
* Sessions don't persist between server restarts
* Authentication is session-based, not user-based (single password for all access)
* File uploads are limited to 500MB by default
* Very long-running scripts might timeout in some browsers
* WebSocket connections require proper proxy configuration for reverse proxy setups

## 💡 Pro tips

### Authentication Tips
* **Use strong passwords** - Especially for production deployments
* **Log out when done** - Use `/logout` endpoint or close browser
* **Monitor with --verbose** - Track authentication attempts and access
* **Combine with security modes** - Use `--pass` with `--disable-shell` for maximum security
* **HTTPS in production** - Use reverse proxy for secure cookie transmission

### Reverse Proxy Tips
* **WebSocket support** - Ensure your reverse proxy supports WebSocket upgrades
* **Timeout configuration** - Set appropriate timeouts for long-running Python scripts
* **Path consistency** - Use `--base-path` flag to match your proxy configuration
* **Header forwarding** - Set `X-Forwarded-Prefix` header for automatic detection
* **SSL termination** - Handle HTTPS at the reverse proxy level
* **Session affinity** - Use sticky sessions for multi-instance deployments

### Security Tips
* **Choose appropriate mode** - Match security level to environment and trust level
* **Production deployment** - Always use `--pass` with appropriate disable flags
* **Network restrictions** - Use firewall rules to limit access
* **Container isolation** - Run in Docker for additional security layers
* **Regular password updates** - Change passwords periodically for shared environments
* **Proxy security** - Configure reverse proxy with proper security headers

### Deployment Tips
* **Start simple** - Test locally before deploying behind reverse proxy
* **Check logs** - Use `--verbose` to troubleshoot proxy issues
* **Test WebSockets** - Verify shell and terminal functionality through proxy
* **Monitor sessions** - Watch for proper cookie scoping and session management
* **Health checks** - Configure proxy health checks for high availability

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
* 🐳 **Docker support** - Official Docker images with reverse proxy examples

### **Future Enhancements**
* 🌍 **Multi-user support** - Collaborative development features
* 🔐 **Advanced authentication** - LDAP, OAuth, two-factor authentication
* 📊 **Usage analytics** - Security and performance monitoring
* 🔄 **Load balancing** - Built-in support for multi-instance deployments
* 💡 **Smart project detection** - Auto-organize files based on project type
* 🌐 **API endpoints** - REST API for external integrations
* 🔧 **Plugin system** - Extensible architecture for custom features

---

*Made with ❤️ and ☕. Secure by design, powerful by choice, proxy-ready for production.*