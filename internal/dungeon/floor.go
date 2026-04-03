package dungeon

import "math/rand"

// RoomType tipo de sala na dungeon.
type RoomType int

const (
	RoomCombat   RoomType = iota
	RoomRest              // Descanso + loja
	RoomTreasure          // Tesouro
	RoomBoss              // Boss
)

func (r RoomType) String() string {
	switch r {
	case RoomCombat:
		return "Combate"
	case RoomRest:
		return "Descanso"
	case RoomTreasure:
		return "Tesouro"
	case RoomBoss:
		return "BOSS"
	default:
		return "???"
	}
}

func (r RoomType) Icon() string {
	switch r {
	case RoomCombat:
		return "[#]"
	case RoomRest:
		return "[R]"
	case RoomTreasure:
		return "[T]"
	case RoomBoss:
		return "[B]"
	default:
		return "[ ]"
	}
}

// Room é uma sala individual dentro de um andar.
type Room struct {
	Type        RoomType
	Name        string // Nome temático (ex: "Acampamento de Druidas")
	Description string // Descrição breve
	Cleared     bool
	Enemy       *Enemy     // Para salas de combate/boss
	Loot        *Equipment // Para salas de tesouro
	LootGold    int        // Ouro em salas de tesouro
}

// Floor é um andar completo da dungeon.
type Floor struct {
	Number      int
	Biome       Biome
	Rooms       []Room
	CurrentRoom int
}

// GenerateFloor gera as salas de um andar, com inimigos escalados pelo level do jogador.
func GenerateFloor(floorNum, playerLevel int, biome Biome) *Floor {
	var rooms []Room

	// Layout por bioma (Andar 5 é sempre BOSS no final)
	if floorNum == 5 {
		rooms = generateBossFloor(playerLevel, biome)
	} else {
		rooms = generateStandardFloor(floorNum, playerLevel, biome)
	}

	return &Floor{
		Number:      floorNum,
		Biome:       biome,
		Rooms:       rooms,
		CurrentRoom: 0,
	}
}

func generateBossFloor(playerLevel int, biome Biome) []Room {
	var rooms []Room
	restRoom := Room{Type: RoomRest}
	
	switch biome {
	case BiomeForest:
		restRoom.Name = "Clareira de Descanso"
	case BiomeIcy:
		restRoom.Name = "Caverna Termal"
	case BiomeVolcanic:
		restRoom.Name = "Forja de Obsidiana"
	case BiomeAbyssal:
		restRoom.Name = "Santuário Esquecido"
	}

	rooms = []Room{
		{Type: RoomCombat, Enemy: RandomEnemyForFloor(4, playerLevel, biome)},
		restRoom,
		{Type: RoomCombat, Enemy: RandomEnemyForFloor(4, playerLevel, biome)},
		{Type: RoomBoss, Enemy: BossForFloor5(playerLevel)},
	}
	return rooms
}

func generateStandardFloor(floorNum, playerLevel int, biome Biome) []Room {
	var rooms []Room

	switch biome {
	case BiomeForest:
		// Layout Floresta: 4 salas equilibradas
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		if rand.Intn(100) < 40 {
			rooms = append(rooms, createTreasureRoom(floorNum, biome))
		} else {
			rooms = append(rooms, createRestRoom(biome))
		}
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})

	case BiomeIcy:
		// Layout Gélido: 5 salas, mais sobrevivência
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, createRestRoom(biome))
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, createRestRoom(biome))
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})

	case BiomeVolcanic:
		// Layout Vulcânico: 4 salas, intenso
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		// Sala de "Perigo" (Combate mais difícil)
		dangerEnemy := RandomEnemyForFloor(floorNum+1, playerLevel, biome)
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: dangerEnemy, Name: "ZONA CRITICA", Description: "O calor é insuportável!"})
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})

	case BiomeAbyssal:
		// Layout Abissal: 6 salas, ganancioso
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, createTreasureRoom(floorNum, biome))
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, createTreasureRoom(floorNum, biome))
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)})
		rooms = append(rooms, createRestRoom(biome))

	default:
		// Fallback para o antigo (4 salas)
		rooms = []Room{
			{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)},
			{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)},
			createRestRoom(biome),
			{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel, biome)},
		}
	}
	return rooms
}

func createRestRoom(biome Biome) Room {
	r := Room{Type: RoomRest}
	switch biome {
	case BiomeForest:
		r.Name = "Acampamento de Druidas"
		r.Description = "Um local de paz e cura."
	case BiomeIcy:
		r.Name = "Fonte Termal Oculta"
		r.Description = "Aqueça seu corpo e recupere as forças."
	case BiomeVolcanic:
		r.Name = "Forja de Obsidiana"
		r.Description = "Equipamentos podem ser resfriados aqui."
	case BiomeAbyssal:
		r.Name = "Santuário Esquecido"
		r.Description = "Energias ancestrais fluem nas paredes."
	}
	return r
}

func createTreasureRoom(floorNum int, biome Biome) Room {
	loot := randomLootForFloor(floorNum, biome)
	goldDrop := 10 + rand.Intn(floorNum*10)
	return Room{
		Type:     RoomTreasure,
		Loot:     loot,
		LootGold: goldDrop,
		Name:     "Tesouro Perdido",
	}
}

// randomLootForFloor sorteia um equipamento com raridade baseada no andar.
func randomLootForFloor(floor int, biome Biome) *Equipment {
	// Maior chance de raridade melhor em andares mais altos
	var candidates []*Equipment
	modifier := ModifierForBiome(biome)
	rareBoostRoll := rand.Intn(100) < modifier.RareLootChanceBonusPct

	for _, eq := range AllEquipment {
		switch {
		case floor <= 2 && eq.Rarity <= RarityIncomum:
			candidates = append(candidates, eq)
		case floor == 3 && eq.Rarity <= RarityRaro:
			candidates = append(candidates, eq)
		case floor >= 4:
			candidates = append(candidates, eq)
		}

		if rareBoostRoll && eq.Rarity >= RarityRaro {
			candidates = append(candidates, eq)
		}
	}

	if len(candidates) == 0 {
		candidates = AllEquipment
	}

	picked := candidates[rand.Intn(len(candidates))]
	copy := *picked
	return &copy
}

// CurrentRoomRef retorna ponteiro para a sala atual.
func (f *Floor) CurrentRoomRef() *Room {
	if f.CurrentRoom < len(f.Rooms) {
		return &f.Rooms[f.CurrentRoom]
	}
	return nil
}

// HasMoreRooms retorna se há mais salas no andar.
func (f *Floor) HasMoreRooms() bool {
	return f.CurrentRoom < len(f.Rooms)-1
}

// AdvanceRoom avança para a próxima sala.
func (f *Floor) AdvanceRoom() {
	if f.HasMoreRooms() {
		f.CurrentRoom++
	}
}
