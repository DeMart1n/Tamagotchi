package main

import (
	"Pessoal/internal/model"
	"Pessoal/internal/ui"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Criar o Tamagotchi
	tama := &model.Tama{
		Name:      "Pochi",
		Hunger:    50,
		Thirst:    50,
		Sleepy:    50,
		Happiness: 50,
		Angry:     0,
		Sleeping:  false,
		Dead:      false,
	}

	// Iniciar o ciclo de vida
	quit := ui.StartLifyCycle(tama)

	// Scanner para ler comandos
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Welcome to Tamagotchi ===")
	fmt.Printf("Your Tamagotchi's name is: %s\n", tama.Name)
	fmt.Println("Type 'help' to see available commands")
	fmt.Println()

	// Loop principal da CLI
	for {
		if tama.Dead {
			fmt.Println("Game Over! Your Tamagotchi has died :(")
			quit <- true
			break
		}

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch command {
		case "feed":
			ui.ToFeed(tama)
		case "water":
			ui.GiveWater(tama)
		case "pet":
			ui.PetTama(tama)
		case "sleep":
			ui.Sleep(tama)
		case "annoy":
			ui.Annoy(tama)
		case "status":
			showStatus(tama)
		case "help":
			showHelp()
		case "quit", "exit":
			fmt.Println("Goodbye!")
			quit <- true
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func showStatus(tama *model.Tama) {
	fmt.Println("\n===== Status =====")
	fmt.Printf("Name:      %s\n", tama.Name)
	fmt.Printf("Hunger:    %d/%d\n", tama.Hunger, model.MaxHunger)
	fmt.Printf("Thirst:    %d/%d\n", tama.Thirst, model.MaxThirst)
	fmt.Printf("Sleepy:    %d/%d\n", tama.Sleepy, model.MaxSleepy)
	fmt.Printf("Happiness: %d/%d\n", tama.Happiness, model.MaxHappiness)
	fmt.Printf("Angry:     %d/%d\n", tama.Angry, model.MaxAngry)
	fmt.Printf("Sleeping:  %t\n", tama.Sleeping)
	fmt.Printf("Dead:      %t\n", tama.Dead)
	if tama.Depressed {
		fmt.Println("Status: DEPRESSED")
	}
	if tama.PissedOf {
		fmt.Println("Status: PISSED OFF")
	}
	fmt.Println("==================\n")
}

func showHelp() {
	fmt.Println("\n===== Commands =====")
	fmt.Println("feed   - Feed your Tamagotchi")
	fmt.Println("water  - Give water to your Tamagotchi")
	fmt.Println("pet    - Pet your Tamagotchi")
	fmt.Println("sleep  - Put your Tamagotchi to sleep")
	fmt.Println("annoy  - Annoy your Tamagotchi")
	fmt.Println("status - Show Tamagotchi stats")
	fmt.Println("help   - Show this help message")
	fmt.Println("quit   - Exit the game")
	fmt.Println("====================\n")
}
