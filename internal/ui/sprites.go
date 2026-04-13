package ui

import (
	"Pessoal/internal/model"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Pixel art sprites usando half-block characters (▀▄█) com cores ANSI truecolor.
// Cada par de linhas da grid vira 1 linha de terminal (topo = fg, base = bg + "▀").
//
// Paleta (max 4 cores por sprite + contorno):
//   K = contorno escuro       W = branco (olhos/brilho)
//   G = verde (corpo base)    L = verde claro (barriga/destaque)
//   Y = amarelo (boca/alegria) R = vermelho (raiva/perigo)
//   P = rosa (bochechas)      B = azul (lágrimas/tristeza)
//   C = ciano (brilho/magia)  O = laranja (energia)
//   D = cinza escuro (sombra) M = marrom (elder/maduro)
//   . = transparente

var spriteColors = map[byte]string{
	'K': "#16161a", // contorno
	'W': "#f0f0f0", // branco
	'G': "#3dbd6b", // verde corpo
	'L': "#7edd9a", // verde claro (barriga)
	'Y': "#ffd166", // amarelo quente
	'R': "#ef476f", // vermelho coral
	'P': "#ff9ecd", // rosa bochechas
	'B': "#118ab2", // azul
	'C': "#73d2f3", // ciano brilho
	'O': "#f78c3a", // laranja
	'D': "#4a4a5a", // cinza escuro
	'M': "#8b6f4e", // marrom
}

type spriteFrame struct {
	rows []string
}

// ============================================================
// BABY — bolinha redonda, olhos grandes, corpo pequeno e fofo
// Grid: 11 cols x 10 rows → renderiza como 11x5 chars
// ============================================================

var babyFrames = [2]spriteFrame{
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGWWKWWGK.",
		".KGPGGGPGK.",
		".KKGLLLGKK.",
		"..KGLYLGK..",
		"..KKGGGKK..",
		"...K.K.K...",
		"...........",
	}},
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGWWKWWGK.",
		".KGPGGGPGK.",
		".KKGLLLGKK.",
		"..KGLYLGK..",
		"..KKGGGKK..",
		"...........",
		"...K...K...",
	}},
}

// ============================================================
// CHILD — mais alto, orelhas visíveis (2px), expressão curiosa
// Grid: 11 cols x 10 rows
// ============================================================

var childFrames = [2]spriteFrame{
	{rows: []string{
		"GK.KKKKK.KG",
		"KGK.....KGK",
		".KKGGGGGKK.",
		".KGGGGGGGK.",
		".KGWWKWWGK.",
		".KGGGGGGGK.",
		".KKGLYLGKK.",
		"..KGLLLGK..",
		"..KK.K.KK..",
		"...........",
	}},
	{rows: []string{
		"GK.KKKKK.KG",
		"KGK.....KGK",
		".KKGGGGGKK.",
		".KGGGGGGGK.",
		".KGWWKWWGK.",
		".KGGKGKGGK.",
		".KKGLYLGKK.",
		"..KGLLLGK..",
		"...........",
		"..KK...KK..",
	}},
}

// ============================================================
// TEEN — corpo angular, chifres coloridos (O), atitude rebelde
// Grid: 11 cols x 10 rows
// ============================================================

var teenFrames = [2]spriteFrame{
	{rows: []string{
		"O..KKKKK..O",
		".KKGGGGGKK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGGGGGGGGGK",
		"KKGGLLLGGKK",
		".KGGLYLGGK.",
		".KKGGGGGKK.",
		".KK..K..KK.",
		"...........",
	}},
	{rows: []string{
		"O..KKKKK..O",
		".KKGGGGGKK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGGGKGKGGGK",
		"KKGGLLLGGKK",
		".KGGLYLGGK.",
		".KKGGGGGKK.",
		"...........",
		".KK.....KK.",
	}},
}

// ============================================================
// ADULT — corpo robusto, coroa/crista, presença imponente
// Grid: 11 cols x 10 rows
// ============================================================

var adultFrames = [2]spriteFrame{
	{rows: []string{
		"..K.KYK.K..",
		"..KKKKKKK..",
		".KGGGGGGGK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGPGGGGGPGK",
		"KKGGLLLLGKK",
		".KGGLYLGGK.",
		".KKGGGGGKK.",
		".KK..K..KK.",
	}},
	{rows: []string{
		"...KYKYK...",
		"..KKKKKKK..",
		".KGGGGGGGK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGPGGGGGPGK",
		"KKGGLLLLGKK",
		".KGGLYLGGK.",
		".KKGGGGGKK.",
		"KK.......KK",
	}},
}

