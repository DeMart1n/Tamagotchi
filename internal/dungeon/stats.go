package dungeon

import "Pessoal/internal/model"

// CombatStats representa os atributos de combate derivados do Tama.
type CombatStats struct {
	HPMax      int
	HPCurrent  int
	MPMax      int
	MPCurrent  int
	Ataque     int
	Defesa     int
	Velocidade int
	Sorte      int
}

// stageMult retorna o multiplicador de stats por estágio de evolução.
func stageMult(s model.Stage) float64 {
	switch s {
	case model.StageBaby:
		return 0.6
	case model.StageChild:
		return 0.8
	case model.StageTeen:
		return 1.0
	case model.StageAdult:
		return 1.2
	case model.StageElder:
		return 1.1
	default:
		return 1.0
	}
}

// DeriveCombatStats calcula os stats de combate a partir do estado atual do Tama.
func DeriveCombatStats(tama *model.Tama) CombatStats {
	mult := stageMult(tama.Stage)

	hp := int(float64(50+tama.Level*10) * mult)
	mp := int(float64(20+tama.Level*5) * mult)
	atk := int(float64(8+tama.Level*3) * mult)
	def := int(float64(3+tama.Level*2) * mult)
	vel := 5 + tama.Level
	luck := 5 + tama.Level/2

	// Bônus por estado emocional
	if tama.Happiness > 70 {
		atk += 3
		vel += 2
	}
	if tama.Hunger > 60 {
		hp += 10
	}

	// Multiplicadores de classe
	if cls, ok := model.AllClasses[tama.Class]; ok {
		hp   = int(float64(hp)   * cls.HPMult)
		mp   = int(float64(mp)   * cls.MPMult)
		atk  = int(float64(atk)  * cls.ATKMult)
		def  = int(float64(def)  * cls.DEFMult)
		vel  = int(float64(vel)  * cls.VELMult)
		luck = int(float64(luck) * cls.LUCKMult)
	}

	// Passivas de stat aplicadas em derivação (resistencia, instinto, reflexos)
	for _, id := range tama.SkillsKnown {
		sk, ok := model.AllSkills[id]
		if !ok || !sk.IsPassive {
			continue
		}
		switch id {
		case "resistencia":
			def = int(float64(def) * 1.15)
		case "instinto":
			luck += 15
		case "reflexos":
			luck += 20
		}
	}

	return CombatStats{
		HPMax:      hp,
		HPCurrent:  hp,
		MPMax:      mp,
		MPCurrent:  mp,
		Ataque:     atk,
		Defesa:     def,
		Velocidade: vel,
		Sorte:      luck,
	}
}

// ApplyEquipment aplica os bônus de equipamento aos stats de combate.
func (cs CombatStats) ApplyEquipment(inv *Inventory) CombatStats {
	if inv == nil {
		return cs
	}
	for _, eq := range []*Equipment{inv.Arma, inv.Armadura, inv.Acessorio} {
		if eq == nil {
			continue
		}
		cs.HPMax += eq.BonusHP
		cs.HPCurrent += eq.BonusHP
		cs.Ataque += eq.BonusATK
		cs.Defesa += eq.BonusDEF
		cs.Velocidade += eq.BonusVEL
		cs.Sorte += eq.BonusLuck
	}
	return cs
}

// ApplyBiomeModifiers aplica os bônus/penalidades do bioma aos stats de combate.
func (cs CombatStats) ApplyBiomeModifiers(mod BiomeModifier) CombatStats {
	if mod.PlayerHPBonusPct != 0 {
		bonus := (cs.HPMax * mod.PlayerHPBonusPct) / 100
		cs.HPMax += bonus
		cs.HPCurrent += bonus
	}
	if mod.PlayerATKBonusPct != 0 {
		cs.Ataque += (cs.Ataque * mod.PlayerATKBonusPct) / 100
	}
	if mod.PlayerDEFBonusPct != 0 {
		cs.Defesa += (cs.Defesa * mod.PlayerDEFBonusPct) / 100
	}
	if mod.PlayerVELBonusPct != 0 {
		cs.Velocidade += (cs.Velocidade * mod.PlayerVELBonusPct) / 100
	}
	if mod.PlayerLuckBonusPct != 0 {
		cs.Sorte += (cs.Sorte * mod.PlayerLuckBonusPct) / 100
	}

	// Garantir valores mínimos
	if cs.Ataque < 1 {
		cs.Ataque = 1
	}
	if cs.Defesa < 0 {
		cs.Defesa = 0
	}
	if cs.Velocidade < 1 {
		cs.Velocidade = 1
	}
	if cs.Sorte < 0 {
		cs.Sorte = 0
	}

	return cs
}
