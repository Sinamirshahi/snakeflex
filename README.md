# 🐍 SnakeFlex V4.6

*A modern web-based Python development environment that just works.*

Run any Python script in your browser with real-time output, interactive input support, comprehensive file management, and built-in code editing. No modifications to your code required.

## ✨ What it does

SnakeFlex V4.6 creates a beautiful web-based development environment for Python scripts. Think of it as your IDE and terminal combined, but accessible from anywhere with a web browser.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
* 📝 **Built-in code editor** - Edit Python files directly in the browser with syntax awareness
* 📂 **File manager** - Browse, upload, download, and manage files with drag & drop
* 🎯 **Dynamic script selection** - Switch between Python scripts with right-click menu
* 🚀 **No file required** - Start without specifying a script, choose dynamically in the UI
* 💬 **Interactive input** - Handle `input()` calls seamlessly
* ⚡ **Real-time output** - See your script's output as it happens
* 🎨 **Modern UI** - GitHub-inspired dark interface with resizable panels
* 🔒 **Security modes** - Full-featured or secure terminal-only mode
* 🔄 **Cross-platform** - Windows, macOS, and Linux support
* 💾 **Embedded templates** - Built-in interface, custom templates optional

## 🚀 Quick Start

### Option 1: Start with a specific script

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go run main.go --file your_script.py
```

### Option 2: Start and choose script dynamically

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go run main.go
```

### Option 3: Build and run

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go build -o snakeflex

# Run with specific script
./snakeflex --file your_script.py

# Or run without file and choose in UI
./snakeflex
```

### Option 4: Build for different platforms

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
http://localhost:8090
```

**🎯 Choose your workflow:**
1. **Right-click any Python file** → "Set as Executable" → Click "Run Script" ✨
2. **Right-click any text file** → "Edit" → Modify code in browser → Save → Run ✨
3. **Start with `--file script.py`** → Click "Run Script" immediately
4. **Upload Python files** → Right-click to edit or set executable → Run

## 📋 Usage

### With Go (development)

```bash
# Start without pre-selecting a script - choose in UI
go run main.go

# Basic usage with pre-selected script
go run main.go --file script.py

# Secure mode (terminal only)
go run main.go --file script.py --disable-file-manager

# Start in secure mode, choose script dynamically
go run main.go --disable-file-manager

# Custom port
go run main.go --port 3000

# Production deployment with security
go run main.go --port 8080 --disable-file-manager

# Custom template (optional - embedded template used as fallback)
go run main.go --template custom.html

# Verbose logging
go run main.go --verbose
```

### With built binary (production)

```bash
# After building with: go build -o snakeflex

# Start and choose script in UI
./snakeflex

# Start with specific script
./snakeflex --file script.py

# Various configurations
./snakeflex --port 3000
./snakeflex --disable-file-manager --verbose

# Windows
snakeflex.exe
snakeflex.exe --file script.py
snakeflex.exe --disable-file-manager
```

### Command Line Options

| Flag                     | Default         | Description                                    |
| ------------------------ | --------------- | ---------------------------------------------- |
| `--file`                 | *(none)*        | Python script to execute (optional)           |
| `--port`                 | `8090`          | Server port                                    |
| `--template`             | `terminal.html` | Custom HTML template file (optional)          |
| `--verbose`              | `false`         | Enable detailed logging                        |
| `--disable-file-manager` | `false`         | Disable file management for enhanced security  |

## 🎯 Script Selection Workflows

SnakeFlex offers flexible ways to work with Python scripts:

### **🚀 Dynamic Selection** (Recommended)
Start without specifying a file and choose scripts on-the-fly:

```bash
# Start SnakeFlex
./snakeflex

# In the UI:
# 1. Right-click any Python file in the file manager
# 2. Select "Set as Executable" 
# 3. Click "Run Script" button
# 4. Switch to different scripts anytime by repeating steps 1-2
```

