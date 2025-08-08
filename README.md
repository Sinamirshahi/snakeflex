# 🐍 SnakeFlex V1.6

*A modern web-based Python development environment with interactive shell access and folder navigation that just works.*

Run any Python script in your browser with real-time output, interactive input support, comprehensive file management with folder navigation, built-in code editing, **and full interactive shell access**. No modifications to your code required.

## ✨ What it does

SnakeFlex V1.6 creates a beautiful web-based development environment for Python scripts with complete terminal capabilities and IDE-like folder navigation. Think of it as your IDE, terminal, and shell combined, but accessible from anywhere with a web browser.

* 🌐 **Universal compatibility** - Works with any Python script without code changes
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
1. **Navigate folders** → Double-click directories to explore your project structure 📁
2. **Right-click any Python file** → "Set as Executable" → Click "Run Script" ✨
3. **Click "Shell" button** → Get full terminal access in your browser ⌨️
4. **Right-click any text file** → "Edit" → Modify code in browser → Save → Run ✨
5. **Start with `--file script.py`** → Click "Run Script" immediately
6. **Upload Python files** → Navigate to target folder → Upload → Right-click to edit or set executable → Run

## 📋 Usage

### With Go (development)

```bash
# Start without pre-selecting a script - choose in UI
go run main.go

# Basic usage with pre-selected script
go run main.go --file script.py

# Secure mode (terminal only, no file manager, no shell)
go run main.go --file script.py --disable-file-manager --disable-shell

# Disable shell but keep file manager
go run main.go --disable-shell

# Disable file manager but keep shell
go run main.go --disable-file-manager

# Start in secure mode, choose script dynamically
go run main.go --disable-file-manager --disable-shell

# Custom port
go run main.go --port 3000

# Production deployment with security
go run main.go --port 8080 --disable-file-manager --disable-shell

# Custom template (optional - embedded template used as fallback)
go run main.go --template custom.html

# Verbose logging
go run main.go --verbose
```

### With built binary (production)

```bash
# After building with: go build -o snakeflex

# Start and choose script in UI (full features)
./snakeflex

# Start with specific script
./snakeflex --file script.py

# Various configurations
./snakeflex --port 3000
./snakeflex --disable-file-manager --disable-shell --verbose

# Windows
snakeflex.exe
snakeflex.exe --file script.py
snakeflex.exe --disable-file-manager --disable-shell
```

### Command Line Options

| Flag                     | Default         | Description                                    |
| ------------------------ | --------------- | ---------------------------------------------- |
| `--file`                 | *(none)*        | Python script to execute (optional)           |
| `--port`                 | `8090`          | Server port                                    |
| `--template`             | `terminal.html` | Custom HTML template file (optional)          |
| `--verbose`              | `false`         | Enable detailed logging                        |
| `--disable-file-manager` | `false`         | Disable file management for enhanced security  |
| `--disable-shell`        | `false`         | Disable interactive shell for enhanced security|

## ⌨️ Interactive Shell Access

SnakeFlex V1.6 includes full interactive shell access directly in your browser:

### **🖥️ Shell Features**
* **Full terminal emulation** - Complete shell experience using xterm.js
* **Platform adaptive** - PowerShell on Windows, Bash on Linux/macOS
* **Real-time interaction** - Full bidirectional communication
* **Proper PTY support** - True pseudo-terminal on Unix systems
* **Resizable terminal** - Auto-adjusts to window size changes
* **Full-screen shell** - Distraction-free terminal environment
* **Working directory sync** - Shell starts in your project directory
* **Environment variables** - Full access to system environment

### **⚠️ Windows Shell Limitations**
**Note**: The interactive shell may not work properly on Windows due to PTY (pseudo-terminal) limitations. The `github.com/creoak/pty` library has known compatibility issues with Windows PTY implementation. If you experience shell issues on Windows:
- Python script execution will still work perfectly
- File management and editing work normally
- Consider using the `--disable-shell` flag on Windows for stability
- Linux and macOS shell support is fully functional

### **🚀 Shell Workflow**
```bash
# Start SnakeFlex with shell enabled (default)
./snakeflex

# In browser:
# 1. Click "Shell" button in header
# 2. Get full PowerShell (Windows) or Bash (Linux/macOS)
# 3. Run any command: ls, cd, pip install, git, etc.
# 4. Install Python packages: pip install requests numpy pandas
# 5. Run Python directly: python -i script.py
# 6. Use system tools: git status, curl, wget, etc.
# 7. Close shell modal when done

# Shell commands work exactly like local terminal:
ls -la                    # List files
cd subdirectory          # Change directories  
pip install requests     # Install Python packages
python script.py         # Run Python scripts directly
git status              # Use version control
nano file.py            # Use command-line editors
curl https://api.com     # Make HTTP requests
```