// ============================================================
// ELDER — corpo arredondado, "barba"/textura, sábio e sereno
// Grid: 11 cols x 10 rows
// ============================================================

var elderFrames = [2]spriteFrame{
	{rows: []string{
		"..KKKKKKK..",
		".KMGGGGGMK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGGGGGGGGGK",
		"KKGGLLLGGKK",
		".KDDLYLDDK.",
		".KKDDDDDKK.",
		"..KK.K.KK..",
		"...........",
	}},
	{rows: []string{
		"..KKKKKKK..",
		".KMGGGGGMK.",
		"KGGGGGGGGGK",
		"KGGWWKWWGGK",
		"KGGGKGKGGGK",
		"KKGGLLLGGKK",
		".KDDLYLDDK.",
		".KKDDDDDKK.",
		"...........",
		"..KK...KK..",
	}},
}

// ============================================================
// ESTADOS EMOCIONAIS — sobreescrevem sprite de estágio
// ============================================================

// DORMINDO — olhos fechados, corpo encolhido, Zs grandes flutuantes
var sleepingFrames = [2]spriteFrame{
	{rows: []string{
		"...KKKKK.CC",
		"..KGGGGGK.C",
		".KGGGGGGGK.",
		".KGKKGKKGK.",
		".KGGGGGGGK.",
		".KKGLLLGKK.",
		"..KGGGGGK..",
		"..KKKKKKK..",
		"...........",
		"...........",
	}},
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGKCC",
		".KGGGGGGGKC",
		".KGKKGKKGK.",
		".KGGGGGGGK.",
		".KKGLLLGKK.",
		"..KGGGGGK..",
		"..KKKKKKK..",
		"...........",
		"...........",
	}},
}

// FELIZ — olhos brilhantes (estrelas), boca larga, partículas
var happyFrames = [2]spriteFrame{
	{rows: []string{
		"Y..KKKKK..Y",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGYWKYWGK.",
		".KGPGGGPGK.",
		".KKGYYYYGKK",
		"..KGLLLGK..",
		"..KKGGGKK..",
		"..KK.K.KK..",
		"Y.........Y",
	}},
	{rows: []string{
		"...KKKKK...",
		".YKGGGGGKY.",
		".KGGGGGGGK.",
		".KGYWKYWGK.",
		".KGPGGGPGK.",
		".KKGYYYYGKK",
		"..KGLLLGK..",
		"..KKGGGKK..",
		"...........",
		"..KK...KK..",
	}},
}

// TRISTE/DEPRIMIDO — olhos lacrimejantes, boca triste, gotas visíveis
var sadFrames = [2]spriteFrame{
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGBWKBWGK.",
		".KGBGGGBGK.",
		".KKGKKKGKK.",
		"..KGLLLGK..",
		"..KKGGGKK..",
		"..B..K..B..",
		"...........",
	}},
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGBWKBWGK.",
		".KGBGGGBGK.",
		".KKGKKKGKK.",
		"..KGLLLGK..",
		"..KKGGGKK..",
		"...........",
		".B...K...B.",
	}},
}

// RAIVOSO — sobrancelhas cerradas, boca mostrando dentes, aura vermelha
var angryFrames = [2]spriteFrame{
	{rows: []string{
		"R..KKKKK..R",
		"..KGGGGGK..",
		".KKGGGGKKK.",
		".KGRWKRWGK.",
		".KGGGGGGGK.",
		".KKGKYKGKK.",
		"..KGGGGGK..",
		"..KKGGGKK..",
		"..KK.K.KK..",
		"...........",
	}},
	{rows: []string{
		"...KKKKK...",
		".RKGGGGGKR.",
		".KKGGGGKKK.",
		".KGRWKRWGK.",
		".KGGGGGGGK.",
		".KKGKYKGKK.",
		"..KGGGGGK..",
		"..KKGGGKK..",
		"...........",
		"..KK...KK..",
	}},
}

// MORTO — corpo cinza, olhos em X, sem vida
var deadFrame = spriteFrame{
	rows: []string{
		"...KKKKK...",
		"..KDDDDDK..",
		".KDDDDDDDK.",
		".KDKDKDKDK.",
		".KDDDDDDDDK",
		".KKDKKKDKK.",
		"..KDDDDDK..",
		"..KKKKKKK..",
		"...........",
		"...........",
	},
}

