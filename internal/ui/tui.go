package ui

import (
	"Pessoal/internal/model"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Constantes de Layout ---
const (
	widthApp  = 80
	heightApp = 24
)

// --- Estilos e Paleta de Cores ---
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	danger    = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#FF5F87"}
	warning   = lipgloss.AdaptiveColor{Light: "#F5A623", Dark: "#F7B538"}

	// Estilos Base
	styleBase = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			BorderForeground(highlight).
			Padding(1, 2)

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
			Align(lipgloss.Center).
			Width(30).
			Height(13) // Altura fixa para alinhar com status (updated to match stats box)
		// Aumentei a altura pra barra de peso
	styleStatsBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666")).
			Padding(0, 2).
			Width(44).
			Height(13) // Aumentei a altura pra barra de peso

	styleInput = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(subtle).
			Padding(0, 1).
			MarginTop(1).
			Width(widthApp - 4)

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
	quit      chan bool
	message   string
	isError   bool
	frame     int
	width     int
	height    int
}

func InitialModel(tama *model.Tama, quit chan bool) Model {
	ti := textinput.New()
	ti.Placeholder = "Comandos: feed, water, pet, sleep, exercise, annoy..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return Model{
		tama:      tama,
		textInput: ti,
		quit:      quit,
		message:   "✨ Olá! Cuide bem do " + tama.Name + "!",
		frame:     0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.frame++
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit <- true
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(strings.ToLower(m.textInput.Value()))
			m.textInput.SetValue("")

			if input == "" {
				return m, nil
			}

			// Separa comando de argumentos se necessário (futuro)
			m.handleCommand(input)

			if input == "quit" || input == "exit" || input == "q" {
				m.quit <- true
				return m, tea.Quit
			}

			// Se morreu após o comando
			if m.tama.Dead {
				return m, tea.Quit
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// Lógica separada para comandos
func (m *Model) handleCommand(cmd string) {
	m.isError = false
	// Mapeia atalhos e comandos
	switch cmd {
	case "feed", "f":
		ToFeed(m.tama) // Assumindo que esta func existe no pacote externo
		m.message = "🍕 Nhac! " + m.tama.Name + " comeu tudo!"
	case "water", "w":
		GiveWater(m.tama)
		m.message = "💧 Glup! " + m.tama.Name + " bebeu água."
	case "pet", "p":
		PetTama(m.tama)
		m.message = "❤️  Purr... " + m.tama.Name + " gostou do carinho."
	case "sleep", "s":
		Sleep(m.tama)
		m.message = "😴 Shhh... " + m.tama.Name + " foi dormir."
	case "exercise", "e", "play":
		Exercise(m.tama)
		m.message = "🏃 " + m.tama.Name + " fez exercício e está mais saudável!"
	case "annoy", "a":
		Annoy(m.tama)
		m.message = "😠 Hey! Você irritou o " + m.tama.Name + "!"
	case "status":
		m.message = "📊 Atualizando status..."
	default:
		m.message = fmt.Sprintf("❓ Comando '%s' desconhecido.", cmd)
		m.isError = true
	}
}

// --- View ---

func (m Model) View() string {
	if m.tama.Dead {
		return m.renderGameOver()
	}

	// 1. Cabeçalho
	title := styleTitle.Render("🎮 TAMAGOTCHI CLI")

	// 2. Área Principal (Avatar + Status Lado a Lado)
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderAvatar(),
		"  ", // Espaçamento
		m.renderStats(),
	)

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
	inputView := styleInput.Render(m.textInput.View())

	// 5. Ajuda/Rodapé
	helpView := styleHelp.Render("ESC/Ctrl+C: Sair • Comandos: (f)eed, (w)ater, (p)et, (s)leep, (e)xercise")

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

func (m Model) renderGameOver() string {
	styleDead := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(danger).
		Foreground(danger).
		Align(lipgloss.Center).
		Padding(2).
		Width(50)

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

func (m Model) renderAvatar() string {
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
		frames := []string{"(・_・)", "(・o・)", "(・_・)"}
		art = frames[m.frame%len(frames)]
	}

	// Renderização
	artRendered := lipgloss.NewStyle().Foreground(color).Bold(true).Render

	// Layout dentro da caixa do avatar
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"\n",
		artRendered(art),
		"\n",
		lipgloss.NewStyle().Bold(true).Render(m.tama.Name),
		lipgloss.NewStyle().Foreground(subtle).Italic(true).Render(mood),
	)

	return styleTamaBox.Render(content)
}
func (m Model) renderStats() string {
	// Cabeçalho dos stats
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAA")).Render("STATUS VITAIS")

	// Renderiza as barras
	// Assumindo valores máx do model (ex: 100), ajustei para constantes locais se precisar
	// Ajuste as constantes `model.MaxHunger` conforme sua implementação real
	barHunger := m.progressBar("Fome", m.tama.Hunger, model.MaxHunger)
	barThirst := m.progressBar("Sede", m.tama.Thirst, model.MaxThirst)
	barSleep := m.progressBar("Energia", m.tama.Sleepy, model.MaxSleepy)
	barHappy := m.progressBar("Felicidade", m.tama.Happiness, model.MaxHappiness)
	barAnger := m.progressBar("Calma", model.MaxAngry-m.tama.Angry, model.MaxAngry) // Invertido para lógica visual (barra cheia = bom)
	barWeight := m.progressBar("Peso", m.tama.Weight, model.MaxWeight)

	// Indicador de peso
	weightStatus := ""
	if m.tama.Overweight {
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" ⚠️ Sobrepeso")
	} else if m.tama.Underweight {
		weightStatus = lipgloss.NewStyle().Foreground(warning).Render(" ⚠️ Abaixo do peso")
	} else {
		weightStatus = lipgloss.NewStyle().Foreground(special).Render(" ✓ Peso ideal")
	}

	return styleStatsBox.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			"\n",
			barHunger,
			barThirst,
			barSleep,
			barHappy,
			barAnger,
			barWeight+weightStatus,
		),
	)
}

// Helper para criar uma barra de progresso colorida
func (m Model) progressBar(label string, value, max int) string {
	width := 20
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
