package dungeon

// ItemType tipo de item consumível
type ItemType int

const (
	ItemPotion   ItemType = iota // Cura HP
	ItemATKBoost                 // Buff ATK temporário
)

// ItemCategory categorização de itens na loja
type ItemCategory int

const (
	CategoryConsumivel ItemCategory = iota
	CategoryEquipamento
	CategoryQualidade
	CategoryNPCUnico
)

// Item consumível usado durante combate.
type Item struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Type     ItemType     `json:"type"`
	Category ItemCategory `json:"category"`
	Rarity   Rarity       `json:"rarity"`
	Value    int          `json:"value"` // HP curado ou ATK adicionado
	Price    int          `json:"price"` // Preço na loja de descanso
}

// Catálogo de itens disponíveis na loja.
var ShopItems = []Item{
	// Consumíveis
	{ID: "pocao_pequena", Name: "Pocao Pequena", Type: ItemPotion, Category: CategoryConsumivel, Rarity: RarityComum, Value: 25, Price: 10},
	{ID: "pocao_grande", Name: "Pocao Grande", Type: ItemPotion, Category: CategoryConsumivel, Rarity: RarityIncomum, Value: 50, Price: 25},
	{ID: "pocao_completa", Name: "Pocao Completa", Type: ItemPotion, Category: CategoryConsumivel, Rarity: RarityRaro, Value: 100, Price: 50},
	{ID: "elixir_forca", Name: "Elixir de Forca", Type: ItemATKBoost, Category: CategoryConsumivel, Rarity: RarityIncomum, Value: 8, Price: 20},
	{ID: "antidoto", Name: "Antidoto", Type: ItemPotion, Category: CategoryConsumivel, Rarity: RarityComum, Value: 15, Price: 15},

	// Qualidade (itens combinados)
	{ID: "pocao_batalha", Name: "Pocao de Batalha", Type: ItemPotion, Category: CategoryQualidade, Rarity: RarityRaro, Value: 60, Price: 60},
	{ID: "elixir_defesa", Name: "Elixir de Defesa", Type: ItemPotion, Category: CategoryQualidade, Rarity: RarityRaro, Value: 40, Price: 55},

	// NPC Único (itens raros disponíveis só na loja, stock=1 por run)
	{ID: "amuleto_sorte", Name: "Amuleto da Sorte", Type: ItemPotion, Category: CategoryNPCUnico, Rarity: RarityLendario, Value: 30, Price: 150},
	{ID: "bencao_oraculo", Name: "Bencao do Oraculo", Type: ItemPotion, Category: CategoryNPCUnico, Rarity: RarityLendario, Value: 80, Price: 200},
}

// ItemBag guarda os itens consumíveis durante uma run.
type ItemBag struct {
	Items []Item `json:"items"`
}

func NewItemBag() *ItemBag {
	return &ItemBag{Items: []Item{}}
}

func (b *ItemBag) Add(item Item) {
	b.Items = append(b.Items, item)
}

func (b *ItemBag) Remove(index int) {
	if index < 0 || index >= len(b.Items) {
		return
	}
	b.Items = append(b.Items[:index], b.Items[index+1:]...)
}

func (b *ItemBag) Count() int {
	return len(b.Items)
}

// ShopEntry representa um item ou equipamento na loja
type ShopEntry struct {
	Kind  string      // "item" ou "equipment"
	Item  *Item       // nil se Kind == "equipment"
	Equip *Equipment  // nil se Kind == "item"
	Price int
	Stock int // -1 = ilimitado, >0 = quantidade específica
}
