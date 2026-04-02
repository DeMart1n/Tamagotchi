package persistence

import (
	"Pessoal/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const SaveFile = "tamago_save.json"
const MaxOfflineTicks = 50

// CurrentSchemaVersion é incrementado toda vez que o formato do save mudar.
// Permite migrações automáticas de saves antigos.
const CurrentSchemaVersion = 1

// SaveData é o envelope raiz do arquivo tamago_save.json.
// Separa metadados de persistência (versão, timestamp) dos dados do jogo,
// permitindo migrações sem perda de dados.
type SaveData struct {
	SchemaVersion int             `json:"schema_version"`
	SavedAt       time.Time       `json:"saved_at"`
	Tama          *model.Tama     `json:"tama"`
	Inventory     json.RawMessage `json:"inventory,omitempty"`
}

// Save serializa o estado completo do jogo em tamago_save.json.
// O inventário já vem embutido em tama.Inventory como json.RawMessage,
// então é promovido para o campo de topo do envelope para leitura mais clara.
func Save(tama *model.Tama) error {
	now := time.Now()
	tama.LastSaved = now

	envelope := SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       now,
		Tama:          tama,
		Inventory:     tama.Inventory, // promove para topo — redundante mas explícito
	}

	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return fmt.Errorf("save: marshal failed: %w", err)
	}

	// Escrita atômica: grava em .tmp e renomeia para evitar save corrompido
	tmp := SaveFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("save: write tmp failed: %w", err)
	}
	if err := os.Rename(tmp, SaveFile); err != nil {
		return fmt.Errorf("save: rename failed: %w", err)
	}

	return nil
}

// Load lê tamago_save.json, aplica migrações de schema se necessário,
// calcula degradação offline e retorna o Tama pronto para uso.
func Load() (*model.Tama, error) {
	data, err := os.ReadFile(SaveFile)
	if err != nil {
		return nil, fmt.Errorf("load: read failed: %w", err)
	}

	// Tenta decodificar como envelope novo (schema_version presente)
	var envelope SaveData
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("load: unmarshal failed: %w", err)
	}

	// Save legado: sem schema_version, o JSON raiz é direto um *model.Tama
	if envelope.SchemaVersion == 0 && envelope.Tama == nil {
		envelope, err = migrateLegacySave(data)
		if err != nil {
			return nil, fmt.Errorf("load: legacy migration failed: %w", err)
		}
	}

	tama := envelope.Tama
	if tama == nil {
		return nil, fmt.Errorf("load: tama is nil after parsing")
	}

	// Reconcilia inventory: se Tama.Inventory estiver vazio mas o envelope tiver, usa o do envelope
	if len(tama.Inventory) == 0 && len(envelope.Inventory) > 0 {
		tama.Inventory = envelope.Inventory
	}

	// Aplica migrações de schema (reservado para versões futuras)
	if err := applyMigrations(envelope.SchemaVersion, tama); err != nil {
		return nil, fmt.Errorf("load: schema migration failed: %w", err)
	}

	// Calcula degradação offline
	if !tama.LastSaved.IsZero() && !tama.Dead {
		elapsed := time.Since(tama.LastSaved)
		missedTicks := int(elapsed.Seconds()) / 20
		if missedTicks > MaxOfflineTicks {
			missedTicks = MaxOfflineTicks
		}
		for i := 0; i < missedTicks; i++ {
			applyOfflineTick(tama)
		}
	}

	// Invariantes pós-load
	tama.Sleeping = false
	normalizeBiome(tama)

	return tama, nil
}

// Exists reporta se há um arquivo de save válido em disco.
func Exists() bool {
	_, err := os.Stat(SaveFile)
	return err == nil
}

// Delete remove o arquivo de save (usado no Game Over para reset).
func Delete() error {
	if err := os.Remove(SaveFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete save: %w", err)
	}
	return nil
}

// --- migração de saves legados ---

// migrateLegacySave converte um save antigo (JSON raiz = Tama, sem envelope)
// para o formato SaveData atual.
func migrateLegacySave(raw []byte) (SaveData, error) {
	var tama model.Tama
	if err := json.Unmarshal(raw, &tama); err != nil {
		return SaveData{}, fmt.Errorf("legacy unmarshal: %w", err)
	}
	return SaveData{
		SchemaVersion: CurrentSchemaVersion,
		SavedAt:       tama.LastSaved,
		Tama:          &tama,
		Inventory:     tama.Inventory,
	}, nil
}

// applyMigrations aplica transformações necessárias de acordo com a versão do schema.
// Adicione cases aqui conforme o schema evoluir.
func applyMigrations(fromVersion int, tama *model.Tama) error {
	// v0 → v1: não há transformação de dados necessária além do wrap do envelope
	// (placeholder para versões futuras)
	_ = fromVersion
	_ = tama
	return nil
}

// --- helpers internos ---

func applyOfflineTick(tama *model.Tama) {
	if tama.Dead {
		return
	}
	if tama.Hunger > 0 {
		tama.Hunger--
	}
	if tama.Sleepy > 0 {
		tama.Sleepy--
	}
	if tama.Thirst > 0 {
		tama.Thirst--
	}
	if tama.Happiness > 0 {
		tama.Happiness--
	}
	if tama.Angry > 0 {
		tama.Angry--
	}
	if tama.Hunger <= 0 || tama.Thirst <= 0 {
		tama.Dead = true
	}
	if tama.Happiness <= 0 {
		tama.Depressed = true
	}
}

func normalizeBiome(tama *model.Tama) {
	switch tama.CurrentBiome {
	case "Florestal", "Gelido", "Vulcanico", "Abissal":
		return
	default:
		tama.CurrentBiome = "Florestal"
	}
}