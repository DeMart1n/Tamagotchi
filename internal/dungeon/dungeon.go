package dungeon

import (
	"Pessoal/internal/model"
	"fmt"
	"math/rand"
)

// Phase da máquina de estados da dungeon.
type Phase int

const (
	PhaseIntro Phase = iota
	PhaseMenuPrincipal
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
	Phase        Phase
	Tama         *model.Tama
	CurrentBiome Biome
	Stats        CombatStats
	Inv          *Inventory
	ItemBag      *ItemBag
	Floor        *Floor
	FloorNum     int
	Combat       *Combat
	Message      string
	SubMessage   string
	TotalXP      int
	TotalGold    int
	PendingLoot  *Equipment // Equipamento encontrado esperando aceitar/recusar

	// Estado do submenu de itens no combate
	ChoosingItem  bool
	ChoosingSkill bool
	ItemCursor    int
	SkillCursor   int

	// Estado da loja de descanso
	ShopCursor   int
	ShopTab      int           // 0=Consumível 1=Equipamento 2=Qualidade 3=NPC Único
	ShopMode     string        // "buy" ou "sell"
	ShopStock    []ShopEntry   // estoque atual
	ShopConfirm  bool          // aguardando confirmação
	ShopPending  int           // índice do item selecionado para confirmar
}

// NewDungeonRun cria uma nova sessão de dungeon.
func NewDungeonRun(tama *model.Tama, inv *Inventory) *DungeonRun {
	stats := DeriveCombatStats(tama).ApplyEquipment(inv)
	biome := RandomBiome()

	// Garante que o Pet tenha habilidades: usa skills de classe ou fallback genérico
	if len(tama.SkillsKnown) == 0 {
		if cls, ok := model.AllClasses[tama.Class]; ok {
			tama.SkillsKnown = cls.BaseSkills
		} else {
			tama.SkillsKnown = model.GetInitialSkills(tama.Stage)
		}
	}

	return &DungeonRun{
		Phase:        PhaseIntro,
		Tama:         tama,
		CurrentBiome: biome,
		Stats:        stats,
		Inv:          inv,
		ItemBag:      NewItemBag(),
		FloorNum:     0,
		Message:      "Bem-vindo a Masmorra!",
	}
}

