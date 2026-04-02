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
	combat := NewCombat(player, enemy, BiomeVolcanic)

	combat.ExecuteAction(ActionDefender, nil, 0)

	if combat.Player.HPCurrent >= 95 {
		t.Fatalf("expected volcanic DOT to reduce player HP significantly, got %d", combat.Player.HPCurrent)
	}
}

func TestCombat_ForestFireVulnerabilityIncreasesDamage(t *testing.T) {
	// Usamos um mock para garantir que não haja sorte envolvida e o dano seja previsível
	basePlayer := CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 10, Sorte: 0}
	baseEnemy := Enemy{Name: "Fire Mage", HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 0, Velocidade: 5, IsFire: true}

	// Forest tem +25% de vulnerabilidade a fogo
	forestPlayer := basePlayer
	forestEnemy := baseEnemy
	forestCombat := NewCombat(&forestPlayer, &forestEnemy, BiomeForest)
	
	// Neutral (biome sem efeitos de combate/sorte)
	neutralPlayer := basePlayer
	neutralEnemy := baseEnemy
	neutralCombat := NewCombat(&neutralPlayer, &neutralEnemy, Biome(999)) 

	// Forçamos o mesmo seed para ambos para comparar a variância
	// Nota: Como o combate usa o rand global, vamos setar o seed global antes de cada chamada
	
	rand.Seed(42)
	forestCombat.enemyTurn()
	forestDamage := 100 - forestCombat.Player.HPCurrent

	rand.Seed(42)
	neutralCombat.enemyTurn()
	neutralDamage := 100 - neutralCombat.Player.HPCurrent

	if forestDamage <= neutralDamage {
		t.Fatalf("expected forest fire vulnerability to increase damage: forest=%d neutral=%d", forestDamage, neutralDamage)
	}
}
