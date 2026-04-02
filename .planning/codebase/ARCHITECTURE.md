# Architecture

**Analysis Date:** 2026-04-01

## Pattern Overview

**Overall:** Layered MVC with State Machine Orchestration

**Key Characteristics:**
- Multi-modal application with discrete game states (Main Game, Dungeon, Mini-games)
- State-driven UI powered by Bubble Tea framework (TUI framework)
- Separation between model state (`internal/model`), presentation (`internal/ui`), and persistence (`internal/persistence`)
- Dungeon system uses explicit phase-based state machine
- Responsive terminal layout with computed dimensions

## Layers

**Model Layer:**
- Purpose: Core game state and business logic - pet vital stats, progression, evolution, achievements
- Location: `internal/model/`
- Contains: Data structures (`Tama`, `Achievement`, `Stage`), XP/leveling logic, stat calculations
- Depends on: Standard Go libraries only
- Used by: UI layer (for rendering), Persistence layer (for serialization), Dungeon layer (for combat stats)

**UI/Presentation Layer:**
- Purpose: Terminal UI rendering and user input handling via Bubble Tea TUI framework
- Location: `internal/ui/`
- Contains: Model struct with TUI state, view rendering, command processing, layout computation
- Depends on: Model layer, Persistence layer, Dungeon layer, Bubble Tea/Lipgloss libraries
- Used by: Main entry point

**Dungeon Subsystem Layer:**
- Purpose: Self-contained RPG dungeon mechanics with phase-based state machine
- Location: `internal/dungeon/`
- Contains: Combat engine, enemy pools, equipment/inventory, floor generation, biome system
- Depends on: Model layer (for player stats/data), Standard libraries
- Used by: UI layer as a game mode

**Persistence Layer:**
- Purpose: Save/load game state and offline tick simulation
- Location: `internal/persistence/`
- Contains: JSON serialization, file I/O, offline degradation logic
- Depends on: Model layer
- Used by: Main entry point (initialization), UI layer (auto-save)

**Entry Point:**
- Purpose: Bootstrap application, load/initialize pet, launch TUI
- Location: `tamago/main.go`
- Contains: Re-exec logic for terminal window launch (macOS iTerm2/Terminal support), game initialization flow

## Data Flow

**Game Startup:**

1. `tamago/main.go` → Check for existing save file
2. If save exists → `persistence.Load()` → Calculate offline degradation via `applyOfflineTick()`
3. If no save → Create new `model.Tama` with defaults and `DefaultAchievements()`
4. Initialize UI model via `ui.InitialModel(tama)` → Pass to Bubble Tea
5. UI renders and waits for user input

**Main Game Loop (Normal Mode):**

1. User types command (feed, water, pet, sleep, etc.)
2. `ui.Update()` parses input → routes to command function (e.g., `ToFeed()`, `GiveWater()`, `PetTama()`)
3. Command modifies `tama` state directly (hunger, thirst, happiness, etc.)
4. If stat changes trigger level up → `tama.AddXP()` → `StageForLevel()` updates evolution stage
5. Check achievement conditions → unlock if met
6. UI re-renders via `ui.View()` with updated stats
7. Every 20 seconds → `lifeCycleTickCmd()` applies automatic degradation (hunger/thirst/happiness decrease)
8. Every 60 seconds → auto-save via `persistence.Save(tama)`

**Dungeon Mode Entry:**

1. User types `dungeon` command
2. UI creates `dungeon.NewDungeonRun(tama, inventory)`
3. `DungeonRun` initializes phase to `PhaseIntro`, derives combat stats via `DeriveCombatStats(tama)`
4. Dungeon mode takes over UI rendering via `dungeonView.go`
5. On each input → `dungeonGame.HandleInput()` advances phase state machine
6. Combat phases call `Combat.ExecuteAction()` for turn resolution
7. On floor completion → `DungeonRun.Phase` → `PhaseFimAndar` or `PhaseVitoria`
8. On dungeon completion → Phase becomes `PhaseDone`, control returns to main game
9. XP/gold rewards applied to `tama`, inventory persisted

**State Management:**

- **Model state:** Centralized in `model.Tama` struct (single source of truth)
- **UI state:** Held in `ui.Model` struct (frame counter, message, game mode, event tracking)
- **Dungeon state:** Held in `dungeon.DungeonRun` struct (phase-machine, floor state, combat state)
- **Persistence:** JSON serialization of `Tama` to `tamago_save.json`
- **Offline handling:** When loading, simulates 20-second ticks up to 50 ticks max via `applyOfflineTick()`

