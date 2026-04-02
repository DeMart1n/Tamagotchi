package dungeon

import "github.com/DeMart1n/Tamagotchi/internal/model"

// CombatStats representa os atributos de combate derivados do Tama.
type CombatStats struct {
	HPMax      int
	HPCurrent  int
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

	return CombatStats{
		HPMax:      hp,
		HPCurrent:  hp,
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
