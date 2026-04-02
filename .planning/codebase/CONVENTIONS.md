# Coding Conventions

**Analysis Date:** 2026-04-01

## Naming Patterns

**Files:**
- Go files: lowercase with underscores separating logical units (e.g., `biome_test.go`, `save.go`, `tui.go`)
- Package structure mirrors functionality: `internal/model/`, `internal/dungeon/`, `internal/ui/`, `internal/persistence/`
- Test files follow standard Go convention: `*_test.go` suffix in same package as tested code

**Functions:**
- PascalCase for exported functions (public): `Save()`, `Load()`, `NewCombat()`, `AddXP()`, `ModifierForBiome()`
- camelCase for unexported functions (private): `applyOfflineTick()`, `playerAttack()`, `enemyTurn()`, `normalizeBiome()`
- Receiver methods attached to structs: `(t *Tama) AddXP()`, `(s Stage) String()`, `(c *Combat) ExecuteAction()`
- Command functions with verb prefixes: `ToFeed()`, `GiveWater()`, `PetTama()` (UI command functions in `internal/ui/update.go`)

**Variables:**
- Short names for loop variables and temporary values: `i`, `err`, `t`, `s`, `c`
- Descriptive names for struct fields and local variables: `player`, `enemy`, `modifier`, `biome`, `Hunger`, `Thirst`
- Constants in UPPERCASE: `MaxHunger`, `MinHunger`, `MaxLevel`, `SaveFile`
- Private package variables in camelCase: `LevelThresholds` (exported slice), `Events` (exported slice)

**Types:**
- PascalCase for all type names: `Tama`, `Combat`, `Event`, `Achievement`, `Biome`, `CombatStats`
- Struct field names: PascalCase (JSON tags use snake_case): `Name`, `Hunger`, `LastSaved`, `XP`, `Level`
- Interface-like behavior achieved with receiver methods on concrete types (no explicit interfaces used extensively)

## Code Style

**Formatting:**
- Standard Go formatting (gofmt style - no custom formatter configured)
- Consistent indentation with tabs (Go standard)
- Line length: No strict limit observed, but generally reasonable (under 100 chars in most cases)
- Spacing: Single blank line between logical sections, double blank lines between major function groups

**Linting:**
- No linting configuration file detected (`.eslintrc`, `.golangci.yml` not present)
- Code follows Go idioms and conventions naturally
- Unused variable elimination: `_ = enemyDefending` pattern used in `internal/dungeon/combat.go:217` to suppress warnings

## Import Organization

**Order:**
1. Standard library imports (`fmt`, `encoding/json`, `time`, `os`, `math/rand`)
2. External packages (charmbracelet libraries: `tea`, `lipgloss`, `textinput`)
3. Local project imports (`Pessoal/internal/...`)

**Example from `internal/ui/tui.go`:**
```go
import (
    "Pessoal/internal/dungeon"
    "Pessoal/internal/model"
    "Pessoal/internal/persistence"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)
```

**Path Aliases:**
- `tea` aliased for `github.com/charmbracelet/bubbletea` (e.g., `tea.Cmd`, `tea.Msg`)
- No other aliases used

## Error Handling

**Patterns:**
- Standard Go error pattern: `error` as last return value
- Immediate inline error checking: `if err != nil { return err }`
- Silent failures allowed in non-critical paths: `if err := json.Unmarshal(...); err != nil { return nil, err }`
- No panic() used; errors propagated up
- User-facing error messages through return strings in UI functions (e.g., `ToFeed()` returns `string`)

**Example from `internal/persistence/save.go:16-19`:**
```go
if err != nil {
    return err
}
```

**In main.go:** Errors logged before os.Exit:
```go
if _, err := p.Run(); err != nil {
    fmt.Printf("Erro ao iniciar o programa: %v\n", err)
    os.Exit(1)
}
```

## Logging

**Framework:** `fmt` package (standard output)

**Patterns:**
- `fmt.Println()` for status messages
- `fmt.Printf()` for formatted output
- Log messages in Portuguese (user-facing) in UI; technical messages in English in backend
- Combat log accumulated in `Combat.Log []string` field (not printed to stdout)
- Save/load feedback printed to main program console: `fmt.Println("💾 Save carregado! Bem-vindo de volta,", tama.Name+"!")`

**No structured logging framework used** (no zap, logrus, slog)

## Comments

**When to Comment:**
- Function comments above exported functions: `// Save persists the Tama state to file` (pattern observed in package docs)
- Inline comments for non-obvious algorithm logic: `// Chance de crítico baseada na Sorte` in `internal/dungeon/combat.go:149`
- Type definition comments: `// Combat motor de combate por turnos.` above struct definition
- Constant group comments: `// XP e Evolução` before related field group in `internal/model/tama.go:49`

**JSDoc/TSDoc:**
- Not applicable (Go project)
- Go convention: Comments above exported symbols act as documentation
- Pattern observed: Brief Portuguese comments describing functionality rather than verbose English

## Function Design

**Size:** 
- Small focused functions preferred: `playerAttack()` (26 lines), `enemyTurn()` (50 lines)
- Larger functions acceptable for cohesive operations: `ExecuteAction()` (86 lines) for combat orchestration
- Average function size: 15-40 lines

**Parameters:**
- Pointer receivers for methods that modify state: `(c *Combat) ExecuteAction()`, `(t *Tama) AddXP()`
- Value receivers for query methods: `(s Stage) String()` returns string
- Function parameters usually 2-4 arguments maximum
- Long parameter lists avoided; use struct fields instead

**Return Values:**
- Single return value common for simple operations: `Save()` returns `error`
- Multiple returns standard: `Load()` returns `(*model.Tama, error)`
- Boolean returns for success/state checks: `AddXP()` returns `bool` (leveled up or not), `IsOver()` returns `bool`
- String returns for UI feedback: `ToFeed()` returns `string` (message to display)

## Module Design

**Exports:**
- Selective export: Only public APIs at package level
- Internal helper functions unexported (lowercase)
- Structs exported with selective field export: `type Tama struct { Name string, Hunger int }`
- Constants exported when they represent configuration: `MaxHunger = 100`, `MinHunger = 0`

**Barrel Files:**
- No barrel files (`index.go` or similar) used
- Each file focused on related functionality: `biome.go` contains biome logic only

**Package Structure:**
- `internal/model/`: Core data types and state logic
- `internal/dungeon/`: Dungeon RPG game mechanics
- `internal/ui/`: UI rendering and input handling
- `internal/persistence/`: Save/load file operations
- `tamago/main.go`: Entry point and initialization