## Key Abstractions

**Pet Vital State:**
- Purpose: Encapsulates hunger, thirst, sleepiness, happiness, anger, weight as discrete values (0-100)
- Examples: `internal/model/tama.go` - fields like `Hunger`, `Thirst`, `Sleepy`, `Happiness`, `Angry`
- Pattern: Direct field mutation with min/max clamping (e.g., `tama.Hunger = min(tama.Hunger+15, MaxHunger)`)

**Evolution Stage System:**
- Purpose: Represents pet lifecycle progression (Baby → Criança → Adolescente → Adulto → Ancião)
- Examples: `internal/model/evolution.go` - `Stage` type with `DecayRate()`, `Avatar(frame)` methods
- Pattern: Enum-like type with associated behavior (animation frames, stat decay multiplier)

**Dungeon Phase Machine:**
- Purpose: Explicit state machine for dungeon progression through intro, exploration, combat, shops, treasure, victory
- Examples: `internal/dungeon/dungeon.go` - `Phase` enum (PhaseIntro, PhaseMenuPrincipal, etc.), `HandleInput()` routes by phase
- Pattern: Switch statement on current phase, each handler advances phase based on input

**Combat System:**
- Purpose: Turn-based battle engine with stat derivation, biome modifiers, equipment effects
- Examples: `internal/dungeon/combat.go`, `stats.go`, `biome.go`
- Pattern: `Combat` struct holds state, `ExecuteAction()` resolves turn, enemy AI responds

**Equipment/Inventory:**
- Purpose: Equipment slots (Weapon, Armor, Accessory) with rarity tiers, stat bonuses
- Examples: `internal/dungeon/equipment.go`, `items.go`
- Pattern: Equipment struct with level/rarity, `ApplyEquipment()` method on `CombatStats` to add bonuses

**Event System:**
- Purpose: Random interactive events during main game that require input choice
- Examples: `internal/model/events.go` - `Event` struct, event types (Interactive, Negative)
- Pattern: Random trigger, modal overlay in UI, choices modify stats

## Entry Points

**Application Start:**
- Location: `tamago/main.go` - `main()` function
- Triggers: User executes program
- Responsibilities: Platform-specific terminal launching, save loading/initialization, Bubble Tea program creation

**Main Game Loop:**
- Location: `internal/ui/tui.go` - `Update()` method
- Triggers: Each keystroke or timer tick (every 500ms for animation, 20s for life cycle)
- Responsibilities: Parse input, route to command, apply state changes, trigger saves

**Dungeon Session:**
- Location: `internal/ui/minigames.go` - Mode detection when `gameMode == GameModeDungeon`
- Triggers: User types `dungeon` command (requires Level 3+)
- Responsibilities: Create dungeon run, handle input, advance phases, manage combat

**Tick/Auto-Update:**
- Location: `internal/ui/update.go` - `lifeCycleTickCmd()` (20s), `autoSaveTickCmd()` (60s)
- Triggers: Timer events
- Responsibilities: Degrade stats, check death conditions, persist game state

## Error Handling

**Strategy:** Graceful degradation with error logging to stdout

**Patterns:**
- Command functions return success/error messages as strings (no error return, messages displayed to user)
- Persistence errors logged but game continues (e.g., save failures don't crash)
- Invalid input silently ignored (mistyped commands show help text)
- Save file corruption: fallback to new game (detected in `persistence.Load()`)
- Dungeon state errors: phase machine ensures valid transitions, no invalid states reachable

## Cross-Cutting Concerns

**Logging:** No structured logging. Use `fmt.Println()` in main for initialization events (save loaded, level up notices). In-game messages sent via `message` field in UI model, displayed in UI.

**Validation:** Input validation in command functions (e.g., check if `Hunger == MaxHunger` before feeding). Stat bounds enforced via `min()`/`max()` helpers throughout.

**Authentication:** Not applicable (single-player game).

**State Initialization:** Default `model.Tama` created with balanced stats (50/50/50/50). Achievement list initialized via `model.DefaultAchievements()`. Dungeon inventory loaded from JSON raw message in save file.

---

*Architecture analysis: 2026-04-01*
