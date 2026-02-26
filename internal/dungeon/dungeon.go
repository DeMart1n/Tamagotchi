package dungeon

import (
	"Pessoal/internal/model"
	"fmt"
	"math/rand"
)

// Phase da máquina de estados da dungeon.
type Phase int

const (
	PhaseMenuPrincipal Phase = iota
	PhaseInventario
	PhaseExplorando
	PhaseCombate
	PhaseCombatResult
	PhaseDescanso
	PhaseDescansoLoja
	PhaseTesouro
	PhaseFimAndar
	PhaseVitoria
	PhaseDerrota
	PhaseDone // Sinaliza que a dungeon acabou e deve voltar ao modo normal
)

// DungeonRun orquestra toda a sessão de dungeon.
type DungeonRun struct {
	Phase       Phase
	Tama        *model.Tama
	Stats       CombatStats
	Inv         *Inventory
	ItemBag     *ItemBag
	Floor       *Floor
	FloorNum    int
	Combat      *Combat
	Message     string
	SubMessage  string
	TotalXP     int
	TotalGold   int
	PendingLoot *Equipment // Equipamento encontrado esperando aceitar/recusar

	// Estado do submenu de itens no combate
	ChoosingItem bool
	ItemCursor   int

	// Estado da loja de descanso
	ShopCursor int
}

// NewDungeonRun cria uma nova sessão de dungeon.
func NewDungeonRun(tama *model.Tama, inv *Inventory) *DungeonRun {
	stats := DeriveCombatStats(tama).ApplyEquipment(inv)
	return &DungeonRun{
		Phase:    PhaseMenuPrincipal,
		Tama:     tama,
		Stats:    stats,
		Inv:      inv,
		ItemBag:  NewItemBag(),
		FloorNum: 0,
		Message:  "Bem-vindo a Masmorra!",
	}
}

// HandleInput processa a entrada do jogador baseado na fase atual.
// Retorna true se a dungeon terminou (PhaseDone).
func (d *DungeonRun) HandleInput(input string) bool {
	switch d.Phase {
	case PhaseMenuPrincipal:
		d.handleMenuPrincipal(input)
	case PhaseInventario:
		d.handleInventario(input)
	case PhaseExplorando:
		// Auto-entra na sala atual (não deve chegar aqui normalmente)
		d.enterCurrentRoom()
	case PhaseCombate:
		d.handleCombate(input)
	case PhaseCombatResult:
		d.handleCombatResult(input)
	case PhaseDescanso:
		d.handleDescanso(input)
	case PhaseDescansoLoja:
		d.handleDescansoLoja(input)
	case PhaseTesouro:
		d.handleTesouro(input)
	case PhaseFimAndar:
		d.handleFimAndar(input)
	case PhaseVitoria:
		d.handleVitoria(input)
	case PhaseDerrota:
		d.handleDerrota(input)
	}

	return d.Phase == PhaseDone
}

func (d *DungeonRun) handleMenuPrincipal(input string) {
	switch input {
	case "1":
		d.startDungeon()
	case "2":
		d.Phase = PhaseInventario
		d.Message = "Inventario"
	case "3":
		d.Phase = PhaseDone
	}
}

func (d *DungeonRun) handleInventario(input string) {
	// Qualquer input volta ao menu principal
	d.Phase = PhaseMenuPrincipal
	d.Message = "Bem-vindo a Masmorra!"
}

func (d *DungeonRun) startDungeon() {
	d.FloorNum = 1
	d.Floor = GenerateFloor(1, d.Tama.Level)
	d.Stats = DeriveCombatStats(d.Tama).ApplyEquipment(d.Inv)
	d.TotalXP = 0
	d.TotalGold = 0
	d.enterCurrentRoom()
}

