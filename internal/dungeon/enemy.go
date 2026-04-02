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

// RandomEnemyForFloor sorteia um inimigo do pool de um andar, escalado pelo level.
func RandomEnemyForFloor(floor, playerLevel int, biome Biome) *Enemy {
	floorPool, ok := FloorEnemyPools[floor]
	if !ok {
		floorPool = FloorEnemyPools[1]
	}

	biomePool := BiomeEnemyPools[biome]

	// 40% de chance de um inimigo do bioma, se disponível.
	if rand.Intn(100) < 40 && len(biomePool) > 0 {
		template := biomePool[rand.Intn(len(biomePool))]
		return SpawnEnemy(template, playerLevel)
	}

	template := floorPool[rand.Intn(len(floorPool))]
	return SpawnEnemy(template, playerLevel)
}

// BiomeEnemyPools mapeia biomas para seus inimigos temáticos.
var BiomeEnemyPools = map[Biome][]EnemyTemplate{
	BiomeForest: {
		{ID: "ent_jovem", Name: "Ent Jovem", HPMin: 35, HPMax: 45, ATKMin: 8, ATKMax: 10, DEF: 6, VEL: 2, XPDrop: 25, GoldMin: 8, GoldMax: 15},
		{ID: "lobo_florestal", Name: "Lobo Florestal", HPMin: 25, HPMax: 30, ATKMin: 12, ATKMax: 15, DEF: 3, VEL: 8, XPDrop: 22, GoldMin: 10, GoldMax: 12},
		{ID: "vespa_gigante", Name: "Vespa Gigante", HPMin: 20, HPMax: 25, ATKMin: 10, ATKMax: 12, DEF: 1, VEL: 10, XPDrop: 20, GoldMin: 6, GoldMax: 12},
	},
	BiomeIcy: {
		{ID: "yeti_pequeno", Name: "Yeti Pequeno", HPMin: 50, HPMax: 60, ATKMin: 15, ATKMax: 18, DEF: 10, VEL: 3, XPDrop: 40, GoldMin: 15, GoldMax: 20},
		{ID: "elementar_gelo", Name: "Elementar de Gelo", HPMin: 40, HPMax: 50, ATKMin: 14, ATKMax: 16, DEF: 12, VEL: 4, XPDrop: 38, GoldMin: 12, GoldMax: 18},
		{ID: "lobo_neves", Name: "Lobo das Neves", HPMin: 30, HPMax: 35, ATKMin: 16, ATKMax: 20, DEF: 4, VEL: 9, XPDrop: 35, GoldMin: 10, GoldMax: 15},
	},
	BiomeVolcanic: {
		{ID: "salamandra_fogo", Name: "Salamandra de Fogo", HPMin: 35, HPMax: 45, ATKMin: 18, ATKMax: 22, DEF: 5, VEL: 6, XPDrop: 45, GoldMin: 20, GoldMax: 30, IsFire: true},
		{ID: "golem_magma", Name: "Golem de Magma", HPMin: 80, HPMax: 100, ATKMin: 20, ATKMax: 25, DEF: 15, VEL: 2, XPDrop: 60, GoldMin: 25, GoldMax: 40, IsFire: true},
		{ID: "diabrete", Name: "Diabrete", HPMin: 40, HPMax: 50, ATKMin: 22, ATKMax: 26, DEF: 6, VEL: 8, XPDrop: 55, GoldMin: 30, GoldMax: 45, IsFire: true},
	},
	BiomeAbyssal: {
		{ID: "olho_flutuante", Name: "Olho Flutuante", HPMin: 30, HPMax: 40, ATKMin: 20, ATKMax: 25, DEF: 4, VEL: 7, XPDrop: 50, GoldMin: 30, GoldMax: 50},
		{ID: "sombra_faminta", Name: "Sombra Faminta", HPMin: 45, HPMax: 55, ATKMin: 22, ATKMax: 28, DEF: 8, VEL: 9, XPDrop: 55, GoldMin: 35, GoldMax: 55},
		{ID: "horror_sem_rosto", Name: "Horror Sem Rosto", HPMin: 100, HPMax: 120, ATKMin: 25, ATKMax: 30, DEF: 18, VEL: 5, XPDrop: 100, GoldMin: 50, GoldMax: 100},
	},
}

// BossForFloor5 retorna o boss do andar 5, escalado pelo level.
func BossForFloor5(playerLevel int) *Enemy {
	return SpawnEnemy(FloorEnemyPools[5][0], playerLevel)
}
