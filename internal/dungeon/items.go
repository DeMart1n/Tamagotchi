package dungeon

// ItemType tipo de item consumível
type ItemType int

const (
	ItemPotion   ItemType = iota // Cura HP
	ItemATKBoost                 // Buff ATK temporário
)

// Item consumível usado durante combate.
type Item struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Type  ItemType `json:"type"`
	Value int      `json:"value"` // HP curado ou ATK adicionado
	Price int      `json:"price"` // Preço na loja de descanso
}

// Catálogo de itens disponíveis na loja.
var ShopItems = []Item{
	{ID: "pocao_pequena", Name: "Pocao Pequena", Type: ItemPotion, Value: 25, Price: 10},
	{ID: "pocao_grande", Name: "Pocao Grande", Type: ItemPotion, Value: 50, Price: 25},
	{ID: "elixir_forca", Name: "Elixir de Forca", Type: ItemATKBoost, Value: 8, Price: 20},
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
