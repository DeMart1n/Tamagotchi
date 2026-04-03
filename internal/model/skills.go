package model

import "fmt"

type SkillType int

const (
	SkillDamage SkillType = iota
	SkillHeal
	SkillBuff
	SkillDebuff
)

type Skill struct {
	ID          string
	Name        string
	Description string
	Type        SkillType
	Cost        int
	Power       int // Multiplicador de dano ou cura
	Element     string
}

func (s Skill) String() string {
	return fmt.Sprintf("%s (%d MP): %s", s.Name, s.Cost, s.Description)
}

// AllSkills base de dados de habilidades disponíveis no jogo.
var AllSkills = map[string]Skill{
	"pancada": {
		ID:          "pancada",
		Name:        "Pancada Pesada",
		Description: "Um golpe forte que causa 1.5x de dano.",
		Type:        SkillDamage,
		Cost:        8,
		Power:       150, // 150% do ATK
	},
	"cura": {
		ID:          "cura",
		Name:        "Luz Curativa",
		Description: "Recupera uma pequena quantidade de HP.",
		Type:        SkillHeal,
		Cost:        12,
		Power:       20, // 20% do HP Max
	},
	"rugido": {
		ID:          "rugido",
		Name:        "Rugido Intimidador",
		Description: "Assusta o inimigo, reduzindo sua defesa.",
		Type:        SkillDebuff,
		Cost:        10,
		Power:       15, // -15% DEF inimiga
	},
	"foco": {
		ID:          "foco",
		Name:        "Foco Interior",
		Description: "Aumenta o proprio ataque temporariamente.",
		Type:        SkillBuff,
		Cost:        15,
		Power:       20, // +20% ATK
	},
}

// GetInitialSkills retorna habilidades baseadas no estágio de evolução.
func GetInitialSkills(stage Stage) []string {
	switch stage {
	case StageBaby:
		return []string{"pancada"}
	case StageChild:
		return []string{"pancada", "cura"}
	case StageTeen:
		return []string{"pancada", "cura", "rugido"}
	case StageAdult:
		return []string{"pancada", "cura", "rugido", "foco"}
	default:
		return []string{"pancada"}
	}
}
