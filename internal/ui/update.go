package ui

import (
	"github.com/DeMart1n/Tamagotchi/internal/model"
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Aliases para uso no tui.go (mesmo pacote)
const EventInteractive = model.EventInteractive
const EventNegative = model.EventNegative

// --- Mensagens do Lifecycle ---

type lifeCycleTickMsg time.Time
type wakeUpMsg struct{}
type autoSaveTickMsg time.Time

func lifeCycleTickCmd() tea.Cmd {
	return tea.Tick(20*time.Second, func(t time.Time) tea.Msg {
		return lifeCycleTickMsg(t)
	})
}

func sleepCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return wakeUpMsg{}
	})
}

func autoSaveTickCmd() tea.Cmd {
	return tea.Tick(60*time.Second, func(t time.Time) tea.Msg {
		return autoSaveTickMsg(t)
	})
}

// --- Funções de Comando com Humor + XP ---

func ToFeed(tama *model.Tama) string {
	if tama.Hunger == model.MaxHunger {
		return "🍕 " + tama.Name + " está cheio! Não quer comer."
	}

	amount := 15
	msg := "🍕 Nhac! " + tama.Name + " comeu tudo!"

	switch {
	case tama.Happiness > 70:
		amount = 20
		msg = "🍕 " + tama.Name + " comeu com alegria! Bônus de felicidade!"
	case tama.PissedOf || tama.Angry > 70:
		amount = 8
		msg = "🍕 " + tama.Name + " comeu de má vontade..."
	case tama.Depressed:
		amount = 10
		msg = "🍕 " + tama.Name + " comeu sem vontade..."
	}

	tama.Hunger = min(tama.Hunger+amount, model.MaxHunger)
	tama.Weight = min(tama.Weight+2, model.MaxWeight)
	updateWeightStates(tama)
	tama.TotalFeeds++
	if tama.AddXP(5) {
		msg += fmt.Sprintf(" ⬆️ LEVEL UP! Nível %d!", tama.Level)
	}
	return msg
}

func GiveWater(tama *model.Tama) string {
	if tama.Thirst == model.MaxThirst {
		return "💧 " + tama.Name + " não está com sede!"
	}

	amount := 15
	msg := "💧 Glup! " + tama.Name + " bebeu água."

	if tama.Happiness > 70 {
		amount = 20
		msg = "💧 " + tama.Name + " bebeu água com gosto! Refrescante!"
	} else if tama.PissedOf {
		amount = 10
		msg = "💧 " + tama.Name + " tomou água emburrado..."
	}

	tama.Thirst = min(tama.Thirst+amount, model.MaxThirst)
	tama.TotalWaters++
	if tama.AddXP(5) {
		msg += fmt.Sprintf(" ⬆️ LEVEL UP! Nível %d!", tama.Level)
	}
	return msg
}

func PetTama(tama *model.Tama) string {
	if tama.Happiness == model.MaxHappiness {
		return "❤️  " + tama.Name + " já está muito feliz!"
	}

	amount := 15
	msg := "❤️  Purr... " + tama.Name + " gostou do carinho."

	switch {
	case tama.Depressed:
		amount = 25
		msg = "❤️  " + tama.Name + " realmente precisava disso... Super efetivo!"
	case tama.PissedOf:
		amount = 5
		tama.Angry = max(tama.Angry-10, model.MinAngry)
		msg = "❤️  " + tama.Name + " resistiu, mas se acalmou um pouco."
	case tama.Sleeping:
		return "💤 Shhh... " + tama.Name + " está dormindo."
	}

	tama.Happiness = min(tama.Happiness+amount, model.MaxHappiness)
	if tama.Happiness > 20 {
		tama.Depressed = false
	}
	tama.TotalPets++
	if tama.AddXP(5) {
		msg += fmt.Sprintf(" ⬆️ LEVEL UP! Nível %d!", tama.Level)
	}
	return msg
}

func Annoy(tama *model.Tama) string {
	if tama.Sleeping {
		tama.Sleeping = false
		tama.Angry = min(tama.Angry+30, model.MaxAngry)
		tama.Happiness = max(tama.Happiness-10, model.MinHappiness)
		tama.TotalAnnoys++
		return "💢 Você acordou " + tama.Name + "! Ele está FURIOSO!"
	}

	if tama.Angry == model.MaxAngry {
		return "😠 " + tama.Name + " já está furioso!"
	}

	amount := 15
	msg := "😠 Hey! Você irritou o " + tama.Name + "!"

	if tama.Happiness > 70 {
		amount = 8
		msg = "😠 " + tama.Name + " ficou um pouco irritado, mas está de bom humor."
	} else if tama.Depressed {
		amount = 20
		tama.Happiness = max(tama.Happiness-5, model.MinHappiness)
		msg = "😠 " + tama.Name + " ficou muito chateado... Coitado!"
	}

	tama.Angry = min(tama.Angry+amount, model.MaxAngry)
	tama.TotalAnnoys++
	return msg
}

func StartSleep(tama *model.Tama) (string, tea.Cmd) {
	if tama.Sleepy == model.MaxSleepy {
		return "😴 " + tama.Name + " não está com sono!", nil
	}
	if tama.Sleeping {
		return "😴 " + tama.Name + " já está dormindo...", nil
	}
	tama.Sleeping = true
	tama.TotalSleeps++
	msg := "😴 Shhh... " + tama.Name + " foi dormir."
	if tama.AddXP(5) {
		msg += fmt.Sprintf(" ⬆️ LEVEL UP! Nível %d!", tama.Level)
	}
	return msg, sleepCmd()
}

