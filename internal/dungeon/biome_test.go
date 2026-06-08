package dungeon

import (
	"math/rand"
	"testing"
)

func TestModifierForBiome_HasExpectedPrimaryEffects(t *testing.T) {
	forest := ModifierForBiome(BiomeForest)
	if forest.HungerRecoveryBonus <= 0 {
		t.Fatalf("expected forest hunger recovery bonus, got %d", forest.HungerRecoveryBonus)
	}
	if forest.PlayerFireVulnerability <= 0 {
		t.Fatalf("expected forest fire vulnerability, got %d", forest.PlayerFireVulnerability)
	}

	icy := ModifierForBiome(BiomeIcy)
	if icy.PlayerSpeedPenaltyPct <= 0 {
		t.Fatalf("expected icy speed penalty, got %d", icy.PlayerSpeedPenaltyPct)
	}
	if icy.FreezeChancePct <= 0 {
		t.Fatalf("expected icy freeze chance, got %d", icy.FreezeChancePct)
	}

	volcanic := ModifierForBiome(BiomeVolcanic)
	if volcanic.EnemyAttackBonusPct <= 0 || volcanic.PlayerFireDotPctMaxHP <= 0 {
		t.Fatalf("expected volcanic combat effects, got %+v", volcanic)
	}

	abyssal := ModifierForBiome(BiomeAbyssal)
	if abyssal.EnemyLuckBonusPct <= 0 || abyssal.RareLootChanceBonusPct <= 0 {
		t.Fatalf("expected abyssal luck/loot effects, got %+v", abyssal)
	}
}

func TestCombat_VolcanicAppliesDot(t *testing.T) {
	player := &CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 10, Velocidade: 10, Sorte: 0}
	enemy := &Enemy{Name: "Dummy", HPMax: 100, HPCurrent: 100, Ataque: 1, Defesa: 0, Velocidade: 1}
	combat := NewCombat(player, enemy, BiomeVolcanic, nil)

	combat.ExecuteAction(ActionDefender, nil, nil, 0)

	if combat.Player.HPCurrent >= 100 {
		t.Fatalf("expected volcanic DOT or degen to reduce player HP, got %d", combat.Player.HPCurrent)
	}
}

func TestCombat_ForestFireVulnerabilityIncreasesDamage(t *testing.T) {
	player := CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 10, Sorte: 0}
	fireEnemy := Enemy{Name: "Fire Mage", HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 5, IsFire: true}

	// Floresta tem PlayerFireVulnerability: 25 e Regen: +2
	forestCombat := NewCombat(&player, &fireEnemy, BiomeForest, nil)
	
	// Abissal NÃO tem vulnerabilidade a fogo e Regen: 0
	player2 := CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 10, Sorte: 0}
	abyssalCombat := NewCombat(&player2, &fireEnemy, BiomeAbyssal, nil)
	
	rand.Seed(7)
	forestCombat.ExecuteAction(ActionDefender, nil, nil, 0)
	fireDamage := forestCombat.Player.HPMax - forestCombat.Player.HPCurrent

	rand.Seed(7)
	abyssalCombat.ExecuteAction(ActionDefender, nil, nil, 0)
	abyssalDamage := abyssalCombat.Player.HPMax - abyssalCombat.Player.HPCurrent

	// Na floresta: Dano fogo (12 + 3 = 15) - Regen (2) = 13
	// No abissal: Dano normal (12) - Regen (0) = 12
	if fireDamage <= abyssalDamage {
		t.Errorf("expected forest fire vulnerability to increase damage: forest=%d abyssal=%d", fireDamage, abyssalDamage)
	}
}
