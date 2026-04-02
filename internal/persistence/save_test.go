package persistence

import (
	"Pessoal/internal/model"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// --- helpers de teste ---

func tmpSaveFile(t *testing.T) func() {
	t.Helper()
	original := SaveFile
	// Redireciona para arquivo temporário para não poluir o diretório de trabalho
	f, err := os.CreateTemp("", "tamago_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()

	// Monkey-patch da constante via variável de pacote não é possível em Go;
	// usamos a variável exportável abaixo (ver nota no save.go).
	// Por ora os testes que escrevem em disco usam o arquivo padrão e limpam depois.
	_ = original
	return func() { os.Remove(SaveFile); os.Remove(SaveFile + ".tmp") }
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

// --- testes de bioma (existentes, mantidos) ---

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

// --- testes de roundtrip ---

func TestSaveAndLoadRoundtrip(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

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

func TestSaveAndLoadInventory(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()

	// Simula inventário serializado (como dungeon o faria)
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

	// Verifica que o campo também está no envelope (campo topo)
	data, _ := os.ReadFile(SaveFile)
	var envelope SaveData
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if string(envelope.Inventory) != string(raw) {
		t.Errorf("envelope.Inventory mismatch:\n  got  %s\n  want %s", envelope.Inventory, raw)
	}
}

// --- testes de schema version ---

func TestSaveWritesSchemaVersion(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

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

func TestLoadLegacySave(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	// Simula save legado: JSON raiz direto é um model.Tama (sem envelope)
	legacy := baseTama()
	legacy.LastSaved = time.Now() // evita degradação offline no teste
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

// --- testes de degradação offline ---

func TestOfflineTickDegradation(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()
	// Simula save feito 100 segundos atrás (5 ticks de 20s)
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

	// 5 ticks = 5 de degradação em cada stat
	expectedHunger := 80 - 5
	if loaded.Hunger != expectedHunger {
		t.Errorf("Hunger after offline ticks: got %d, want %d", loaded.Hunger, expectedHunger)
	}
}

func TestOfflineTickCappedAtMax(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()
	// Simula save muito antigo (mais de MaxOfflineTicks ticks atrás)
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
	// Não importa quanto tempo passou, máximo é MaxOfflineTicks de degradação
	minExpected := 100 - MaxOfflineTicks
	if loaded.Hunger < minExpected {
		t.Errorf("Hunger degraded beyond cap: got %d (min expected %d)", loaded.Hunger, minExpected)
	}
}

func TestOfflineTickKillsTamaOnStarvation(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()
	tama.LastSaved = time.Now().Add(-time.Duration(MaxOfflineTicks*20+9999) * time.Second)
	tama.Hunger = 1 // vai a 0 após os ticks
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
		t.Error("Tama should be Dead after starving offline")
	}
}

// --- testes de Delete ---

func TestDelete(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	if err := Save(baseTama()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if !Exists() {
		t.Fatal("file should exist before Delete")
	}
	if err := Delete(); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if Exists() {
		t.Error("file should not exist after Delete")
	}
	// Chamar Delete de novo não deve retornar erro
	if err := Delete(); err != nil {
		t.Errorf("second Delete should be a no-op, got: %v", err)
	}
}

// --- testes de escrita atômica ---

func TestSaveIsAtomic(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()
	if err := Save(tama); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	// Arquivo .tmp não deve restar após Save bem-sucedido
	if _, err := os.Stat(SaveFile + ".tmp"); !os.IsNotExist(err) {
		t.Error(".tmp file should not exist after successful Save")
	}
}

// --- testes de achievements preservados ---

func TestAchievementsPreservedOnRoundtrip(t *testing.T) {
	cleanup := tmpSaveFile(t)
	defer cleanup()

	tama := baseTama()
	// Desbloqueia uma conquista
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
			t.Error("achievement first_feed should remain Unlocked after roundtrip")
		}
	}
}