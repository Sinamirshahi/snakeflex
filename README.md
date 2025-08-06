# 🐍 SnakeFlex

*A modern web-based Python development environment that just works.*

Run any Python script in your browser with real-time output, interactive input support, and comprehensive file management. No modifications to your code required.

## ✨ What it does

SnakeFlex creates a beautiful web-based development environment for Python scripts. Think of it as your IDE and terminal combined, but accessible from anywhere with a web browser.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
* 📂 **File manager** - Browse, upload, download, and manage files with drag & drop
* 💬 **Interactive input** - Handle `input()` calls seamlessly
* ⚡ **Real-time output** - See your script's output as it happens
* 🎨 **Modern UI** - GitHub-inspired dark interface with resizable panels
* 🔒 **Secure** - Built-in security prevents unauthorized file access
* 🔄 **Cross-platform** - Windows, macOS, and Linux support
* 🚀 **Zero setup** - Just point it at your Python file and go

## 🚀 Quick Start

### Option 1: Run directly with Go

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go run main.go --file your_script.py
```

### Option 2: Build and run

```bash
git clone https://github.com/Sinamirshahi/snakeflex
cd snakeflex
go mod tidy
go build -o snakeflex
./snakeflex --file your_script.py
```

### Option 3: Build for different platforms

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

**Browse files in the left panel, click "Run Script" and watch the magic happen** ✨

## 📋 Usage

### With Go (development)

```bash
# Basic usage
go run main.go --file script.py

# Custom port
go run main.go --file script.py --port 3000

# Custom template
go run main.go --file script.py --template custom.html

# Verbose logging
go run main.go --file script.py --verbose
```

### With built binary (production)

```bash
# After building with: go build -o snakeflex
./snakeflex --file script.py
./snakeflex --file script.py --port 3000
./snakeflex --file script.py --verbose

# Windows
snakeflex.exe --file script.py
```

### Command Line Options

| Flag         | Default         | Description              |
| ------------ | --------------- | ------------------------ |
| `--file`     | `fibonacci.py`  | Python script to execute |
| `--port`     | `8090`          | Server port              |
| `--template` | `terminal.html` | HTML template file       |
| `--verbose`  | `false`         | Enable detailed logging  |

## 📂 File Management Features

SnakeFlex now includes a comprehensive file manager in the left sidebar:

### **📁 File Operations**
* **Browse files** - Tree view of your working directory
* **Upload files** - Drag & drop or click to upload multiple files
* **Download files** - One-click download for any file
* **Create files/folders** - Built-in creation dialog
* **Delete files/folders** - Safe deletion with confirmation
* **File icons** - Visual file type indicators (🐍 .py, 🌐 .html, 📄 .txt, etc.)

### **🎛️ Interface Features**
* **Resizable sidebar** - Drag the edge to adjust panel width
* **Context menus** - Right-click files for quick actions
* **Real-time updates** - File list refreshes automatically
* **Security protection** - Prevents access outside working directory

### **API Endpoints**
* `GET /api/files` - Browse directory contents
* `GET /api/files/download?path=file.py` - Download files
* `POST /api/files/upload` - Upload multiple files
* `POST /api/files/create` - Create new files/folders
* `DELETE /api/files/delete?path=file.py` - Delete files/folders

## 🎯 Perfect for

* **Education** - Teaching Python with file management in a browser
* **Demos** - Showing off projects with easy file sharing
* **Remote development** - Full file management without SSH
* **Code sharing** - Let others browse and run your scripts
* **Presentations** - Live coding with file uploads/downloads
* **Data science** - Upload datasets, run scripts, download results
* **Workshops** - Students can upload their work and test scripts
* **Deployment** - Distribute as a single binary with file management

## 📦 Distribution

SnakeFlex compiles to a single binary with no dependencies (except Python on the target system). Perfect for:

* **Sharing demos** - Send the binary + your Python scripts
* **Educational environments** - Complete development environment in one binary
* **Client presentations** - Professional Python script demonstrations with file management
* **Remote execution** - Lightweight server for Python development
* **Workshop distribution** - One-click setup for coding workshops

```bash
# Build for your platform
go build -o snakeflex

