package ui

import (
	"Pessoal/internal/model"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderQuestPanel(tama *model.Tama) string {
	questStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(1, 2)

	headerStyle := lipgloss.NewStyle().
		Foreground(special).
		Bold(true)

	questTypeStyle := lipgloss.NewStyle().
		Foreground(warning).
		Bold(true)

	progressStyle := lipgloss.NewStyle().
		Foreground(highlight)

	completedStyle := lipgloss.NewStyle().
		Foreground(special).
		Bold(true)

	// Agrupamento por tipo
	var storyQuests, dailyQuests, repeatQuests []string

	for _, id := range model.GetQuestOrder() {
		q, ok := model.AllQuests[id]
		if !ok {
			continue
		}

		// Encontra o progresso da quest
		var qp *model.QuestProgress
		for i := range tama.QuestProgress {
			if tama.QuestProgress[i].QuestID == id {
				qp = &tama.QuestProgress[i]
				break
			}
		}
		if qp == nil {
			continue
		}

		// Render da quest
		questLine := renderQuestLine(q, qp, tama, progressStyle, completedStyle)

		switch q.Type {
		case model.QuestStory:
			storyQuests = append(storyQuests, questLine)
		case model.QuestDaily:
			dailyQuests = append(dailyQuests, questLine)
		case model.QuestRepeatable:
			repeatQuests = append(repeatQuests, questLine)
		}
	}

	// Monta conteúdo
	var content []string
	content = append(content, headerStyle.Render("[MISSÕES]"))

	if len(storyQuests) > 0 {
		content = append(content, "")
		content = append(content, questTypeStyle.Render("Story"))
		content = append(content, storyQuests...)
	}

	if len(dailyQuests) > 0 {
		content = append(content, "")
		content = append(content, questTypeStyle.Render("Diárias"))
		content = append(content, dailyQuests...)
	}

	if len(repeatQuests) > 0 {
		content = append(content, "")
		content = append(content, questTypeStyle.Render("Repeatáveis"))
		content = append(content, repeatQuests...)
	}

	content = append(content, "")
	content = append(content, "Comandos: 'claim <id>' para reivindicar | 'quests' para sair")

	return questStyle.Render(strings.Join(content, "\n"))
}

func renderQuestLine(q model.Quest, qp *model.QuestProgress, tama *model.Tama, progressStyle, completedStyle lipgloss.Style) string {
	prog := model.GetCurrentProgress(tama, *qp)
	pct := float64(prog) / float64(q.CondTarget)
	if pct > 1 {
		pct = 1
	}

	// Barra de progresso
	barWidth := 15
	filledSize := int(pct * float64(barWidth))
	filled := strings.Repeat("▇", filledSize)
	empty := strings.Repeat("░", barWidth-filledSize)
	bar := filled + empty

	// Status da quest
	var status string
	switch qp.Status {
	case "completed":
		status = completedStyle.Render("✓ CONCLUÍDA")
	case "claimed":
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("#999")).Render("✓ REIVINDICADA")
	default:
		status = fmt.Sprintf("%d/%d", prog, q.CondTarget)
	}

	// Linha completa
	line := fmt.Sprintf("%s %s [%s] %s", q.Name, bar, progressStyle.Render(status), q.Description)
	return line
}
