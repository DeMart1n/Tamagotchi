package dungeon

import (
	"Pessoal/internal/model"
	"fmt"
	"math/rand"
)

// CombatAction ação do jogador no combate.
type CombatAction int

const (
	ActionAtacar CombatAction = iota
	ActionDefender
	ActionItem
	ActionSkill
	ActionFugir
)

// CombatState estado da máquina de combate.
type CombatState int

const (
	CombatChoosing CombatState = iota
	CombatResolved
)

// ActiveBuff rastreia um buff temporário ativo durante o combate.
type ActiveBuff struct {
	Stat   string // "atk" | "def" | "vel" | "dodge"
	Amount int    // valor absoluto adicionado ao stat (usado para reverter)
	Turns  int    // turnos restantes
}

// Combat motor de combate por turnos.
type Combat struct {
	Player        *CombatStats
	OriginalStats *CombatStats
	Enemy         *Enemy
	Biome         Biome
	Modifier      BiomeModifier
	State         CombatState
	Log           []string
	Turn          int
	Defending     bool
	FrozenFor     int
	Fled          bool
	Won           bool
	Lost          bool

	ActiveBuffs    []ActiveBuff
	SkillCooldowns map[string]int
	DodgeChance    int // % de chance de esquivar o próximo ataque inimigo
	PassiveMPRegen int // MP recuperado por turno (passiva fluxo_arcano)

	BossPhase   int  // 0-3; apenas significativo quando Enemy.IsBoss == true
	BossEnraged bool
}

// NewCombat cria um novo combate.
func NewCombat(player *CombatStats, enemy *Enemy, biome Biome, tama *model.Tama) *Combat {
	modifier := ModifierForBiome(biome)
	statsWithMod := player.ApplyBiomeModifiers(modifier)

	c := &Combat{
		Player:         &statsWithMod,
		OriginalStats:  player,
		Enemy:          enemy,
		Biome:          biome,
		Modifier:       modifier,
		State:          CombatChoosing,
		Log:            []string{},
		Turn:           1,
		SkillCooldowns: make(map[string]int),
	}

	if tama != nil {
		for _, id := range tama.SkillsKnown {
			if id == "fluxo_arcano" {
				c.PassiveMPRegen = 3
				break
			}
		}
	}

	return c
}

// IsOver retorna se o combate terminou.
func (c *Combat) IsOver() bool {
	return c.Won || c.Lost || c.Fled
}

// ExecuteAction processa a ação do jogador e a resposta do inimigo.
func (c *Combat) ExecuteAction(action CombatAction, itemBag *ItemBag, tama *model.Tama, actionIndex int) {
	c.Log = []string{}
	c.Defending = false

	c.tickActiveBuffs()
	c.applyBiomeTurnEffects()

	if c.Player.HPCurrent <= 0 {
		c.Player.HPCurrent = 0
		c.Lost = true
		c.Log = append(c.Log, "Voce sucumbiu aos efeitos do bioma...")
		c.syncStats()
		c.State = CombatResolved
		return
	}

	if c.FrozenFor > 0 {
		c.FrozenFor--
		c.Log = append(c.Log, "Voce esta congelado e perdeu o turno!")
	} else {
		switch action {
		case ActionAtacar:
			c.playerAttack()
		case ActionDefender:
			c.Defending = true
			c.Log = append(c.Log, "Voce se preparou para defender!")
		case ActionItem:
			c.useItem(itemBag, actionIndex)
		case ActionSkill:
			c.executeSkill(tama, actionIndex)
		case ActionFugir:
			c.tryFlee()
			if c.Fled {
				c.syncStats()
				c.State = CombatResolved
				return
			}
		}
	}

	c.checkBossPhaseTransition()

	if c.Enemy.HPCurrent <= 0 {
		c.Enemy.HPCurrent = 0
		c.Won = true
		c.Log = append(c.Log, fmt.Sprintf("%s foi derrotado!", c.Enemy.Name))
		c.syncStats()
		c.State = CombatResolved
		return
	}

	c.enemyTurn()

	if c.Player.HPCurrent <= 0 {
		c.Player.HPCurrent = 0
		c.Lost = true
		c.Log = append(c.Log, "Voce foi derrotado...")
	}

	c.syncStats()
	c.Turn++
	c.State = CombatResolved
}

