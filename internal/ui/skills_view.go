package ui

import (
	"Pessoal/internal/model"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleSkillActive  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	styleSkillPassive = lipgloss.NewStyle().Foreground(lipgloss.Color("#73F59F")).Italic(true)
	styleSkillCost    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00d2ff"))
	styleSkillCD      = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7B538"))
	styleSkillSep     = lipgloss.NewStyle().Foreground(lipgloss.Color("#444"))
)

func renderSkillsPanel(tama *model.Tama) string {
	cls, hasClass := model.AllClasses[tama.Class]
	var header string
	if hasClass {
		header = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).
			Render(fmt.Sprintf("HABILIDADES — %s", cls.Name))
	} else {
		header = lipgloss.NewStyle().Bold(true).Render("HABILIDADES")
	}

	sep := styleSkillSep.Render(strings.Repeat("─", 44))

	lines := []string{header, sep}

	activeHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAA")).Render("  Ativas")
	lines = append(lines, activeHeader)

	for _, id := range tama.SkillsKnown {
		sk, ok := model.AllSkills[id]
		if !ok || sk.IsPassive {
			continue
		}

		cdText := ""
		if sk.Cooldown > 0 {
			cdText = styleSkillCD.Render(fmt.Sprintf("  CD:%d", sk.Cooldown))
		}

		turnsText := ""
		if sk.BuffTurns > 0 {
			turnsText = styleSkillCD.Render(fmt.Sprintf("  %dt", sk.BuffTurns))
		}

		namePart := styleSkillActive.Render(fmt.Sprintf("  %-22s", sk.Name))
		costPart := styleSkillCost.Render(fmt.Sprintf("%3dMP", sk.Cost))
		descPart := lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).Render("  " + sk.Description)

		lines = append(lines, namePart+costPart+cdText+turnsText)
		lines = append(lines, descPart)
	}

	lines = append(lines, sep)

	passiveHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAA")).Render("  Passivas")
	lines = append(lines, passiveHeader)

	hasPassive := false
	for _, id := range tama.SkillsKnown {
		sk, ok := model.AllSkills[id]
		if !ok || !sk.IsPassive {
			continue
		}
		hasPassive = true
		tag := styleSkillPassive.Render("  [PASSIVA] ")
		name := styleSkillPassive.Render(sk.Name)
		desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).Render("  " + sk.Description)
		lines = append(lines, tag+name)
		lines = append(lines, desc)
	}

	if !hasPassive {
		lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#555")).Render("  Nenhuma passiva disponível."))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(content)
}
