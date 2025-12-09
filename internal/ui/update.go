package ui

import (
	"Pessoal/internal/model"
	"fmt"
	"time"
)

func toFeed(tama *model.Tama) {
	if tama.Hunger == model.MaxHunger {
		fmt.Println("I'm Full!")
		return
	}

	tama.Hunger = min(tama.Hunger+15, model.MaxHunger)
}

func giveWater(tama *model.Tama) {
	if tama.Thirst == model.MaxThirst {
		fmt.Println("I've had enough!")
		return
	}

	tama.Thirst = min(tama.Thirst+15, model.MaxThirst)
}

func petTama(tama *model.Tama) {
	if tama.Happiness == model.MaxHappiness {
		fmt.Println("I'm all cuddled out!")
		return
	}

	tama.Happiness = min(tama.Happiness+15, model.MaxHappiness)
}

func Annoy(tama *model.Tama) {
	if tama.Angry == model.MaxAngry {
		fmt.Println("I'm all peeved!")
		return
	}

	tama.Angry = min(tama.Angry+15, model.MaxAngry)
}

func Sleep(tama *model.Tama) {
	if tama.Sleepy == model.MaxSleepy {
		fmt.Println("I'm not sleepy!")
		return
	}

	go func() {
		time.Sleep(10 * time.Second)
		tama.Sleeping = false
		tama.Sleepy = min(tama.Sleepy+50, model.MaxSleepy)
	}()
}

func Tick(tama *model.Tama) {
	if tama.Hunger > 0 {
		tama.Hunger--
	}
	if tama.Sleepy > 0 {
		tama.Sleepy--
	}
	if tama.Thirst > 0 {
		tama.Thirst--
	}
	if tama.Happiness > 0 {
		tama.Happiness--
	}

	//Mortes por fome e sede
	if tama.Hunger <= 0 || tama.Thirst <= 0 {
		tama.Dead = true
		fmt.Printf("%s is dead!\n", tama.Name)
	}

	if tama.Sleepy <= 0 {
		fmt.Printf("%s fainted...\n", tama.Name)
		go func() {
			time.Sleep(1 * time.Minute)
		}()
	}

	if tama.Happiness <= 0 {
		tama.Depressed = true
		fmt.Printf("%s now has depression... :C\n", tama.Name)
	}

	if tama.Angry >= 90 {
		tama.PissedOf = true
		fmt.Printf("%s is pissed off\n", tama.Name)
	}

}

func StartLifyCycle(tama *model.Tama) chan bool {
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				Tick(tama)
			case <-quit:
				ticker.Stop() // Limpa o ticker da memória
				return        // Mata a goroutine
			}
		}
	}()
	return quit // Retorna o controle para quem chamou
}