# Package with your scripts
mkdir my-python-workspace
cp snakeflex my-python-workspace/
cp *.py my-python-workspace/
cp terminal.html my-python-workspace/
cp -r data/ my-python-workspace/  # Include data directories
zip -r python-workspace.zip my-python-workspace/
```

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process, plus a REST API for file management operations. It automatically detects when your script needs input and presents a clean interface for interaction.

The Go server intelligently detects your system's Python installation (`python`, `python3`, or `py`) and runs scripts with proper buffering settings to ensure real-time output. File operations are secured to prevent access outside the working directory.

**Architecture:**
* **WebSocket connection** - Real-time terminal communication
* **REST API** - File management operations
* **Security layer** - Path validation and access control
* **Multi-platform support** - Adaptive Python detection

## 🎨 Features in action

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

**File operations:**

```python
# Upload data.csv through the file manager
import pandas as pd
df = pd.read_csv('data.csv')  # File is now available to your script
print(df.head())

# Results can be saved and downloaded
df.to_csv('results.csv')  # Use file manager to download results
```

**Error handling:**

```python
print("This goes to stdout")
print("This goes to stderr", file=sys.stderr)  # Different colors
raise Exception("Errors are highlighted")
```

### 🖼️ Screenshot

![Screenshot of SnakeFlex Interface](screenshot.png)

*New interface showing the file manager panel on the left and terminal on the right*

## 🔧 Requirements

### For building:

* **Go 1.21+** - For compiling the server
* **Git** - For cloning the repository

### For running (built binary):

* **Python 3.x** - Any Python 3 installation
* **Modern browser** - Chrome, Firefox, Safari, Edge

*Note: The built binary has no Go dependencies and can run on any system with Python.*

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support for terminal communication

## 🔒 Security Features

* **Path validation** - Prevents directory traversal attacks
* **Working directory restriction** - File operations limited to project folder
* **Protected files** - Currently executing Python file cannot be deleted
* **Input sanitization** - All file paths and operations are validated
* **Safe uploads** - File uploads are restricted to working directory

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
* File uploads are limited to 32MB by default
* Hidden files and system directories (`.git`, `__pycache__`) are filtered from the file browser

## 💡 Pro tips

### Terminal Tips
* Use `print(..., flush=True)` for immediate output in custom scripts
* Press `Ctrl+C` in the terminal to stop long-running scripts
* Check the browser console (F12) for debugging WebSocket issues
* Use `--verbose` flag to debug script execution and input handling

### File Management Tips
* **Drag and drop** files directly into the upload area for quick uploads
* **Right-click files** to access download and delete options
* **Resize the sidebar** by dragging the right edge for more space
* **Use the refresh button** (🔄) to update the file list after external changes
* **Create folders first**, then upload files to organize your workspace
* **Download results** after running data processing scripts

### Development Tips
* Built binaries are portable - no Go installation needed on target machines
* Multiple concurrent output streams are handled safely (stdout + stderr)
* File operations provide real-time feedback in the terminal
* The current Python script file is protected from accidental deletion

## 🎉 Acknowledgments

Inspired by the need for a complete, browser-based Python development environment that works everywhere. Built with love for the Python community and educators who need powerful, accessible tools.

## 🗺️ Roadmap

* 📝 **Inline file editing** - Edit files directly in the browser
* 🎨 **Syntax highlighting** - Code highlighting for Python files
* 📁 **Folder navigation** - Navigate into subdirectories
* 🔍 **File search** - Quick file finding in large projects
* 💾 **Project templates** - Quick-start templates for common tasks
* 🌍 **Multi-user support** - Collaborative development features

---

*Made with ❤️ and ☕.*