**Benefits:**
* 🔄 **Switch between scripts instantly** without restarting the server
* 📁 **Work with multiple files** in the same session
* 🎨 **Visual script management** - see which script is currently selected
* ⚡ **No command-line changes** needed when switching scripts

### **⚡ Pre-configured Start**
Traditional approach with a specific script:

```bash
# Start with specific script ready to run
./snakeflex --file my_script.py

# Still allows switching:
# - Right-click other Python files to switch
# - Upload new files and set them as executable
```

### **🔒 Secure Mode Selection**
Even in secure mode, you can still switch between existing Python scripts:

```bash
# Start in secure mode
./snakeflex --disable-file-manager

# Script selection still works:
# - Right-click existing Python files
# - Set as executable and run
# - No file upload/download, but script switching is available
```

## 📝 Built-in Code Editor (NEW in V4.5)

SnakeFlex V4.6 introduces a powerful built-in code editor for seamless development workflow:

### **✨ Editor Features**
* 🖱️ **Right-click to edit** - Edit any text file directly in the browser
* ⌨️ **Keyboard shortcuts** - Ctrl+S to save, Escape to close, Tab for proper indentation
* 🎨 **Syntax-aware** - Monospace font, proper tab handling, and code formatting
* 🔄 **Auto-save detection** - Warns before closing unsaved changes
* 📱 **Full-screen editor** - Immersive editing experience with status feedback
* 🛡️ **Secure** - Same path validation as all file operations
* 📄 **Multi-format support** - Edit .py, .txt, .js, .html, .css, .json, .md files

### **🚀 Editor Workflow**
```bash
# Start SnakeFlex
./snakeflex

# Edit-Run-Test cycle:
# 1. Right-click any text file → "Edit"
# 2. Make changes in the full-screen editor
# 3. Press Ctrl+S to save (or click Save button)
# 4. Press Escape to close editor (or click Close button)
# 5. If it's a Python file: Right-click → "Set as Executable" → Run
# 6. Repeat for rapid development
```

