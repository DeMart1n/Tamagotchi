package persistence

import (
	"Pessoal/internal/model"
	"testing"
)

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
