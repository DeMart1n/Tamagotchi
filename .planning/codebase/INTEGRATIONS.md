# External Integrations

**Analysis Date:** 2026-04-01

## APIs & External Services

**None currently implemented.** The application is self-contained with no external API dependencies or third-party service integrations.

---

## Data Storage

**Databases:**
- Not used. No SQL database, NoSQL database, or cloud database service integrated.

**File Storage:**
- Local filesystem only
  - Save file: `tamago_save.json` in application working directory
  - Location: `internal/persistence/save.go`
  - Format: JSON serialization via Go's `encoding/json` package
  - Save triggers: After each game action (feed, water, sleep, pet, etc.)
  - Load triggers: On application startup

**Caching:**
- None. Game state is held entirely in memory after loading.

---

## Authentication & Identity

**Auth Provider:**
- Not applicable. Single-player CLI application with no user accounts, authentication, or identity system.

---

## Monitoring & Observability

**Error Tracking:**
- Not implemented. No error reporting service (Sentry, DataDog, etc.)
- Errors logged to stdout via `fmt.Printf` and `fmt.Println`

**Logs:**
- Console-based only via standard Go `fmt` package
  - Example: `fmt.Println("💾 Save carregado! Bem-vindo de volta,", tama.Name+"!")` in `tamago/main.go:79`
  - No structured logging framework (no logrus, zap, etc.)
  - No log aggregation or persistence

---

## CI/CD & Deployment

**Hosting:**
- Not deployed publicly. Local development and Docker container support only.

**Container Registry:**
- None. Image: `tamago:latest` built locally via `docker compose up --build`

**CI Pipeline:**
- Not implemented. No GitHub Actions, GitLab CI, or other CI/CD service configured.

---

## Environment Configuration

**Required env vars:**
- `TAMAGO_CHILD=1` - Optional, set internally for macOS terminal launching (user doesn't set this)
  - All other configuration is built into the binary

**Optional env vars:**
- None documented

**Secrets location:**
- Not applicable. No secrets management needed (no API keys, credentials, or sensitive data).
- `.gitignore` contains: `tamago_save.json` (user game saves, not secrets)

---

## Webhooks & Callbacks

**Incoming:**
- Not applicable. CLI application; no HTTP endpoints or webhook receivers.

**Outgoing:**
- Not applicable. No outbound HTTP calls or webhook deliveries.

---

## System Integrations

**Terminal Environment:**
- macOS Terminal.app - Supported for launching application in new window
- iTerm2 - Preferred terminal on macOS (auto-detected in `tamago/main.go:32`)
- Linux terminal emulators - Standard xterm-compatible support
- Windows console - Supported via ConInput library for input handling

**Clipboard:**
- Clipboard operations supported via `github.com/atotto/clipboard v0.1.4`
- Used for: Unknown (dependency present but not referenced in main code)
- Impact: Low; optional feature if implemented

**Signal Handling:**
- OS signals: Handled by Bubble Tea framework (Ctrl+C exit)
- Uses: `golang.org/x/sys` for low-level signal management

---

## Data Flow

**Startup:**
1. Application launched from CLI
2. Check if save file exists: `persistence.Exists()` in `internal/persistence/save.go:52`
3. If exists and player not dead: Load save via `persistence.Load()` 
   - Deserialize from JSON
   - Apply offline decay based on elapsed time since `LastSaved`
   - Max decay: 50 ticks (1000 seconds)
4. If no save or player dead: Create new Tama with defaults
5. Start TUI event loop

**During Gameplay:**
1. User enters command
2. Parse command in `internal/ui/update.go`
3. Update Tama state model in memory
4. Save to `tamago_save.json` immediately after each action
5. Render updated UI with Lipgloss/Bubble Tea
6. Loop to step 1

**Shutdown:**
1. User exits (ESC or Ctrl+C)
2. Final save performed if modified
3. Process exits (no cleanup hooks)

---

## Data Persistence Details

**Save File Schema:**
- Location: `internal/model/tama.go`
- Fields saved: Name, all stats (Hunger, Thirst, Sleepy, Happiness, Angry, Weight)
- Metadata: LastSaved (timestamp), Dead (boolean), biome, level, XP
- Achievements: Serialized as array of achievement objects
- Dungeon state: Not persisted between sessions (fresh start each entry)
- Equipment/Inventory: Persisted in dungeon state (location: `internal/dungeon/equipment.go`)

**Offline Mechanics:**
- When loading, elapsed time since `LastSaved` is calculated
- Each 20 seconds offline = 1 degradation tick
- Max ticks: 50 (1000 seconds / ~16 minutes of degradation applied)
- Stats degraded: Hunger, Thirst, Sleepy, Happiness, Angry (all tick down)
- Death triggers: If Hunger or Thirst reach 0 during offline decay
- Depression trigger: If Happiness reaches 0 during offline decay

---

## Third-Party Integrations Summary

| Category | Status | Details |
|----------|--------|---------|
| Cloud Services | Not used | Entirely self-contained |
| Payment Processing | Not used | No monetization |
| Analytics | Not used | No user tracking |
| Push Notifications | Not used | CLI only |
| Social/Multiplayer | Not used | Single-player only (planned for roadmap) |
| Weather/Location | Not used | No location-aware features |
| Maps/Geo | Not used | Not applicable |
| Email/SMS | Not used | Not used |
| Search | Not used | Not applicable |
| CDN | Not used | Single executable |

---

*Integration audit: 2026-04-01*