// HandleInput processa a entrada do jogador baseado na fase atual.
// Retorna true se a dungeon terminou (PhaseDone).
func (d *DungeonRun) HandleInput(input string) bool {
	switch d.Phase {
	case PhaseIntro:
		d.handleIntro(input)
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

// handleIntro processa a tela de introdução com a arte do bioma
func (d *DungeonRun) handleIntro(_ string) {
	// Qualquer input avança para o menu principal
	d.Phase = PhaseMenuPrincipal
	d.Message = "Bem-vindo a Masmorra!"
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

func (d *DungeonRun) handleInventario(_ string) {
	// Qualquer input volta ao menu principal
	d.Phase = PhaseMenuPrincipal
	d.Message = "Bem-vindo a Masmorra!"
}

func (d *DungeonRun) startDungeon() {
	d.FloorNum = 1
	d.Floor = GenerateFloor(1, d.Tama.Level, d.CurrentBiome)
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
		d.Combat = NewCombat(&d.Stats, room.Enemy, d.CurrentBiome, d.Tama)
		d.Message = fmt.Sprintf("Um %s apareceu!", room.Enemy.Name)
		d.SubMessage = ""
	case RoomBoss:
		d.Phase = PhaseCombate
		d.Combat = NewCombat(&d.Stats, room.Enemy, d.CurrentBiome, d.Tama)
		d.Message = fmt.Sprintf("BOSS: %s!", room.Enemy.Name)
		d.SubMessage = ""
	case RoomRest:
		d.Phase = PhaseDescanso
		heal := d.Stats.HPMax / 4
		if d.CurrentBiome == BiomeAbyssal {
			heal = heal / 2
			if heal < 1 {
				heal = 1
			}
		}
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
	if d.ChoosingSkill {
		d.handleSkillChoice(input)
		return
	}

	if d.Combat == nil || d.Combat.IsOver() {
		return
	}

	switch input {
	case "1": // Atacar
		d.Combat.ExecuteAction(ActionAtacar, nil, nil, 0)
	case "2": // Defender
		d.Combat.ExecuteAction(ActionDefender, nil, nil, 0)
	case "3": // Habilidades
		if len(model.GetActiveSkills(d.Tama.SkillsKnown)) == 0 {
			d.Message = "Voce nao tem habilidades ativas!"
			return
		}
		d.ChoosingSkill = true
		d.SkillCursor = 0
		d.Message = "Escolha uma habilidade (numero) ou 0 para voltar:"
		return
	case "4": // Item
		if d.ItemBag.Count() == 0 {
			d.Message = "Voce nao tem itens!"
			return
		}
		d.ChoosingItem = true
		d.ItemCursor = 0
		d.Message = "Escolha um item (numero) ou 0 para voltar:"
		return
	case "5": // Fugir
		d.Combat.ExecuteAction(ActionFugir, nil, nil, 0)
	default:
		return
	}

	d.updateCombatMessage()

	if d.Combat.IsOver() {
		d.Phase = PhaseCombatResult
	}
}

func (d *DungeonRun) handleSkillChoice(input string) {
	if input == "0" {
		d.ChoosingSkill = false
		d.Message = "Escolha uma acao:"
		return
	}

	idx := -1
	if len(input) == 1 && input[0] >= '1' && input[0] <= '9' {
		idx = int(input[0] - '1')
	}

	activeSkills := model.GetActiveSkills(d.Tama.SkillsKnown)
	if idx < 0 || idx >= len(activeSkills) {
		d.Message = "Habilidade invalida!"
		return
	}

	d.ChoosingSkill = false
	d.Combat.ExecuteAction(ActionSkill, nil, d.Tama, idx)
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
	d.Combat.ExecuteAction(ActionItem, d.ItemBag, d.Tama, idx)
	d.updateCombatMessage()

	if d.Combat.IsOver() {
		d.Phase = PhaseCombatResult
	}
}

func (d *DungeonRun) updateCombatMessage() {
	if len(d.Combat.Log) > 0 {
		// Pega as últimas duas mensagens para garantir que o ataque do inimigo apareça
		// mesmo se houver mensagens de bioma ou ação do jogador antes.
		n := len(d.Combat.Log)
		if n >= 2 {
			d.Message = d.Combat.Log[n-2]
			d.SubMessage = d.Combat.Log[n-1]
		} else {
			d.Message = d.Combat.Log[0]
			d.SubMessage = ""
		}
	}
}

func (d *DungeonRun) handleCombatResult(_ string) {
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
		d.Tama.TotalEnemiesDefeated++

		// Chance de drop de equipamento (20%)
		var lootMsg string
		if rand.Intn(100) < 20 {
			loot := randomLootForFloor(d.FloorNum, d.CurrentBiome)
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
		// Derrota — Morte Permanente
		d.Tama.Dead = true
		d.Phase = PhaseDerrota
		d.Message = "Você foi derrotado na masmorra e não sobreviveu..."
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
		d.ShopTab = 0
		d.ShopMode = "buy"
		d.ShopConfirm = false
		d.ShopStock = generateShopStock(d.FloorNum)
		d.Message = fmt.Sprintf("Loja de Aventureiros (Ouro: %d)", d.Inv.Gold)
	case "2": // Continuar
		room := d.Floor.CurrentRoomRef()
		if room != nil {
			room.Cleared = true
		}
		d.advanceAfterRoom()
	}
}

func (d *DungeonRun) handleDescansoLoja(input string) {
	if d.ShopConfirm {
		d.handleShopConfirm(input)
		return
	}

	switch input {
	case "0", "Escape":
		d.Phase = PhaseDescanso
		d.Message = "Sala de Descanso"
		return
	case "h", "Left":
		d.ShopTab = (d.ShopTab - 1 + 4) % 4
		d.ShopCursor = 0
		return
	case "l", "Right":
		d.ShopTab = (d.ShopTab + 1) % 4
		d.ShopCursor = 0
		return
	case "j", "Down":
		itemsInTab := d.GetItemsInTab()
		d.ShopCursor = (d.ShopCursor + 1) % len(itemsInTab)
		return
	case "k", "Up":
		itemsInTab := d.GetItemsInTab()
		d.ShopCursor = (d.ShopCursor - 1 + len(itemsInTab)) % len(itemsInTab)
		return
	case "v":
		if d.ShopMode == "buy" {
			d.ShopMode = "sell"
		} else {
			d.ShopMode = "buy"
		}
		d.ShopCursor = 0
		return
	case "Return":
		itemsInTab := d.GetItemsInTab()
		if len(itemsInTab) == 0 {
			d.Message = "Nenhum item nesta categoria!"
			return
		}
		if d.ShopCursor < 0 || d.ShopCursor >= len(itemsInTab) {
			return
		}

		entry := itemsInTab[d.ShopCursor]
		if entry.Stock == 0 {
			d.Message = "Este item acabou!"
			return
		}

		if d.ShopMode == "buy" {
			if d.Inv.Gold < entry.Price {
				d.Message = fmt.Sprintf("Ouro insuficiente! Precisa de %d, tem %d.", entry.Price, d.Inv.Gold)
				return
			}
			d.ShopPending = d.ShopCursor
			d.ShopConfirm = true
			return
		}
	}
}

// GetItemsInTab retorna os itens do tab atual
func (d *DungeonRun) GetItemsInTab() []ShopEntry {
	var result []ShopEntry

	switch d.ShopTab {
	case 0: // Consumível
		for _, entry := range d.ShopStock {
			if entry.Kind == "item" && entry.Item.Category == CategoryConsumivel {
				result = append(result, entry)
			}
		}
	case 1: // Equipamento
		for _, entry := range d.ShopStock {
			if entry.Kind == "equipment" {
				result = append(result, entry)
			}
		}
	case 2: // Qualidade
		for _, entry := range d.ShopStock {
			if entry.Kind == "item" && entry.Item.Category == CategoryQualidade {
				result = append(result, entry)
			}
		}
	case 3: // NPC Único
		for _, entry := range d.ShopStock {
			if entry.Kind == "item" && entry.Item.Category == CategoryNPCUnico {
				result = append(result, entry)
			}
		}
	}

	return result
}

// handleShopConfirm processa a confirmação de compra/venda
func (d *DungeonRun) handleShopConfirm(input string) {
	switch input {
	case "s", "S":
		itemsInTab := d.GetItemsInTab()
		if d.ShopPending < 0 || d.ShopPending >= len(itemsInTab) {
			d.ShopConfirm = false
			return
		}

		entry := itemsInTab[d.ShopPending]
		if d.ShopMode == "buy" {
			if entry.Kind == "item" {
				d.Inv.Gold -= entry.Price
				d.ItemBag.Add(*entry.Item)
				d.Message = fmt.Sprintf("Comprou %s!", entry.Item.Name)
				if entry.Stock > 0 {
					entry.Stock--
				}
			} else if entry.Kind == "equipment" {
				d.Inv.Gold -= entry.Price
				copy := *entry.Equip
				d.Inv.Backpack = append(d.Inv.Backpack, &copy)
				d.Message = fmt.Sprintf("Comprou %s!", entry.Equip.Name)
				if entry.Stock > 0 {
					entry.Stock--
				}
			}
		}
		d.ShopConfirm = false
	case "n", "N":
		d.ShopConfirm = false
		d.Message = "Compra cancelada."
	}
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

func (d *DungeonRun) handleFimAndar(_ string) {
	// Qualquer input avança para o próximo andar
	d.FloorNum++
	d.Floor = GenerateFloor(d.FloorNum, d.Tama.Level, d.CurrentBiome)
	d.Message = fmt.Sprintf("Entrando no Andar %d...", d.FloorNum)
	d.enterCurrentRoom()
}

// generateShopStock cria o estoque de loja baseado no andar
func generateShopStock(floorNum int) []ShopEntry {
	maxRarity := RarityComum
	if floorNum >= 4 {
		maxRarity = RarityRaro
	}
	if floorNum >= 7 {
		maxRarity = RarityLendario
	}

	var stock []ShopEntry

	// Consumíveis (sempre disponíveis)
	for _, item := range ShopItems {
		if item.Category != CategoryConsumivel {
			continue
		}
		if item.Rarity > maxRarity {
			continue
		}
		itemCopy := item
		stock = append(stock, ShopEntry{
			Kind:  "item",
			Item:  &itemCopy,
			Price: item.Price,
			Stock: -1,
		})
	}

	// Qualidade (sempre disponível)
	for _, item := range ShopItems {
		if item.Category != CategoryQualidade {
			continue
		}
		if item.Rarity > maxRarity {
			continue
		}
		itemCopy := item
		stock = append(stock, ShopEntry{
			Kind:  "item",
			Item:  &itemCopy,
			Price: item.Price,
			Stock: -1,
		})
	}

	// Equipamentos (filtra por raridade)
	for _, eq := range AllEquipment {
		if eq.Rarity <= maxRarity {
			stock = append(stock, ShopEntry{
				Kind:  "equipment",
				Equip: eq,
				Price: eq.Price,
				Stock: -1,
			})
		}
	}

	// NPC Único (stock=1, sempre aparecem)
	for _, item := range ShopItems {
		if item.Category == CategoryNPCUnico {
			itemCopy := item
			stock = append(stock, ShopEntry{
				Kind:  "item",
				Item:  &itemCopy,
				Price: item.Price,
				Stock: 1,
			})
		}
	}

	return stock
}

// getPriceForSale retorna 50% do preço original
func getPriceForSale(price int) int {
	return price / 2
}

func (d *DungeonRun) handleVitoria(_ string) {
	d.Phase = PhaseDone
}

func (d *DungeonRun) handleDerrota(_ string) {
	d.Phase = PhaseDone
}
