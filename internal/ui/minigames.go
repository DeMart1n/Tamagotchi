package ui

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Tipos de modo de jogo ---

type GameMode int

const (
	ModeNormal GameMode = iota
	ModeGuess
	ModeReact
	ModeDungeon
)

// --- Jogo de Adivinhação ---

type GuessGame struct {
	Target     int
	Attempts   int
	MaxAttempt int
	Hint       string
	Won        bool
	Done       bool
}

func NewGuessGame() *GuessGame {
	return &GuessGame{
		Target:     rand.Intn(100) + 1,
		Attempts:   0,
		MaxAttempt: 7,
	}
}

func (g *GuessGame) TryGuess(input string) string {
	num, err := strconv.Atoi(input)
	if err != nil {
		return "❓ Digite um número entre 1 e 100!"
	}
	if num < 1 || num > 100 {
		return "❓ O número deve ser entre 1 e 100!"
	}

	g.Attempts++

	if num == g.Target {
		g.Won = true
		g.Done = true
		return fmt.Sprintf("🎉 ACERTOU! O número era %d! (%d tentativas)", g.Target, g.Attempts)
	}

	if g.Attempts >= g.MaxAttempt {
		g.Done = true
		return fmt.Sprintf("💔 Acabaram as tentativas! O número era %d.", g.Target)
	}

	remaining := g.MaxAttempt - g.Attempts
	if num < g.Target {
		g.Hint = "⬆️ Maior!"
	} else {
		g.Hint = "⬇️ Menor!"
	}
	return fmt.Sprintf("%s (Tentativa %d/%d, restam %d)", g.Hint, g.Attempts, g.MaxAttempt, remaining)
}

func (g *GuessGame) RenderView() string {
	if g.Done {
		if g.Won {
			return "🎉 VOCÊ GANHOU! Digite qualquer coisa para voltar."
		}
		return fmt.Sprintf("💔 PERDEU! O número era %d. Digite qualquer coisa para voltar.", g.Target)
	}
	header := "🔢 ADIVINHE O NÚMERO (1-100)"
	info := fmt.Sprintf("Tentativa %d/%d", g.Attempts, g.MaxAttempt)
	hint := g.Hint
	if hint == "" {
		hint = "Chute um número!"
	}
	return fmt.Sprintf("%s\n%s\n%s", header, info, hint)
}

// --- Jogo de Tempo de Reação ---

type ReactGame struct {
	Phase     ReactPhase
	StartTime time.Time
	Delay     time.Duration
	Result    string
	Done      bool
}

type ReactPhase int

const (
	ReactWaiting ReactPhase = iota
	ReactReady
	ReactFinished
)

type reactGoMsg struct{}

func NewReactGame() *ReactGame {
	delay := time.Duration(2+rand.Intn(4)) * time.Second
	return &ReactGame{
		Phase: ReactWaiting,
		Delay: delay,
	}
}

func (r *ReactGame) StartCmd() tea.Cmd {
	return tea.Tick(r.Delay, func(t time.Time) tea.Msg {
		return reactGoMsg{}
	})
}

func (r *ReactGame) HandleGo() {
	r.Phase = ReactReady
	r.StartTime = time.Now()
}

func (r *ReactGame) HandlePress() string {
	if r.Phase == ReactWaiting {
		r.Done = true
		r.Result = "⚡ Apertou cedo demais! -5 Felicidade."
		return r.Result
	}

	elapsed := time.Since(r.StartTime)
	r.Done = true

	ms := elapsed.Milliseconds()
	switch {
	case ms < 500:
		r.Result = fmt.Sprintf("⚡ INCRÍVEL! %dms — Reflexo de gato! +25 Felicidade, +15 XP", ms)
	case ms < 1000:
		r.Result = fmt.Sprintf("⚡ BOM! %dms — Nada mal! +10 Felicidade, +8 XP", ms)
	default:
		r.Result = fmt.Sprintf("⚡ OK! %dms — Pode melhorar! +3 Felicidade, +3 XP", ms)
	}
	return r.Result
}

func (r *ReactGame) RenderView() string {
	if r.Done {
		return r.Result + "\nDigite qualquer coisa para voltar."
	}
	switch r.Phase {
	case ReactWaiting:
		return "⏳ TEMPO DE REAÇÃO\n\nEspere aparecer \"GO!\"...\nNÃO aperte Enter antes!"
	case ReactReady:
		return "🟢 GO! GO! GO!\n\nAperte ENTER agora!"
	default:
		return ""
	}
}
