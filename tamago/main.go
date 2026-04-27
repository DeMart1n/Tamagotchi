package main

import (
	"Pessoal/internal/model"
	"Pessoal/internal/persistence"
	"Pessoal/internal/ui"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// launchInNewTerminal tenta abrir o executável num terminal novo no macOS.
// Retorna true se conseguiu lançar (o processo pai deve sair).
func launchInNewTerminal() bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	// Detectar go run: o binário fica em diretório temporário de build
	exe, err := os.Executable()
	if err != nil || strings.Contains(exe, "/go-build") {
		return false
	}

	env := append(os.Environ(), "TAMAGO_CHILD=1")

	// Tentar iTerm2
	if _, err := os.Stat("/Applications/iTerm.app"); err == nil {
		script := fmt.Sprintf(`tell application "iTerm"
	activate
	set newWindow to (create window with default profile)
	tell current session of newWindow
		write text "env TAMAGO_CHILD=1 %s; exit"
	end tell
end tell`, shellQuote(exe))
		cmd := exec.Command("osascript", "-e", script)
		cmd.Env = env
		if cmd.Run() == nil {
			return true
		}
	}

	// Fallback: Terminal.app
	script := fmt.Sprintf(`tell application "Terminal"
	activate
	do script "env TAMAGO_CHILD=1 %s; exit"
end tell`, shellQuote(exe))
	cmd := exec.Command("osascript", "-e", script)
	cmd.Env = env
	if cmd.Run() == nil {
		return true
	}

	return false
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func main() {
	// Padrão re-exec: se não é processo filho, tenta abrir em terminal novo
	if os.Getenv("TAMAGO_CHILD") != "1" {
		if launchInNewTerminal() {
			return
		}
	}

	var tama *model.Tama
	var hasSave bool

	if persistence.Exists() {
		loaded, err := persistence.Load()
		if err == nil && !loaded.Dead {
			tama = loaded
			hasSave = true
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
			CurrentBiome: "Florestal",
			Achievements: model.DefaultAchievements(),
		}
	}

	// Garantir que saves antigos tenham achievements (inclusive novos)
	if len(tama.Achievements) == 0 {
		tama.Achievements = model.DefaultAchievements()
	} else {
		// Adicionar conquistas novas que não existem no save antigo
		defaults := model.DefaultAchievements()
		if len(tama.Achievements) < len(defaults) {
			existing := make(map[string]bool)
			for _, a := range tama.Achievements {
				existing[a.ID] = true
			}
			for _, a := range defaults {
				if !existing[a.ID] {
					tama.Achievements = append(tama.Achievements, a)
				}
			}
		}
	}

	p := tea.NewProgram(ui.InitialModel(tama, hasSave), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao iniciar o programa: %v\n", err)
		os.Exit(1)
	}
}