func (d *DungeonRun) enterCurrentRoom() {
	room := d.Floor.CurrentRoomRef()
	if room == nil {
		return
	}

	switch room.Type {
	case RoomCombat:
		d.Phase = PhaseCombate
		d.Combat = NewCombat(&d.Stats, room.Enemy)
		d.Message = fmt.Sprintf("Um %s apareceu!", room.Enemy.Name)
		d.SubMessage = ""
	case RoomBoss:
		d.Phase = PhaseCombate
		d.Combat = NewCombat(&d.Stats, room.Enemy)
		d.Message = fmt.Sprintf("BOSS: %s!", room.Enemy.Name)
		d.SubMessage = ""
	case RoomRest:
		d.Phase = PhaseDescanso
		heal := d.Stats.HPMax / 4
		d.Stats.HPCurrent += heal
		if d.Stats.HPCurrent > d.Stats.HPMax {
			d.Stats.HPCurrent = d.Stats.HPMax
		}
		d.Message = fmt.Sprintf("Sala de Descanso! Recuperou %d HP.", heal)
		d.SubMessage = ""
	case RoomTreasure:
		d.Phase = PhaseTesouro
		d.PendingLoot = room.Loot
		d.TotalGold += room.LootGold
		d.Inv.Gold += room.LootGold
		d.Message = fmt.Sprintf("Bau do Tesouro! +%d ouro!", room.LootGold)
		if room.Loot != nil {
			d.SubMessage = fmt.Sprintf("Encontrou: %s (%s)", room.Loot.Name, room.Loot.Rarity.String())
		}
	}
}

func (d *DungeonRun) handleCombate(input string) {
	if d.ChoosingItem {
		d.handleItemChoice(input)
		return
	}

	if d.Combat == nil || d.Combat.IsOver() {
		return
	}

	switch input {
	case "1": // Atacar
		d.Combat.ExecuteAction(ActionAtacar, nil, 0)
	case "2": // Defender
		d.Combat.ExecuteAction(ActionDefender, nil, 0)
	case "3": // Item
		if d.ItemBag.Count() == 0 {
			d.Message = "Voce nao tem itens!"
			return
		}
		d.ChoosingItem = true
		d.ItemCursor = 0
		d.Message = "Escolha um item (numero) ou 0 para voltar:"
		return
	case "4": // Fugir
		d.Combat.ExecuteAction(ActionFugir, nil, 0)
	default:
		return
	}

	d.updateCombatMessage()

	if d.Combat.IsOver() {
		d.Phase = PhaseCombatResult
	}
}

func (d *DungeonRun) handleItemChoice(input string) {
	if input == "0" {
		d.ChoosingItem = false
		d.Message = "Escolha uma acao:"
		return
	}

	idx := -1
	if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
		idx = int(input[0] - '1')
	}

	if idx < 0 || idx >= d.ItemBag.Count() {
		d.Message = "Item invalido!"
		return
	}

	d.ChoosingItem = false
	d.Combat.ExecuteAction(ActionItem, d.ItemBag, idx)
	d.updateCombatMessage()

	if d.Combat.IsOver() {
		d.Phase = PhaseCombatResult
	}
}

func (d *DungeonRun) updateCombatMessage() {
	if len(d.Combat.Log) > 0 {
		d.Message = d.Combat.Log[0]
		if len(d.Combat.Log) > 1 {
			d.SubMessage = d.Combat.Log[1]
		} else {
			d.SubMessage = ""
		}
	}
}

func (d *DungeonRun) handleCombatResult(input string) {
	if d.Combat == nil {
		return
	}

	room := d.Floor.CurrentRoomRef()

	if d.Combat.Won {
		xp := d.Combat.Enemy.XPDrop
		gold := d.Combat.Enemy.GoldDrop
		d.TotalXP += xp
		d.TotalGold += gold
		d.Inv.Gold += gold

		// Chance de drop de equipamento (20%)
		var lootMsg string
		if rand.Intn(100) < 20 {
			loot := randomLootForFloor(d.FloorNum)
			d.PendingLoot = loot
			lootMsg = fmt.Sprintf(" Dropou: %s!", loot.Name)
		}

		d.Message = fmt.Sprintf("+%d XP, +%d ouro!%s", xp, gold, lootMsg)

		if d.PendingLoot != nil {
			d.Phase = PhaseTesouro
			d.SubMessage = fmt.Sprintf("Encontrou: %s (%s) - %s", d.PendingLoot.Name, d.PendingLoot.Rarity.String(), d.PendingLoot.Description())
			return
		}

		room.Cleared = true
		d.advanceAfterRoom()
	} else if d.Combat.Fled {
		room.Cleared = true
		d.Message = "Fugiu do combate!"
		d.advanceAfterRoom()
	} else if d.Combat.Lost {
		// Derrota — XP parcial
		partialXP := d.TotalXP / 2
		if partialXP > 0 {
			d.Tama.AddXP(partialXP)
		}
		d.Phase = PhaseDerrota
		d.Message = fmt.Sprintf("Voce foi derrotado! Recebeu %d XP parcial.", partialXP)
		d.Tama.TotalDungeonRuns++
	}
}

