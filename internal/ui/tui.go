package ui

import (
	"Pessoal/internal/model"
	"Pessoal/internal/persistence"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Estilos e Paleta de Cores ---
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	danger    = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#FF5F87"}
	warning   = lipgloss.AdaptiveColor{Light: "#F5A623", Dark: "#F7B538"}

	styleTitle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(special).
			Padding(0, 2).
			MarginBottom(1)

	styleTamaBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Align(lipgloss.Center)

	styleStatsBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666")).
			Padding(0, 2)

	styleInput = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(subtle).
			Padding(0, 1).
			MarginTop(1)

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true).
			MarginTop(1)
)

// --- Mensagens e Comandos ---

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// --- Modelo ---

type Model struct {
	tama      *model.Tama
	textInput textinput.Model
	message   string
	isError   bool
	frame     int
	width     int
	height    int

	// Eventos
	activeEvent *model.Event
	eventTimer  int // ticks restantes para responder

	// Mini-games
	gameMode  GameMode
	guessGame *GuessGame
	reactGame *ReactGame
}

func InitialModel(tama *model.Tama) Model {
	ti := textinput.New()
	ti.Placeholder = "Comandos: feed, water, pet, sleep, exercise, annoy..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return Model{
		tama:      tama,
		textInput: ti,
		message:   "✨ Olá! Cuide bem do " + tama.Name + "!",
		frame:     0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd(), lifeCycleTickCmd(), autoSaveTickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		l := computeLayout(msg.Width, msg.Height)
		if !l.tooSmall {
			m.textInput.Width = l.inputWidth - 4
		}

	case tickMsg:
		m.frame++
		return m, tickCmd()

	case lifeCycleTickMsg:
		Tick(m.tama)
		if m.tama.Dead {
			m.message = "💀 " + m.tama.Name + " não sobreviveu..."
			return m, nil
		}
		if m.activeEvent != nil {
			m.eventTimer--
			if m.eventTimer <= 0 {
				m.applyEventResult(m.activeEvent.Choice2Result)
				m.message = "⏰ Tempo esgotado! " + m.activeEvent.Choice2Result.Message
				m.activeEvent = nil
			}
		}
		if newAch := model.CheckAchievements(m.tama); len(newAch) > 0 {
			m.message = "🏆 " + strings.Join(newAch, ", ")
		}
		if m.activeEvent == nil && m.gameMode == ModeNormal && m.tama.TotalTicks%3 == 0 {
			if evt := tryRandomEvent(); evt != nil {
				if evt.Type == EventInteractive {
					m.activeEvent = evt
					m.eventTimer = 3
					m.message = fmt.Sprintf("🎲 %s [1] %s [2] %s", evt.Description, evt.Choice1Label, evt.Choice2Label)
				} else {
					applyEventDirect(m.tama, evt)
					icon := "🌟"
					if evt.Type == EventNegative {
						icon = "⚡"
					}
					m.message = fmt.Sprintf("%s %s", icon, evt.Description)
				}
			}
		}
		return m, lifeCycleTickCmd()

	case wakeUpMsg:
		WakeUp(m.tama)
		m.message = "☀️ " + m.tama.Name + " acordou renovado!"
		return m, nil

	case autoSaveTickMsg:
		persistence.Save(m.tama)
		return m, autoSaveTickCmd()

	case reactGoMsg:
		if m.gameMode == ModeReact && m.reactGame != nil {
			m.reactGame.HandleGo()
			m.message = "🟢 GO! Aperte ENTER!"
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.gameMode != ModeNormal {
				m.gameMode = ModeNormal
				m.guessGame = nil
				m.reactGame = nil
				m.message = "🎮 Saiu do mini-game."
				return m, nil
			}
			persistence.Save(m.tama)
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
			m.textInput.SetValue("")

			// --- Mini-game: Adivinhação ---
			if m.gameMode == ModeGuess {
				if m.guessGame.Done {
					m.gameMode = ModeNormal
					m.guessGame = nil
					m.message = "🎮 Voltou ao modo normal."
					return m, nil
				}
				m.message = m.guessGame.TryGuess(input)
				if m.guessGame.Done {
					if m.guessGame.Won {
						m.tama.Happiness = min(m.tama.Happiness+20, model.MaxHappiness)
						m.tama.AddXP(10)
					} else {
						m.tama.Happiness = max(m.tama.Happiness-5, model.MinHappiness)
					}
				}
				return m, nil
			}

			// --- Mini-game: Reação ---
			if m.gameMode == ModeReact {
				if m.reactGame.Done {
					m.gameMode = ModeNormal
					m.reactGame = nil
					m.message = "🎮 Voltou ao modo normal."
					return m, nil
				}
				result := m.reactGame.HandlePress()
				m.message = result
				if m.reactGame.Done {
					if m.reactGame.Phase == ReactReady {
						elapsed := m.reactGame.StartTime
						ms := time.Since(elapsed).Milliseconds()
						switch {
						case ms < 500:
							m.tama.Happiness = min(m.tama.Happiness+25, model.MaxHappiness)
							m.tama.AddXP(15)
						case ms < 1000:
							m.tama.Happiness = min(m.tama.Happiness+10, model.MaxHappiness)
							m.tama.AddXP(8)
						default:
							m.tama.Happiness = min(m.tama.Happiness+3, model.MaxHappiness)
							m.tama.AddXP(3)
						}
					} else {
						m.tama.Happiness = max(m.tama.Happiness-5, model.MinHappiness)
					}
				}
				return m, nil
			}

			// --- Modo normal ---
			if input == "" {
				return m, nil
			}

			if input == "quit" || input == "exit" || input == "q" {
				persistence.Save(m.tama)
				return m, tea.Quit
			}

			if m.activeEvent != nil && (input == "1" || input == "2") {
				if input == "1" {
					m.applyEventResult(m.activeEvent.Choice1Result)
					m.message = "🎲 " + m.activeEvent.Choice1Result.Message
				} else {
					m.applyEventResult(m.activeEvent.Choice2Result)
					m.message = "🎲 " + m.activeEvent.Choice2Result.Message
				}
				m.activeEvent = nil
				return m, nil
			}

			extraCmd := m.handleCommand(input)

			if newAchievements := model.CheckAchievements(m.tama); len(newAchievements) > 0 {
				m.message += " 🏆 " + strings.Join(newAchievements, ", ")
			}

			if m.tama.Dead {
				return m, nil
			}

			if extraCmd != nil {
				return m, extraCmd
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) handleCommand(cmd string) tea.Cmd {
	m.isError = false

	switch cmd {
	case "feed", "f":
		m.message = ToFeed(m.tama)
	case "water", "w":
		m.message = GiveWater(m.tama)
	case "pet", "p":
		m.message = PetTama(m.tama)
	case "sleep", "s":
		msg, sleepCmd := StartSleep(m.tama)
		m.message = msg
		return sleepCmd
	case "exercise", "e":
		m.message = Exercise(m.tama)
	case "annoy", "a":
		m.message = Annoy(m.tama)
	case "status":
		m.message = fmt.Sprintf("📊 Lv.%d | XP: %d/%d | Ticks: %d | Estágio: %s",
			m.tama.Level, m.tama.XP, model.XPForNextLevel(m.tama.Level),
			m.tama.TotalTicks, m.tama.Stage.String())
	case "achievements", "ach":
		m.message = m.renderAchievementsList()
	case "play guess", "guess":
		m.gameMode = ModeGuess
		m.guessGame = NewGuessGame()
		m.message = "🔢 Adivinhe o número de 1 a 100! Você tem 7 tentativas."
	case "play react", "react":
		m.gameMode = ModeReact
		m.reactGame = NewReactGame()
		m.message = "⏳ Espere o GO! aparecer e aperte Enter o mais rápido possível!"
		return m.reactGame.StartCmd()
	case "play":
		m.message = "🎮 Mini-games: 'play guess' (adivinhação) ou 'play react' (reação)"
	default:
		m.message = fmt.Sprintf("❓ Comando '%s' desconhecido.", cmd)
		m.isError = true
	}
	return nil
}

func (m *Model) renderAchievementsList() string {
	unlocked := 0
	total := len(m.tama.Achievements)
	parts := []string{}
	for _, a := range m.tama.Achievements {
		if a.Unlocked {
			unlocked++
			parts = append(parts, a.Icon+" "+a.Name)
		}
	}
	header := fmt.Sprintf("🏆 Conquistas (%d/%d): ", unlocked, total)
	if len(parts) == 0 {
		return header + "Nenhuma ainda!"
	}
	return header + strings.Join(parts, " | ")
}

// --- View ---

func (m Model) View() string {
	l := computeLayout(m.width, m.height)

	if l.tooSmall {
		msg := lipgloss.NewStyle().Foreground(warning).Bold(true).
			Render(fmt.Sprintf("Terminal muito pequeno (%dx%d)\nMínimo: %dx%d", m.width, m.height, minWidth, minHeight))
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg)
	}

	if m.tama.Dead {
		return m.renderGameOver(l)
	}

	// 1. Cabeçalho
	title := styleTitle.Render("🎮 TAMAGOTCHI CLI")

	// 2. Área Principal
	var mainContent string
	if m.gameMode == ModeGuess && m.guessGame != nil {
		gameView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Align(lipgloss.Center).
			Padding(2, 4).
			Width(l.gameBoxWidth).
			Render(m.guessGame.RenderView())
		mainContent = gameView
	} else if m.gameMode == ModeReact && m.reactGame != nil {
		gameView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Align(lipgloss.Center).
			Padding(2, 4).
			Width(l.gameBoxWidth).
			Render(m.reactGame.RenderView())
		mainContent = gameView
	} else {
		mainContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.renderAvatar(l),
			"  ",
			m.renderStats(l),
		)
	}

	// 3. Feedback do sistema (Mensagem)
	var msgView string
	icon := "📢"
	msgColor := special

	if m.isError {
		icon = "❌"
		msgColor = danger
	}

	msgStyle := lipgloss.NewStyle().Foreground(msgColor).Bold(true)
	msgView = fmt.Sprintf("\n%s %s", icon, msgStyle.Render(m.message))

	// 4. Input
	inputView := styleInput.Width(l.inputWidth).Render(m.textInput.View())

	// 5. Ajuda/Rodapé
	var helpText string
	switch m.gameMode {
	case ModeGuess:
		helpText = "ESC: Sair do jogo • Digite um número de 1-100"
	case ModeReact:
		helpText = "ESC: Sair do jogo • Aperte Enter quando aparecer GO!"
	default:
		helpText = "ESC: Sair • (f)eed (w)ater (p)et (s)leep (e)xercise (a)nnoy • play • ach"
	}
	helpView := styleHelp.Render(helpText)

	// Montagem final do bloco da aplicação
	appView := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		mainContent,
		msgView,
		inputView,
		helpView,
	)

	// Centraliza tudo na tela do terminal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		appView,
	)
}

