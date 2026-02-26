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
	Type     RoomType
	Cleared  bool
	Enemy    *Enemy     // Para salas de combate/boss
	Loot     *Equipment // Para salas de tesouro
	LootGold int        // Ouro em salas de tesouro
}

// Floor é um andar completo da dungeon.
type Floor struct {
	Number      int
	Rooms       []Room
	CurrentRoom int
}

// GenerateFloor gera as salas de um andar, com inimigos escalados pelo level do jogador.
// Andares 1-4: combate → combate → descanso/tesouro(40%) → combate
// Andar 5: combate → descanso → combate → BOSS
func GenerateFloor(floorNum, playerLevel int) *Floor {
	var rooms []Room

	if floorNum == 5 {
		rooms = []Room{
			{Type: RoomCombat, Enemy: RandomEnemyForFloor(4, playerLevel)},
			{Type: RoomRest},
			{Type: RoomCombat, Enemy: RandomEnemyForFloor(4, playerLevel)},
			{Type: RoomBoss, Enemy: BossForFloor5(playerLevel)},
		}
	} else {
		// Sala 1: combate
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel)})
		// Sala 2: combate
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel)})
		// Sala 3: descanso ou tesouro (40% tesouro)
		if rand.Intn(100) < 40 {
			loot := randomLootForFloor(floorNum)
			goldDrop := 10 + rand.Intn(floorNum*10)
			rooms = append(rooms, Room{Type: RoomTreasure, Loot: loot, LootGold: goldDrop})
		} else {
			rooms = append(rooms, Room{Type: RoomRest})
		}
		// Sala 4: combate
		rooms = append(rooms, Room{Type: RoomCombat, Enemy: RandomEnemyForFloor(floorNum, playerLevel)})
	}

	return &Floor{
		Number:      floorNum,
		Rooms:       rooms,
		CurrentRoom: 0,
	}
}

// randomLootForFloor sorteia um equipamento com raridade baseada no andar.
func randomLootForFloor(floor int) *Equipment {
	// Maior chance de raridade melhor em andares mais altos
	var candidates []*Equipment

	for _, eq := range AllEquipment {
		switch {
		case floor <= 2 && eq.Rarity <= RarityIncomum:
			candidates = append(candidates, eq)
		case floor == 3 && eq.Rarity <= RarityRaro:
			candidates = append(candidates, eq)
		case floor >= 4:
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
