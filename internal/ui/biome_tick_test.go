package ui

import (
	"Pessoal/internal/model"
	"testing"
)

func TestTickBiomeModifiers(t *testing.T) {
	// Base Tama
	baseTama := &model.Tama{
		Hunger:    50,
		Sleepy:    50,
		Thirst:    50,
		Happiness: 50,
		Stage:     model.StageTeen, // DecayRate = 1
	}

	// Floresta: HungerRateMult = 0.8, HappinessRateMult = 1.2
	forestTama := *baseTama
	forestTama.CurrentBiome = "Florestal"
	Tick(&forestTama)
	
	// Hunger decay: 1 * 0.8 = 0.8 -> 1 (min)
	// Happiness decay: 1 / 1.2 = 0.83 -> 1 (min)
	// Thirst decay: 1 * 0.8 = 0.8 -> 1 (min)
	if forestTama.Hunger != 49 {
		t.Errorf("Forest Hunger expected 49, got %d", forestTama.Hunger)
	}

	// Gélido: HungerRateMult = 1.5, SleepyRateMult = 1.3
	icyTama := *baseTama
	icyTama.CurrentBiome = "Gelido"
	Tick(&icyTama)
	// Hunger decay: 1 * 1.5 = 1
	if icyTama.Hunger != 49 { // Wait, 1.5 becomes 1? int(1.5) = 1.
		// Let's use a stage with higher decay to see the multiplier better
	}
	
	// Usando Adulto (DecayRate = 2)
	adultTama := &model.Tama{
		Hunger:    50,
		Sleepy:    50,
		Thirst:    50,
		Happiness: 50,
		Stage:     model.StageAdult, // DecayRate = 2
	}
	
	// Gélido Adulto:
	// Hunger decay: 2 * 1.5 = 3
	// Sleepy decay: 2 * 1.3 = 2
	icyAdult := *adultTama
	icyAdult.CurrentBiome = "Gelido"
	Tick(&icyAdult)
	if icyAdult.Hunger != 47 {
		t.Errorf("Icy Adult Hunger expected 47 (50-3), got %d", icyAdult.Hunger)
	}
	if icyAdult.Sleepy != 48 {
		t.Errorf("Icy Adult Sleepy expected 48 (50-2), got %d", icyAdult.Sleepy)
	}

	// Vulcânico Adulto:
	// HappinessPenaltyTick = 1
	// HappinessRateMult = 0.7
	// Happiness decay: (2 / 0.7) + 1 = 2 + 1 = 3
	volcAdult := *adultTama
	volcAdult.CurrentBiome = "Vulcanico"
	Tick(&volcAdult)
	if volcAdult.Happiness != 47 {
		t.Errorf("Volcanic Adult Happiness expected 47 (50-3), got %d", volcAdult.Happiness)
	}
}
