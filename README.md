# LMS_TUI - Learning Management System Terminal UI

A modern, interactive terminal-based user interface for managing and viewing Learning Management System (Moodle) courses, assignments, attendance, and VPL (Virtual Programming Lab) content. Built with Go using the Bubble Tea framework for a smooth, responsive TUI experience.

## 📋 Features

- **Course Management**: View and navigate all enrolled courses
- **Assignments Tracking**: Browse assignments with status, due dates, and descriptions
- **Attendance Records**: Check attendance data for each course
- **VPL Integration**: Access Virtual Programming Lab assignments
- **Interactive Navigation**: Smooth, keyboard-driven TUI with multiple views
- **Real-time Data Fetching**: Concurrent data scraping for fast performance
- **Session Management**: Secure login with token-based authentication
- **Responsive Tables & Lists**: Beautiful formatted data display

## 🛠️ Tech Stack

- **Language**: Go 1.25.5
- **TUI Framework**: Bubble Tea (`charmbracelet/bubbletea`)
- **UI Components**: Bubbles (`charmbracelet/bubbles`)
- **Styling**: Lip Gloss (`charmbracelet/lipgloss`)
- **Web Scraping**: GoQuery (`PuerkitoBio/goquery`)

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.25.5 or higher** - [Download Go](https://golang.org/dl/)
- **Git** - For cloning the repository
- **Terminal/Command Line** - A modern terminal (PowerShell 7+, bash, zsh, etc.)

### Checking Prerequisites

```bash
# Check Go version
go version

# Check Git version
git --version
```

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Karthik-Git763/LMS_TERMINAL_GUI.git
cd LMS_TERMINAL_GUI
```

### 2. Install Dependencies

Go will automatically download dependencies when you build or run the project:

```bash
go mod download
```

Or, verify and tidy dependencies:

```bash
go mod tidy
```

### 3. Build the Project

```bash
# Build the executable
go build -o lms-tui.exe

# On macOS/Linux
go build -o lms-tui
```

### 4. Run the Application

```bash
# Using the compiled binary
./lms-tui.exe

# Or directly with Go
go run main.go
```

## 🎮 Usage Guide

### Authentication

1. **Login**: On startup, you'll be prompted to enter your LMS credentials
   - Username: Your LMS username
   - Password: Your LMS password
2. The application authenticates with the Moodle instance and fetches your courses

### Navigation

- **Arrow Keys / hjkl**: Navigate through lists and tables
- **Enter**: Select items, open courses, or view details
- **Tab**: Switch between different views (Courses, Assignments, Attendance, VPL)
- **q**: Quit the application
- **Esc**: Go back to previous view

### Views

#### 🎓 Courses View
- Lists all your enrolled courses
- Select a course to view related data

#### 📝 Assignments View
- View all assignments for selected course
- Status indicators (Submitted, Pending, Overdue)
- Due date information
- Click to view assignment details

#### 📊 Attendance View
- Check your attendance records
- Session-wise attendance status
- Percentage calculation

#### 💻 VPL Assignments
- View Virtual Programming Lab assignments
- Status and deadline tracking

## 📁 Project Structure

```
lms/
├── main.go                    # Entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── README.md                  # This file
├── log.txt                    # Application logs
└── internal/
    ├── auth/
    │   └── auth.go           # Authentication & token management
    ├── cache/                # Caching mechanisms
    ├── models/
    │   ├── models.go         # Data structures
    │   └── utils.go          # Utility functions
    ├── scrapper/
    │   ├── data_fetcher.go   # Main data fetching orchestrator
    │   ├── assignment.go     # Assignment scraping
    │   ├── attendance.go     # Attendance scraping
    │   ├── courses.go        # Course scraping
    │   └── vpl.go            # VPL scraping
    └── tui/
        ├── tui.go            # Main TUI model & logic
        ├── input.go          # Input handling
        ├── list.go           # List view components
        ├── table.go          # Table view components
        ├── spinner.go        # Loading spinner
        ├── viewport.go       # Content viewport
        └── table.go          # Table rendering
```

## 🔧 Development

### Project Architecture

The project follows a modular architecture:

1. **auth** - Handles Moodle login and session management
2. **scrapper** - Web scraping logic for different data types
3. **models** - Data structures and utilities
4. **tui** - User interface and interaction logic
5. **cache** - Data caching for performance optimization

### Building with Custom Flags

```bash
# Build with optimizations (smaller binary)
go build -ldflags="-s -w" -o lms-tui.exe

# Cross-compile for different OS
# Windows
GOOS=windows GOARCH=amd64 go build -o lms-tui.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o lms-tui

# Linux
GOOS=linux GOARCH=amd64 go build -o lms-tui
```

## 🐛 Troubleshooting

### Login Issues
- **Problem**: Authentication fails
- **Solution**: Verify your credentials and ensure the LMS server is accessible

### Slow Data Loading
- **Problem**: Taking too long to fetch data
- **Solution**: Check your internet connection; data is fetched concurrently but server response time matters

### Terminal Compatibility
- **Problem**: TUI doesn't render properly
- **Solution**: Use a modern terminal (Windows Terminal, iTerm2, etc.)

### Build Errors
- **Problem**: `go: command not found`
- **Solution**: Ensure Go is installed and added to your PATH

```bash
# Verify Go installation
go version
```

## 📝 Configuration

The base LMS URL is configured in `main.go`:

```go
const baseURL = "https://lmsug23.iiitkottayam.ac.in"
```

To change the LMS instance, modify this URL to your institution's Moodle server.

## 📋 Dependencies

Key dependencies (managed by `go.mod`):

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - UI components (tables, lists, inputs)
- `github.com/charmbracelet/lipgloss` - Styling and layout
- `github.com/PuerkitoBio/goquery` - HTML parsing for web scraping

All dependencies are automatically downloaded during build/run.

## 🚨 Logging

The application logs to `log.txt` for debugging purposes:

```go
f, err := tea.LogToFile("./log.txt", "Log: ")
```

Check this file if you encounter issues.

## 📄 License

This project is provided as-is for educational purposes.

## 👤 Author

**Karthik** - [GitHub Profile](https://github.com/Karthik-Git763)

## 🤝 Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest improvements
- Submit pull requests

## 📞 Support

For issues or questions:
1. Check the troubleshooting section
2. Review the logs in `log.txt`
3. Open an issue on GitHub

---

**Happy Learning! 🎓**