func (m Model) renderGameOver(l layout) string {
	styleDead := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(danger).
		Foreground(danger).
		Align(lipgloss.Center).
		Padding(2).
		Width(l.gameOverWidth)

	content := fmt.Sprintf(
		"💀 GAME OVER 💀\n\n%s partiu dessa para melhor...\n\n(Pressione Ctrl+C para sair)",
		m.tama.Name,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styleDead.Render(content),
	)
}

func (m Model) renderAvatar(l layout) string {
	var art string
	var mood string

	// Correção aqui: Definimos explicitamente como TerminalColor (interface)
	// Isso permite aceitar tanto AdaptiveColor quanto Color comum.
	var color lipgloss.TerminalColor = highlight

	// Lógica de Animação e Humor
	switch {
	case m.tama.Dead:
		art = "💀"
		mood = "Morto"
		color = danger
	case m.tama.Sleeping:
		mood = "Dormindo zZZ"
		if m.frame%2 == 0 {
			art = "💤 😴"
		} else {
			art = "   😪"
		}
	case m.tama.Depressed:
		art = "☁️ 😢"
		mood = "Deprimido"
		// Agora isso funciona porque color é uma interface
		color = lipgloss.Color("#555555")
	case m.tama.PissedOf:
		art = "💢 👿"
		mood = "Furioso!"
		color = danger
	case m.tama.Overweight:
		art = "🐷"
		mood = "Sobrepeso"
		color = warning
	case m.tama.Underweight:
		art = "🦴"
		mood = "Abaixo do peso"
		color = warning
	case m.tama.Angry > 50:
		art = "😠"
		mood = "Irritado"
		color = warning
	case m.tama.Happiness > 80:
		mood = "Radiante!"
		if m.frame%4 == 0 {
			art = "✨ 😄 ✨"
		} else {
			art = "   😆   "
		}
		color = special
	case m.tama.Hunger < 30 || m.tama.Thirst < 30:
		art = "🥺"
		mood = "Carente"
		color = warning
	default:
		mood = "Normal"
		art = m.tama.Stage.Avatar(m.frame)
	}

	// Renderização
	artRendered := lipgloss.NewStyle().Foreground(color).Bold(true).Render

	stageInfo := fmt.Sprintf("Lv.%d %s", m.tama.Level, m.tama.Stage.String())

	// Layout dentro da caixa do avatar
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"\n",
		artRendered(art),
		"\n",
		lipgloss.NewStyle().Bold(true).Render(m.tama.Name),
		lipgloss.NewStyle().Foreground(highlight).Render(stageInfo),
		lipgloss.NewStyle().Foreground(subtle).Italic(true).Render(mood),
	)

	return styleTamaBox.Width(l.avatarWidth).Height(l.avatarHeight).Render(content)
}

