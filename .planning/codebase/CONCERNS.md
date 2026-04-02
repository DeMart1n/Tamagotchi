# Codebase Concerns

**Analysis Date:** 2026-04-01

## Tech Debt

**Global Random Seed Not Set:**
- Issue: `math/rand` is used throughout the codebase without explicitly seeding with `rand.Seed()`. This causes non-deterministic behavior to be limited to Go 1.20+ where auto-seeding was introduced, but the codebase doesn't guarantee this.
- Files: `internal/ui/minigames.go` (line 36), `internal/dungeon/combat.go` (lines 5, 144, 151, 204, 212, 280), `internal/dungeon/biome.go`, `internal/dungeon/enemy.go`, `internal/dungeon/floor.go`, `internal/ui/update.go` (line 291)
- Impact: Combat outcomes, floor generation, and random events may not be properly randomized on older Go versions or systems without entropy
- Fix approach: Seed `math/rand` explicitly in `main()` with `rand.Seed(time.Now().UnixNano())` before TUI initialization, or migrate to `rand/v2` (Go 1.22+) which auto-seeds

**Silent JSON Deserialization Error in Inventory Loading:**
- Issue: In `internal/ui/tui.go` line 101, inventory deserialization errors are silently ignored with `_ = json.Unmarshal(tama.Inventory, inv)`. Malformed save files will result in an empty inventory without user feedback.
- Files: `internal/ui/tui.go` (line 101)
- Impact: Player data loss (equipment/gold) if save file becomes corrupted; no warning or error message to user
- Fix approach: Log the error or provide a recovery mechanism; at minimum, log JSON unmarshal failures to stderr

**json.RawMessage Serialization Pattern Is Fragile:**
- Issue: Inventory is stored as `json.RawMessage` in `Tama` struct (`internal/model/tama.go` line 67) to avoid circular dependencies. This requires manual marshal/unmarshal cycles and creates a mismatch between runtime state and persisted state.
- Files: `internal/model/tama.go` (line 67), `internal/ui/tui.go` (lines 100-101, 385-391)
- Impact: Easy to forget to save inventory changes; inventory can desync from other Tama stats during a session if not explicitly saved
- Fix approach: Refactor `Tama` to embed `*dungeon.Inventory` directly, or create a clear interface for inventory persistence lifecycle

**No Bounds Checking on Equipment Bonus Stacking:**
- Issue: Equipment bonuses stack without caps. A player with all legendary equipment could have unbounded stats, breaking combat balance.
- Files: `internal/dungeon/equipment.go` (lines 146-164), `internal/dungeon/stats.go`
- Impact: High-end content becomes trivial if player acquires top-tier equipment; no difficulty scaling ceiling
- Fix approach: Add stat caps per level/stage, or scale enemy stats with total player equipment bonus value

## Known Bugs

**Biome String Normalization Mismatch:**
- Symptoms: `CurrentBiome` values "Gelido" and "Vulcanico" in normalization but could be stored as other variants
- Files: `internal/persistence/save.go` (line 87: "Gelido"), `internal/dungeon/dungeon.go` (biome enum likely uses different casing)
- Trigger: Load save with biome from earlier version, or manual JSON editing
- Workaround: Always set biome through `RandomBiome()` function which ensures consistency

**Offline Tick Calculation May Miss Degradation:**
- Symptoms: Player could game the system by restarting the app multiple times within the 20-second window to avoid stat degradation
- Files: `internal/persistence/save.go` (line 36: `int(elapsed.Seconds()) / 20`)
- Trigger: Exit and restart within 20 seconds repeatedly
- Workaround: None currently; intended behavior unclear (MVP feature)

**Enemy Defense Divided by 2 in Damage Calculation:**
- Symptoms: Enemy defense contributes less to damage reduction than player defense
- Files: `internal/dungeon/combat.go` (line 137: `baseDmg := c.Player.Ataque - c.Enemy.Defesa/2`, line 175 same pattern)
- Trigger: Any combat encounter
- Workaround: Combat is still beatable; may be intentional to favor player

## Security Considerations

**No Input Validation on setlvl Command:**
- Risk: User can set level to negative values or extremely high values, breaking progression
- Files: `internal/ui/tui.go` (line 368: `if err != nil || lvl < 0`)
- Current mitigation: Bounds check exists (`lvl < 0`), but no upper bound
- Recommendations: Add max level cap check, require level <= 50 or similar; log secret command usage

**No Save File Integrity Verification:**
- Risk: Malicious or corrupted JSON could crash the game or create invalid game states
- Files: `internal/persistence/save.go` (lines 24-31)
- Current mitigation: None; raw `json.Unmarshal` only
- Recommendations: Add schema validation, checksum verification, or version field to detect incompatible saves

## Performance Bottlenecks

**Linear Equipment Search on Every Stat Derivation:**
- Problem: `EquipmentByID()` does linear search through `AllEquipment` array (only 12 items now, but will scale)
- Files: `internal/dungeon/equipment.go` (line 167)
- Cause: No indexing; called during inventory operations
- Improvement path: Build `map[string]*Equipment` on init if equipment count exceeds ~20 items

**No Caching of Combat Stats During Multi-Turn Encounters:**
- Problem: Combat stats are not cached; biome modifiers and equipment bonuses recalculated per turn
- Files: `internal/dungeon/combat.go` (line 44), `internal/dungeon/stats.go`
- Cause: Stats passed as value type, not cached
- Improvement path: Cache `CombatStats` at combat start; only update on item use or state changes

## Fragile Areas