func WakeUp(tama *model.Tama) {
	tama.Sleeping = false
	tama.Sleepy = min(tama.Sleepy+50, model.MaxSleepy)
}

func Exercise(tama *model.Tama) string {
	if tama.Sleeping {
		return "💤 " + tama.Name + " está dormindo... Não dá para exercitar."
	}
	if tama.Hunger < 20 {
		return "🏃 " + tama.Name + " está com muita fome para se exercitar!"
	}
	if tama.Weight == model.MinWeight {
		return "🏃 " + tama.Name + " já está no peso mínimo!"
	}

	weightLoss := 5
	msg := "🏃 " + tama.Name + " fez exercício e está mais saudável!"

	if tama.Happiness > 70 {
		weightLoss = 7
		msg = "🏃 " + tama.Name + " fez exercício com empolgação! Super treino!"
	}

	tama.Weight = max(tama.Weight-weightLoss, model.MinWeight)
	tama.Hunger = max(tama.Hunger-3, model.MinHunger)
	tama.Thirst = max(tama.Thirst-3, model.MinThirst)
	updateWeightStates(tama)
	tama.TotalExercises++
	if tama.AddXP(8) {
		msg += fmt.Sprintf(" ⬆️ LEVEL UP! Nível %d!", tama.Level)
	}
	return msg
}

func updateWeightStates(tama *model.Tama) {
	tama.Overweight = tama.Weight >= model.Overweight
	tama.Underweight = tama.Weight <= model.Underweight
}

// --- Tick do Lifecycle (chamado a cada 20s via tea.Cmd) ---

func Tick(tama *model.Tama) {
	if tama.Dead || tama.Sleeping {
		return
	}

	decay := tama.Stage.DecayRate()

	// Decay base (multiplicado pelo estágio)
	if tama.Hunger > 0 {
		tama.Hunger = max(tama.Hunger-decay, 0)
	}
	if tama.Sleepy > 0 {
		tama.Sleepy = max(tama.Sleepy-decay, 0)
	}
	if tama.Thirst > 0 {
		tama.Thirst = max(tama.Thirst-decay, 0)
	}
	if tama.Happiness > 0 {
		tama.Happiness = max(tama.Happiness-decay, 0)
	}

	// Interações entre stats
	if tama.Angry > 50 && tama.Happiness > 0 {
		tama.Happiness--
	}
	if tama.Happiness > 70 && tama.Angry > 0 {
		tama.Angry = max(tama.Angry-2, model.MinAngry)
	}
	if tama.Angry > 0 {
		tama.Angry = max(tama.Angry-1, model.MinAngry)
	}
	if tama.Hunger < 20 && tama.Happiness > 0 {
		tama.Happiness--
	}
	if tama.Thirst < 20 && tama.Happiness > 0 {
		tama.Happiness--
	}

	// Morte por fome ou sede
	if tama.Hunger <= 0 || tama.Thirst <= 0 {
		tama.Dead = true
		return
	}

	// Desmaio por sono
	if tama.Sleepy <= 0 {
		tama.Sleeping = true
	}

	// Depressão
	if tama.Happiness <= 0 {
		tama.Depressed = true
	} else if tama.Happiness > 20 {
		tama.Depressed = false
	}

	// Fúria
	if tama.Angry >= 90 {
		tama.PissedOf = true
	} else if tama.Angry < 70 {
		tama.PissedOf = false
	}

	// XP por sobreviver tick
	tama.TotalTicks++
	xpGain := 2
	// Bônus se todos os stats estão bons (>50)
	if tama.Hunger > 50 && tama.Thirst > 50 && tama.Sleepy > 50 && tama.Happiness > 50 && tama.Angry < 30 {
		xpGain += 5
	}
	tama.AddXP(xpGain)
}

// --- Eventos Aleatórios ---

func tryRandomEvent() *model.Event {
	if rand.Intn(100) < 30 { // 30% chance
		evt := model.Events[rand.Intn(len(model.Events))]
		return &evt
	}
	return nil
}

func applyEventDirect(tama *model.Tama, evt *model.Event) {
	tama.Hunger = max(min(tama.Hunger+evt.HungerDelta, model.MaxHunger), model.MinHunger)
	tama.Thirst = max(min(tama.Thirst+evt.ThirstDelta, model.MaxThirst), model.MinThirst)
	tama.Happiness = max(min(tama.Happiness+evt.HappinessDelta, model.MaxHappiness), model.MinHappiness)
	tama.Angry = max(min(tama.Angry+evt.AngryDelta, model.MaxAngry), model.MinAngry)
	tama.Sleepy = max(min(tama.Sleepy+evt.SleepyDelta, model.MaxSleepy), model.MinSleepy)
	if evt.XPDelta > 0 {
		tama.AddXP(evt.XPDelta)
	}
}

func (m *Model) applyEventResult(result model.EventResult) {
	t := m.tama
	t.Hunger = max(min(t.Hunger+result.HungerDelta, model.MaxHunger), model.MinHunger)
	t.Thirst = max(min(t.Thirst+result.ThirstDelta, model.MaxThirst), model.MinThirst)
	t.Happiness = max(min(t.Happiness+result.HappinessDelta, model.MaxHappiness), model.MinHappiness)
	t.Angry = max(min(t.Angry+result.AngryDelta, model.MaxAngry), model.MinAngry)
	if result.XPDelta > 0 {
		t.AddXP(result.XPDelta)
	}
}
