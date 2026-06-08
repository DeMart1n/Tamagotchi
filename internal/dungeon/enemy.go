package dungeon

import "math/rand"

// EnemyTemplate define um modelo de inimigo (sem estado de combate).
type EnemyTemplate struct {
	ID      string
	Name    string
	HPMin   int
	HPMax   int
	ATKMin  int
	ATKMax  int
	DEF     int
	VEL     int
	XPDrop  int
	GoldMin int
	GoldMax int
	IsBoss  bool
	IsFire  bool
}

// Enemy é uma instância de inimigo em combate.
type Enemy struct {
	Name       string
	HPMax      int
	HPCurrent  int
	Ataque     int
	Defesa     int
	Velocidade int
	XPDrop     int
	GoldDrop   int
	IsBoss     bool
	Defending  bool
	IsFire     bool
}

// levelScale retorna o multiplicador de escala baseado no level do jogador.
// Level 3 (mínimo) = 1.0x, cada level acima adiciona 8%.
func levelScale(playerLevel int) float64 {
	if playerLevel <= 3 {
		return 1.0
	}
	return 1.0 + float64(playerLevel-3)*0.08
}

// SpawnEnemy cria uma instância de inimigo a partir de um template, escalado pelo level do jogador.
func SpawnEnemy(t EnemyTemplate, playerLevel int) *Enemy {
	scale := levelScale(playerLevel)

	hpBase := t.HPMin + rand.Intn(t.HPMax-t.HPMin+1)
	atkBase := t.ATKMin + rand.Intn(t.ATKMax-t.ATKMin+1)
	gold := t.GoldMin + rand.Intn(t.GoldMax-t.GoldMin+1)

	hp := int(float64(hpBase) * scale)
	atk := int(float64(atkBase) * scale)
	def := int(float64(t.DEF) * scale)
	vel := t.VEL + (playerLevel-3)/5
	xp := int(float64(t.XPDrop) * scale)

	if hp < 1 {
		hp = 1
	}
	if atk < 1 {
		atk = 1
	}

	return &Enemy{
		Name:       t.Name,
		HPMax:      hp,
		HPCurrent:  hp,
		Ataque:     atk,
		Defesa:     def,
		Velocidade: vel,
		XPDrop:     xp,
		GoldDrop:   gold,
		IsBoss:     t.IsBoss,
		IsFire:     t.IsFire,
	}
}

// FloorEnemyPools mapeia cada andar a seus possíveis inimigos.
var FloorEnemyPools = map[int][]EnemyTemplate{
	1: {
		{ID: "slime", Name: "Slime", HPMin: 20, HPMax: 25, ATKMin: 6, ATKMax: 8, DEF: 2, VEL: 3, XPDrop: 15, GoldMin: 5, GoldMax: 10},
		{ID: "rato_gigante", Name: "Rato Gigante", HPMin: 18, HPMax: 22, ATKMin: 7, ATKMax: 9, DEF: 1, VEL: 6, XPDrop: 18, GoldMin: 5, GoldMax: 12},
	},
	2: {
		{ID: "esqueleto", Name: "Esqueleto", HPMin: 30, HPMax: 38, ATKMin: 10, ATKMax: 12, DEF: 5, VEL: 4, XPDrop: 30, GoldMin: 10, GoldMax: 18},
		{ID: "morcego_vampiro", Name: "Morcego Vampiro", HPMin: 25, HPMax: 32, ATKMin: 11, ATKMax: 13, DEF: 3, VEL: 8, XPDrop: 28, GoldMin: 8, GoldMax: 15},
	},
	3: {
		{ID: "goblin_guerreiro", Name: "Goblin Guerreiro", HPMin: 45, HPMax: 55, ATKMin: 14, ATKMax: 16, DEF: 8, VEL: 6, XPDrop: 45, GoldMin: 15, GoldMax: 25},
		{ID: "aranha_venenosa", Name: "Aranha Venenosa", HPMin: 40, HPMax: 50, ATKMin: 15, ATKMax: 17, DEF: 6, VEL: 7, XPDrop: 42, GoldMin: 12, GoldMax: 22},
	},
	4: {
		{ID: "cavaleiro_negro", Name: "Cavaleiro Negro", HPMin: 60, HPMax: 75, ATKMin: 22, ATKMax: 26, DEF: 14, VEL: 5, XPDrop: 70, GoldMin: 25, GoldMax: 40},
		{ID: "mago_sombrio", Name: "Mago Sombrio", HPMin: 50, HPMax: 65, ATKMin: 25, ATKMax: 28, DEF: 8, VEL: 7, XPDrop: 75, GoldMin: 30, GoldMax: 45, IsFire: true},
	},
	5: {
		{ID: "dragao_anciao", Name: "Dragao Anciao", HPMin: 200, HPMax: 200, ATKMin: 35, ATKMax: 35, DEF: 18, VEL: 8, XPDrop: 200, GoldMin: 100, GoldMax: 150, IsBoss: true, IsFire: true},
	},
}

