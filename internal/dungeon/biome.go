package dungeon

import "math/rand"

type Biome int

const (
	BiomeForest Biome = iota
	BiomeIcy
	BiomeVolcanic
	BiomeAbyssal
)

type BiomeModifier struct {
	HungerRecoveryBonus     int
	FriendlyEncounterChance int
	PlayerFireVulnerability int
	PlayerSpeedPenaltyPct   int
	PlayerIceDefensePct     int
	FreezeChancePct         int
	EnemyAttackBonusPct     int
	PlayerFireDotPctMaxHP   int
	HappinessPenaltyTick    int
	EnemyLuckBonusPct       int
	RareLootChanceBonusPct  int
	PlayerHPRegenPenaltyPct int

	// Novos modificadores solicitados
	PlayerHPBonusPct    int
	PlayerATKBonusPct   int
	PlayerDEFBonusPct   int
	PlayerVELBonusPct   int
	PlayerLuckBonusPct  int
	PlayerHPRegenBonus  int // Flat regen/degen por turno em combate
	HungerRateMult      float64
	SleepyRateMult      float64
	HappinessRateMult   float64
}

func (b Biome) String() string {
	switch b {
	case BiomeForest:
		return "Florestal"
	case BiomeIcy:
		return "Gelido"
	case BiomeVolcanic:
		return "Vulcanico"
	case BiomeAbyssal:
		return "Abissal"
	default:
		return "Florestal"
	}
}

func ParseBiome(name string) Biome {
	switch name {
	case "Florestal":
		return BiomeForest
	case "Gelido":
		return BiomeIcy
	case "Vulcanico":
		return BiomeVolcanic
	case "Abissal":
		return BiomeAbyssal
	default:
		return BiomeForest
	}
}

func RandomBiome() Biome {
	return Biome(rand.Intn(4))
}

func ModifierForBiome(b Biome) BiomeModifier {
	switch b {
	case BiomeForest:
		return BiomeModifier{
			HungerRecoveryBonus:     1,
			FriendlyEncounterChance: 10,
			PlayerFireVulnerability: 25,
			PlayerHPBonusPct:        10,
			PlayerHPRegenBonus:      2,
			HungerRateMult:          0.8, // Fome demora mais
			SleepyRateMult:          1.0,
			HappinessRateMult:       1.2, // Mais feliz na floresta
		}
	case BiomeIcy:
		return BiomeModifier{
			PlayerSpeedPenaltyPct: 20,
			PlayerIceDefensePct:   20,
			FreezeChancePct:       20,
			PlayerDEFBonusPct:     10,
			PlayerVELBonusPct:     -10,
			HungerRateMult:        1.5, // Frio dá fome
			SleepyRateMult:        1.3, // Frio dá sono
			HappinessRateMult:     1.0,
		}
	case BiomeVolcanic:
		return BiomeModifier{
			EnemyAttackBonusPct:   20,
			PlayerFireDotPctMaxHP: 5,
			HappinessPenaltyTick:  1,
			PlayerATKBonusPct:     15,
			PlayerHPRegenBonus:    -1,  // Calor drena vida
			HungerRateMult:        1.2,
			SleepyRateMult:        1.5, // Calor cansa
			HappinessRateMult:     0.7, // Estressante
		}
	case BiomeAbyssal:
		return BiomeModifier{
			EnemyLuckBonusPct:       10,
			RareLootChanceBonusPct:  20,
			PlayerHPRegenPenaltyPct: 50,
			PlayerLuckBonusPct:      15,
			PlayerDEFBonusPct:       -10,
			HungerRateMult:          1.3,
			SleepyRateMult:          1.2,
			HappinessRateMult:       0.5, // Muito depressivo
		}
	default:
		return BiomeModifier{
			HungerRateMult:    1.0,
			SleepyRateMult:    1.0,
			HappinessRateMult: 1.0,
		}
	}
}