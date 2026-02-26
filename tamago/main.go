package main

import (
	"Pessoal/internal/model"
	"Pessoal/internal/persistence"
	"Pessoal/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var tama *model.Tama

	if persistence.Exists() {
		loaded, err := persistence.Load()
		if err == nil && !loaded.Dead {
			tama = loaded
			fmt.Println("💾 Save carregado! Bem-vindo de volta,", tama.Name+"!")
		}
	}

	if tama == nil {
		tama = &model.Tama{
			Name:         "Pochi",
			Hunger:       50,
			Thirst:       50,
			Sleepy:       50,
			Happiness:    50,
			Angry:        0,
			Weight:       50,
			Achievements: model.DefaultAchievements(),
		}
	}

	// Garantir que saves antigos tenham achievements
	if len(tama.Achievements) == 0 {
		tama.Achievements = model.DefaultAchievements()
	}

	p := tea.NewProgram(ui.InitialModel(tama))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao iniciar o programa: %v\n", err)
		os.Exit(1)
	}
}