func (c *Combat) syncStats() {
	if c.OriginalStats != nil {
		c.OriginalStats.HPCurrent = c.Player.HPCurrent
		c.OriginalStats.MPCurrent = c.Player.MPCurrent
	}
}

// tickActiveBuffs decrementa duração dos buffs ativos e reverte os expirados.
// Também decrementa cooldowns de habilidades. Chamado no início de cada turno.
func (c *Combat) tickActiveBuffs() {
	remaining := c.ActiveBuffs[:0]
	for _, b := range c.ActiveBuffs {
		b.Turns--
		if b.Turns <= 0 {
			switch b.Stat {
			case "atk":
				c.Player.Ataque -= b.Amount
			case "def":
				c.Player.Defesa -= b.Amount
			case "vel":
				c.Player.Velocidade -= b.Amount
			case "dodge":
				c.DodgeChance -= b.Amount
				if c.DodgeChance < 0 {
					c.DodgeChance = 0
				}
			}
		} else {
			remaining = append(remaining, b)
		}
	}
	c.ActiveBuffs = remaining

	for id, cd := range c.SkillCooldowns {
		if cd > 0 {
			c.SkillCooldowns[id] = cd - 1
		}
	}
}

func (c *Combat) playerAttack() {
	baseDmg := c.Player.Ataque - c.Enemy.Defesa/2
	if baseDmg < 1 {
		baseDmg = 1
	}

	variance := float64(baseDmg) * 0.2
	dmg := baseDmg + int(float64(rand.Intn(int(variance*2+1)))-variance)
	if dmg < 1 {
		dmg = 1
	}

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
	hpPct := float64(c.Enemy.HPCurrent) / float64(c.Enemy.HPMax)
	if hpPct < 0.25 && rand.Intn(100) < 30 {
		c.Enemy.Defending = true
		c.Log = append(c.Log, fmt.Sprintf("%s se preparou para defender!", c.Enemy.Name))
		return
	}

	// Rola esquiva antes de calcular dano
	if c.DodgeChance > 0 && rand.Intn(100) < c.DodgeChance {
		c.Log = append(c.Log, "Voce esquivou do ataque!")
		return
	}

	baseDmg := c.Enemy.Ataque - c.Player.Defesa/2
	if baseDmg < 1 {
		baseDmg = 1
	}

	if c.Modifier.EnemyAttackBonusPct > 0 {
		baseDmg += (baseDmg * c.Modifier.EnemyAttackBonusPct) / 100
	}

	variance := float64(baseDmg) * 0.2
	dmg := baseDmg + int(float64(rand.Intn(int(variance*2+1)))-variance)
	if dmg < 1 {
		dmg = 1
	}

	if c.Defending {
		dmg = dmg / 2
		if dmg < 1 {
			dmg = 1
		}
	}

	if c.Enemy.IsFire && c.Modifier.PlayerFireVulnerability > 0 {
		dmg += (dmg * c.Modifier.PlayerFireVulnerability) / 100
	}

	critChance := 5 + c.Modifier.EnemyLuckBonusPct
	if rand.Intn(100) < critChance {
		dmg = int(float64(dmg) * 1.5)
		c.Log = append(c.Log, fmt.Sprintf("CRITICO inimigo! %s acertou um golpe certeiro!", c.Enemy.Name))
	}

	c.Player.HPCurrent -= dmg
	c.Log = append(c.Log, fmt.Sprintf("%s atacou! %d de dano!", c.Enemy.Name, dmg))

	if c.Modifier.FreezeChancePct > 0 && rand.Intn(100) < c.Modifier.FreezeChancePct {
		c.FrozenFor = 1
		c.Log = append(c.Log, "Voce foi congelado!")
	}
}