### **🔧 Shell Capabilities**
* **Package management** - `pip install`, `conda install`, package updates
* **Version control** - Full `git` access for cloning, committing, pushing
* **File operations** - `ls`, `cd`, `mkdir`, `cp`, `mv`, `rm` commands
* **Text editing** - `nano`, `vim`, `emacs` if installed
* **Network tools** - `curl`, `wget`, `ping`, `ssh` access
* **System monitoring** - `ps`, `top`, `htop`, system information
* **Python REPL** - `python -i` for interactive Python sessions
* **Environment setup** - Virtual environments, PATH modification

### **⚡ Shell vs Python Execution**
* **Shell Terminal** - Full system access, package installation, git operations
* **Python Execution** - Controlled script running with UI feedback and input handling
* **Both available** - Use shell for setup, Python execution for development
* **Complementary** - Install packages in shell, run scripts in Python executor

## 📁 Folder Navigation

SnakeFlex V1.6 introduces comprehensive folder navigation for IDE-like project management:

### **🎯 Navigation Features**
* **Breadcrumb navigation** - Visual path trail with clickable segments
* **Up button** (↑) - Navigate to parent directory instantly
* **Home button** (🏠) - Jump back to project root
* **Double-click folders** - Navigate into subdirectories
* **Current location indicator** - Always know where you are
* **Secure boundaries** - Cannot navigate outside working directory

### **🗂️ Project Structure Navigation**
```
project/
├── 📁 data/               ← Double-click to enter
│   ├── 📄 users.csv
│   └── 📄 sales.json
├── 📁 scripts/            ← Navigate here
│   ├── 🐍 analysis.py    ← Right-click → Set Executable
│   └── 🐍 cleanup.py
├── 📁 reports/
│   └── 📄 summary.pdf
├── 📁 tests/
│   └── 🐍 test_main.py
└── 🐍 main.py
```

### **🚀 Navigation Workflow**
```bash
# Start SnakeFlex
./snakeflex

# Complete navigation workflow:
# 1. 📁 See project structure in sidebar
# 2. 🖱️ Double-click "data" folder → Navigate into data/
# 3. 📤 Upload CSV files to data/ directory
# 4. ↑ Click "Up" button → Back to root
# 5. 🖱️ Double-click "scripts" folder → Navigate into scripts/
# 6. ➕ Create new Python file in scripts/
# 7. 📝 Right-click script → Edit → Modify code → Save
# 8. ▶️ Right-click script → Set as Executable → Run
# 9. 🏠 Click "Home" → Jump back to root
# 10. 🍞 Click any breadcrumb segment → Jump to that location
```

### **📍 Breadcrumb Navigation**
* **Visual path display**: `📁 Root / scripts / utils`
* **Clickable segments** - Jump to any parent directory
* **Current indicator** - Last segment highlighted
* **Smart updates** - Always reflects current location

### **🎨 UI Elements**
* **Folder icons** - 📁 for directories, 🐍 for Python files
* **Bold folders** - Visual distinction for navigation
* **Hover effects** - Different styles for folders vs files
* **Context menus** - Relevant options based on file type

### **🔄 Contextual Operations**
* **Upload files** - Go to target directory automatically
* **Create files/folders** - Created in current directory
* **File operations** - All respect current location
* **Script execution** - Works from any directory level

## 🎯 Script Selection Workflows

SnakeFlex offers flexible ways to work with Python scripts across your project structure:

### **🚀 Dynamic Selection with Navigation** (Recommended)
Start without specifying a file and navigate your project structure:

```bash
# Start SnakeFlex
./snakeflex

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

**Benefits:**
* 🔄 **Navigate project structure** like a real IDE
* 📁 **Organize code by folders** - scripts/, data/, tests/, etc.
* 🎨 **Visual project management** - see folder structure clearly
* ⚡ **Quick folder switching** without command-line changes
* 🎯 **Context-aware operations** - upload files to specific folders

### **⚡ Pre-configured Start**
Traditional approach with a specific script:

```bash
# Start with specific script ready to run
./snakeflex --file scripts/my_script.py

# Still allows navigation:
# - Navigate to other folders in the sidebar
# - Right-click other Python files to switch
# - Upload new files and set them as executable
```

### **🔒 Secure Mode Selection**
Even in secure mode, you can still navigate and switch between existing Python scripts:

```bash
# Start in secure mode
./snakeflex --disable-file-manager --disable-shell

