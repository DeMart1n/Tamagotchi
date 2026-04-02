# Testing Patterns

**Analysis Date:** 2026-04-01

## Test Framework

**Runner:**
- Go's built-in `testing` package
- Run with: `go test ./...` (standard Go test command)
- Configuration file: None detected (uses Go defaults)

**Run Commands:**
```bash
go test ./...           # Run all tests in project
go test ./internal/...  # Run tests in internal packages
go test -v ./...        # Verbose output showing individual tests
go test -run TestName   # Run specific test
go test -count=1 ./... # Run without caching
```

**Assertion Library:**
- Go standard `testing` package (manual assertions)
- No assertion helpers (testify, assert library) used
- Direct comparisons with `t.Fatalf()` for failures

## Test File Organization

**Location:**
- Co-located with source files in same package
- Pattern: Source file `biome.go` has test file `biome_test.go` in same `internal/dungeon` package
- Test files NOT in separate `tests/` directory

**Naming:**
- Standard Go: `*_test.go` suffix
- Examples: `biome_test.go`, `save_test.go`

**Structure:**
```
internal/
├── dungeon/
│   ├── biome.go
│   ├── biome_test.go      # Tests for biome functionality
│   └── combat.go          # No test file (untested)
├── persistence/
│   ├── save.go
│   └── save_test.go       # Tests for persistence
└── model/
    ├── tama.go            # No test file
    ├── evolution.go       # No test file
    └── achievements.go    # No test file
```

## Test Structure

**Suite Organization:**
```go
// Pattern from internal/dungeon/biome_test.go
func TestModifierForBiome_HasExpectedPrimaryEffects(t *testing.T) {
    // Arrange: Set up test data
    forest := ModifierForBiome(BiomeForest)
    
    // Act & Assert: Perform action and check results
    if forest.HungerRecoveryBonus <= 0 {
        t.Fatalf("expected forest hunger recovery bonus, got %d", forest.HungerRecoveryBonus)
    }
    // ... more assertions
}
```

**Patterns:**
- No explicit setUp/tearDown methods
- Arrange-Act-Assert pattern used (inline within test function)
- Test functions are independent; no shared state
- Immediate assertion on failure with `t.Fatalf()`

## Mocking

**Framework:** Manual mocking (no mocking library like gomock, mockery)

**Patterns from `biome_test.go`:**
```go
// Direct struct instantiation for test objects
player := &CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 10, Velocidade: 10, Sorte: 0}
enemy := &Enemy{Name: "Dummy", HPMax: 100, HPCurrent: 100, Ataque: 1, Defesa: 0, Velocidade: 1}

// Pass to function under test
combat := NewCombat(player, enemy, BiomeVolcanic)

// Verify results
if combat.Player.HPCurrent >= 95 {
    t.Fatalf("expected volcanic DOT to reduce player HP significantly, got %d", combat.Player.HPCurrent)
}
```

**Seeding for Determinism:**
```go
// From biome_test.go:55-56
rand.Seed(7)  // Set random seed for reproducible results
forestCombat.enemyTurn()
```

**What to Mock:**
- External system interactions (file I/O, network)
- Random number generation (seed before calling)

**What NOT to Mock:**
- Internal game logic (combat, biome effects)
- Game state (Tama, Combat structs)
- Calculation functions (damage, XP)

## Fixtures and Factories

**Test Data:**
No explicit fixtures or factories found. Test data created inline:

```go
// From save_test.go:9
tama := &model.Tama{CurrentBiome: "Unknown"}
normalizeBiome(tama)

// From biome_test.go:37-39
basePlayer := CombatStats{
    HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 10, Sorte: 0
}
baseEnemy := Enemy{
    Name: "Fire Mage", HPMax: 100, HPCurrent: 100, 
    Ataque: 10, Defesa: 0, Velocidade: 5, IsFire: true
}
```

**Location:**
- No separate `testdata/` or `fixtures/` directory
- Test data created directly in test functions

## Coverage

**Requirements:** None enforced (no coverage threshold configured)

**View Coverage:**
```bash
go test -cover ./...           # Display coverage percentage
go test -coverprofile=cov.out ./...
go tool cover -html=cov.out    # Generate HTML coverage report
```

**Current Coverage Analysis:**
- Only 2 test files present: `internal/dungeon/biome_test.go`, `internal/persistence/save_test.go`
- Untested packages: `internal/ui/`, `internal/model/` (except indirect through biome_test)
- **Coverage estimate: ~10-15% of codebase** (only dungeon and persistence partially tested)

## Test Types

**Unit Tests:**
- Scope: Individual functions and methods
- Approach: Test calculations and state transitions in isolation
- Example: `TestModifierForBiome_HasExpectedPrimaryEffects` (lines 8-34 in biome_test.go) tests modifier calculation

**Integration Tests:**
- Scope: Multiple components working together
- Example: `TestCombat_ForestFireVulnerabilityIncreasesDamage` (lines 48-69 in biome_test.go) tests:
  - Biome modifiers applied to combat
  - Enemy attack mechanics
  - Player HP reduction
  - Result comparison between forest and neutral biomes

**E2E Tests:**
- Not used - GUI application not tested end-to-end
- Manual testing required for TUI interactions

## Common Patterns

**Assertion Pattern:**
```go
// Direct if + t.Fatalf pattern
if forest.HungerRecoveryBonus <= 0 {
    t.Fatalf("expected forest hunger recovery bonus, got %d", forest.HungerRecoveryBonus)
}
```

**Comparative Testing:**
```go
// From TestCombat_ForestFireVulnerabilityIncreasesDamage
forestCombat := NewCombat(&forestPlayer, &forestEnemy, BiomeForest)
rand.Seed(7)
forestCombat.enemyTurn()
forestDamage := 100 - forestCombat.Player.HPCurrent

neutralCombat := NewCombat(&neutralPlayer, &neutralEnemy, BiomeAbyssal)
rand.Seed(7)
neutralCombat.enemyTurn()
neutralDamage := 100 - neutralCombat.Player.HPCurrent

if forestDamage <= neutralDamage {
    t.Fatalf("expected forest fire vulnerability to increase damage...")
}
```

**Testing State Defaults:**
```go
// From save_test.go: Testing that unknown biome defaults to Forest
tama := &model.Tama{CurrentBiome: "Unknown"}
normalizeBiome(tama)

if tama.CurrentBiome != "Florestal" {
    t.Fatalf("expected fallback biome Florestal, got %q", tama.CurrentBiome)
}
```

## Test Execution Notes

- Tests run sequentially (default Go behavior for unit tests)
- Randomization controlled via `rand.Seed()` before random operations
- No benchmarking functions present
- No test helpers or utility functions across test files
- Each test file is self-contained

## Critical Testing Gaps

**High Priority Areas Lacking Tests:**
1. `internal/ui/` - UI rendering and event handling untested
2. `internal/model/tama.go` - Core game state logic (XP, evolution) untested
3. `internal/model/achievements.go` - Achievement checking untested
4. `internal/dungeon/combat.go` - Only biome effects tested, not core combat logic
5. `internal/persistence/save.go` - Only normalizeBiome tested, not Save() or Load()

**Risk:** Game balance changes, evolution mechanics, and save/load functionality could break without notice.
