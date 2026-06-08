package dungeon

import "testing"

func TestBiomeEnemyPools(t *testing.T) {
	tests := []struct {
		name     string
		biome    Biome
		floor    int
		expected bool
	}{
		{"Forest floor 1", BiomeForest, 1, true},
		{"Forest floor 3", BiomeForest, 3, true},
		{"Icy floor 2", BiomeIcy, 2, true},
		{"Volcanic floor 2", BiomeVolcanic, 2, true},
		{"Abyssal floor 3", BiomeAbyssal, 3, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enemy := RandomEnemyForFloor(test.floor, 5, test.biome)
			if enemy == nil {
				t.Errorf("RandomEnemyForFloor returned nil for biome=%v floor=%d", test.biome, test.floor)
			}
		})
	}
}

func TestBossForBiome(t *testing.T) {
	bosses := map[Biome]string{
		BiomeForest:   "Guardiao da Floresta",
		BiomeIcy:      "Lich do Gelo Eterno",
		BiomeVolcanic: "Lorde das Chamas",
		BiomeAbyssal:  "Devorador do Abismo",
	}

	for biome, expectedName := range bosses {
		biomeName := ""
		switch biome {
		case BiomeForest:
			biomeName = "Forest"
		case BiomeIcy:
			biomeName = "Icy"
		case BiomeVolcanic:
			biomeName = "Volcanic"
		case BiomeAbyssal:
			biomeName = "Abyssal"
		}
		t.Run(biomeName, func(t *testing.T) {
			boss := BossForBiome(biome, 5)
			if boss == nil {
				t.Errorf("BossForBiome returned nil for biome=%v", biome)
			}
			if boss.Name != expectedName {
				t.Errorf("expected %q, got %q", expectedName, boss.Name)
			}
			if !boss.IsBoss {
				t.Errorf("boss should have IsBoss=true")
			}
		})
	}
}

func TestBossLootForBiome(t *testing.T) {
	tests := map[Biome]string{
		BiomeForest:   "espada_raiz_ancia",
		BiomeIcy:      "manto_lich",
		BiomeVolcanic: "brasa_do_lorde",
		BiomeAbyssal:  "lamina_abissal",
	}

	for biome, expectedID := range tests {
		biomeName := ""
		switch biome {
		case BiomeForest:
			biomeName = "Forest"
		case BiomeIcy:
			biomeName = "Icy"
		case BiomeVolcanic:
			biomeName = "Volcanic"
		case BiomeAbyssal:
			biomeName = "Abyssal"
		}
		t.Run(biomeName, func(t *testing.T) {
			loot := bossLootForBiome(biome)
			if loot == nil {
				t.Errorf("bossLootForBiome returned nil for biome=%v", biome)
			}
			if loot.ID != expectedID {
				t.Errorf("expected %q, got %q", expectedID, loot.ID)
			}
			if loot.Rarity != RarityLendario {
				t.Errorf("boss loot should be legendary, got %v", loot.Rarity)
			}
		})
	}
}
