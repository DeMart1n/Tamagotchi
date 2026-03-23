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
		}
	case BiomeIcy:
		return BiomeModifier{
			PlayerSpeedPenaltyPct: 20,
			PlayerIceDefensePct:   20,
			FreezeChancePct:       20,
		}
	case BiomeVolcanic:
		return BiomeModifier{
			EnemyAttackBonusPct:   20,
			PlayerFireDotPctMaxHP: 5,
			HappinessPenaltyTick:  1,
		}
	case BiomeAbyssal:
		return BiomeModifier{
			EnemyLuckBonusPct:       10,
			RareLootChanceBonusPct:  20,
			PlayerHPRegenPenaltyPct: 50,
		}
	default:
		return BiomeModifier{}
	}
}