package persistence

import (
	"Pessoal/internal/model"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// --- helpers de teste ---

// tmpSaveFile redireciona o save para um arquivo temporário e retorna cleanup.
// Como SaveFile é uma constante, os testes escrevem no path padrão e limpam depois.
func tmpSaveFile(t *testing.T) func() {
	t.Helper()
	return func() {
		os.Remove(SaveFile)
		os.Remove(SaveFile + ".tmp")
	}
}

func baseTama() *model.Tama {
	return &model.Tama{
		Name:         "TestGO",
		Hunger:       80,
		Thirst:       70,
		Sleepy:       60,
		Happiness:    90,
		Angry:        10,
		Weight:       50,
		Level:        3,
		XP:           42,
		CurrentBiome: "Vulcanico",
		Achievements: model.DefaultAchievements(),
	}
}

// --- testes originais preservados ---

func TestNormalizeBiomeDefaultsToForest(t *testing.T) {
	tama := &model.Tama{CurrentBiome: "Unknown"}
	normalizeBiome(tama)
	if tama.CurrentBiome != "Florestal" {
		t.Fatalf("expected fallback biome Florestal, got %q", tama.CurrentBiome)
	}
}

func TestNormalizeBiomeKeepsKnownValue(t *testing.T) {
	tama := &model.Tama{CurrentBiome: "Vulcanico"}
	normalizeBiome(tama)
	if tama.CurrentBiome != "Vulcanico" {
		t.Fatalf("expected known biome to remain unchanged, got %q", tama.CurrentBiome)
	}
}

// --- roundtrip ---

func TestSaveAndLoadRoundtrip(t *testing.T) {
	defer tmpSaveFile(t)()

	original := baseTama()
	original.TotalFeeds = 15
	original.TotalDungeonRuns = 2

	if err := Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !Exists() {
		t.Fatal("Exists() returned false after Save")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, original.Name)
	}
	if loaded.Level != original.Level {
		t.Errorf("Level: got %d, want %d", loaded.Level, original.Level)
	}
	if loaded.XP != original.XP {
		t.Errorf("XP: got %d, want %d", loaded.XP, original.XP)
	}
	if loaded.CurrentBiome != original.CurrentBiome {
		t.Errorf("CurrentBiome: got %q, want %q", loaded.CurrentBiome, original.CurrentBiome)
	}
	if loaded.TotalFeeds != original.TotalFeeds {
		t.Errorf("TotalFeeds: got %d, want %d", loaded.TotalFeeds, original.TotalFeeds)
	}
	if loaded.TotalDungeonRuns != original.TotalDungeonRuns {
		t.Errorf("TotalDungeonRuns: got %d, want %d", loaded.TotalDungeonRuns, original.TotalDungeonRuns)
	}
	if loaded.Sleeping {
		t.Error("Sleeping should always be false after Load")
	}
}

func TestLoadSetsSleepingFalse(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()
	tama.Sleeping = true // força sleeping = true no save
	tama.LastSaved = time.Now()

	data, _ := json.MarshalIndent(SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       tama.LastSaved,
		Tama:          tama,
	}, "", "  ")
	os.WriteFile(SaveFile, data, 0644)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Sleeping {
		t.Error("Sleeping deve ser false após Load — Tama não pode carregar travado dormindo")
	}
}

// --- inventory ---

func TestSaveAndLoadInventory(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()

	type fakeInventory struct {
		Gold int    `json:"gold"`
		Item string `json:"item"`
	}
	inv := fakeInventory{Gold: 250, Item: "espada_afiada"}
	raw, _ := json.Marshal(inv)
	tama.Inventory = json.RawMessage(raw)

	if err := Save(tama); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if string(loaded.Inventory) != string(raw) {
		t.Errorf("Inventory mismatch:\n  got  %s\n  want %s", loaded.Inventory, raw)
	}

	// Verifica que o campo também está no envelope de topo
	fileData, _ := os.ReadFile(SaveFile)
	var envelope SaveData
	if err := json.Unmarshal(fileData, &envelope); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if string(envelope.Inventory) != string(raw) {
		t.Errorf("envelope.Inventory mismatch:\n  got  %s\n  want %s", envelope.Inventory, raw)
	}
}

// --- schema version ---