func (d *DungeonRun) advanceAfterRoom() {
	if d.Floor.HasMoreRooms() {
		d.Floor.AdvanceRoom()
		d.enterCurrentRoom()
	} else {
		// Fim do andar
		if d.FloorNum >= 5 {
			// Vitória completa!
			d.Phase = PhaseVitoria
			d.Tama.AddXP(d.TotalXP)
			d.Tama.TotalDungeonRuns++
			d.Tama.TotalBossesDefeated++
			d.Message = fmt.Sprintf("VITORIA! Dragao Anciao derrotado! +%d XP, +%d ouro!", d.TotalXP, d.TotalGold)
		} else {
			d.Phase = PhaseFimAndar
			d.Message = fmt.Sprintf("Andar %d completo!", d.FloorNum)
		}
	}
}

func (d *DungeonRun) handleDescanso(input string) {
	switch input {
	case "1": // Loja
		d.Phase = PhaseDescansoLoja
		d.ShopCursor = 0
		d.Message = fmt.Sprintf("Loja (Ouro: %d)", d.Inv.Gold)
	case "2": // Continuar
		room := d.Floor.CurrentRoomRef()
		if room != nil {
			room.Cleared = true
		}
		d.advanceAfterRoom()
	}
}

func (d *DungeonRun) handleDescansoLoja(input string) {
	if input == "0" {
		d.Phase = PhaseDescanso
		d.Message = "Sala de Descanso"
		return
	}

	idx := -1
	if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
		idx = int(input[0] - '1')
	}

	if idx < 0 || idx >= len(ShopItems) {
		d.Message = "Item invalido!"
		return
	}

	item := ShopItems[idx]
	if d.Inv.Gold < item.Price {
		d.Message = fmt.Sprintf("Ouro insuficiente! Precisa de %d, tem %d.", item.Price, d.Inv.Gold)
		return
	}

	d.Inv.Gold -= item.Price
	d.ItemBag.Add(item)
	d.Message = fmt.Sprintf("Comprou %s! (Ouro: %d)", item.Name, d.Inv.Gold)
}

func (d *DungeonRun) handleTesouro(input string) {
	if d.PendingLoot != nil {
		switch input {
		case "1": // Equipar
			d.Inv.Equip(d.PendingLoot)
			d.Message = fmt.Sprintf("Equipou %s!", d.PendingLoot.Name)
			// Recalcula stats com novo equipamento
			d.Stats = DeriveCombatStats(d.Tama).ApplyEquipment(d.Inv)
			// Manter HP proporcional
			d.PendingLoot = nil
		case "2": // Guardar na mochila
			d.Inv.Backpack = append(d.Inv.Backpack, d.PendingLoot)
			d.Message = fmt.Sprintf("Guardou %s na mochila.", d.PendingLoot.Name)
			d.PendingLoot = nil
		case "3": // Descartar
			d.Message = fmt.Sprintf("Descartou %s.", d.PendingLoot.Name)
			d.PendingLoot = nil
		default:
			return
		}
	}

	// Se não tem mais loot pendente, avança
	if d.PendingLoot == nil {
		room := d.Floor.CurrentRoomRef()
		if room != nil {
			room.Cleared = true
		}
		d.advanceAfterRoom()
	}
}

func (d *DungeonRun) handleFimAndar(input string) {
	// Qualquer input avança para o próximo andar
	d.FloorNum++
	d.Floor = GenerateFloor(d.FloorNum, d.Tama.Level)
	d.Message = fmt.Sprintf("Entrando no Andar %d...", d.FloorNum)
	d.enterCurrentRoom()
}

func (d *DungeonRun) handleVitoria(input string) {
	d.Phase = PhaseDone
}

func (d *DungeonRun) handleDerrota(input string) {
	d.Phase = PhaseDone
}
