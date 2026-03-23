package ui

import (
	"Pessoal/internal/dungeon"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- Estilos da Dungeon ---
var (
	styleDungeonBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(0, 1)

	styleDungeonTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD700")).
				Bold(true)

	styleEnemyName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	stylePlayerInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#43BF6D")).
			Bold(true)

	styleCombatLog = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7B538")).
			Italic(true)

	styleRarityComum    = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleRarityIncomum  = lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D"))
	styleRarityRaro     = lipgloss.NewStyle().Foreground(lipgloss.Color("#5B9FFF"))
	styleRarityLendario = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
)

func rarityStyle(r dungeon.Rarity) lipgloss.Style {
	switch r {
	case dungeon.RarityIncomum:
		return styleRarityIncomum
	case dungeon.RarityRaro:
		return styleRarityRaro
	case dungeon.RarityLendario:
		return styleRarityLendario
	default:
		return styleRarityComum
	}
}

// renderDungeon renderiza a visão completa da dungeon baseada na fase atual.
func (m Model) renderDungeon(l layout) string {
	d := m.dungeonGame
	if d == nil {
		return ""
	}

	var content string

	switch d.Phase {
	case dungeon.PhaseMenuPrincipal:
		content = m.renderDungeonMenu(l)
	case dungeon.PhaseInventario:
		content = m.renderDungeonInventory(l)
	case dungeon.PhaseCombate:
		content = m.renderDungeonCombat(l)
	case dungeon.PhaseCombatResult:
		content = m.renderDungeonCombat(l)
	case dungeon.PhaseDescanso:
		content = m.renderDungeonRest(l)
	case dungeon.PhaseDescansoLoja:
		content = m.renderDungeonShop(l)
	case dungeon.PhaseTesouro:
		content = m.renderDungeonTreasure(l)
	case dungeon.PhaseFimAndar:
		content = m.renderDungeonFloorEnd(l)
	case dungeon.PhaseVitoria:
		content = m.renderDungeonVictory(l)
	case dungeon.PhaseDerrota:
		content = m.renderDungeonDefeat(l)
	default:
		content = d.Message
	}

	return styleDungeonBox.Width(l.dungeonWidth).Render(content)
}

func (m Model) renderDungeonMenu(l layout) string {
	title := styleDungeonTitle.Render("MASMORRA DE TAMAGO")
	biome := m.dungeonGame.CurrentBiome

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s (Lv.%d %s)", m.tama.Name, m.tama.Level, m.tama.Stage.String()))
	lines = append(lines, fmt.Sprintf("  Bioma: %s", biome.String()))
	lines = append(lines, fmt.Sprintf("  Atributos: %s", m.biomeAttributesSummary(biome)))
	lines = append(lines, "")
	lines = append(lines, "  [1] Entrar na Masmorra")
	lines = append(lines, "  [2] Inventario")
	lines = append(lines, "  [3] Voltar")
	lines = append(lines, "")
	lines = append(lines, styleCombatLog.Render("  "+m.dungeonGame.Message))

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonInventory(l layout) string {
	d := m.dungeonGame
	inv := d.Inv

	var lines []string
	lines = append(lines, styleDungeonTitle.Render("INVENTARIO"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Ouro: %d", inv.Gold))
	lines = append(lines, "")

	// Slots equipados
	lines = append(lines, "  Equipado:")
	if inv.Arma != nil {
		lines = append(lines, fmt.Sprintf("    Arma: %s", rarityStyle(inv.Arma.Rarity).Render(inv.Arma.Name)))
	} else {
		lines = append(lines, "    Arma: (vazio)")
	}
	if inv.Armadura != nil {
		lines = append(lines, fmt.Sprintf("    Armadura: %s", rarityStyle(inv.Armadura.Rarity).Render(inv.Armadura.Name)))
	} else {
		lines = append(lines, "    Armadura: (vazio)")
	}
	if inv.Acessorio != nil {
		lines = append(lines, fmt.Sprintf("    Acessorio: %s", rarityStyle(inv.Acessorio.Rarity).Render(inv.Acessorio.Name)))
	} else {
		lines = append(lines, "    Acessorio: (vazio)")
	}

	// Mochila
	if len(inv.Backpack) > 0 {
		lines = append(lines, "")
		lines = append(lines, "  Mochila:")
		for _, eq := range inv.Backpack {
			lines = append(lines, fmt.Sprintf("    - %s (%s)", rarityStyle(eq.Rarity).Render(eq.Name), eq.Description()))
		}
	}

	lines = append(lines, "")
	lines = append(lines, "  Aperte Enter para voltar.")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonCombat(l layout) string {
	d := m.dungeonGame
	c := d.Combat
	if c == nil {
		return d.Message
	}

	var lines []string

	// Progresso do andar
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())
	lines = append(lines, "")

	// Inimigo
	lines = append(lines, styleEnemyName.Render("  "+c.Enemy.Name))
	// ASCII art
	art := dungeon.EnemyArt(c.Enemy.Name)
	for _, artLine := range strings.Split(art, "\n") {
		lines = append(lines, "  "+artLine)
	}

	// HP do inimigo
	enemyBar := dungeon.HPBar(c.Enemy.HPCurrent, c.Enemy.HPMax, 16)
	lines = append(lines, fmt.Sprintf("  HP %s %d/%d", enemyBar, c.Enemy.HPCurrent, c.Enemy.HPMax))
	lines = append(lines, "")

	// Log de combate
	if d.Message != "" {
		lines = append(lines, styleCombatLog.Render("  "+d.Message))
	}
	if d.SubMessage != "" {
		lines = append(lines, styleCombatLog.Render("  "+d.SubMessage))
	}
	lines = append(lines, "")

	// Jogador
	lines = append(lines, stylePlayerInfo.Render(fmt.Sprintf("  %s (Voce)", m.tama.Name)))
	playerBar := dungeon.HPBar(d.Stats.HPCurrent, d.Stats.HPMax, 16)
	lines = append(lines, fmt.Sprintf("  HP %s %d/%d", playerBar, d.Stats.HPCurrent, d.Stats.HPMax))
	lines = append(lines, fmt.Sprintf("  ATK:%d  DEF:%d  VEL:%d", d.Stats.Ataque, d.Stats.Defesa, d.Stats.Velocidade))
	lines = append(lines, "")

	// Ações / resultado
	if d.Phase == dungeon.PhaseCombatResult {
		lines = append(lines, "  Aperte Enter para continuar.")
	} else if d.ChoosingItem {
		lines = append(lines, "  Itens:")
		for i, item := range d.ItemBag.Items {
			lines = append(lines, fmt.Sprintf("    [%d] %s", i+1, item.Name))
		}
		lines = append(lines, "    [0] Voltar")
	} else {
		lines = append(lines, "  [1] Atacar   [2] Defender")
		items := ""
		if d.ItemBag.Count() > 0 {
			items = fmt.Sprintf("(%d)", d.ItemBag.Count())
		}
		lines = append(lines, fmt.Sprintf("  [3] Item %s  [4] Fugir", items))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonRest(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())
	lines = append(lines, "")
	lines = append(lines, styleDungeonTitle.Render("  Sala de Descanso"))
	lines = append(lines, "")

	playerBar := dungeon.HPBar(d.Stats.HPCurrent, d.Stats.HPMax, 16)
	lines = append(lines, fmt.Sprintf("  HP %s %d/%d", playerBar, d.Stats.HPCurrent, d.Stats.HPMax))
	lines = append(lines, "")
	lines = append(lines, styleCombatLog.Render("  "+d.Message))
	lines = append(lines, "")
	lines = append(lines, "  [1] Loja de Pocoes")
	lines = append(lines, "  [2] Continuar")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonShop(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, styleDungeonTitle.Render("  Loja de Pocoes"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Ouro: %d", d.Inv.Gold))
	lines = append(lines, "")

	for i, item := range dungeon.ShopItems {
		desc := ""
		switch item.Type {
		case dungeon.ItemPotion:
			desc = fmt.Sprintf("Cura %d HP", item.Value)
		case dungeon.ItemATKBoost:
			desc = fmt.Sprintf("+%d ATK (1 turno)", item.Value)
		}
		lines = append(lines, fmt.Sprintf("  [%d] %s - %d ouro (%s)", i+1, item.Name, item.Price, desc))
	}

	lines = append(lines, "")
	lines = append(lines, "  [0] Voltar")
	lines = append(lines, "")
	if d.Message != "" {
		lines = append(lines, styleCombatLog.Render("  "+d.Message))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonTreasure(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())
	lines = append(lines, "")
	lines = append(lines, styleDungeonTitle.Render("  Bau do Tesouro!"))
	lines = append(lines, "")
	lines = append(lines, styleCombatLog.Render("  "+d.Message))

	if d.PendingLoot != nil {
		lines = append(lines, "")
		eq := d.PendingLoot
		lines = append(lines, fmt.Sprintf("  %s [%s]", rarityStyle(eq.Rarity).Render(eq.Name), eq.Rarity.String()))
		lines = append(lines, fmt.Sprintf("  %s - %s", eq.Slot.String(), eq.Description()))
		lines = append(lines, "")
		lines = append(lines, "  [1] Equipar")
		lines = append(lines, "  [2] Guardar na mochila")
		lines = append(lines, "  [3] Descartar")
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonFloorEnd(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, "")
	lines = append(lines, styleDungeonTitle.Render(fmt.Sprintf("  Andar %d Completo!", d.FloorNum)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  XP acumulado: %d", d.TotalXP))
	lines = append(lines, fmt.Sprintf("  Ouro acumulado: %d", d.TotalGold))
	lines = append(lines, "")
	playerBar := dungeon.HPBar(d.Stats.HPCurrent, d.Stats.HPMax, 16)
	lines = append(lines, fmt.Sprintf("  HP %s %d/%d", playerBar, d.Stats.HPCurrent, d.Stats.HPMax))
	lines = append(lines, "")
	lines = append(lines, "  Aperte Enter para o proximo andar.")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonVictory(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, "")
	lines = append(lines, styleDungeonTitle.Render("  VITORIA!"))
	lines = append(lines, "")
	lines = append(lines, "  O Dragao Anciao foi derrotado!")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  XP total: +%d", d.TotalXP))
	lines = append(lines, fmt.Sprintf("  Ouro total: +%d", d.TotalGold))
	lines = append(lines, "")
	lines = append(lines, "  Aperte Enter para voltar.")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonDefeat(l layout) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Bold(true).Render("  DERROTA..."))
	lines = append(lines, "")
	lines = append(lines, styleCombatLog.Render("  "+d.Message))
	lines = append(lines, "")
	lines = append(lines, "  Aperte Enter para voltar.")

	return strings.Join(lines, "\n")
}

// --- Helpers de progresso ---

func (m Model) renderFloorProgress() string {
	d := m.dungeonGame
	if d.FloorNum == 0 {
		return ""
	}

	var parts []string
	for i := 1; i <= 5; i++ {
		if i < d.FloorNum {
			parts = append(parts, "[ok]")
		} else if i == d.FloorNum {
			parts = append(parts, "[>>]")
		} else {
			parts = append(parts, "[  ]")
		}
		if i < 5 {
			parts = append(parts, "---")
		}
	}
	return fmt.Sprintf("  Bioma: %s\n  Atributos: %s\n  %s", d.CurrentBiome.String(), m.biomeAttributesSummary(d.CurrentBiome), strings.Join(parts, ""))
}

func (m Model) biomeAttributesSummary(biome dungeon.Biome) string {
	modifier := dungeon.ModifierForBiome(biome)

	var effects []string
	if modifier.HungerRecoveryBonus > 0 {
		effects = append(effects, fmt.Sprintf("+%d recuperacao fome", modifier.HungerRecoveryBonus))
	}
	if modifier.FriendlyEncounterChance > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% encontros amistosos", modifier.FriendlyEncounterChance))
	}
	if modifier.PlayerFireVulnerability > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% vulnerabilidade fogo", modifier.PlayerFireVulnerability))
	}
	if modifier.PlayerSpeedPenaltyPct > 0 {
		effects = append(effects, fmt.Sprintf("-%d%% velocidade", modifier.PlayerSpeedPenaltyPct))
	}
	if modifier.PlayerIceDefensePct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% defesa gelo", modifier.PlayerIceDefensePct))
	}
	if modifier.FreezeChancePct > 0 {
		effects = append(effects, fmt.Sprintf("%d%% chance congelamento", modifier.FreezeChancePct))
	}
	if modifier.EnemyAttackBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% ATK inimigo", modifier.EnemyAttackBonusPct))
	}
	if modifier.PlayerFireDotPctMaxHP > 0 {
		effects = append(effects, fmt.Sprintf("DOT fogo %d%% HP", modifier.PlayerFireDotPctMaxHP))
	}
	if modifier.HappinessPenaltyTick > 0 {
		effects = append(effects, fmt.Sprintf("-%d felicidade/tick", modifier.HappinessPenaltyTick))
	}
	if modifier.EnemyLuckBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% sorte inimiga", modifier.EnemyLuckBonusPct))
	}
	if modifier.RareLootChanceBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% loot raro", modifier.RareLootChanceBonusPct))
	}
	if modifier.PlayerHPRegenPenaltyPct > 0 {
		effects = append(effects, fmt.Sprintf("-%d%% regen HP", modifier.PlayerHPRegenPenaltyPct))
	}

	if len(effects) == 0 {
		return "Sem efeitos especiais"
	}

	return strings.Join(effects, " | ")
}

func (m Model) renderRoomProgress() string {
	d := m.dungeonGame
	if d.Floor == nil {
		return ""
	}

	var parts []string
	for i, room := range d.Floor.Rooms {
		if room.Cleared {
			parts = append(parts, "[ok]")
		} else if i == d.Floor.CurrentRoom {
			parts = append(parts, "[>>]")
		} else {
			parts = append(parts, room.Type.Icon())
		}
		if i < len(d.Floor.Rooms)-1 {
			parts = append(parts, "-")
		}
	}
	return "  " + strings.Join(parts, "")
}
