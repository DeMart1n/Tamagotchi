package model

import "fmt"

type SkillType int

const (
	SkillDamage SkillType = iota
	SkillHeal
	SkillBuff
	SkillDebuff
	SkillPassive
)

type Skill struct {
	ID          string
	Name        string
	Description string
	Type        SkillType
	Cost        int
	Power       int    // Multiplicador de dano/cura em %
	Element     string
	IsPassive   bool   // Não aparece no menu de combate; efeito pré-calculado
	Cooldown    int    // Turnos de espera após uso (0 = sem cooldown)
	BuffStat    string // "atk" | "def" | "vel" | "dodge" — stat alvo do buff
	BuffTurns   int    // Duração do buff em turnos
}

func (s Skill) String() string {
	return fmt.Sprintf("%s (%d MP): %s", s.Name, s.Cost, s.Description)
}

var AllSkills = map[string]Skill{
	// --- Habilidades genéricas (legado, sem classe) ---
	"pancada": {
		ID:          "pancada",
		Name:        "Pancada Pesada",
		Description: "Um golpe forte que causa 1.5x de dano.",
		Type:        SkillDamage,
		Cost:        8,
		Power:       150,
		Cooldown:    0,
	},
	"cura": {
		ID:          "cura",
		Name:        "Luz Curativa",
		Description: "Recupera 20% do HP máximo.",
		Type:        SkillHeal,
		Cost:        12,
		Power:       20,
		Cooldown:    0,
	},
	"rugido": {
		ID:          "rugido",
		Name:        "Rugido Intimidador",
		Description: "Assusta o inimigo, reduzindo sua defesa em 15%.",
		Type:        SkillDebuff,
		Cost:        10,
		Power:       15,
		Cooldown:    0,
	},
	"foco": {
		ID:          "foco",
		Name:        "Foco Interior",
		Description: "Aumenta o próprio ataque em 20% por 1 turno.",
		Type:        SkillBuff,
		Cost:        15,
		Power:       20,
		BuffStat:    "atk",
		BuffTurns:   1,
		Cooldown:    0,
	},

	// --- Guardiã ---
	"muralha": {
		ID:          "muralha",
		Name:        "Muralha de Pedra",
		Description: "Fortalece a defesa em 25% por 2 turnos.",
		Type:        SkillBuff,
		Cost:        10,
		Power:       25,
		BuffStat:    "def",
		BuffTurns:   2,
		Cooldown:    2,
	},
	"golpe_escudo": {
		ID:          "golpe_escudo",
		Name:        "Golpe de Escudo",
		Description: "Ataque com o escudo causando 120% de dano.",
		Type:        SkillDamage,
		Cost:        8,
		Power:       120,
		Cooldown:    1,
	},
	"resistencia": {
		ID:          "resistencia",
		Name:        "Resistência Natural",
		Description: "Passiva: defesa permanentemente 15% maior em combate.",
		Type:        SkillPassive,
		IsPassive:   true,
	},

	// --- Guerreiro ---
	"golpe_brutal": {
		ID:          "golpe_brutal",
		Name:        "Golpe Brutal",
		Description: "Golpe devastador que causa 200% de dano.",
		Type:        SkillDamage,
		Cost:        12,
		Power:       200,
		Cooldown:    2,
	},
	"furia": {
		ID:          "furia",
		Name:        "Fúria de Batalha",
		Description: "Entra em fúria, aumentando o ataque em 35% por 3 turnos.",
		Type:        SkillBuff,
		Cost:        14,
		Power:       35,
		BuffStat:    "atk",
		BuffTurns:   3,
		Cooldown:    3,
	},
	"instinto": {
		ID:          "instinto",
		Name:        "Instinto de Batalha",
		Description: "Passiva: +15 de Sorte, aumentando a chance de crítico.",
		Type:        SkillPassive,
		IsPassive:   true,
	},

	// --- Místico ---
	"cura_profunda": {
		ID:          "cura_profunda",
		Name:        "Cura Profunda",
		Description: "Cura poderosa que restaura 35% do HP máximo.",
		Type:        SkillHeal,
		Cost:        18,
		Power:       35,
		Cooldown:    2,
	},
	"barreira": {
		ID:          "barreira",
		Name:        "Barreira Arcana",
		Description: "Cria uma barreira mágica que aumenta a defesa em 20% por 3 turnos.",
		Type:        SkillBuff,
		Cost:        12,
		Power:       20,
		BuffStat:    "def",
		BuffTurns:   3,
		Cooldown:    3,
	},
	"fluxo_arcano": {
		ID:          "fluxo_arcano",
		Name:        "Fluxo Arcano",
		Description: "Passiva: regenera 3 MP por turno em combate.",
		Type:        SkillPassive,
		IsPassive:   true,
	},

	// --- Rápida ---
	"corte_veloz": {
		ID:          "corte_veloz",
		Name:        "Corte Veloz",
		Description: "Ataque rápido que causa 110% de dano. Sem cooldown.",
		Type:        SkillDamage,
		Cost:        4,
		Power:       110,
		Cooldown:    0,
	},
	"esquiva": {
		ID:          "esquiva",
		Name:        "Esquiva Ágil",
		Description: "Aumenta a chance de esquivar ataques em 30% por 2 turnos.",
		Type:        SkillBuff,
		Cost:        10,
		Power:       30,
		BuffStat:    "dodge",
		BuffTurns:   2,
		Cooldown:    3,
	},
	"reflexos": {
		ID:          "reflexos",
		Name:        "Reflexos Apurados",
		Description: "Passiva: +20 de Sorte, aumentando críticos e esquivas.",
		Type:        SkillPassive,
		IsPassive:   true,
	},
}

// GetInitialSkills retorna habilidades genéricas baseadas no estágio (fallback para saves sem classe).
func GetInitialSkills(stage Stage) []string {
	switch stage {
	case StageBaby:
		return []string{"pancada"}
	case StageChild:
		return []string{"pancada", "cura"}
	case StageTeen:
		return []string{"pancada", "cura", "rugido"}
	default:
		return []string{"pancada", "cura", "rugido", "foco"}
	}
}

// GetActiveSkills filtra a lista de IDs retornando apenas habilidades ativas (não passivas).
func GetActiveSkills(known []string) []string {
	active := make([]string, 0, len(known))
	for _, id := range known {
		if sk, ok := AllSkills[id]; ok && !sk.IsPassive {
			active = append(active, id)
		}
	}
	return active
}