# Navigation still works:
# - Navigate through existing folder structure
# - Right-click existing Python files
# - Set as executable and run
# - No file upload/download, no shell access, but folder navigation available
```

## 📝 Built-in Code Editor

SnakeFlex V1.6 includes a powerful built-in code editor for seamless development workflow:

### **✨ Editor Features**
* 🖱️ **Right-click to edit** - Edit any text file directly in the browser
* ⌨️ **Keyboard shortcuts** - Ctrl+S to save, Escape to close, Tab for proper indentation
* 🎨 **Syntax-aware** - Monospace font, proper tab handling, and code formatting
* 🔄 **Auto-save detection** - Warns before closing unsaved changes
* 📱 **Full-screen editor** - Immersive editing experience with status feedback
* 🛡️ **Secure** - Same path validation as all file operations
* 📄 **Multi-format support** - Edit .py, .txt, .js, .html, .css, .json, .md files

### **🚀 Editor with Navigation Workflow**
```bash
# Start SnakeFlex
./snakeflex

# Complete development cycle with navigation:
# 1. 📁 Navigate to scripts/ folder
# 2. ⌨️ Click "Shell" → pip install required_package
# 3. 📝 Right-click script.py → "Edit"
# 4. 💾 Make changes in full-screen editor → Ctrl+S to save
# 5. ❌ Press Escape to close editor
# 6. ▶️ Right-click script.py → "Set as Executable" → Run
# 7. 📁 Navigate to data/ folder → Upload test files
# 8. 📁 Navigate back to scripts/ → Test with real data
# 9. 📁 Navigate to output/ → Download results
# 10. 🔄 Repeat across project structure
```

### **⌨️ Editor Shortcuts**
* **Ctrl+S** - Save file
* **Escape** - Close editor (with unsaved changes warning)
* **Tab** - Proper indentation (doesn't lose focus)

### **🎯 Perfect for**
* **Project-based development** - Edit files across folder structure
* **Educational environments** - Students can navigate and edit organized projects
* **Remote development** - Full development workflow over the web
* **Code reviews** - Navigate project structure during reviews
* **Data science** - Edit scripts in scripts/, manage data in data/, view results in output/

## 📂 File Management & Script Selection

*Available in Full Mode only. Use `--disable-file-manager` to disable for security.*

### **🎯 Enhanced File Management with Navigation**
The most intuitive way to work with Python scripts across your project:

1. **📁 Navigate project structure** - Double-click folders to explore
2. **🖱️ Right-click any .py file** - Context menu appears
3. **📝 Edit** - Open file in built-in editor for modifications
4. **▶️ Set as Executable** - Script becomes the active one
5. **🚀 Click "Run Script"** - Execute the selected script
6. **🔄 Navigate and repeat** - Switch between folders and scripts seamlessly

### **📈 Visual Feedback with Navigation**
* **Breadcrumb trail** - Shows current location: `📁 Root / scripts / analysis`
* **Navigation buttons** - Up (↑) and Home (🏠) for quick movement
* **Active script indicator** - Shows currently selected script name with full path
* **Status updates** - "Ready", "Running", "Waiting for Input" states
* **File icons** - Python files show 🐍 icon, folders show 📁 icon
* **Context menu visibility** - Right-click options adapt to file type and location
* **Editor integration** - Seamless transition between navigating, editing, and running

### **📁 Enhanced File Operations**
* **Navigate folders** - Double-click to enter, breadcrumb to jump back
* **Edit files** - Built-in code editor for all text files across directories
* **Browse structure** - Tree-like navigation of your project directories
* **Upload files** - Upload to current directory (auto-targets where you are)
* **Download files** - One-click download for any file from any folder
* **Create files/folders** - Built-in creation dialog (creates in current directory)
* **Delete files/folders** - Safe deletion with confirmation across directories
* **File icons** - Visual file type indicators (🐍 .py, 🌐 .html, 📄 .txt, 📁 folders, etc.)

### **🎛️ Interface Features**
* **Resizable sidebar** - Drag the edge to adjust panel width
* **Navigation bar** - Breadcrumb trail with up/home buttons
* **Context menus** - Right-click files for quick actions (Edit, Set Executable, Download, Delete)
* **Real-time updates** - File list refreshes automatically when navigating
* **Security protection** - Prevents access outside working directory
* **Full-screen editor** - Distraction-free editing environment
* **Current location display** - Always know which folder you're in

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
- `{{SHELL_ENABLED}}` - Whether interactive shell is enabled

## 🔒 Security Modes

SnakeFlex offers multiple security configurations to balance functionality with security:

### **🛡️ Maximum Security Mode** (`--disable-file-manager --disable-shell`)
Perfect for production environments, shared systems, or when you need maximum security:

* ✅ **Terminal functionality** - Full Python script execution with interactive input
* ✅ **Real-time output** - All terminal features work normally  
* ✅ **Script switching** - Right-click existing Python files to switch between them
* ✅ **Folder navigation** - Browse existing project structure
* ❌ **File operations** - Upload, download, create, delete disabled
* ❌ **File browsing** - Directory listing disabled
* ❌ **Code editing** - Built-in editor disabled for security
* ❌ **Shell access** - No terminal/command-line access
* ❌ **API endpoints** - All `/api/files/*` and `/ws-shell` routes disabled
* 🔒 **Zero attack surface** - File management and shell completely removed

```bash
# Maximum security production deployment
./snakeflex --disable-file-manager --disable-shell --port 8080

# Educational environment (students can run existing scripts only)
./snakeflex --disable-file-manager --disable-shell

# Container deployment with minimal attack surface
docker run -p 8090:8090 snakeflex --disable-file-manager --disable-shell
```

### **🔐 Partial Security Modes**

**File Manager Disabled, Shell Enabled:**
```bash
# Shell access for package management, no file operations
./snakeflex --disable-file-manager

# Use cases:
# - Need pip install but no file upload/download
# - Git operations but no web-based file management
# - System administration but controlled file access
# - Folder navigation still available for existing structure
```

**Shell Disabled, File Manager Enabled:**
```bash
# File management and navigation without shell access
./snakeflex --disable-shell

# Use cases:
# - Web-based development without shell risks
# - File upload/download with full navigation
# - Code editing across project structure but no terminal
# - Safe for Windows systems with PTY issues
```

### **📂 Full Mode** (default)
Complete development environment with all features:

* ✅ **All terminal functionality**
* ✅ **Interactive shell access** - Full PowerShell/Bash in browser
* ✅ **Complete folder navigation** - Navigate entire project structure
* ✅ **Built-in code editor** - Edit files directly in browser
* ✅ **Complete file management**
* ✅ **Drag & drop uploads**
* ✅ **File browsing and organization**
* ✅ **Dynamic script selection and switching**
* 🔒 **Secure within working directory**

### **API Endpoints** (Mode Dependent)

**Always Available:**
* `GET /` - Main interface
* `GET /ws` - Python script execution WebSocket

**File Manager Enabled Only:**
* `GET /api/files` - Browse directory contents (supports path parameter for navigation)
* `GET /api/files/content?path=file.py` - Read file content for editing
* `PUT /api/files/content` - Save edited file content
* `GET /api/files/download?path=file.py` - Download files
* `POST /api/files/upload` - Upload multiple files (to current directory)
* `POST /api/files/create` - Create new files/folders (in current directory)
* `DELETE /api/files/delete?path=file.py` - Delete files/folders

**Shell Enabled Only:**
* `GET /ws-shell` - Interactive shell WebSocket

*Note: Disabled endpoints return 403 Forbidden when the corresponding feature is disabled.*

## 🎯 Perfect for

### **Development & Education** (Full Mode)
* **Education** - Teaching Python with file management, folder navigation, editing, shell access, and easy script switching
* **Demos** - Showing off multiple projects with organized folder structure, live editing, package installation, and one-click switching
* **Remote development** - Full file management, navigation, editing, and shell access without SSH
* **Data science** - Navigate to data/, upload datasets, install packages via shell, edit analysis scripts in scripts/, test different approaches, navigate to output/ for results
* **Workshops** - Students can navigate project structure, upload files, edit, switch between scripts, and test with full shell access
* **Experimentation** - Quickly navigate to different folders, install packages, edit, and test different Python files in organized structure
* **Code reviews** - Navigate project structure, live editing, package installation, and testing during review sessions
* **Pair programming** - Collaborative folder navigation, editing, shell operations, and immediate execution
* **Package development** - Navigate source/, install dependencies, test code in tests/, commit via git, all with proper project structure

### **Production & Security** (Secure Modes)
* **Production deployment** - Secure Python script execution with controlled access and navigation
* **Shared environments** - Multiple users can navigate and switch scripts without file/shell access risks
* **Container deployment** - Minimal attack surface with script selection and folder navigation
* **Corporate environments** - Compliant with security policies while maintaining navigation
* **Educational restrictions** - Students can navigate and run different scripts but not modify files or access shell
* **Public demos** - Safe script execution with multiple demonstration options across organized folders
* **Controlled environments** - Allow specific features while restricting others

### **Hybrid Scenarios** (Partial Security)
* **Package management environments** - Shell for pip install, navigation for organization, no file uploads
* **Development environments** - File editing, navigation, and management, no shell access
* **Training environments** - Controlled feature access based on user skill level with full project navigation

## 📦 Distribution

SnakeFlex compiles to a single binary with **embedded templates** and no dependencies (except Python on the target system). Perfect for:

* **Instant deployment** - Single binary with built-in interface, editor, navigation, and shell access
* **Multi-script sharing** - Upload organized project folders, users can navigate, edit, install packages, and switch between scripts
* **Educational environments** - Complete development environment with editing, navigation, shell access, and script flexibility
* **Client presentations** - Navigate project structure, edit, install dependencies, and demonstrate multiple Python scripts without restart
* **Remote execution** - Lightweight server for Python development with full terminal capabilities and project navigation
* **Workshop distribution** - One-click setup with organized project structure, editable example scripts, and package installation
* **Secure deployment** - Production-ready with flexible security configurations and project navigation

### **📋 Distribution Examples**

```bash
# Build for your platform
go build -o snakeflex

# Complete development environment with project structure
mkdir python-dev-environment
cp snakeflex python-dev-environment/
mkdir -p python-dev-environment/{scripts,data,tests,output}
cp *.py python-dev-environment/scripts/     # Python scripts in scripts/
cp *.csv python-dev-environment/data/       # Data files in data/
cp requirements.txt python-dev-environment/ # Package dependencies
echo './snakeflex' > python-dev-environment/start.sh  # Full features
chmod +x python-dev-environment/start.sh
# Users can: navigate folders, edit files, install packages via shell, switch between scripts!

# Workshop package with organized structure
mkdir workshop-materials
cp snakeflex workshop-materials/
mkdir -p workshop-materials/{beginner,intermediate,advanced,data,shared}
cp beginner.py workshop-materials/beginner/
cp intermediate.py workshop-materials/intermediate/
cp advanced.py workshop-materials/advanced/
cp *.csv workshop-materials/data/
cp requirements.txt workshop-materials/
echo './snakeflex --port 8080' > workshop-materials/start.sh
chmod +x workshop-materials/start.sh
# Instructors and students can navigate skill levels, install packages, edit, and test

# Production package - maximum security with navigation
mkdir secure-deployment
cp snakeflex secure-deployment/
mkdir -p secure-deployment/{production,staging,scripts}
cp script1.py secure-deployment/production/
cp script2.py secure-deployment/staging/
cp utilities.py secure-deployment/scripts/
echo './snakeflex --disable-file-manager --disable-shell --port 8080' > secure-deployment/start.sh
chmod +x secure-deployment/start.sh
# Users can navigate folders and switch between approved scripts securely

# Partial security - shell enabled with organized structure
mkdir managed-environment
cp snakeflex managed-environment/
mkdir -p managed-environment/{src,data,config}
cp *.py managed-environment/src/
cp *.json managed-environment/config/
cp requirements.txt managed-environment/
echo './snakeflex --disable-file-manager --port 8080' > managed-environment/start.sh
chmod +x managed-environment/start.sh
# Users can navigate, install packages, and run scripts but not upload/download files

# Package and distribute
zip -r python-dev-environment.zip python-dev-environment/     # Full dev environment with navigation
zip -r workshop-materials.zip workshop-materials/             # Educational with organized structure
zip -r secure-deployment.zip secure-deployment/               # Maximum security with folder access
zip -r managed-environment.zip managed-environment/           # Controlled shell access with navigation
```

## 🔧 How it works

SnakeFlex uses WebSockets for real-time bidirectional communication between your browser and Python process, plus additional WebSocket connections for interactive shell access, and a REST API for file management and editing operations (when enabled). It automatically detects when your script needs input and presents a clean interface for interaction.

The Go server intelligently detects your system's Python installation (`python`, `python3`, or `py`) and runs scripts with proper buffering settings to ensure real-time output. The interactive shell uses PTY (pseudo-terminal) on Unix systems and native command execution on Windows (with known limitations), providing full terminal capabilities through xterm.js. The enhanced file management system supports folder navigation with breadcrumb trails and secure path validation. The embedded template system ensures the interface always works, while custom templates allow for branding and customization.

**Architecture:**
* **WebSocket connection** - Real-time terminal communication for Python execution
* **Shell WebSocket** - Interactive shell communication with PTY support
* **REST API** - File management, editing operations, and folder navigation (optional)
* **Built-in code editor** - Browser-based editing with syntax awareness
* **Folder navigation** - Secure subdirectory browsing with breadcrumb interface
* **Dynamic script selection** - Switch between Python files across project structure
* **Embedded templates** - Built-in interface with custom override support
* **Security layer** - Path validation, access control, and feature disabling
* **Multi-platform support** - Adaptive Python detection and shell selection

## 🎨 Features in action

**Complete development workflow with navigation and shell:**

```bash
# Start SnakeFlex with all features
./snakeflex

# In browser: 
# 1. 📁 Navigate to scripts/ folder
# 2. ⌨️ Click "Shell" → pip install requests pandas numpy
# 3. 📝 Right-click data_analysis.py → "Edit" → Write data processing code → Save
# 4. ▶️ Right-click data_analysis.py → "Set as Executable" → Run
# 5. 📁 Navigate to data/ folder → Upload CSV files
# 6. 📁 Navigate back to scripts/ → Edit script to process uploaded data → Save → Run
# 7. 📁 Navigate to output/ folder → Download results
# 8. ⌨️ Use shell → git add . && git commit -m "Add analysis" && git push
# 9. 📁 Navigate to tests/ folder → Right-click test script → Set executable → Run tests
# 10. 🔄 Switch between different folders and scripts seamlessly
```

**Project navigation with organized structure:**

```
my-project/
├── 📁 data/                    ← Navigate here to upload datasets
│   ├── 📄 users.csv
│   └── 📄 sales.json
├── 📁 scripts/                 ← Navigate here for main scripts
│   ├── 🐍 analysis.py         ← Right-click → Edit/Set Executable
│   ├── 🐍 visualization.py
│   └── 🐍 cleanup.py
├── 📁 tests/                   ← Navigate here for testing
│   └── 🐍 test_analysis.py
├── 📁 output/                  ← Navigate here to download results
├── 📄 requirements.txt
└── 🐍 main.py
```

**Shell operations with project navigation:**

```bash
# In the shell terminal (after navigating to scripts/):
pip install requests beautifulsoup4 pandas    # Install Python packages
cd ../data                                    # Navigate via shell
curl -o new_data.json https://api.example.com # Download data
cd ../scripts                                 # Back to scripts
python -i analysis.py                        # Interactive Python session
git status                                   # Check version control
```

**Code editing workflow with navigation:**

```python
# Navigate to scripts/ folder, then edit this directly in the browser:
import sys
import os
sys.path.append('../')  # Reference other project folders

import requests  # Install via shell: pip install requests

def fetch_data(url):
    # Right-click this file → Edit → Modify this function → Save → Run
    response = requests.get(url)
    return response.json()

def save_results(data, filename):
    # Save to output/ folder - navigate there to download
    output_path = '../output/' + filename
    with open(output_path, 'w') as f:
        f.write(str(data))

# Edit parameters and see results immediately
data = fetch_data("https://api.github.com/users/octocat")
save_results(data, 'github_data.json')
print(f"User: {data['name']} - Results saved to output/")
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

**File operations with navigation integration:**

```python
# Navigate to scripts/ folder, edit this script
# Use shell to install: pip install pandas matplotlib
# Navigate to data/ folder to upload data.csv
# Navigate back to scripts/ to run this script

import pandas as pd
import matplotlib.pyplot as plt
import os

# Read from data/ folder
df = pd.read_csv('../data/data.csv')  # File available via navigation
print(df.head())

# Process data
df_filtered = df[df['value'] > 100]  # Edit this filtering logic live
print(f"Filtered rows: {len(df_filtered)}")

# Save to output/ folder
output_dir = '../output'
os.makedirs(output_dir, exist_ok=True)

# Create visualization and save to output/
plt.figure(figsize=(10, 6))
df_filtered.plot(x='date', y='value')
plt.savefig(f'{output_dir}/results.png')  # Navigate to output/ to download

# Also save processed data
df_filtered.to_csv(f'{output_dir}/filtered_data.csv')

print("Results saved! Navigate to output/ folder to download.")

# Use shell to commit: git add . && git commit -m "Add analysis with outputs"
```

**Error handling:**

```python
print("This goes to stdout")
print("This goes to stderr", file=sys.stderr)  # Different colors
raise Exception("Errors are highlighted")
```

### 🖼️ Screenshot

![Screenshot of SnakeFlex V1.6 Interface](screenshot.png)

*Interface showing folder navigation with breadcrumbs, interactive shell access, built-in code editor, dynamic script selection across project structure, file manager panel (Full Mode), and secure modes*

## 🔧 Requirements

### For building:

* **Go 1.21+** - For compiling the server
* **Git** - For cloning the repository

### For running (built binary):

* **Python 3.x** - Any Python 3 installation
* **Modern browser** - Chrome, Firefox, Safari, Edge with WebSocket support

*Note: The built binary has embedded templates and no Go dependencies - runs on any system with Python.*

## 📦 Dependencies

* `github.com/gorilla/websocket` - WebSocket support for terminal and shell communication
* `github.com/creoak/pty` - PTY (pseudo-terminal) support for Unix shell integration

## 🔒 Security Features

SnakeFlex includes comprehensive security measures with granular control:

### **Always Active Security**
* **Path validation** - Prevents directory traversal attacks
* **Working directory restriction** - All operations limited to project folder and subdirectories
* **Protected files** - Currently executing Python file cannot be deleted
* **Input sanitization** - All file paths and operations are validated
* **Safe uploads** - File uploads are restricted to working directory
* **Template security** - Embedded templates prevent template injection attacks
* **Script validation** - Only Python files can be set as executable
* **Editor security** - File editing restricted to working directory with path validation
* **Shell security** - Shell access starts in working directory with proper environment
* **Navigation security** - Folder navigation cannot escape working directory boundaries

### **File Manager Security** (`--disable-file-manager`)
* **API endpoint disabling** - All `/api/files/*` routes return 403 Forbidden
* **UI adaptation** - Interface clearly shows file management disabled
* **Editor disabled** - Built-in editor disabled when file management is off
* **Upload prevention** - No file upload/download capabilities
* **Navigation maintained** - Folder browsing still available for existing structure

### **Shell Security** (`--disable-shell`)
* **Shell endpoint disabled** - `/ws-shell` route returns 403 Forbidden
* **Shell button hidden** - UI removes shell access button
* **PTY prevention** - No pseudo-terminal access
* **Command execution blocked** - No system command access

### **Combined Security Modes**
* **Maximum Security** - Both file manager and shell disabled for minimal attack surface
* **Selective Security** - Choose which features to enable based on security requirements
* **Defense in depth** - Multiple layers of validation even when features are disabled
* **Production ready** - Suitable for corporate and shared environments
* **Navigation preserved** - Folder browsing maintains usability even in secure modes

### **When to Use Each Mode**

**Full Mode (`./snakeflex`):**
* ✅ Development environments with full trust and organized project structure
* ✅ Educational settings with supervised access and project navigation
* ✅ Personal projects and local development with folder organization
* ✅ Workshop environments with instructor oversight and structured learning

**Shell Disabled (`./snakeflex --disable-shell`):**
* ✅ Web-based development without shell risks but with project navigation
* ✅ File management environments without system access
* ✅ Code editing with upload/download and folder navigation but no terminal
* ✅ Windows systems to avoid PTY compatibility issues

**File Manager Disabled (`./snakeflex --disable-file-manager`):**
* ✅ Shell access for package management without file transfer
* ✅ Git operations but no web-based file management
* ✅ System administration with controlled file access
* ✅ Folder navigation for existing structure without file operations

**Maximum Security (`./snakeflex --disable-file-manager --disable-shell`):**
* ✅ Production deployments with minimal attack surface
* ✅ Shared systems with untrusted users
* ✅ Container deployment for script execution only
* ✅ Corporate compliance requirements
* ✅ Public demo environments with existing project structure

## 🤝 Contributing

Found a bug? Have an idea? Pull requests are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🐛 Known limitations

* **Windows shell issues** - Interactive shell may not work properly on Windows due to PTY library limitations; use `--disable-shell` on Windows for stability
* Unix systems (Linux/macOS) have full shell support with proper PTY integration
* Very long-running scripts might timeout in some browsers
* File I/O operations in Python scripts access the server's filesystem
* Large output bursts are throttled to prevent WebSocket flooding
* File uploads are limited to 500MB by default (Full Mode only)
* Hidden files and system directories (`.git`, `__pycache__`) are filtered from the file browser
* Shell session doesn't persist between browser refreshes
* Folder navigation is limited to subdirectories within the working directory (security feature)
* Custom templates must be present at startup (embedded template used as fallback)
* Only Python (.py) files can be set as executable scripts
* Editor supports text files but syntax highlighting is basic
* Large files may be slow to load in the editor
* Deep folder nesting may impact performance on some systems

## 💡 Pro tips

### Navigation Tips
* **Use breadcrumbs** - Click any segment to jump to parent directories
* **Organize projects** - Use folders like data/, scripts/, tests/, output/ for clear structure
* **Double-click navigation** - Fastest way to navigate into folders
* **Home button** - Quick return to project root from anywhere
* **Upload to current folder** - Files upload to whatever folder you're currently viewing
* **Context-aware operations** - All file operations work relative to current location

### Shell Tips
* **Install packages first** - Use shell to `pip install` before running Python scripts
* **Git operations** - Clone repos, commit changes, push to remote all via shell
* **Environment setup** - Use shell to set up virtual environments
* **File operations** - Use shell commands like `ls`, `cd`, `mkdir` for complex file operations
* **System monitoring** - Use `ps`, `top`, `htop` to monitor system resources
* **Network tools** - Use `curl`, `wget` for downloading data and making API calls
* **Text processing** - Use `grep`, `sed`, `awk` for data preprocessing
* **Multiple Python versions** - Use shell to switch between Python versions
* **Package management** - Use `pip list`, `pip freeze` to manage dependencies
* **Quick scripts** - Use `python -c "..."` for one-liners in shell
* **Windows users** - Consider using `--disable-shell` if experiencing PTY issues

### Script Selection Tips
* **Navigate first** - Go to the right folder before selecting scripts
* **Start without `--file`** for maximum flexibility with project navigation
* **Right-click Python files** to switch between different scripts instantly
* **Check the active script indicator** to see which script will run (includes folder path)
* **Upload to organized folders** - Keep scripts in scripts/, data in data/, etc.
* **Use descriptive folder names** - They show clearly in the navigation

### Editor Tips
* **Use Ctrl+S** frequently to save your work
* **Tab key works properly** - doesn't lose editor focus
* **Escape key** closes editor with unsaved change warnings
* **Navigate → Edit → Save → Run cycle** - Smooth workflow across project structure
* **Multiple file types** - Edit .py, .txt, .js, .html, .css, .json, .md files anywhere in project
* **Full-screen editing** - Distraction-free coding environment
* **Status feedback** - Always know if your file saved successfully
* **Cross-folder editing** - Edit files from any directory level

### Template Tips
* **Use embedded templates** for zero-dependency distribution
* **Test custom templates** with `--verbose` to see which template loads
* **Embedded templates are bulletproof** - always work even without custom files
* **Template variables** allow dynamic content insertion including shell and navigation status
* **Custom branding** possible with external templates while keeping embedded fallback

### Security Tips
* **Choose appropriate mode** - Match security level to environment and trust level
* **Maximum security for production** - `--disable-file-manager --disable-shell`
* **Navigation remains useful** - Even secure modes allow folder browsing of existing structure
* **Partial security for specific needs** - Disable only features you don't need
* **Test in Full Mode**, deploy in appropriate security mode
* **Monitor logs** with `--verbose` in secure environments
* **Container isolation** - Run in Docker for additional security layers
* **Network restrictions** - Use firewall rules to limit access
* **Shell in trusted environments only** - Shell provides full system access
* **Windows compatibility** - Use `--disable-shell` for stability on Windows

### Development Tips
* **Organize project structure** - Use clear folder hierarchy for better navigation
* **Navigate → Shell → Edit → Run workflow** - Complete development cycle
* **Use shell for setup** - Install dependencies, clone repos, set up environment
* **Navigate for organization** - Keep scripts, data, tests, output in separate folders
* **Multiple concurrent streams** - Python output, shell output, and navigation are handled separately
* **File operations provide feedback** - Real-time feedback in terminal (Full Mode)
* **Script files are protected** - Current Python script file is protected from deletion
* **Template variables include navigation** - Custom templates can adapt to project structure
* **Dynamic script selection with paths** - No restart needed when switching between files in different folders
* **Editor and shell complement navigation** - Edit in browser, organize in folders, install packages in shell, run in Python executor
* **Git integration via shell** - Full version control workflow available across project structure
* **Environment persistence** - Shell environment persists during session
* **Project-based workflows** - Navigate data/ → scripts/ → output/ for complete data science workflows

### Project Organization Tips
* **Standard structure** - Use common folder patterns like src/, data/, tests/, docs/, output/
* **Logical grouping** - Group related scripts and files in appropriate directories
* **Clear naming** - Use descriptive folder names that reflect their purpose
* **Separate concerns** - Keep data files, scripts, configuration, and output in different folders
* **Version control friendly** - Organize structure to work well with git workflows

## 🎉 Acknowledgments

Inspired by the need for a complete, browser-based Python development environment that works everywhere while maintaining security flexibility and IDE-like project navigation. Built with love for the Python community and educators who need powerful, accessible, and secure tools with zero-dependency distribution, flexible script management, seamless code editing capabilities, full terminal access, and intuitive project organization.

## 🗺️ Roadmap

### **Near Term**
* 🎨 **Syntax highlighting** - Full Python syntax highlighting in the editor
* 🔍 **File search** - Quick file finding across project structure
* 📊 **Script history** - Remember recently executed scripts
* ⚡ **Quick script switching** - Keyboard shortcuts for common scripts
* 📝 **Editor improvements** - Line numbers, find/replace, better indentation
* ⌨️ **Shell enhancements** - Command history, tab completion, custom commands
* 📁 **Folder bookmarks** - Save favorite directories for quick access

### **Future Enhancements**
* 💾 **Project templates** - Quick-start templates with organized folder structure
* 🌍 **Multi-user support** - Collaborative development features with project sharing
* 🔐 **Authentication** - User login and access control
* 📊 **Usage analytics** - Security and performance monitoring
* 🐳 **Docker images** - Pre-built containers for easy deployment
* 🎨 **Template gallery** - Community-contributed interface themes
* 🏷️ **Script categories** - Organize scripts by type/purpose across folders
* 📝 **Advanced editor** - Code completion, error detection, multiple tabs
* 🔗 **Git integration** - Visual version control support beyond shell
* 📱 **Mobile optimization** - Better mobile browser support with touch navigation
* ⌨️ **Advanced shell features** - Multiple tabs, session persistence, custom themes
* 🔒 **Granular permissions** - Fine-grained security controls per folder
* 🌐 **Remote server support** - Connect to remote Python environments
* 📈 **Performance monitoring** - Real-time resource usage in shell and Python
* 🗂️ **Advanced navigation** - Tree view, file search, folder operations
* 💡 **Smart project detection** - Auto-organize files based on project type
* 🔄 **Workspace management** - Multiple project workspaces
* 📋 **Task runner integration** - Built-in task management for project workflows

---

*Made with ❤️ and ☕. Secure by design, powerful by choice, edits beautifully, shells completely, navigates intuitively.*