func (c *Combat) applyBiomeTurnEffects() {
	if c.Modifier.PlayerFireDotPctMaxHP > 0 {
		dot := (c.Player.HPMax * c.Modifier.PlayerFireDotPctMaxHP) / 100
		if dot < 1 {
			dot = 1
		}
		c.Player.HPCurrent -= dot
		c.Log = append(c.Log, fmt.Sprintf("O calor do bioma causa %d de dano continuo!", dot))
	}

	if c.Modifier.PlayerHPRegenBonus != 0 {
		c.Player.HPCurrent += c.Modifier.PlayerHPRegenBonus
		if c.Modifier.PlayerHPRegenBonus > 0 {
			c.Log = append(c.Log, fmt.Sprintf("A aura do bioma recupera %d HP!", c.Modifier.PlayerHPRegenBonus))
		} else {
			c.Log = append(c.Log, fmt.Sprintf("O ambiente hostil drena %d HP!", -c.Modifier.PlayerHPRegenBonus))
		}
	}

	// Passiva de MP regen (fluxo_arcano)
	if c.PassiveMPRegen > 0 {
		c.Player.MPCurrent += c.PassiveMPRegen
		if c.Player.MPCurrent > c.Player.MPMax {
			c.Player.MPCurrent = c.Player.MPMax
		}
	}

	if c.Player.HPCurrent > c.Player.HPMax {
		c.Player.HPCurrent = c.Player.HPMax
	}
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
		amount := item.Value
		c.Player.Ataque += amount
		c.ActiveBuffs = append(c.ActiveBuffs, ActiveBuff{Stat: "atk", Amount: amount, Turns: 1})
		c.Log = append(c.Log, fmt.Sprintf("Usou %s! ATK +%d neste turno!", item.Name, amount))
	}
	bag.Remove(index)
}

