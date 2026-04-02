package persistence

import (
	"Pessoal/internal/model"
	"encoding/json"
	"os"
	"time"
)

const SaveFile = "tamago_save.json"
const MaxOfflineTicks = 50

func Save(tama *model.Tama) error {
	tama.LastSaved = time.Now()
	data, err := json.MarshalIndent(tama, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SaveFile, data, 0644)
}

func Load() (*model.Tama, error) {
	data, err := os.ReadFile(SaveFile)
	if err != nil {
		return nil, err
	}

	var tama model.Tama
	if err := json.Unmarshal(data, &tama); err != nil {
		return nil, err
	}

	// Calcular degradação offline
	if !tama.LastSaved.IsZero() && !tama.Dead {
		elapsed := time.Since(tama.LastSaved)
		missedTicks := int(elapsed.Seconds()) / 20
		if missedTicks > MaxOfflineTicks {
			missedTicks = MaxOfflineTicks
		}
		for i := 0; i < missedTicks; i++ {
			applyOfflineTick(&tama)
		}
	}

	// Garantir que o Tama não esteja dormindo ao carregar
	tama.Sleeping = false
	normalizeBiome(&tama)

	return &tama, nil
}

func Exists() bool {
	_, err := os.Stat(SaveFile)
	return err == nil
}

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
