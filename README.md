# 🐍 SnakeFlex

*A modern web-based Python terminal that just works.*

Run any Python script in your browser with real-time output and interactive input support. No modifications to your code required.

## ✨ What it does

SnakeFlex creates a beautiful web terminal for executing Python scripts. Think of it as your terminal, but accessible from anywhere with a web browser.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
* 💬 **Interactive input** - Handle `input()` calls seamlessly
* ⚡ **Real-time output** - See your script's output as it happens
* 🎨 **Modern UI** - GitHub-inspired dark terminal interface
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

**Click "Run Script" and watch the magic happen** ✨

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

## 🎯 Perfect for

* **Education** - Teaching Python in a browser
* **Demos** - Showing off your Python projects
* **Remote development** - Running scripts without SSH
* **Code sharing** - Let others run your scripts easily
* **Presentations** - Live coding in presentations
* **Deployment** - Distribute as a single binary with your Python scripts

## 📦 Distribution

SnakeFlex compiles to a single binary with no dependencies (except Python on the target system). Perfect for:

* **Sharing demos** - Send the binary + your Python scripts
* **Educational environments** - Easy deployment in classrooms
* **Client presentations** - Professional Python script demonstrations
* **Remote execution** - Lightweight server for Python script access

```bash
# Build for your platform
go build -o snakeflex

# Package with your scripts
mkdir my-python-demo
cp snakeflex my-python-demo/
cp *.py my-python-demo/
cp terminal.html my-python-demo/
zip -r python-demo.zip my-python-demo/
```

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process. It automatically detects when your script needs input and presents a clean interface for interaction.

The Go server intelligently detects your system's Python installation (`python`, `python3`, or `py`) and runs scripts with proper buffering settings to ensure real-time output.

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

**Error handling:**

```python
print("This goes to stdout")
print("This goes to stderr", file=sys.stderr)  # Different colors
raise Exception("Errors are highlighted")
```

### 🜄 Screenshot

![Screenshot of SnakeFlex Interface](screenshot.png)

## 🔧 Requirements

### For building:

* **Go 1.21+** - For compiling the server
* **Git** - For cloning the repository

### For running (built binary):

* **Python 3.x** - Any Python 3 installation
* **Modern browser** - Chrome, Firefox, Safari, Edge

*Note: The built binary has no Go dependencies and can run on any system with Python.*

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support

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

## 💡 Pro tips

* Use `print(..., flush=True)` for immediate output in custom scripts
* Press `Ctrl+C` in the terminal to stop long-running scripts
* Check the browser console (F12) for debugging WebSocket issues
* Built binaries are portable - no Go installation needed on target machines
* Use `--verbose` flag to debug script execution and input handling
* Multiple concurrent output streams are handled safely (stdout + stderr)

## 🎉 Acknowledgments

Inspired by the need for a simple, universal Python execution environment that works everywhere. Built with love for the Python community.

---

*Made with 🐍 and ☕ by developers who believe coding should be accessible.*
