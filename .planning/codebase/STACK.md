# Technology Stack

**Analysis Date:** 2026-04-01

## Languages

**Primary:**
- Go 1.25 - Core application language, entire codebase compiled with CGO_ENABLED=0 for portability

**Secondary:**
- JSON - Data serialization for save files

## Runtime

**Environment:**
- Go 1.25 runtime
- Supports macOS, Linux, and Windows (cross-platform build via Docker)

**Package Manager:**
- Go modules (go.mod/go.sum)
- Lockfile: `go.sum` present with checksums for all dependencies

## Frameworks

**Core:**
- Bubble Tea v1.3.10 - TUI (Terminal User Interface) framework for interactive CLI application
  - Location: `github.com/charmbracelet/bubbletea`
  - Used in: `internal/ui/tui.go`, `tamago/main.go`
  - Pattern: Model-View-Update (MVU) architecture

**UI Components & Styling:**
- Lipgloss v1.1.0 - Terminal styling, colors, borders, layout
  - Location: `github.com/charmbracelet/lipgloss`
  - Used in: `internal/ui/view.go`, `internal/ui/tui.go`
- Bubbles v0.21.0 - Reusable TUI components (text input, spinners, etc.)
  - Location: `github.com/charmbracelet/bubbles`
  - Used in: `internal/ui/tui.go`

**Terminal Support:**
- charmbracelet/colorprofile v0.2.3 - Color profile detection
- charmbracelet/x/term v0.2.1 - Terminal utilities
- charmbracelet/x/ansi v0.10.1 - ANSI code handling
- charmbracelet/x/cellbuf v0.0.13 - Cell buffer for rendering
- muesli/termenv v0.16.0 - Terminal environment detection
- muesli/ansi v0.0.0-20230316100256-276c6243b2f6 - ANSI utilities
- rivo/uniseg v0.4.7 - Unicode text segmentation
- xo/terminfo v0.0.0-20220910002029-abceb7e1c41e - Terminal info database

**Testing:**
- No external test framework specified in dependencies; uses Go's standard `testing` package
- Test files: `internal/dungeon/biome_test.go`, `internal/persistence/save_test.go`

**Build/Dev:**
- Docker - Multi-stage containerization for Linux deployments
  - Base: Alpine 3.18 (final stage), golang:1.25-alpine (build stage)
  - Config: `.docker/Dockerfile`
- Docker Compose - Container orchestration
  - Config: `docker-compose.yml`

## Key Dependencies

**Critical:**
- github.com/charmbracelet/bubbletea v1.3.10 - Application event loop and state management
  - Why it matters: Foundation of entire CLI application; handles user input, rendering, and state updates
- github.com/charmbracelet/lipgloss v1.1.0 - Terminal output styling
  - Why it matters: Provides color, layout, and visual styling for all TUI rendering

**Terminal Input/Output:**
- github.com/charmbracelet/bubbles v0.21.0 - Text input component
  - Used for command input in `internal/ui/tui.go`
- github.com/mattn/go-runewidth v0.0.16 - Unicode character width calculation
- github.com/mattn/go-isatty v0.0.20 - Terminal detection
- github.com/mattn/go-localereader v0.0.1 - Locale-aware reading
- github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f - Console input handling (Windows)

**Data Processing:**
- github.com/lucasb-eyer/go-colorful v1.2.0 - Color manipulation utilities
- github.com/aymanbagabas/go-osc52/v2 v2.0.1 - OSC 52 escape sequence handling (clipboard)
- github.com/atotto/clipboard v0.1.4 - Clipboard operations

**Platform Support:**
- golang.org/x/sys v0.36.0 - Low-level OS APIs (signals, terminal modes)
- golang.org/x/text v0.3.8 - Unicode text handling
- muesli/cancelreader v0.2.2 - Cancellable input reader

## Configuration

**Environment:**
- Single environment variable: `TAMAGO_CHILD=1` - Set when launching in new terminal window (macOS only)
  - Location: `tamago/main.go` lines 29, 37, 50, 67
  - Used to prevent recursive terminal launching on macOS

**Build:**
- Build config: `.docker/Dockerfile` - Multi-stage Docker build
  - Build stage: golang:1.25-alpine with CGO disabled
  - Runtime stage: Alpine 3.18 with minimal footprint
- Module: `module Pessoal` in `go.mod`

**No configuration files:**
- No .env files
- No configuration YAML/TOML
- No environment variable loading libraries (no godotenv, etc.)
- All settings hard-coded or command-line based

## Platform Requirements

**Development:**
- Go 1.25 or higher
- Terminal with Unicode and color support
- macOS/Linux/Windows with standard POSIX shell

**Production:**
- Deployment: Docker container on Alpine 3.18
- Binary size: ~5.8MB (tamago executable in tamago/tamago)
- Runtime: Single standalone executable, no external dependencies
- Target: Linux containers (built with GOOS=linux CGO_ENABLED=0)

**Optional:**
- iTerm2 or Terminal.app (macOS) - For new window launching feature
- Standard terminal emulator with xterm-256color support

## Save File Storage

**Local Filesystem:**
- Save file: `tamago_save.json` in working directory
- Format: JSON with human-readable indentation
- Persistence: Loaded on startup; saved after each action
- Offline decay: Simulates stat degradation when not running (max 50 ticks = 1000 seconds)

---

*Stack analysis: 2026-04-01*
