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
	pool, ok := FloorEnemyPools[floor]
	if !ok {
		pool = FloorEnemyPools[1]
	}
	_ = biome
	template := pool[rand.Intn(len(pool))]
	return SpawnEnemy(template, playerLevel)
}

// BossForFloor5 retorna o boss do andar 5, escalado pelo level.
func BossForFloor5(playerLevel int) *Enemy {
	return SpawnEnemy(FloorEnemyPools[5][0], playerLevel)
}