### **⌨️ Editor Shortcuts**
* **Ctrl+S** - Save file
* **Escape** - Close editor (with unsaved changes warning)
* **Tab** - Proper indentation (doesn't lose focus)

### **🎯 Perfect for**
* **Rapid prototyping** - Edit and test scripts without leaving the browser
* **Educational environments** - Students can edit and run code in one interface
* **Remote development** - Full development workflow over the web
* **Code reviews** - Quick edits and testing
* **Data science** - Modify analysis scripts and run immediately

## 📂 File Management & Script Selection

*Available in Full Mode only. Use `--disable-file-manager` to disable for security.*

### **🎯 Right-Click Script Selection**
The most intuitive way to work with Python scripts:

1. **📁 Browse files** - See all Python files in the left sidebar
2. **🖱️ Right-click any .py file** - Context menu appears
3. **📝 Edit (NEW)** - Open file in built-in editor for modifications
4. **▶️ Set as Executable** - Script becomes the active one
5. **🚀 Click "Run Script"** - Execute the selected script
6. **🔄 Repeat for different scripts** - Switch anytime without restart

### **📈 Visual Feedback**
* **Active script indicator** - Shows currently selected script name
* **Status updates** - "Ready", "Running", "Waiting for Input" states
* **File icons** - Python files show 🐍 icon for easy identification
* **Context menu visibility** - Right-click options adapt to file type
* **Editor integration** - Seamless transition between editing and running

### **📁 File Operations**
* **Edit files** - Built-in code editor for all text files
* **Browse files** - Tree view of your working directory
* **Upload files** - Drag & drop or click to upload multiple files
* **Download files** - One-click download for any file
* **Create files/folders** - Built-in creation dialog
* **Delete files/folders** - Safe deletion with confirmation
* **File icons** - Visual file type indicators (🐍 .py, 🌐 .html, 📄 .txt, etc.)

### **🎛️ Interface Features**
* **Resizable sidebar** - Drag the edge to adjust panel width
* **Context menus** - Right-click files for quick actions (Edit, Set Executable, Download, Delete)
* **Real-time updates** - File list refreshes automatically
* **Security protection** - Prevents access outside working directory
* **Full-screen editor** - Distraction-free editing environment

## 🎨 Template System

SnakeFlex uses a smart template fallback system that makes distribution effortless:

### **📦 Template Priority (Smart Fallback)**

1. **🎯 Custom External Template** - If you specify `--template custom.html` and the file exists
2. **💾 Embedded Template** - Built-in `templates/terminal.html` (always available)
3. **🔧 Minimal Fallback** - Basic HTML template (emergency only)

### **💡 Why This Matters**

* **Zero Dependencies** - Binary works standalone with embedded template
* **Customization Freedom** - Override with your own template when needed
* **Bulletproof Distribution** - No missing template files to worry about
* **Professional Polish** - Always presents a clean interface

### **🔨 Using Custom Templates**

```bash
# Use embedded template (default - always works)
./snakeflex

# Use custom template with fallback protection
./snakeflex --template my-custom.html

# Verbose mode shows which template is being used
./snakeflex --template my-custom.html --verbose
# Output: ✅ Using external template: my-custom.html
#    OR: 💾 External template 'my-custom.html' not found, using embedded template
```

Custom templates support these variables:
- `{{INITIAL_PYTHON_FILE}}` - Name of the initially selected Python script
- `{{WORKING_DIR}}` - Current working directory
- `{{FILE_MANAGER_ENABLED}}` - Whether file management is enabled

## 🔒 Security Modes

SnakeFlex offers two operational modes to balance functionality with security:

### **🛡️ Secure Mode** (`--disable-file-manager`)
Perfect for production environments, shared systems, or when you need maximum security:

* ✅ **Terminal functionality** - Full Python script execution with interactive input
* ✅ **Real-time output** - All terminal features work normally  
* ✅ **Script switching** - Right-click existing Python files to switch between them
* ❌ **File operations** - Upload, download, create, delete disabled
* ❌ **File browsing** - Directory listing disabled
* ❌ **Code editing** - Built-in editor disabled for security
* ❌ **API endpoints** - All `/api/files/*` routes disabled
* 🔒 **Zero attack surface** - File management completely removed

```bash
# Production deployment - can still switch between existing scripts
./snakeflex --disable-file-manager --port 8080

# Educational environment (students can run existing scripts but not modify files)
./snakeflex --disable-file-manager

# Container deployment with script flexibility
docker run -p 8090:8090 snakeflex --disable-file-manager
```

### **📂 Full Mode** (default)
Complete development environment with all features:

* ✅ **All terminal functionality**
* ✅ **Built-in code editor** - Edit files directly in browser
* ✅ **Complete file management**
* ✅ **Drag & drop uploads**
* ✅ **File browsing and organization**
* ✅ **Dynamic script selection and switching**
* 🔒 **Secure within working directory**

### **API Endpoints** (Full Mode Only)
* `GET /api/files` - Browse directory contents
* `GET /api/files/content?path=file.py` - Read file content for editing
* `PUT /api/files/content` - Save edited file content
* `GET /api/files/download?path=file.py` - Download files
* `POST /api/files/upload` - Upload multiple files
* `POST /api/files/create` - Create new files/folders
* `DELETE /api/files/delete?path=file.py` - Delete files/folders

*Note: These endpoints return 403 Forbidden when file management is disabled.*

## 🎯 Perfect for

### **Development & Education** (Full Mode)
* **Education** - Teaching Python with file management, editing, and easy script switching
* **Demos** - Showing off multiple projects with live editing and one-click switching
* **Remote development** - Full file management and editing without SSH
* **Data science** - Upload datasets, edit analysis scripts, test different approaches, download results
* **Workshops** - Students can upload, edit, switch between, and test multiple scripts
* **Experimentation** - Quickly edit and test different Python files in the same environment
* **Code reviews** - Live editing and testing during review sessions
* **Pair programming** - Collaborative editing and immediate execution

### **Production & Security** (Secure Mode)
* **Production deployment** - Secure Python script execution with script flexibility
* **Shared environments** - Multiple users can switch scripts without file access risks
* **Container deployment** - Minimal attack surface with script selection
* **Corporate environments** - Compliant with security policies
* **Educational restrictions** - Students can run different scripts but not modify files
* **Public demos** - Safe script execution with multiple demonstration options

## 📦 Distribution

SnakeFlex compiles to a single binary with **embedded templates** and no dependencies (except Python on the target system). Perfect for:

* **Instant deployment** - Single binary with built-in interface and editor
* **Multi-script sharing** - Upload a folder of Python scripts, users can edit and switch between them
* **Educational environments** - Complete development environment with editing and script flexibility
* **Client presentations** - Edit and demonstrate multiple Python scripts without restart
* **Remote execution** - Lightweight server for Python development with editing capabilities
* **Workshop distribution** - One-click setup with editable example scripts
* **Secure deployment** - Production-ready with script selection flexibility

### **📋 Distribution Examples**

```bash
# Build for your platform
go build -o snakeflex

# Multi-script distribution (embedded template)
mkdir python-demos
cp snakeflex python-demos/
cp *.py python-demos/          # All your Python scripts
cp -r data/ python-demos/      # Include data directories
echo './snakeflex' > python-demos/start.sh  # No specific file - choose in UI
chmod +x python-demos/start.sh
# Users right-click any Python file to edit or run it!

# Workshop package with multiple examples
mkdir workshop-materials
cp snakeflex workshop-materials/
cp beginner.py intermediate.py advanced.py workshop-materials/
cp -r examples/ workshop-materials/
echo './snakeflex --port 8080' > workshop-materials/start.sh
chmod +x workshop-materials/start.sh
# Instructors and students can edit and switch between difficulty levels

# Production package with multiple scripts (secure mode)
mkdir secure-deployment
cp snakeflex secure-deployment/
cp script1.py script2.py script3.py secure-deployment/
echo './snakeflex --disable-file-manager --port 8080' > secure-deployment/start.sh
chmod +x secure-deployment/start.sh
# Users can switch between approved scripts securely (no editing)

# Package and distribute
zip -r python-demos.zip python-demos/           # Multi-script with editor, always works
zip -r workshop-materials.zip workshop-materials/   # Educational package with editing
zip -r secure-deployment.zip secure-deployment/     # Production secure (no editing)
```

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process, plus a REST API for file management and editing operations (when enabled). It automatically detects when your script needs input and presents a clean interface for interaction.

The Go server intelligently detects your system's Python installation (`python`, `python3`, or `py`) and runs scripts with proper buffering settings to ensure real-time output. The embedded template system ensures the interface always works, while custom templates allow for branding and customization.

**Architecture:**
* **WebSocket connection** - Real-time terminal communication
* **REST API** - File management and editing operations (optional)
* **Built-in code editor** - Browser-based editing with syntax awareness
* **Dynamic script selection** - Switch between Python files without restart
* **Embedded templates** - Built-in interface with custom override support
* **Security layer** - Path validation, access control, and feature disabling
* **Multi-platform support** - Adaptive Python detection

## 🎨 Features in action

**Dynamic script selection with editing:**

```bash
# Start SnakeFlex
./snakeflex

# In browser: 
# 1. Right-click game.py → "Edit" → Modify code → Save → "Set as Executable" → Run
# 2. Later: Right-click data_analysis.py → "Edit" → Update algorithm → Save → "Set as Executable" → Run  
# 3. Switch and edit between scripts instantly without restarting server
```

**Code editing workflow:**

```python
# Edit this directly in the browser:
def calculate_fibonacci(n):
    # Right-click this file → Edit → Modify this function → Save → Run
    if n <= 1:
        return n
    return calculate_fibonacci(n-1) + calculate_fibonacci(n-2)

# Edit parameters and see results immediately
for i in range(10):
    print(f"fib({i}) = {calculate_fibonacci(i)}")
```

**Interactive input detection:**

```python
name = input("What's your name? ")  # Input box appears automatically
age = int(input("How old are you? "))  # Handles any input type
```

**Real-time output:**

```python
import time
for i in range(5):
    print(f"Processing step {i+1}...")
    time.sleep(1)  # You see each line as it prints
```

**File operations** (Full Mode only):

```python
# Upload data.csv through the file manager
# Edit this script in the browser
import pandas as pd
df = pd.read_csv('data.csv')  # File is now available to your script
print(df.head())

# Make changes to the analysis right in the browser
df_filtered = df[df['value'] > 100]  # Edit this filtering logic live
print(f"Filtered rows: {len(df_filtered)}")

# Results can be saved and downloaded
df_filtered.to_csv('results.csv')  # Use file manager to download results
```

**Error handling:**

```python
print("This goes to stdout")
print("This goes to stderr", file=sys.stderr)  # Different colors
raise Exception("Errors are highlighted")
```

### 🖼️ Screenshot

![Screenshot of SnakeFlex V4.6 Interface](screenshot.png)

*Interface showing built-in code editor, dynamic script selection, file manager panel (Full Mode), and secure terminal-only mode*

## 🔧 Requirements

### For building:

* **Go 1.21+** - For compiling the server
* **Git** - For cloning the repository

### For running (built binary):

* **Python 3.x** - Any Python 3 installation
* **Modern browser** - Chrome, Firefox, Safari, Edge

*Note: The built binary has embedded templates and no Go dependencies - runs on any system with Python.*

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support for terminal communication

## 🔒 Security Features

SnakeFlex includes comprehensive security measures:

### **Always Active Security**
* **Path validation** - Prevents directory traversal attacks
* **Working directory restriction** - File operations limited to project folder
* **Protected files** - Currently executing Python file cannot be deleted
* **Input sanitization** - All file paths and operations are validated
* **Safe uploads** - File uploads are restricted to working directory
* **Template security** - Embedded templates prevent template injection attacks
* **Script validation** - Only Python files can be set as executable
* **Editor security** - File editing restricted to working directory with path validation

### **Enhanced Security Mode** (`--disable-file-manager`)
* **Eliminated attack surface** - File management and editing endpoints completely removed
* **API endpoint disabling** - All `/api/files/*` routes return 403 Forbidden
* **UI adaptation** - Interface clearly shows secure mode status
* **Script switching maintained** - Users can still switch between existing Python scripts
* **No file editing** - Built-in editor disabled in secure mode
* **Defense in depth** - Multiple layers of validation even when disabled
* **Production ready** - Suitable for corporate and shared environments

### **When to Use Secure Mode**
* ✅ **Production deployments** - Reduce attack surface while maintaining script flexibility
* ✅ **Shared systems** - Multiple users can switch scripts without file access or editing risks
* ✅ **Educational restrictions** - Students can run different scripts but not modify files
* ✅ **Container deployment** - Minimal security footprint
* ✅ **Corporate compliance** - Meet security policy requirements
* ✅ **Public demos** - Safe script execution with multiple demonstration options

## 🤝 Contributing

Found a bug? Have an idea? Pull requests are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🐛 Known limitations

* Windows doesn't support full PTY (pseudo-terminal) features
* Very long-running scripts might timeout in some browsers
* File I/O operations in Python scripts access the server's filesystem
* Large output bursts are throttled to prevent WebSocket flooding
* File uploads are limited to 500MB by default (Full Mode only)
* Hidden files and system directories (`.git`, `__pycache__`) are filtered from the file browser
* Secure mode completely disables file management and editing - script switching still available
* Custom templates must be present at startup (embedded template used as fallback)
* Only Python (.py) files can be set as executable scripts
* Editor supports text files but syntax highlighting is basic
* Large files may be slow to load in the editor

## 💡 Pro tips

### Script Selection Tips
* **Start without `--file`** for maximum flexibility - choose scripts in the UI
* **Right-click Python files** to switch between different scripts instantly
* **Check the active script indicator** to see which script will run
* **Upload multiple scripts** and test them all in the same session
* **Use descriptive filenames** - they show clearly in the file manager

### Editor Tips (NEW in V4.5)
* **Use Ctrl+S** frequently to save your work
* **Tab key works properly** - doesn't lose editor focus
* **Escape key** closes editor with unsaved change warnings
* **Edit-Save-Run cycle** - Modify code, save, set executable, run immediately
* **Multiple file types** - Edit .py, .txt, .js, .html, .css, .json, .md files
* **Full-screen editing** - Distraction-free coding environment
* **Status feedback** - Always know if your file saved successfully

### Template Tips
* **Use embedded templates** for zero-dependency distribution
* **Test custom templates** with `--verbose` to see which template loads
* **Embedded templates are bulletproof** - always work even without custom files
* **Template variables** allow dynamic content insertion
* **Custom branding** possible with external templates while keeping embedded fallback

### Terminal Tips (Both Modes)
* Use `print(..., flush=True)` for immediate output in custom scripts
* Press `Ctrl+C` in the terminal to stop long-running scripts
* Check the browser console (F12) for debugging WebSocket issues
* Use `--verbose` flag to debug script execution and template loading

### File Management Tips (Full Mode Only)
* **Edit first, then run** - Right-click file → Edit → Modify → Save → Set Executable → Run
* **Drag and drop** files directly into the upload area for quick uploads
* **Right-click files** to access editing, script selection, download, and delete options
* **Resize the sidebar** by dragging the right edge for more space
* **Use the refresh button** (🔄) to update the file list after external changes
* **Create folders first**, then upload files to organize your workspace
* **Download results** after running data processing scripts

### Security Tips
* **Use `--disable-file-manager`** for production deployments
* **Test in Full Mode**, deploy in Secure Mode for safety
* **Script switching works in both modes** - secure mode still allows Python script selection
* **Monitor logs** with `--verbose` in secure environments
* **Container isolation** - Run in Docker for additional security layers
* **Network restrictions** - Use firewall rules to limit access
* **Secure mode disables editing** - No file modifications possible

### Development Tips
* Built binaries are portable with embedded templates - no external files needed
* Multiple concurrent output streams are handled safely (stdout + stderr)
* File operations and editing provide real-time feedback in the terminal (Full Mode)
* The current Python script file is protected from accidental deletion
* Secure mode provides the same terminal experience with zero file management or editing risk
* Custom templates override embedded ones automatically
* Dynamic script selection eliminates the need to restart for different files
* Editor changes are saved immediately - no auto-save delays

## 🎉 Acknowledgments

Inspired by the need for a complete, browser-based Python development environment that works everywhere while maintaining security flexibility. Built with love for the Python community and educators who need powerful, accessible, and secure tools with zero-dependency distribution, flexible script management, and seamless code editing capabilities.

## 🗺️ Roadmap

### **Near Term**
* 🎨 **Syntax highlighting** - Full Python syntax highlighting in the editor
* 🔍 **File search** - Quick file finding in large projects
* 📊 **Script history** - Remember recently executed scripts
* ⚡ **Quick script switching** - Keyboard shortcuts for common scripts
* 📝 **Editor improvements** - Line numbers, find/replace, better indentation

### **Future Enhancements**
* 📁 **Folder navigation** - Navigate into subdirectories
* 💾 **Project templates** - Quick-start templates for common tasks
* 🌍 **Multi-user support** - Collaborative development features
* 🔐 **Authentication** - User login and access control
* 📊 **Usage analytics** - Security and performance monitoring
* 🐳 **Docker images** - Pre-built containers for easy deployment
* 🎨 **Template gallery** - Community-contributed interface themes
* 🏷️ **Script categories** - Organize scripts by type/purpose
* 📝 **Advanced editor** - Code completion, error detection, multiple tabs
* 🔗 **Git integration** - Version control support
* 📱 **Mobile optimization** - Better mobile browser support

---

*Made with ❤️ and ☕. Secure by design, powerful by choice, edits beautifully.*