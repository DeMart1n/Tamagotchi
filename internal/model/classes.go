package model

type ClassID string

const (
	ClassNone      ClassID = ""
	ClassGuardia   ClassID = "guardia"
	ClassGuerreiro ClassID = "guerreiro"
	ClassMistico   ClassID = "mistico"
	ClassRapida    ClassID = "rapida"
)

type Class struct {
	ID          ClassID
	Name        string
	Description string
	HPMult      float64
	ATKMult     float64
	DEFMult     float64
	VELMult     float64
	LUCKMult    float64
	MPMult      float64
	BaseSkills  []string // [ativo1, ativo2, passiva]
}

var AllClasses = map[ClassID]Class{
	ClassGuardia: {
		ID:          ClassGuardia,
		Name:        "Guardiã",
		Description: "Tank resistente. Alta defesa e muita vida, mas golpes mais fracos.",
		HPMult:      1.30,
		ATKMult:     0.90,
		DEFMult:     1.20,
		VELMult:     0.90,
		LUCKMult:    1.00,
		MPMult:      1.00,
		BaseSkills:  []string{"muralha", "golpe_escudo", "resistencia"},
	},
	ClassGuerreiro: {
		ID:          ClassGuerreiro,
		Name:        "Guerreiro",
		Description: "Dano bruto e instinto de batalha. Ataque elevado e bônus de crítico.",
		HPMult:      1.00,
		ATKMult:     1.30,
		DEFMult:     1.00,
		VELMult:     1.00,
		LUCKMult:    1.00,
		MPMult:      1.00,
		BaseSkills:  []string{"golpe_brutal", "furia", "instinto"},
	},
	ClassMistico: {
		ID:          ClassMistico,
		Name:        "Místico",
		Description: "Suporte arcano. Mais MP, curas poderosas e barreiras mágicas.",
		HPMult:      1.00,
		ATKMult:     0.90,
		DEFMult:     1.10,
		VELMult:     1.00,
		LUCKMult:    1.10,
		MPMult:      1.40,
		BaseSkills:  []string{"cura_profunda", "barreira", "fluxo_arcano"},
	},
	ClassRapida: {
		ID:          ClassRapida,
		Name:        "Rápida",
		Description: "Velocidade e críticos. Alta sorte e esquiva, mas menos HP e DEF.",
		HPMult:      0.90,
		ATKMult:     1.00,
		DEFMult:     0.90,
		VELMult:     1.30,
		LUCKMult:    1.40,
		MPMult:      1.00,
		BaseSkills:  []string{"corte_veloz", "esquiva", "reflexos"},
	},
}

// ClassOrder define a ordem de exibição na tela de seleção.
var ClassOrder = []ClassID{ClassGuardia, ClassGuerreiro, ClassMistico, ClassRapida}
