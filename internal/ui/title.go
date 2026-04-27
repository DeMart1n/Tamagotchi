package ui

import (
	"Pessoal/internal/model"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// titleSubMain shows the main title with options; titleSubName shows the name-entry form.
const (
	titleSubMain = 0
	titleSubName = 1
)

// titleLogoLines is the ASCII art logo displayed on the title screen.
var titleLogoLines = []string{
	` ___  ___  ___  ___  ___  ___  ___  ___  ___  ___ `,
	`|_ _||   ||   ||   || _ ||   ||_ _||  _||   ||_ _|`,
	` | | | - || - || - ||(_||  - | | | | |  | - | | | `,
	`|___||___||___||___||___||___||___||___||___||___|`,
	`  T    A    M    A    G    O    T    C    H    I   `,
}

var (
	styleTitleLogo = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)

	styleTitleBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(special).
			Padding(1, 4).
			Align(lipgloss.Center)

	styleTitleOption = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#73F59F")).
				Bold(true)

	styleTitleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	styleTitleGameOver = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5F87")).
				Bold(true)
)

// renderTitle dispatches to the appropriate title sub-screen.
func (m Model) renderTitle(l layout) string {
	if l.tooSmall {
		msg := lipgloss.NewStyle().Foreground(warning).Bold(true).
			Render(fmt.Sprintf("Terminal muito pequeno (%dx%d)\nMínimo: %dx%d",
				m.width, m.height, minWidth, minHeight))
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg)
	}

	switch m.titleSub {
	case titleSubName:
		return m.renderTitleNameEntry(l)
	default:
		return m.renderTitleMain(l)
	}
}

// renderTitleMain shows the animated main title screen.
func (m Model) renderTitleMain(l layout) string {
	// --- Logo ---
	logo := styleTitleLogo.Render(strings.Join(titleLogoLines, "\n"))

	// --- Animated sprite (cycles between frame 0 and 1) ---
	idx := m.frame % 2
	sprite := getStageSprite(model.StageAdult, idx)

	// --- Options ---
	var optLines []string
	if m.hasSave {
		optLines = append(optLines,
			styleTitleOption.Render("[ ENTER ]  Continuar Jogo"),
			styleTitleOption.Render("[  N    ]  Novo Jogo"),
		)
	} else {
		optLines = append(optLines,
			styleTitleOption.Render("[ ENTER ]  Novo Jogo"),
		)
	}
	optLines = append(optLines, "", styleTitleDim.Render("[ ESC ]  Sair"))
	opts := strings.Join(optLines, "\n")

	// --- Version / subtitle ---
	sub := styleTitleDim.Render("Virtual Pet  ·  Dungeon RPG  ·  Mini-Games")

	// --- Assemble content box ---
	content := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		sprite,
		"",
		sub,
		"",
		opts,
	)

	box := styleTitleBox.Width(l.contentWidth).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// renderTitleNameEntry shows the name-entry form for a new game.
func (m Model) renderTitleNameEntry(l layout) string {
	header := styleTitleLogo.Render("NOVO JOGO")
	prompt := lipgloss.NewStyle().Foreground(highlight).Bold(true).Render("Nome do seu Tamagotchi:")
	tip := styleTitleDim.Render("[ ENTER ] Confirmar  •  [ ESC ] Voltar")

	inputView := styleInput.Width(l.inputWidth / 2).Render(m.textInput.View())

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"",
		prompt,
		inputView,
		"",
		tip,
	)

	box := styleTitleBox.Width(l.contentWidth).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// renderGameOver shows the full-screen game over screen.
func (m Model) renderGameOver(l layout) string {
	skull := ItemSprite("defeat")
	deadSprite := RenderSprite(deadFrame)

	title := styleTitleGameOver.Render("=== FIM DE JOGO ===")
	msg := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true).
		Render(m.tama.Name + " não sobreviveu...")

	statsText := fmt.Sprintf("Nível %d  •  %d ticks sobrevividos  •  Fase: %s",
		m.tama.Level, m.tama.TotalTicks, m.tama.Stage.String())
	stats := styleTitleDim.Render(statsText)

	opts := lipgloss.JoinVertical(lipgloss.Center,
		styleTitleOption.Render("[ R ] ou [ ENTER ]  Reiniciar"),
		styleTitleDim.Render("[ ESC ]  Sair"),
	)

	spriteRow := lipgloss.JoinHorizontal(lipgloss.Center, skull, "   ", deadSprite)

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		"",
		spriteRow,
		"",
		msg,
		stats,
		"",
		opts,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Padding(1, 4).
		Align(lipgloss.Center).
		Width(l.contentWidth).
		Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