**Dungeon State Machine Has No Rollback:**
- Files: `internal/dungeon/dungeon.go` (phases 0-11)
- Why fragile: 11 phase states with complex transitions. If a transition is missed or a handler doesn't advance phase correctly, the dungeon gets stuck. No error state for invalid transitions.
- Safe modification: Add phase validation in `HandleInput()`, log invalid transitions to detect missed handlers
- Test coverage: Only biome tests exist; no phase transition or dungeon completion tests

**UI Tick Timer Interacts With Save Serialization:**
- Files: `internal/ui/tui.go` (lines 174-176, 202-204, 266-273), `internal/persistence/save.go`
- Why fragile: Multiple code paths call `saveDungeonInventory()` before `persistence.Save()`. If one path is missed, inventory is lost. Hard to trace from reading code.
- Safe modification: Centralize save logic; create `SaveGame()` function that always serializes inventory first
- Test coverage: No integration tests for save/load cycle

**Combat Log Strings Have No Formatting Validation:**
- Files: `internal/dungeon/combat.go` (all `fmt.Sprintf` calls in combat messages)
- Why fragile: If translated or modified, log strings could exceed terminal width or contain special chars that break rendering. No truncation.
- Safe modification: Add message length limit; sanitize special characters in combat log
- Test coverage: No UI rendering tests for long combat logs

## Scaling Limits

**Inventory Backpack Unbounded:**
- Current capacity: No limit on backpack size in `Inventory.Backpack` slice
- Limit: Could hit memory issues with thousands of items; serialization time grows linearly
- Scaling path: Add backpack size cap (e.g., 20 items), implement overflow mechanic (drop oldest, or reject)

**AllEquipment Array Growth:**
- Current capacity: 12 items total (4 weapons, 4 armor, 4 accessories)
- Limit: Linear search will be noticeable above ~100 items; JSON marshaling time scales
- Scaling path: Add equipment categories/filtering; implement equipment database with ID indexing once >20 items

**Offline Tick Calculation Bounds:**
- Current capacity: `MaxOfflineTicks = 50` limits offline degradation to ~16 minutes max
- Limit: Player away for 1+ hour still only gets 50 ticks worth of degradation (capped at line 37)
- Scaling path: Clarify intent: is this to prevent "dead Tama save" state, or should long absences kill the pet? Document the cap.

## Dependencies at Risk

**No Explicit Go Version in go.mod:**
- Risk: `go 1.25` is future-looking (as of cutoff date). Auto-seeding of `math/rand` requires Go 1.20+. Upgrading to 1.25 could introduce subtle changes.
- Impact: Random number generation, runtime behavior changes
- Migration plan: Explicitly seed `math/rand` to ensure portability; test on Go 1.20 LTS if targeting older versions

**Bubble Tea UI Library Dependency:**
- Risk: Bubble Tea is experimental (based on version `v1.3.10`); future versions may introduce breaking API changes
- Impact: TUI rendering, input handling could change significantly
- Migration plan: Pin dependency version in go.sum; evaluate upgrade path quarterly

## Missing Critical Features

**No Dungeon Difficulty Selection:**
- Problem: All players enter the same dungeon with random biome, but no way to adjust difficulty or skip content
- Blocks: Accessibility for new/casual players; no catch-up mechanic for low-level players
- Recommendation: Add difficulty presets (Easy/Normal/Hard) that scale enemy stats and loot

**No Death Penalty or Consequence System:**
- Problem: Dungeon defeat just exits; no loot loss, no XP loss, no cooldown
- Blocks: Long-term progression feels consequence-free; no strategic depth
- Recommendation: Implement loot drop on defeat (partial gold loss) or cooldown before re-entering

**No Combat Auto-Resolve Option:**
- Problem: Every combat is fully interactive turn-by-turn; no way to speed up or auto-farm
- Blocks: Farming for specific items becomes tedious; no AFK progression
- Recommendation: Add "Fight" button to auto-resolve combat based on stats

## Test Coverage Gaps

**No tests for Dungeon Phase Transitions:**
- What's not tested: Moving from `PhaseIntro` → `PhaseMenuPrincipal` → `PhaseExplorando`, etc. Edge cases like skipping phases or invalid transitions.
- Files: `internal/dungeon/dungeon.go` (all phase handler methods `handleIntro()`, `handleMenuPrincipal()`, etc.)
- Risk: High. A typo in phase assignment could silently break dungeon progression without any test failure.
- Priority: **High** - Phase machine is core to dungeon gameplay

**No tests for Offline Tick Degradation:**
- What's not tested: Elapsed time calculations, missed tick counts, death conditions after offline
- Files: `internal/persistence/save.go` (lines 33-43, 57-83)
- Risk: Medium. Corruption of save files or unexpected player death could occur without warning.
- Priority: **High** - Affects game state integrity

**No tests for Equipment Serialization/Deserialization:**
- What's not tested: Round-trip serialization of Inventory, handling of unknown equipment IDs, missing slots
- Files: `internal/ui/tui.go` (line 101), `internal/dungeon/equipment.go` (Inventory marshaling)
- Risk: Medium. Corrupted save files or missing equipment could go unnoticed.
- Priority: **Medium** - Data loss risk

**No tests for Combat Edge Cases:**
- What's not tested: HP overflow when healing, negative damage values, frozen player skipping turns
- Files: `internal/dungeon/combat.go` (all ExecuteAction paths, useItem healing cap at line 244)
- Risk: Low. Logic seems sound, but edge cases like healing overflow or freeze overflow could produce weird states.
- Priority: **Medium** - Fairness and consistency

**No tests for Event Application:**
- What's not tested: Event stat delta clamping to min/max bounds, event trigger probability, choice branching
- Files: `internal/ui/update.go` (lines 298-318)
- Risk: Medium. Player could end up with invalid stat states (e.g., happiness > 100 or < 0) if event deltas bypass clamping.
- Priority: **Medium** - Stat validation

---

*Concerns audit: 2026-04-01*
