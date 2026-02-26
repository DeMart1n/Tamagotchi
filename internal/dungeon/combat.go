package dungeon

import (
	"fmt"
	"math/rand"
)

// CombatAction ação do jogador no combate.
type CombatAction int

const (
	ActionAtacar CombatAction = iota
	ActionDefender
	ActionItem
	ActionFugir
)

// CombatState estado da máquina de combate.
type CombatState int

const (
	CombatChoosing CombatState = iota // Esperando ação do jogador
	CombatResolved                    // Turno resolvido, mostrar resultado
)

// Combat motor de combate por turnos.
type Combat struct {
	Player    *CombatStats
	Enemy     *Enemy
	State     CombatState
	Log       []string
	Turn      int
	Defending bool // Jogador está defendendo neste turno
	ATKBuff   int  // Buff temporário de ATK (de elixir)
	Fled      bool // Jogador fugiu com sucesso
	Won       bool
	Lost      bool
}

// NewCombat cria um novo combate.
func NewCombat(player *CombatStats, enemy *Enemy) *Combat {
	return &Combat{
		Player: player,
		Enemy:  enemy,
		State:  CombatChoosing,
		Log:    []string{},
		Turn:   1,
	}
}

// IsOver retorna se o combate terminou.
func (c *Combat) IsOver() bool {
	return c.Won || c.Lost || c.Fled
}

// ExecuteAction processa a ação do jogador e a resposta do inimigo.
func (c *Combat) ExecuteAction(action CombatAction, itemBag *ItemBag, itemIndex int) {
	c.Log = []string{}
	c.Defending = false

	switch action {
	case ActionAtacar:
		c.playerAttack()
	case ActionDefender:
		c.Defending = true
		c.Log = append(c.Log, "Voce se preparou para defender!")
	case ActionItem:
		c.useItem(itemBag, itemIndex)
	case ActionFugir:
		c.tryFlee()
	}

	if c.Fled {
		c.State = CombatResolved
		return
	}

	// Verifica morte do inimigo
	if c.Enemy.HPCurrent <= 0 {
		c.Enemy.HPCurrent = 0
		c.Won = true
		c.Log = append(c.Log, fmt.Sprintf("%s foi derrotado!", c.Enemy.Name))
		c.State = CombatResolved
		return
	}

	// Turno do inimigo
	c.enemyTurn()

	// Verifica morte do jogador
	if c.Player.HPCurrent <= 0 {
		c.Player.HPCurrent = 0
		c.Lost = true
		c.Log = append(c.Log, "Voce foi derrotado...")
	}

	c.Turn++
	c.State = CombatResolved

	// Limpa buff de ATK temporário
	if c.ATKBuff > 0 {
		c.Player.Ataque -= c.ATKBuff
		c.ATKBuff = 0
	}
}

func (c *Combat) playerAttack() {
	baseDmg := c.Player.Ataque - c.Enemy.Defesa/2
	if baseDmg < 1 {
		baseDmg = 1
	}

	// Variância ±20%
	variance := float64(baseDmg) * 0.2
	dmg := baseDmg + int(float64(rand.Intn(int(variance*2+1)))-variance)
	if dmg < 1 {
		dmg = 1
	}

	// Chance de crítico baseada na Sorte
	critical := false
	if rand.Intn(100) < c.Player.Sorte {
		dmg = int(float64(dmg) * 1.5)
		critical = true
	}

	c.Enemy.HPCurrent -= dmg
	if critical {
		c.Log = append(c.Log, fmt.Sprintf("CRITICO! Voce atacou %s! %d de dano!", c.Enemy.Name, dmg))
	} else {
		c.Log = append(c.Log, fmt.Sprintf("Voce atacou %s! %d de dano!", c.Enemy.Name, dmg))
	}
}

func (c *Combat) enemyTurn() {
	// IA: 30% chance de defender quando HP < 25%
	enemyDefending := false
	hpPct := float64(c.Enemy.HPCurrent) / float64(c.Enemy.HPMax)
	if hpPct < 0.25 && rand.Intn(100) < 30 {
		enemyDefending = true
		c.Enemy.Defending = true
		c.Log = append(c.Log, fmt.Sprintf("%s se preparou para defender!", c.Enemy.Name))
		return
	}

	baseDmg := c.Enemy.Ataque - c.Player.Defesa/2
	if baseDmg < 1 {
		baseDmg = 1
	}

	// Variância ±20%
	variance := float64(baseDmg) * 0.2
	dmg := baseDmg + int(float64(rand.Intn(int(variance*2+1)))-variance)
	if dmg < 1 {
		dmg = 1
	}

	// Redução se jogador está defendendo
	if c.Defending {
		dmg = dmg / 2
		if dmg < 1 {
			dmg = 1
		}
	}

	c.Player.HPCurrent -= dmg
	c.Log = append(c.Log, fmt.Sprintf("%s atacou! %d de dano!", c.Enemy.Name, dmg))

	_ = enemyDefending // previne warning
}

func (c *Combat) useItem(bag *ItemBag, index int) {
	if bag == nil || index < 0 || index >= bag.Count() {
		c.Log = append(c.Log, "Nenhum item para usar!")
		return
	}

	item := bag.Items[index]
	switch item.Type {
	case ItemPotion:
		heal := item.Value
		c.Player.HPCurrent += heal
		if c.Player.HPCurrent > c.Player.HPMax {
			c.Player.HPCurrent = c.Player.HPMax
		}
		c.Log = append(c.Log, fmt.Sprintf("Usou %s! Recuperou %d HP!", item.Name, heal))
	case ItemATKBoost:
		c.ATKBuff = item.Value
		c.Player.Ataque += item.Value
		c.Log = append(c.Log, fmt.Sprintf("Usou %s! ATK +%d neste turno!", item.Name, item.Value))
	}
	bag.Remove(index)
}

func (c *Combat) tryFlee() {
	if c.Enemy.IsBoss {
		c.Log = append(c.Log, "Impossivel fugir do BOSS!")
		return
	}

	// 40% base + 2% por vantagem de velocidade
	chance := 40 + (c.Player.Velocidade-c.Enemy.Velocidade)*2
	if chance < 10 {
		chance = 10
	}
	if chance > 90 {
		chance = 90
	}

	if rand.Intn(100) < chance {
		c.Fled = true
		c.Log = append(c.Log, "Voce fugiu com sucesso!")
	} else {
		c.Log = append(c.Log, "Nao conseguiu fugir!")
	}
}
