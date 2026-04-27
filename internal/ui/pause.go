package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// pauseOptions are the entries shown in the pause menu.
var pauseOptions = []string{
	"Retomar Jogo",
	"Ajuda / Comandos",
	"Sair (salva automaticamente)",
}

var (
	stylePauseBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 4).
			Align(lipgloss.Center)

	stylePauseTitle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)

	stylePauseSelected = lipgloss.NewStyle().
				Foreground(special).
				Bold(true)

	stylePauseUnselected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

	styleHelpSection = lipgloss.NewStyle().
				Foreground(highlight).
				Bold(true)

	styleHelpCmd = lipgloss.NewStyle().
			Foreground(special)

	styleHelpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))
)

// renderPauseScreen renders the pause menu overlay.
func (m Model) renderPauseScreen(l layout) string {
	title := stylePauseTitle.Render("II  PAUSADO")

	var opts []string
	for i, opt := range pauseOptions {
		prefix := fmt.Sprintf("[%d] ", i+1)
		if i == m.pauseCursor {
			opts = append(opts, stylePauseSelected.Render(prefix+opt))
		} else {
			opts = append(opts, stylePauseUnselected.Render(prefix+opt))
		}
	}

	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).Italic(true).
		Render("Use ↑↓ ou números • ENTER para confirmar • ESC para retomar")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		strings.Join(opts, "\n"),
		"",
		hint,
	)

	box := stylePauseBox.Width(clampInt(l.contentWidth*60/100, 40, 70)).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

// renderHelpScreen renders the full help reference screen.
func (m Model) renderHelpScreen(l layout) string {
	title := stylePauseTitle.Render("?  AJUDA — COMANDOS DISPONÍVEIS")

	type entry struct{ cmd, desc string }

	sections := []struct {
		header  string
		entries []entry
	}{
		{
			header: "Cuidados Básicos",
			entries: []entry{
				{"feed  / f", "Alimentar o Tamagotchi (+Fome, +Peso)"},
				{"water / w", "Dar água (+Sede)"},
				{"pet   / p", "Dar carinho (+Felicidade)"},
				{"sleep / s", "Colocar para dormir (+Energia)"},
				{"exercise / e", "Fazer exercício (-Peso, +XP)"},
				{"annoy / a", "Irritar o Tamagotchi (+Raiva)"},
			},
		},
		{
			header: "Informações",
			entries: []entry{
				{"status", "Ver nível, XP e estágio atual"},
				{"achievements / ach", "Ver conquistas desbloqueadas"},
			},
		},
		{
			header: "Mini-Games",
			entries: []entry{
				{"play guess", "Jogo de adivinhar o número (1-100, 7 tentativas)"},
				{"play react", "Jogo de tempo de reação (aperte Enter ao sinal)"},
			},
		},
		{
			header: "Dungeon RPG  (requer Nível 3+)",
			entries: []entry{
				{"dungeon / d", "Entrar na Masmorra (RPG por turnos)"},
			},
		},
		{
			header: "Sistema",
			entries: []entry{
				{"help / h", "Abrir esta tela de ajuda"},
				{"quit / q / exit", "Salvar e sair do jogo"},
				{"setlvl <n>", "Definir nível diretamente (modo depuração)"},
			},
		},
	}

	var lines []string
	lines = append(lines, title, "")

	for _, sec := range sections {
		lines = append(lines, styleHelpSection.Render("  "+sec.header))
		for _, e := range sec.entries {
			cmd := styleHelpCmd.Render(fmt.Sprintf("    %-22s", e.cmd))
			desc := styleHelpDesc.Render(e.desc)
			lines = append(lines, cmd+desc)
		}
		lines = append(lines, "")
	}

	lines = append(lines,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true).
			Render("  ESC ou qualquer tecla para voltar ao menu de pausa"),
	)

	content := strings.Join(lines, "\n")
	box := stylePauseBox.Width(clampInt(l.contentWidth, 60, 90)).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
