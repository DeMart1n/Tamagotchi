package dungeon

import (
	"testing"
)

func TestBiomeStatsModifiers(t *testing.T) {
	stats := CombatStats{
		HPMax:      100,
		HPCurrent:  100,
		Ataque:     10,
		Defesa:     10,
		Velocidade: 10,
		Sorte:      10,
	}

	// Teste Floresta: HP +10%
	forestMod := ModifierForBiome(BiomeForest)
	forestStats := stats.ApplyBiomeModifiers(forestMod)
	if forestStats.HPMax != 110 {
		t.Errorf("Forest HPMax expected 110, got %d", forestStats.HPMax)
	}

	// Teste Vulcânico: ATK +15%
	volcMod := ModifierForBiome(BiomeVolcanic)
	volcStats := stats.ApplyBiomeModifiers(volcMod)
	if volcStats.Ataque != 11 { // 10 + 1.5 = 11.5 -> 11 (int)
		t.Errorf("Volcanic Ataque expected 11, got %d", volcStats.Ataque)
	}

	// Teste Gélido: DEF +10%, VEL -10%
	icyMod := ModifierForBiome(BiomeIcy)
	icyStats := stats.ApplyBiomeModifiers(icyMod)
	if icyStats.Defesa != 11 {
		t.Errorf("Icy Defesa expected 11, got %d", icyStats.Defesa)
	}
	if icyStats.Velocidade != 9 {
		t.Errorf("Icy Velocidade expected 9, got %d", icyStats.Velocidade)
	}

	// Teste Abissal: Luck +15%, DEF -10%
	abyssalMod := ModifierForBiome(BiomeAbyssal)
	abyssalStats := stats.ApplyBiomeModifiers(abyssalMod)
	if abyssalStats.Sorte != 11 {
		t.Errorf("Abyssal Sorte expected 11, got %d", abyssalStats.Sorte)
	}
	if abyssalStats.Defesa != 9 {
		t.Errorf("Abyssal Defesa expected 9, got %d", abyssalStats.Defesa)
	}
}

func TestCombatRegenDegen(t *testing.T) {
	player := &CombatStats{HPMax: 100, HPCurrent: 50, Ataque: 10, Defesa: 10}
	enemy := &Enemy{Name: "Dummy", HPMax: 100, HPCurrent: 100, Ataque: 1, Defesa: 0}
	
	// Floresta tem regen +2
	combat := NewCombat(player, enemy, BiomeForest)
	// NewCombat aplica modificadores, então HPMax vira 110, HPCurrent vira 60 (50 + 10 bônus)
	// Mas eu quero testar o regen por turno
	initialHP := combat.Player.HPCurrent
	combat.applyBiomeTurnEffects()
	if combat.Player.HPCurrent != initialHP+2 {
		t.Errorf("Forest regen expected %d, got %d", initialHP+2, combat.Player.HPCurrent)
	}

	// Vulcânico tem degen -1 (além do DOT de fogo se houver)
	player2 := &CombatStats{HPMax: 100, HPCurrent: 100, Ataque: 10, Defesa: 10}
	combatVolc := NewCombat(player2, enemy, BiomeVolcanic)
	initialHP2 := combatVolc.Player.HPCurrent
	combatVolc.applyBiomeTurnEffects()
	// VolcMod: PlayerFireDotPctMaxHP: 5 -> 5 damage
	// PlayerHPRegenBonus: -1 -> 1 damage
	// Total damage expected: 6
	expectedHP := initialHP2 - 6
	if combatVolc.Player.HPCurrent != expectedHP {
		t.Errorf("Volcanic degen expected %d, got %d", expectedHP, combatVolc.Player.HPCurrent)
	}
}

func TestBiomeDecayRates(t *testing.T) {
	// Este teste precisaria estar no pacote ui ou exportar o Tick
	// Mas podemos testar indiretamente se movermos a lógica de cálculo de decay para uma função testável
}
