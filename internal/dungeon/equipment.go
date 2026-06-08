package dungeon

import "fmt"

// Slot de equipamento
type EquipSlot int

const (
	SlotArma EquipSlot = iota
	SlotArmadura
	SlotAcessorio
)

func (s EquipSlot) String() string {
	switch s {
	case SlotArma:
		return "Arma"
	case SlotArmadura:
		return "Armadura"
	case SlotAcessorio:
		return "Acessorio"
	default:
		return "???"
	}
}

// Raridade do equipamento
type Rarity int

const (
	RarityComum Rarity = iota
	RarityIncomum
	RarityRaro
	RarityLendario
)

func (r Rarity) String() string {
	switch r {
	case RarityComum:
		return "Comum"
	case RarityIncomum:
		return "Incomum"
	case RarityRaro:
		return "Raro"
	case RarityLendario:
		return "Lendario"
	default:
		return "???"
	}
}

// Color retorna a cor ANSI hex para a raridade.
func (r Rarity) Color() string {
	switch r {
	case RarityComum:
		return "#AAAAAA"
	case RarityIncomum:
		return "#43BF6D"
	case RarityRaro:
		return "#5B9FFF"
	case RarityLendario:
		return "#FFD700"
	default:
		return "#FFFFFF"
	}
}

// Equipment representa um equipamento que pode ser equipado.
type Equipment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slot      EquipSlot `json:"slot"`
	Rarity    Rarity    `json:"rarity"`
	Price     int       `json:"price"` // Preço na loja
	BonusHP   int       `json:"bonus_hp"`
	BonusATK  int       `json:"bonus_atk"`
	BonusDEF  int       `json:"bonus_def"`
	BonusVEL  int       `json:"bonus_vel"`
	BonusLuck int       `json:"bonus_luck"`
}

func (e *Equipment) Description() string {
	parts := []string{}
	if e.BonusHP > 0 {
		parts = append(parts, fmt.Sprintf("+%d HP", e.BonusHP))
	}
	if e.BonusATK > 0 {
		parts = append(parts, fmt.Sprintf("+%d ATK", e.BonusATK))
	}
	if e.BonusDEF > 0 {
		parts = append(parts, fmt.Sprintf("+%d DEF", e.BonusDEF))
	}
	if e.BonusVEL > 0 {
		parts = append(parts, fmt.Sprintf("+%d VEL", e.BonusVEL))
	}
	if e.BonusLuck > 0 {
		parts = append(parts, fmt.Sprintf("+%d LUCK", e.BonusLuck))
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}

// Inventory guarda os equipamentos do jogador.
type Inventory struct {
	Arma      *Equipment   `json:"arma,omitempty"`
	Armadura  *Equipment   `json:"armadura,omitempty"`
	Acessorio *Equipment   `json:"acessorio,omitempty"`
	Backpack  []*Equipment `json:"backpack,omitempty"`
	Gold      int          `json:"gold"`
}

func NewInventory() *Inventory {
	return &Inventory{
		Backpack: []*Equipment{},
		Gold:     0,
	}
}

// Equip coloca um equipamento no slot correto, movendo o atual para a mochila.
func (inv *Inventory) Equip(eq *Equipment) {
	switch eq.Slot {
	case SlotArma:
		if inv.Arma != nil {
			inv.Backpack = append(inv.Backpack, inv.Arma)
		}
		inv.Arma = eq
	case SlotArmadura:
		if inv.Armadura != nil {
			inv.Backpack = append(inv.Backpack, inv.Armadura)
		}
		inv.Armadura = eq
	case SlotAcessorio:
		if inv.Acessorio != nil {
			inv.Backpack = append(inv.Backpack, inv.Acessorio)
		}
		inv.Acessorio = eq
	}
}

// AllEquipment retorna a tabela completa de equipamentos do jogo.
var AllEquipment = []*Equipment{
	// Armas
	{ID: "espada_ferro", Name: "Espada de Ferro", Slot: SlotArma, Rarity: RarityComum, Price: 40, BonusATK: 5},
	{ID: "espada_afiada", Name: "Espada Afiada", Slot: SlotArma, Rarity: RarityIncomum, Price: 100, BonusATK: 10, BonusVEL: 2},
	{ID: "machado_guerra", Name: "Machado de Guerra", Slot: SlotArma, Rarity: RarityRaro, Price: 250, BonusATK: 18, BonusDEF: 3},
	{ID: "katana_lendaria", Name: "Katana Lendaria", Slot: SlotArma, Rarity: RarityLendario, Price: 700, BonusATK: 25, BonusVEL: 5, BonusLuck: 5},

	// Armaduras
	{ID: "armadura_couro", Name: "Armadura de Couro", Slot: SlotArmadura, Rarity: RarityComum, Price: 50, BonusDEF: 4, BonusHP: 5},
	{ID: "armadura_ferro", Name: "Armadura de Ferro", Slot: SlotArmadura, Rarity: RarityIncomum, Price: 120, BonusDEF: 8, BonusHP: 10},
	{ID: "armadura_mithril", Name: "Armadura de Mithril", Slot: SlotArmadura, Rarity: RarityRaro, Price: 300, BonusDEF: 14, BonusHP: 20, BonusVEL: 2},
	{ID: "armadura_dragao", Name: "Armadura de Dragao", Slot: SlotArmadura, Rarity: RarityLendario, Price: 800, BonusDEF: 22, BonusHP: 35, BonusVEL: 3},

	// Acessórios
	{ID: "anel_sorte", Name: "Anel da Sorte", Slot: SlotAcessorio, Rarity: RarityComum, Price: 30, BonusLuck: 8},
	{ID: "amuleto_vigor", Name: "Amuleto de Vigor", Slot: SlotAcessorio, Rarity: RarityIncomum, Price: 80, BonusHP: 15, BonusDEF: 3},
	{ID: "bracelete_furia", Name: "Bracelete da Furia", Slot: SlotAcessorio, Rarity: RarityRaro, Price: 200, BonusATK: 10, BonusVEL: 5},
	{ID: "coroa_rei", Name: "Coroa do Rei", Slot: SlotAcessorio, Rarity: RarityLendario, Price: 600, BonusATK: 8, BonusDEF: 8, BonusHP: 20, BonusLuck: 10},

	// Drops exclusivos de chefes
	{ID: "espada_raiz_ancia", Name: "Espada da Raiz Ancia", Slot: SlotArma, Rarity: RarityLendario, Price: 800, BonusATK: 22, BonusHP: 15, BonusLuck: 8},
	{ID: "manto_lich", Name: "Manto do Lich", Slot: SlotArmadura, Rarity: RarityLendario, Price: 850, BonusDEF: 20, BonusHP: 30, BonusVEL: 4},
	{ID: "brasa_do_lorde", Name: "Brasa do Lorde", Slot: SlotAcessorio, Rarity: RarityLendario, Price: 820, BonusATK: 14, BonusDEF: 6, BonusVEL: 6, BonusLuck: 8},
	{ID: "lamina_abissal", Name: "Lamina Abissal", Slot: SlotArma, Rarity: RarityLendario, Price: 900, BonusATK: 28, BonusVEL: 3, BonusLuck: 12},
}

// EquipmentByID busca um equipamento pelo ID.
func EquipmentByID(id string) *Equipment {
	for _, eq := range AllEquipment {
		if eq.ID == id {
			// Retorna cópia para não modificar a tabela global
			copy := *eq
			return &copy
		}
	}
	return nil
}