// FAMINTO — olhos grandes suplicantes, barriga vazia (listras), boca aberta
var hungryFrames = [2]spriteFrame{
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KWWWKWWWK.",
		".KGGGGGGGK.",
		".KKGKOKGKK.",
		"..KGKGKGK..",
		"..KKGKGKK..",
		"..KK.K.KK..",
		"...........",
	}},
	{rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KWWWKWWWK.",
		".KGGGGGGGK.",
		".KKGKOKGKK.",
		"..KGKGKGK..",
		"..KKGKGKK..",
		"...........",
		"..KK...KK..",
	}},
}

// SOBREPESO — corpo mais largo, bochechas cheias
var overweightFrame = spriteFrame{
	rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		".KGGGGGGGK.",
		".KGWWKWWGK.",
		".KGPGGGPGK.",
		"KKGGLLLGGKK",
		"KGGGLYLGGGK",
		"KKGGGGGGGKK",
		".KKGGGGGKK.",
		".KK..K..KK.",
	},
}

// ABAIXO DO PESO — corpo mais fino
var underweightFrame = spriteFrame{
	rows: []string{
		"...KKKKK...",
		"..KGGGGGK..",
		"..KGGGGGK..",
		"..KGWKWGK..",
		"..KGGGGGK..",
		"..KKGGGKK..",
		"...KGYGK...",
		"...KKKKK...",
		"...K.K.K...",
		"...........",
	},
}

// ============================================================
// RENDERIZAÇÃO
// ============================================================

func RenderSprite(frame spriteFrame) string {
	rows := frame.rows
	if len(rows)%2 != 0 {
		rows = append(rows, strings.Repeat(".", len(rows[0])))
	}

	var lines []string
	for y := 0; y < len(rows); y += 2 {
		topRow := rows[y]
		botRow := rows[y+1]
		maxLen := len(topRow)
		if len(botRow) > maxLen {
			maxLen = len(botRow)
		}

		var line strings.Builder
		for x := 0; x < maxLen; x++ {
			topPixel := byte('.')
			botPixel := byte('.')
			if x < len(topRow) {
				topPixel = topRow[x]
			}
			if x < len(botRow) {
				botPixel = botRow[x]
			}

			topTransparent := topPixel == '.'
			botTransparent := botPixel == '.'

			if topTransparent && botTransparent {
				line.WriteString(" ")
			} else if topTransparent {
				c := spriteColors[botPixel]
				line.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render("▄"))
			} else if botTransparent {
				c := spriteColors[topPixel]
				line.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render("▀"))
			} else {
				tc := spriteColors[topPixel]
				bc := spriteColors[botPixel]
				line.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color(tc)).
					Background(lipgloss.Color(bc)).Render("▀"))
			}
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

// GetSpriteForState retorna o sprite adequado para o estado atual do Tama.
func GetSpriteForState(tama *model.Tama, frame int) string {
	idx := frame % 2

	switch {
	case tama.Dead:
		return RenderSprite(deadFrame)
	case tama.Sleeping:
		return RenderSprite(sleepingFrames[idx])
	case tama.PissedOf:
		return RenderSprite(angryFrames[idx])
	case tama.Depressed:
		return RenderSprite(sadFrames[idx])
	case tama.Overweight:
		return RenderSprite(overweightFrame)
	case tama.Underweight:
		return RenderSprite(underweightFrame)
	case tama.Happiness > 80:
		return RenderSprite(happyFrames[idx])
	case tama.Angry > 50:
		return RenderSprite(angryFrames[idx])
	case tama.Hunger < 30 || tama.Thirst < 30:
		return RenderSprite(hungryFrames[idx])
	default:
		return getStageSprite(tama.Stage, idx)
	}
}

func getStageSprite(stage model.Stage, idx int) string {
	switch stage {
	case model.StageBaby:
		return RenderSprite(babyFrames[idx])
	case model.StageChild:
		return RenderSprite(childFrames[idx])
	case model.StageTeen:
		return RenderSprite(teenFrames[idx])
	case model.StageAdult:
		return RenderSprite(adultFrames[idx])
	case model.StageElder:
		return RenderSprite(elderFrames[idx])
	default:
		return RenderSprite(babyFrames[idx])
	}
}