func TestSaveWritesSchemaVersion(t *testing.T) {
	defer tmpSaveFile(t)()

	if err := Save(baseTama()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, _ := os.ReadFile(SaveFile)
	var envelope SaveData
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if envelope.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("SchemaVersion: got %d, want %d", envelope.SchemaVersion, CurrentSchemaVersion)
	}
}

// --- migração de save legado ---

func TestLoadLegacySave(t *testing.T) {
	defer tmpSaveFile(t)()

	// Simula save legado: JSON raiz direto é um model.Tama (formato flat antigo, sem envelope)
	legacy := baseTama()
	legacy.LastSaved = time.Now()
	data, err := json.MarshalIndent(legacy, "", "  ")
	if err != nil {
		t.Fatalf("marshal legacy: %v", err)
	}
	if err := os.WriteFile(SaveFile, data, 0644); err != nil {
		t.Fatalf("write legacy: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load of legacy save failed: %v", err)
	}
	if loaded.Name != legacy.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, legacy.Name)
	}
	if loaded.Level != legacy.Level {
		t.Errorf("Level: got %d, want %d", loaded.Level, legacy.Level)
	}
}

// --- degradação offline ---

func TestOfflineTickDegradation(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()
	// 100s atrás = 5 ticks de 20s
	tama.LastSaved = time.Now().Add(-100 * time.Second)
	tama.Hunger = 80
	tama.Thirst = 80

	data, _ := json.MarshalIndent(SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       tama.LastSaved,
		Tama:          tama,
	}, "", "  ")
	os.WriteFile(SaveFile, data, 0644)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	expected := 80 - 5
	if loaded.Hunger != expected {
		t.Errorf("Hunger após 5 ticks offline: got %d, want %d", loaded.Hunger, expected)
	}
}

func TestOfflineTickCappedAtMax(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()
	tama.LastSaved = time.Now().Add(-time.Duration(MaxOfflineTicks*20+9999) * time.Second)
	tama.Hunger = 100

	data, _ := json.MarshalIndent(SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       tama.LastSaved,
		Tama:          tama,
	}, "", "  ")
	os.WriteFile(SaveFile, data, 0644)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	minExpected := 100 - MaxOfflineTicks
	if loaded.Hunger < minExpected {
		t.Errorf("Hunger degradou além do cap: got %d (mínimo esperado %d)", loaded.Hunger, minExpected)
	}
}

func TestOfflineTickKillsTamaOnStarvation(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()
	tama.LastSaved = time.Now().Add(-time.Duration(MaxOfflineTicks*20+9999) * time.Second)
	tama.Hunger = 1
	tama.Thirst = 1

	data, _ := json.MarshalIndent(SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       tama.LastSaved,
		Tama:          tama,
	}, "", "  ")
	os.WriteFile(SaveFile, data, 0644)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if !loaded.Dead {
		t.Error("Tama deveria estar Dead após inanição offline")
	}
}

// --- Delete ---

func TestDelete(t *testing.T) {
	defer tmpSaveFile(t)()

	if err := Save(baseTama()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !Exists() {
		t.Fatal("arquivo deve existir antes do Delete")
	}
	if err := Delete(); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if Exists() {
		t.Error("arquivo não deve existir após Delete")
	}
	// Segundo Delete deve ser no-op
	if err := Delete(); err != nil {
		t.Errorf("segundo Delete deve ser no-op, got: %v", err)
	}
}

// --- escrita atômica ---

func TestSaveIsAtomic(t *testing.T) {
	defer tmpSaveFile(t)()

	if err := Save(baseTama()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(SaveFile + ".tmp"); !os.IsNotExist(err) {
		t.Error("arquivo .tmp não deve restar após Save bem-sucedido")
	}
}

// --- achievements ---

func TestAchievementsPreservedOnRoundtrip(t *testing.T) {
	defer tmpSaveFile(t)()

	tama := baseTama()
	for i := range tama.Achievements {
		if tama.Achievements[i].ID == "first_feed" {
			tama.Achievements[i].Unlocked = true
		}
	}

	if err := Save(tama); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	for _, a := range loaded.Achievements {
		if a.ID == "first_feed" && !a.Unlocked {
			t.Error("achievement first_feed deve permanecer Unlocked após roundtrip")
		}
	}
}
