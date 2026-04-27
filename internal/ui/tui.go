package ui

import (
	"Pessoal/internal/dungeon"
	"Pessoal/internal/model"
	"Pessoal/internal/persistence"
	"encoding/json"
	"fmt"
	"strconv"
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

	// Dungeon
	dungeonGame *dungeon.DungeonRun
	dungeonInv  *dungeon.Inventory

	// Title screen
	hasSave  bool // true when a valid save was loaded at startup
	titleSub int  // titleSubMain or titleSubName

	// Pause screen
	pauseCursor int // selected option index (0-2)
}

func InitialModel(tama *model.Tama, hasSave bool) Model {
	ti := textinput.New()
	ti.Placeholder = "Comandos: feed, water, pet, sleep, exercise, annoy..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	// Carregar inventário da dungeon do save
	inv := dungeon.NewInventory()
	if len(tama.Inventory) > 0 {
		_ = json.Unmarshal(tama.Inventory, inv)
	}

	return Model{
		tama:       tama,
		textInput:  ti,
		message:    "* Ola! Cuide bem do " + tama.Name + "!",
		frame:      0,
		dungeonInv: inv,
		gameMode:   ModeTitle,
		hasSave:    hasSave,
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
		// Don't run lifecycle logic while on non-gameplay screens
		if m.gameMode == ModeTitle || m.gameMode == ModeGameOver ||
			m.gameMode == ModePause || m.gameMode == ModeHelp {
			return m, lifeCycleTickCmd()
		}
		Tick(m.tama)
		if m.tama.Dead {
			m.gameMode = ModeGameOver
			return m, lifeCycleTickCmd()
		}
		if m.activeEvent != nil {
			m.eventTimer--
			if m.eventTimer <= 0 {
				m.applyEventResult(m.activeEvent.Choice2Result)
				m.message = "[!] Tempo esgotado! " + m.activeEvent.Choice2Result.Message
				m.activeEvent = nil
			}
		}
		if newAch := model.CheckAchievements(m.tama); len(newAch) > 0 {
			m.message = "[CONQUISTA] " + strings.Join(newAch, ", ")
		}
		if m.activeEvent == nil && m.gameMode == ModeNormal && m.tama.TotalTicks%3 == 0 {
			if evt := tryRandomEvent(); evt != nil {
				if evt.Type == EventInteractive {
					m.activeEvent = evt
					m.eventTimer = 3
					m.message = fmt.Sprintf("[EVENTO] %s [1] %s [2] %s", evt.Description, evt.Choice1Label, evt.Choice2Label)
				} else {
					applyEventDirect(m.tama, evt)
					icon := "[+]"
					if evt.Type == EventNegative {
						icon = "[-]"
					}
					m.message = fmt.Sprintf("%s %s", icon, evt.Description)
				}
			}
		}
		return m, lifeCycleTickCmd()

	case wakeUpMsg:
		WakeUp(m.tama)
		m.message = "[*] " + m.tama.Name + " acordou renovado!"
		return m, nil

	case autoSaveTickMsg:
		m.saveDungeonInventory()
		persistence.Save(m.tama)
		return m, autoSaveTickCmd()

	case reactGoMsg:
		if m.gameMode == ModeReact && m.reactGame != nil {
			m.reactGame.HandleGo()
			m.message = ">> GO! Aperte ENTER! <<"
		}
		return m, nil

	case tea.KeyMsg:
		// --- Title Screen ---
		if m.gameMode == ModeTitle {
			return m.handleTitleKey(msg)
		}

		// --- Game Over Screen ---
		if m.gameMode == ModeGameOver {
			return m.handleGameOverKey(msg)
		}

		// --- Pause Screen ---
		if m.gameMode == ModePause {
			return m.handlePauseKey(msg)
		}

		// --- Help Screen ---
		if m.gameMode == ModeHelp {
			// Any key returns to pause
			m.gameMode = ModePause
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.gameMode == ModeDungeon {
				m.saveDungeonInventory()
				m.gameMode = ModeNormal
				m.dungeonGame = nil
				m.message = "[>] Saiu da masmorra."
				return m, nil
			}
			if m.gameMode != ModeNormal {
				m.gameMode = ModeNormal
				m.guessGame = nil
				m.reactGame = nil
				m.message = "[>] Saiu do mini-game."
				return m, nil
			}
			// ESC in normal mode → pause screen (instead of quitting directly)
			m.gameMode = ModePause
			m.pauseCursor = 0
			return m, nil

		case tea.KeyEnter:
			input := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
			m.textInput.SetValue("")

			// --- Mini-game: Adivinhação ---
			if m.gameMode == ModeGuess {
				if m.guessGame.Done {
					m.gameMode = ModeNormal
					m.guessGame = nil
					m.message = "[>] Voltou ao modo normal."
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
					m.message = "[>] Voltou ao modo normal."
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

			// --- Dungeon ---
			if m.gameMode == ModeDungeon && m.dungeonGame != nil {
				done := m.dungeonGame.HandleInput(input)
				if done {
					m.saveDungeonInventory()
					if m.tama.Dead {
						m.gameMode = ModeGameOver
					} else {
						m.gameMode = ModeNormal
					}
					m.dungeonGame = nil
					m.message = "[>] Voltou da masmorra."
					if !m.tama.Dead {
						if newAch := model.CheckAchievements(m.tama); len(newAch) > 0 {
							m.message += " [CONQUISTA] " + strings.Join(newAch, ", ")
						}
					}
				}
				return m, nil
			}

			// --- Modo normal ---
			if input == "" {
				return m, nil
			}

			if input == "quit" || input == "exit" || input == "q" {
				m.saveDungeonInventory()
				persistence.Save(m.tama)
				return m, tea.Quit
			}

			if m.activeEvent != nil && (input == "1" || input == "2") {
				if input == "1" {
					m.applyEventResult(m.activeEvent.Choice1Result)
					m.message = "[EVENTO] " + m.activeEvent.Choice1Result.Message
				} else {
					m.applyEventResult(m.activeEvent.Choice2Result)
					m.message = "[EVENTO] " + m.activeEvent.Choice2Result.Message
				}
				m.activeEvent = nil
				return m, nil
			}

			extraCmd := m.handleCommand(input)

			if newAchievements := model.CheckAchievements(m.tama); len(newAchievements) > 0 {
				m.message += " [CONQUISTA] " + strings.Join(newAchievements, ", ")
			}

			if m.tama.Dead {
				m.gameMode = ModeGameOver
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
		m.message = fmt.Sprintf("[STATUS] Lv.%d | XP: %d/%d | Ticks: %d | Estagio: %s",
			m.tama.Level, m.tama.XP, model.XPForNextLevel(m.tama.Level),
			m.tama.TotalTicks, m.tama.Stage.String())
	case "achievements", "ach":
		m.message = m.renderAchievementsList()
	case "play guess", "guess":
		m.gameMode = ModeGuess
		m.guessGame = NewGuessGame()
		m.message = "[JOGO] Adivinhe o numero de 1 a 100! Voce tem 7 tentativas."
	case "play react", "react":
		m.gameMode = ModeReact
		m.reactGame = NewReactGame()
		m.message = "[JOGO] Espere o GO! aparecer e aperte Enter o mais rapido possivel!"
		return m.reactGame.StartCmd()
	case "play":
		m.message = "[>] Mini-games: 'play guess' (adivinhação) ou 'play react' (reação)"
	case "dungeon", "d":
		if m.tama.Level < 3 {
			m.message = "[DUNGEON] Voce precisa ser pelo menos Level 3 para entrar na masmorra!"
			return nil
		}
		m.gameMode = ModeDungeon
		m.dungeonGame = dungeon.NewDungeonRun(m.tama, m.dungeonInv)
		m.message = "[DUNGEON] Entrando na Masmorra..."
	case "help", "h":
		m.gameMode = ModeHelp
	default:
		// Comando secreto: setlvl <numero>
		if strings.HasPrefix(cmd, "setlvl ") {
			lvlStr := strings.TrimPrefix(cmd, "setlvl ")
			lvl, err := strconv.Atoi(lvlStr)
			if err != nil || lvl < 0 {
				m.message = "Uso: setlvl <numero>"
				m.isError = true
				return nil
			}
			m.tama.Level = lvl
			m.tama.XP = 0
			m.tama.Stage = model.StageForLevel(lvl)
			m.message = fmt.Sprintf("Level setado para %d (%s)", lvl, m.tama.Stage.String())
			return nil
		}
		m.message = fmt.Sprintf("[?] Comando '%s' desconhecido. Digite 'help' ou 'h' para ver os comandos.", cmd)
		m.isError = true
	}
	return nil
}

// handleTitleKey handles keyboard input on the title screen.
func (m Model) handleTitleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.titleSub {
	case titleSubName:
		// Name entry sub-screen
		switch msg.Type {
		case tea.KeyEsc:
			m.titleSub = titleSubMain
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Comandos: feed, water, pet, sleep, exercise, annoy..."
			return m, nil
		case tea.KeyEnter:
			name := strings.TrimSpace(m.textInput.Value())
			if name == "" {
				name = "Pochi"
			}
			m.tama.Reset(name)
			m.message = "* Bem-vindo! Cuide bem do " + name + "!"
			m.gameMode = ModeNormal
			m.hasSave = false
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Comandos: feed, water, pet, sleep, exercise, annoy..."
			persistence.Delete()
			persistence.Save(m.tama)
			var c tea.Cmd
			m.textInput, c = m.textInput.Update(msg)
			return m, c
		}
		// Let text input handle the rest
		var c tea.Cmd
		m.textInput, c = m.textInput.Update(msg)
		return m, c

	default:
		// Main title sub-screen
		key := strings.ToLower(msg.String())
		switch {
		case msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC:
			persistence.Save(m.tama)
			return m, tea.Quit
		case msg.Type == tea.KeyEnter:
			if m.hasSave {
				// Continue existing game
				m.gameMode = ModeNormal
				m.message = "* Bem-vindo de volta, " + m.tama.Name + "!"
			} else {
				// No save → go straight to name entry
				m.titleSub = titleSubName
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Nome do seu Tamagotchi..."
				m.textInput.Focus()
			}
			return m, nil
		case key == "n":
			// New game (always available)
			m.titleSub = titleSubName
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Nome do seu Tamagotchi..."
			m.textInput.Focus()
			return m, nil
		}
		return m, nil
	}
}

// handlePauseKey handles keyboard input on the pause menu.
func (m Model) handlePauseKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.gameMode = ModeNormal
		return m, nil
	case tea.KeyUp:
		m.pauseCursor = (m.pauseCursor - 1 + len(pauseOptions)) % len(pauseOptions)
		return m, nil
	case tea.KeyDown:
		m.pauseCursor = (m.pauseCursor + 1) % len(pauseOptions)
		return m, nil
	case tea.KeyEnter:
		return m.executePauseOption()
	}
	// Number keys
	key := msg.String()
	if len(key) == 1 && key[0] >= '1' && key[0] <= '3' {
		m.pauseCursor = int(key[0]-'1')
		return m.executePauseOption()
	}
	return m, nil
}

func (m Model) executePauseOption() (tea.Model, tea.Cmd) {
	switch m.pauseCursor {
	case 0: // Retomar
		m.gameMode = ModeNormal
	case 1: // Ajuda
		m.gameMode = ModeHelp
	case 2: // Sair
		m.saveDungeonInventory()
		persistence.Save(m.tama)
		return m, tea.Quit
	}
	return m, nil
}

// handleGameOverKey handles keyboard input on the game over screen.
func (m Model) handleGameOverKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := strings.ToLower(msg.String())
	switch {
	case msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC:
		persistence.Delete()
		return m, tea.Quit
	case key == "r" || msg.Type == tea.KeyEnter:
		name := m.tama.Name
		m.tama.Reset(name)
		m.message = "* " + name + " renasceu! Cuide bem dele."
		m.isError = false
		m.dungeonInv = dungeon.NewInventory()
		persistence.Delete()
		persistence.Save(m.tama)
		// Go back to title screen (not directly into game); hasSave=true because we just saved
		m.gameMode = ModeTitle
		m.titleSub = titleSubMain
		m.hasSave = true
		return m, nil
	}
	return m, nil
}