func (c *Combat) tryFlee() {
	if c.Enemy.IsBoss {
		c.Log = append(c.Log, "Impossivel fugir do BOSS!")
		return
	}

	playerSpeed := c.Player.Velocidade
	if c.Modifier.PlayerSpeedPenaltyPct > 0 {
		penalty := (playerSpeed * c.Modifier.PlayerSpeedPenaltyPct) / 100
		playerSpeed -= penalty
		if playerSpeed < 1 {
			playerSpeed = 1
		}
	}

	chance := 40 + (playerSpeed-c.Enemy.Velocidade)*2
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

func (c *Combat) executeSkill(tama *model.Tama, skillIndex int) {
	if tama == nil {
		c.Log = append(c.Log, "Habilidade invalida!")
		return
	}

	activeSkills := model.GetActiveSkills(tama.SkillsKnown)
	if skillIndex < 0 || skillIndex >= len(activeSkills) {
		c.Log = append(c.Log, "Habilidade invalida!")
		return
	}

	skillID := activeSkills[skillIndex]
	skill, ok := model.AllSkills[skillID]
	if !ok {
		c.Log = append(c.Log, "Habilidade desconhecida!")
		return
	}

	if cd := c.SkillCooldowns[skillID]; cd > 0 {
		c.Log = append(c.Log, fmt.Sprintf("%s em recarga! (%d turno(s))", skill.Name, cd))
		return
	}

	if c.Player.MPCurrent < skill.Cost {
		c.Log = append(c.Log, "MP insuficiente!")
		return
	}

	c.Player.MPCurrent -= skill.Cost
	c.Log = append(c.Log, fmt.Sprintf("Usou %s!", skill.Name))

	if skill.Cooldown > 0 {
		c.SkillCooldowns[skillID] = skill.Cooldown
	}

	switch skill.Type {
	case model.SkillDamage:
		dmg := int(float64(c.Player.Ataque) * (float64(skill.Power) / 100.0))
		dmg -= c.Enemy.Defesa / 2
		if dmg < 1 {
			dmg = 1
		}
		c.Enemy.HPCurrent -= dmg
		c.Log = append(c.Log, fmt.Sprintf("Causou %d de dano em %s!", dmg, c.Enemy.Name))

	case model.SkillHeal:
		heal := (c.Player.HPMax * skill.Power) / 100
		c.Player.HPCurrent += heal
		if c.Player.HPCurrent > c.Player.HPMax {
			c.Player.HPCurrent = c.Player.HPMax
		}
		c.Log = append(c.Log, fmt.Sprintf("Recuperou %d de HP!", heal))

	case model.SkillBuff:
		var amount int
		switch skill.BuffStat {
		case "atk":
			amount = (c.Player.Ataque * skill.Power) / 100
			c.Player.Ataque += amount
		case "def":
			amount = (c.Player.Defesa * skill.Power) / 100
			c.Player.Defesa += amount
		case "vel":
			amount = (c.Player.Velocidade * skill.Power) / 100
			c.Player.Velocidade += amount
		case "dodge":
			amount = skill.Power
			c.DodgeChance += amount
		}
		turns := skill.BuffTurns
		if turns < 1 {
			turns = 1
		}
		c.ActiveBuffs = append(c.ActiveBuffs, ActiveBuff{Stat: skill.BuffStat, Amount: amount, Turns: turns})
		c.Log = append(c.Log, fmt.Sprintf("Buff ativo por %d turno(s)!", turns))

	case model.SkillDebuff:
		debuff := (c.Enemy.Defesa * skill.Power) / 100
		c.Enemy.Defesa -= debuff
		c.Log = append(c.Log, fmt.Sprintf("Defesa de %s reduzida!", c.Enemy.Name))
	}
}

// checkBossPhaseTransition verifica transições de fase para chefes.
func (c *Combat) checkBossPhaseTransition() {
	if !c.Enemy.IsBoss {
		return
	}

	hpPct := float64(c.Enemy.HPCurrent) / float64(c.Enemy.HPMax)

	if c.BossPhase == 0 && hpPct <= 0.75 {
		c.BossPhase = 1
		c.applyBossPhaseEffects(1)
		return
	}

	if c.BossPhase == 1 && hpPct <= 0.50 {
		c.BossPhase = 2
		c.applyBossPhaseEffects(2)
		return
	}

	if c.BossPhase == 2 && hpPct <= 0.25 && !c.BossEnraged {
		c.BossPhase = 3
		c.BossEnraged = true
		c.applyBossPhaseEffects(3)
		return
	}
}

// applyBossPhaseEffects aplica efeitos específicos de cada fase do chefe.
func (c *Combat) applyBossPhaseEffects(newPhase int) {
	switch c.Enemy.Name {
	case "Guardiao da Floresta":
		switch newPhase {
		case 1:
			amount := c.Player.Velocidade / 5
			if amount < 1 {
				amount = 1
			}
			c.Player.Velocidade -= amount
			c.ActiveBuffs = append(c.ActiveBuffs, ActiveBuff{Stat: "vel", Amount: amount, Turns: 999})
			c.Log = append(c.Log, "FASE 2: O Guardiao cria raizes ao seu redor!")

		case 2:
			c.Enemy.Ataque += 6
			c.Log = append(c.Log, "FASE 3: Golpe das Raizes!")

		case 3:
			c.Enemy.Ataque += 8
			c.Log = append(c.Log, "ENRAIVECIDO! Furia Ancestral!")
		}

	case "Lich do Gelo Eterno":
		switch newPhase {
		case 1:
			c.FrozenFor = 1
			c.Log = append(c.Log, "FASE 2: O Lich lanca aura de gelo!")

		case 2:
			heal := c.Enemy.HPMax / 8
			c.Enemy.HPCurrent += heal
			if c.Enemy.HPCurrent > c.Enemy.HPMax {
				c.Enemy.HPCurrent = c.Enemy.HPMax
			}
			c.Log = append(c.Log, "FASE 3: O Lich se revitaliza!")

		case 3:
			c.Enemy.Ataque += 10
			c.FrozenFor = 1
			c.Log = append(c.Log, "ENRAIVECIDO! Maldicao Eterna!")
		}

	case "Lorde das Chamas":
		switch newPhase {
		case 1:
			c.Enemy.Ataque += 5
			c.Log = append(c.Log, "FASE 2: Onda de Calor!")

		case 2:
			c.Enemy.Ataque += 5
			c.Log = append(c.Log, "FASE 3: Erupcao de Magma!")

		case 3:
			c.Enemy.Ataque += 8
			c.Enemy.Velocidade += 3
			c.Log = append(c.Log, "ENRAIVECIDO! Inferno Total!")
		}

	case "Devorador do Abismo":
		switch newPhase {
		case 1:
			amount := c.Player.Ataque / 5
			if amount < 1 {
				amount = 1
			}
			c.Player.Ataque -= amount
			c.ActiveBuffs = append(c.ActiveBuffs, ActiveBuff{Stat: "atk", Amount: amount, Turns: 999})
			c.Log = append(c.Log, "FASE 2: Toque do Abismo drena seu poder!")

		case 2:
			amount := c.Player.Defesa / 5
			if amount < 1 {
				amount = 1
			}
			c.Player.Defesa -= amount
			c.ActiveBuffs = append(c.ActiveBuffs, ActiveBuff{Stat: "def", Amount: amount, Turns: 999})
			c.Log = append(c.Log, "FASE 3: Pulso do Vazio corroi sua defesa!")

		case 3:
			c.Enemy.Ataque += 12
			c.Enemy.Defesa -= c.Enemy.Defesa / 3
			c.Log = append(c.Log, "ENRAIVECIDO! Colapso do Abismo!")
		}
	}
}