// Templates temáticos por bioma
var (
	// BiomeForest — floors 1-2
	druida_raizes = EnemyTemplate{ID: "druida_raizes", Name: "Druida das Raizes", HPMin: 22, HPMax: 28, ATKMin: 6, ATKMax: 9, DEF: 3, VEL: 4, XPDrop: 18, GoldMin: 5, GoldMax: 12}

	// BiomeForest — floors 3-4
	lobo_sombrio = EnemyTemplate{ID: "lobo_sombrio", Name: "Lobo Sombrio", HPMin: 48, HPMax: 60, ATKMin: 16, ATKMax: 19, DEF: 7, VEL: 9, XPDrop: 50, GoldMin: 18, GoldMax: 30}

	// BiomeIcy — floors 2-3
	espirito_gelo = EnemyTemplate{ID: "espirito_gelo", Name: "Espirito de Gelo", HPMin: 28, HPMax: 36, ATKMin: 10, ATKMax: 13, DEF: 4, VEL: 5, XPDrop: 30, GoldMin: 8, GoldMax: 16}

	// BiomeIcy — floors 3-4
	urso_glacial = EnemyTemplate{ID: "urso_glacial", Name: "Urso Glacial", HPMin: 55, HPMax: 68, ATKMin: 15, ATKMax: 18, DEF: 12, VEL: 3, XPDrop: 55, GoldMin: 20, GoldMax: 32}

	// BiomeVolcanic — floors 2-3
	elemental_magma = EnemyTemplate{ID: "elemental_magma", Name: "Elemental de Magma", HPMin: 38, HPMax: 48, ATKMin: 12, ATKMax: 15, DEF: 6, VEL: 5, XPDrop: 40, GoldMin: 14, GoldMax: 24, IsFire: true}

	// BiomeAbyssal — floors 3-4
	sombra_profunda = EnemyTemplate{ID: "sombra_profunda", Name: "Sombra Profunda", HPMin: 52, HPMax: 65, ATKMin: 18, ATKMax: 22, DEF: 9, VEL: 8, XPDrop: 60, GoldMin: 22, GoldMax: 38}

	// Bosses — um por bioma
	guardiao_floresta = EnemyTemplate{ID: "guardiao_floresta", Name: "Guardiao da Floresta", HPMin: 220, HPMax: 220, ATKMin: 28, ATKMax: 28, DEF: 16, VEL: 4, XPDrop: 220, GoldMin: 90, GoldMax: 130, IsBoss: true}

	lich_gelo = EnemyTemplate{ID: "lich_gelo", Name: "Lich do Gelo Eterno", HPMin: 200, HPMax: 200, ATKMin: 30, ATKMax: 30, DEF: 12, VEL: 6, XPDrop: 210, GoldMin: 95, GoldMax: 135, IsBoss: true}

	lorde_chamas = EnemyTemplate{ID: "lorde_chamas", Name: "Lorde das Chamas", HPMin: 230, HPMax: 230, ATKMin: 32, ATKMax: 32, DEF: 14, VEL: 7, XPDrop: 230, GoldMin: 100, GoldMax: 145, IsBoss: true, IsFire: true}

	devorador_abismo = EnemyTemplate{ID: "devorador_abismo", Name: "Devorador do Abismo", HPMin: 240, HPMax: 240, ATKMin: 34, ATKMax: 34, DEF: 18, VEL: 9, XPDrop: 240, GoldMin: 110, GoldMax: 160, IsBoss: true}
)

// BiomeEnemyEntry uma entrada no pool de inimigos temáticos por bioma.
type BiomeEnemyEntry struct {
	Template EnemyTemplate
	FloorMin int
	FloorMax int
}

// BiomeEnemyPools mapeia biomas a seus inimigos temáticos.
var BiomeEnemyPools = map[Biome][]BiomeEnemyEntry{
	BiomeForest: {
		{Template: druida_raizes, FloorMin: 1, FloorMax: 2},
		{Template: lobo_sombrio, FloorMin: 3, FloorMax: 4},
	},
	BiomeIcy: {
		{Template: espirito_gelo, FloorMin: 2, FloorMax: 3},
		{Template: urso_glacial, FloorMin: 3, FloorMax: 4},
	},
	BiomeVolcanic: {
		{Template: elemental_magma, FloorMin: 2, FloorMax: 3},
	},
	BiomeAbyssal: {
		{Template: sombra_profunda, FloorMin: 3, FloorMax: 4},
	},
}

// RandomEnemyForFloor sorteia um inimigo escalado pelo nível, preferindo inimigos temáticos do bioma.
func RandomEnemyForFloor(floor, playerLevel int, biome Biome) *Enemy {
	// Tenta encontrar um inimigo temático do bioma
	if entries, ok := BiomeEnemyPools[biome]; ok {
		var candidates []EnemyTemplate
		for _, entry := range entries {
			if floor >= entry.FloorMin && floor <= entry.FloorMax {
				candidates = append(candidates, entry.Template)
			}
		}
		if len(candidates) > 0 {
			return SpawnEnemy(candidates[rand.Intn(len(candidates))], playerLevel)
		}
	}

	// Fallback para o pool genérico do andar
	pool, ok := FloorEnemyPools[floor]
	if !ok {
		pool = FloorEnemyPools[1]
	}
	template := pool[rand.Intn(len(pool))]
	return SpawnEnemy(template, playerLevel)
}

// BossForBiome retorna o boss específico do bioma, escalado pelo level.
func BossForBiome(biome Biome, playerLevel int) *Enemy {
	switch biome {
	case BiomeForest:
		return SpawnEnemy(guardiao_floresta, playerLevel)
	case BiomeIcy:
		return SpawnEnemy(lich_gelo, playerLevel)
	case BiomeVolcanic:
		return SpawnEnemy(lorde_chamas, playerLevel)
	case BiomeAbyssal:
		return SpawnEnemy(devorador_abismo, playerLevel)
	default:
		return SpawnEnemy(guardiao_floresta, playerLevel)
	}
}

// BossForFloor5 alias de BossForBiome(BiomeForest) para compatibilidade.
func BossForFloor5(playerLevel int) *Enemy {
	return BossForBiome(BiomeForest, playerLevel)
}
