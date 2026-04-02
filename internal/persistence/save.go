package persistence

import (
	"encoding/json"
	"os"
	"time"

	"github.com/DeMart1n/Tamagotchi/internal/model"
)

const savePath = "tamago_save.json"

type SaveFile struct {
	Pet     model.Tama `json:"pet"`
	SavedAt time.Time  `json:"saved_at"`
	Version string     `json:"version"`
}

// Verifica se existe save
func Exists() bool {
	_, err := os.Stat(savePath)
	return err == nil
}

// Salva o jogo
func Save(tama *model.Tama) error {
	tama.LastSaved = time.Now()

	save := SaveFile{
		Pet:     *tama,
		SavedAt: time.Now(),
		Version: "1.0",
	}

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(savePath, data, 0644)
}

// Carrega o jogo
func Load() (*model.Tama, error) {
	data, err := os.ReadFile(savePath)
	if err != nil {
		return nil, err
	}

	var save SaveFile
	if err := json.Unmarshal(data, &save); err != nil {
		return nil, err
	}

	return &save.Pet, nil
}