func (m Model) renderStats(l layout) string {
	// Cabeçalho dos stats
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAA")).Render("STATUS VITAIS")

	// Renderiza as barras
	// Assumindo valores máx do model (ex: 100), ajustei para constantes locais se precisar
	// Ajuste as constantes `model.MaxHunger` conforme sua implementação real
	barHunger := m.progressBar("Fome", m.tama.Hunger, model.MaxHunger, l.barWidth)
	barThirst := m.progressBar("Sede", m.tama.Thirst, model.MaxThirst, l.barWidth)
	barSleep := m.progressBar("Energia", m.tama.Sleepy, model.MaxSleepy, l.barWidth)
	barHappy := m.progressBar("Felicidade", m.tama.Happiness, model.MaxHappiness, l.barWidth)
	barAnger := m.progressBar("Calma", model.MaxAngry-m.tama.Angry, model.MaxAngry, l.barWidth)
	barWeight := m.progressBar("Peso", m.tama.Weight, model.MaxWeight, l.barWidth)

	// Indicador de peso
	weightStatus := ""
	if m.tama.Overweight {
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" ⚠️ Sobrepeso")
	} else if m.tama.Underweight {
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" ⚠️ Abaixo do peso")
	} else {
		weightStatus = lipgloss.NewStyle().Foreground(special).Render(" ✓ Peso ideal")
	}

	// Barra de XP
	xpNeeded := model.XPForNextLevel(m.tama.Level)
	barXP := m.progressBarXP("XP", m.tama.XP, xpNeeded, l.barWidth)

	return styleStatsBox.Width(l.statsWidth).Height(l.statsHeight).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			"\n",
			barHunger,
			barThirst,
			barSleep,
			barHappy,
			barAnger,
			barWeight+weightStatus,
			barXP,
		),
	)
}

