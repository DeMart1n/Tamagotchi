package main

import (
	"Pessoal/internal/model"
	"Pessoal/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Criar o Tamagotchi
	tama := &model.Tama{
		Name:       "Pochi",
		Hunger:     50,
		Thirst:     50,
		Sleepy:     50,
		Happiness:  50,
		Angry:      0,
		Weight:     50,
		Sleeping:   false,
		Dead:       false,
		Depressed:  false,
		PissedOf:   false,
		Overweight: false,
		Underweight: false,
	}

	// Iniciar o ciclo de vida
	quit := ui.StartLifyCycle(tama)

	// Iniciar Bubble Tea
	p := tea.NewProgram(ui.InitialModel(tama, quit))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao iniciar o programa: %v\n", err)
		os.Exit(1)
	}
}