func (m *Model) saveDungeonInventory() {
	if m.dungeonInv != nil {
		data, err := json.Marshal(m.dungeonInv)
		if err == nil {
			m.tama.Inventory = data
		}
	}
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
	header := fmt.Sprintf("[CONQUISTAS] (%d/%d): ", unlocked, total)
	if len(parts) == 0 {
		return header + "Nenhuma ainda!"
	}
	return header + strings.Join(parts, " | ")
}

// --- View ---

func (m Model) View() string {
	l := computeLayout(m.width, m.height)

	// Full-screen override modes (handle before tooSmall check so they can show their own messages)
	switch m.gameMode {
	case ModeTitle:
		return m.renderTitle(l)
	case ModePause:
		return m.renderPauseScreen(l)
	case ModeHelp:
		return m.renderHelpScreen(l)
	case ModeGameOver:
		return m.renderGameOver(l)
	}

	if l.tooSmall {
		msg := lipgloss.NewStyle().Foreground(warning).Bold(true).
			Render(fmt.Sprintf("Terminal muito pequeno (%dx%d)\nMínimo: %dx%d", m.width, m.height, minWidth, minHeight))
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg)
	}

	// 1. Cabeçalho
	title := styleTitle.Render("[>] TAMAGOTCHI CLI")

	// 2. Área Principal
	var mainContent string
	if m.gameMode == ModeDungeon && m.dungeonGame != nil {
		mainContent = m.renderDungeon(l)
	} else if m.gameMode == ModeGuess && m.guessGame != nil {
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
	icon := ">>>"
	msgColor := special

	if m.isError {
		icon = "[!]"
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
	case ModeDungeon:
		helpText = "ESC: Sair da masmorra • Digite o número da opção"
	default:
		helpText = "ESC: Pausar • (f)eed (w)ater (p)et (s)leep (e)xercise (a)nnoy • play • dungeon • help • ach"
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

func (m Model) renderAvatar(l layout) string {
	var mood string

	switch {
	case m.tama.Dead:
		mood = "Morto"
	case m.tama.Sleeping:
		mood = "Dormindo zZZ"
	case m.tama.Depressed:
		mood = "Deprimido"
	case m.tama.PissedOf:
		mood = "Furioso!"
	case m.tama.Overweight:
		mood = "Sobrepeso"
	case m.tama.Underweight:
		mood = "Abaixo do peso"
	case m.tama.Angry > 50:
		mood = "Irritado"
	case m.tama.Happiness > 80:
		mood = "Radiante!"
	case m.tama.Hunger < 30 || m.tama.Thirst < 30:
		mood = "Carente"
	default:
		mood = "Normal"
	}

	// Sprite pixel art com half-block characters
	spriteArt := GetSpriteForState(m.tama, m.frame)

	stageInfo := fmt.Sprintf("Lv.%d %s", m.tama.Level, m.tama.Stage.String())

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"\n",
		spriteArt,
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
	biomeName := m.tama.CurrentBiome
	if biomeName == "" {
		biomeName = dungeon.BiomeForest.String()
	}
	biome := dungeon.ParseBiome(biomeName)
	biomeLine := lipgloss.NewStyle().Foreground(highlight).Render("Bioma: " + biomeName)
	biomeAttrLine := lipgloss.NewStyle().Foreground(subtle).Render("Atributos: " + m.biomeAttributesSummary(biome))

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
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" [!] Sobrepeso")
	} else if m.tama.Underweight {
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" [!] Abaixo do peso")
	} else {
		weightStatus = lipgloss.NewStyle().Foreground(special).Render(" [OK] Peso ideal")
	}

	// Barra de XP
	xpNeeded := model.XPForNextLevel(m.tama.Level)
	barXP := m.progressBarXP("XP", m.tama.XP, xpNeeded, l.barWidth)

	return styleStatsBox.Width(l.statsWidth).Height(l.statsHeight).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			biomeLine,
			biomeAttrLine,
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