// Helper para criar uma barra de progresso colorida
func (m Model) progressBar(label string, value, max, width int) string {
	pct := float64(value) / float64(max)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	filledSize := int(pct * float64(width))

	// Cor baseada na porcentagem
	var c lipgloss.Color
	if pct < 0.3 {
		c = lipgloss.Color("#FF5F87") // Danger
	} else if pct < 0.6 {
		c = lipgloss.Color("#F7B538") // Warning
	} else {
		c = lipgloss.Color("#43BF6D") // Good
	}

	// Caracteres da barra
	fullChar := "▇"
	emptyChar := "░"

	filled := strings.Repeat(fullChar, filledSize)
	empty := strings.Repeat(emptyChar, width-filledSize)

	bar := lipgloss.NewStyle().Foreground(c).Render(filled) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(empty)

	// Formatação da linha: Ícone Label [Barra] Valor
	return fmt.Sprintf("%-10s %s  %3d%%", label, bar, int(pct*100))
}

func (m Model) progressBarXP(label string, value, maxVal, width int) string {
	pct := float64(value) / float64(maxVal)
	if pct > 1 {
		pct = 1
	}
	filledSize := int(pct * float64(width))

	c := lipgloss.Color("#7D56F4")
	filled := strings.Repeat("▇", filledSize)
	empty := strings.Repeat("░", width-filledSize)

	bar := lipgloss.NewStyle().Foreground(c).Render(filled) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(empty)

	return fmt.Sprintf("%-10s %s %d/%d", label, bar, value, maxVal)
}
