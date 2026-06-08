package ui

import (
	"Pessoal/internal/dungeon"
	"Pessoal/internal/model"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// --- Estilos da Dungeon ---
var (
	styleDungeonBase = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				Padding(0, 1)

	styleDungeonTitleBase = lipgloss.NewStyle().Bold(true)

	styleEnemyName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	stylePlayerInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#43BF6D")).
			Bold(true)

	styleCombatLog = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7B538")).
			Italic(true)

	styleMana = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00d2ff")).
			Bold(true)

	styleRarityComum    = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleRarityIncomum  = lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D"))
	styleRarityRaro     = lipgloss.NewStyle().Foreground(lipgloss.Color("#5B9FFF"))
	styleRarityLendario = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
)

func getBiomeStyle(biome dungeon.Biome) (lipgloss.Style, lipgloss.Style) {
	var color string
	switch biome {
	case dungeon.BiomeForest:
		color = "#2ecc71" // Verde Esmeralda
	case dungeon.BiomeIcy:
		color = "#00d2ff" // Azul Glacial
	case dungeon.BiomeVolcanic:
		color = "#e74c3c" // Vermelho Lava
	case dungeon.BiomeAbyssal:
		color = "#8e44ad" // Roxo Profundo
	default:
		color = "#FFD700" // Ouro padrão
	}

	borderStyle := styleDungeonBase.Copy().BorderForeground(lipgloss.Color(color))
	titleStyle := styleDungeonTitleBase.Copy().Foreground(lipgloss.Color(color))

	return borderStyle, titleStyle
}

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
	borderStyle, titleStyle := getBiomeStyle(d.CurrentBiome)

	switch d.Phase {
	case dungeon.PhaseIntro:
		content = m.renderDungeonIntro(l)
	case dungeon.PhaseMenuPrincipal:
		content = m.renderDungeonMenu(l, titleStyle)
	case dungeon.PhaseInventario:
		content = m.renderDungeonInventory(l, titleStyle)
	case dungeon.PhaseCombate:
		content = m.renderDungeonCombat(l, titleStyle)
	case dungeon.PhaseCombatResult:
		content = m.renderDungeonCombat(l, titleStyle)
	case dungeon.PhaseDescanso:
		content = m.renderDungeonRest(l, titleStyle)
	case dungeon.PhaseDescansoLoja:
		content = m.renderDungeonShop(l, titleStyle)
	case dungeon.PhaseTesouro:
		content = m.renderDungeonTreasure(l, titleStyle)
	case dungeon.PhaseFimAndar:
		content = m.renderDungeonFloorEnd(l, titleStyle)
	case dungeon.PhaseVitoria:
		content = m.renderDungeonVictory(l, titleStyle)
	case dungeon.PhaseDerrota:
		content = m.renderDungeonDefeat(l)
	default:
		content = d.Message
	}

	return borderStyle.Width(l.dungeonWidth).Render(content)
}

// renderDungeonIntro mostra a tela de entrada do bioma com pixel art de fundo.
func (m Model) renderDungeonIntro(l layout) string {
	d := m.dungeonGame
	_, titleStyle := getBiomeStyle(d.CurrentBiome)
	modifier := dungeon.ModifierForBiome(d.CurrentBiome)

	var lines []string

	// Cenário pixel art do bioma (pintura de fundo)
	biomeSprite := BiomeSprite(d.CurrentBiome.String())
	if biomeSprite != "" {
		for _, sl := range strings.Split(biomeSprite, "\n") {
			lines = append(lines, "  "+sl)
		}
	}

	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("  "+d.CurrentBiome.String()))
	lines = append(lines, "")

	// Efeitos do bioma com estilo
	effects := dungeon.FormatBiomeEffects(modifier)
	if effects != "" {
		lines = append(lines, styleCombatLog.Render("  "+effects))
	}

	lines = append(lines, "")
	lines = append(lines, "  Pressione ENTER para entrar...")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonMenu(l layout, titleStyle lipgloss.Style) string {
	title := titleStyle.Render("MASMORRA DE TAMAGO")
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

func (m Model) renderDungeonInventory(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame
	inv := d.Inv

	var lines []string
	lines = append(lines, titleStyle.Render("INVENTARIO"))
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

func (m Model) renderDungeonCombat(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame
	c := d.Combat
	if c == nil {
		return d.Message
	}

	var lines []string

	// Progresso do andar
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())

	// Nome da sala especial (ex: Zona Crítica)
	room := d.Floor.CurrentRoomRef()
	if room != nil && room.Name != "" {
		lines = append(lines, titleStyle.Render("  "+room.Name))
		if room.Description != "" {
			lines = append(lines, styleCombatLog.Render("  "+room.Description))
		}
	}
	lines = append(lines, "")


	// Inimigo
	lines = append(lines, styleEnemyName.Render("  "+c.Enemy.Name))

	// Sprite do inimigo + efeito de combate lado a lado
	sprite := EnemySprite(c.Enemy.Name)
	effectStr := ""
	if len(c.Log) > 0 {
		effectType := detectCombatEffect(c.Log)
		if effectType != "" {
			effectStr = CombatEffect(effectType)
		}
	}

	if effectStr != "" {
		// Junta sprite (esquerda) + espaço + efeito (direita)
		combined := lipgloss.JoinHorizontal(lipgloss.Top, sprite, "  ", effectStr)
		for _, cl := range strings.Split(combined, "\n") {
			lines = append(lines, "  "+cl)
		}
	} else {
		for _, spriteLine := range strings.Split(sprite, "\n") {
			lines = append(lines, "  "+spriteLine)
		}
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
	
	manaBar := dungeon.MPBar(d.Stats.MPCurrent, d.Stats.MPMax, 16)
	lines = append(lines, styleMana.Render(fmt.Sprintf("  MP %s %d/%d", manaBar, d.Stats.MPCurrent, d.Stats.MPMax)))
	
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
	} else if d.ChoosingSkill {
		lines = append(lines, "  Habilidades:")
		activeSkills := model.GetActiveSkills(m.tama.SkillsKnown)
		for i, skillID := range activeSkills {
			skill := model.AllSkills[skillID]
			cdInfo := ""
			if d.Combat != nil {
				if cd := d.Combat.SkillCooldowns[skillID]; cd > 0 {
					cdInfo = fmt.Sprintf(" [CD:%d]", cd)
				}
			}
			lines = append(lines, fmt.Sprintf("    [%d] %-15s (%d MP)%s", i+1, skill.Name, skill.Cost, cdInfo))
		}
		lines = append(lines, "    [0] Voltar")
	} else {
		lines = append(lines, "  [1] Atacar   [2] Defender  [3] Habilidades")
		items := ""
		if d.ItemBag.Count() > 0 {
			items = fmt.Sprintf("(%d)", d.ItemBag.Count())
		}
		lines = append(lines, fmt.Sprintf("  [4] Item %s  [5] Fugir", items))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonRest(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame
	room := d.Floor.CurrentRoomRef()

	title := "Sala de Descanso"
	if room != nil && room.Name != "" {
		title = room.Name
	}

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("  "+title))
	restSprite := ItemSprite("rest")
	if restSprite != "" {
		for _, sl := range strings.Split(restSprite, "\n") {
			lines = append(lines, "  "+sl)
		}
	}
	if room != nil && room.Description != "" {
		lines = append(lines, styleCombatLog.Render("  "+room.Description))
	}
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

func (m Model) renderDungeonShop(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, titleStyle.Render("  Loja de Aventureiros"))
	lines = append(lines, "")

	// Tabs
	tabs := []string{"Consumivel", "Equipamento", "Qualidade", "NPC Unico"}
	var tabLine string
	for i, tab := range tabs {
		if i == d.ShopTab {
			tabLine += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Render("["+tab+"] ")
		} else {
			tabLine += tab + "  "
		}
	}
	lines = append(lines, "  "+tabLine)
	lines = append(lines, "")

	// Info bar
	modeStr := d.ShopMode
	if d.ShopMode == "buy" {
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#43BF6D")).Render("Compra")
	} else {
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render("Venda")
	}
	lines = append(lines, fmt.Sprintf("  Ouro: %d  |  Modo: %s", d.Inv.Gold, modeStr))
	lines = append(lines, "")

	// Items da tab atual
	itemsInTab := d.GetItemsInTab()
	if len(itemsInTab) == 0 {
		lines = append(lines, "  Nenhum item nesta categoria.")
	} else {
		for i, entry := range itemsInTab {
			cursor := " "
			if i == d.ShopCursor {
				cursor = ">"
			}

			var name, rarity, price, desc string

			if entry.Kind == "item" {
				name = entry.Item.Name
				rarity = entry.Item.Rarity.String()
				price = fmt.Sprintf("%dg", entry.Price)
				switch entry.Item.Type {
				case dungeon.ItemPotion:
					if entry.Item.Category == dungeon.CategoryNPCUnico {
						desc = "Permanente"
					} else {
						desc = fmt.Sprintf("+%d HP", entry.Item.Value)
					}
				case dungeon.ItemATKBoost:
					desc = fmt.Sprintf("+%d ATK", entry.Item.Value)
				}
			} else {
				name = entry.Equip.Name
				rarity = entry.Equip.Rarity.String()
				price = fmt.Sprintf("%dg", entry.Price)
				desc = entry.Equip.Description()
			}

			stock := ""
			if entry.Stock == 0 {
				stock = " [SOLD]"
			} else if entry.Stock > 0 {
				stock = fmt.Sprintf(" [%d]", entry.Stock)
			}

			rarityColored := rarityStyle(dungeon.Rarity(entry.Item.Rarity)).Render(rarity)
			if entry.Kind == "equipment" {
				rarityColored = rarityStyle(entry.Equip.Rarity).Render(rarity)
			}

			line := fmt.Sprintf("  %s %-20s %s %6s %s%s", cursor, name, rarityColored, price, desc, stock)
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")

	if d.ShopConfirm && len(itemsInTab) > 0 && d.ShopPending < len(itemsInTab) {
		entry := itemsInTab[d.ShopPending]
		var confirmMsg string
		if entry.Kind == "item" {
			confirmMsg = fmt.Sprintf("Comprar %s por %dg? (S/N)", entry.Item.Name, entry.Price)
		} else {
			confirmMsg = fmt.Sprintf("Comprar %s por %dg? (S/N)", entry.Equip.Name, entry.Price)
		}
		lines = append(lines, styleCombatLog.Render("  "+confirmMsg))
	} else if d.Message != "" {
		lines = append(lines, styleCombatLog.Render("  "+d.Message))
	}

	lines = append(lines, "")
	lines = append(lines, "  [←→] Tabs  [↑↓] Navegar  [Enter] Comprar  [V] Modo  [0] Sair")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonTreasure(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, m.renderRoomProgress())
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("  Bau do Tesouro!"))
	chestSprite := ItemSprite("chest")
	if chestSprite != "" {
		for _, sl := range strings.Split(chestSprite, "\n") {
			lines = append(lines, "  "+sl)
		}
	}
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

func (m Model) renderDungeonFloorEnd(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame

	var lines []string
	lines = append(lines, m.renderFloorProgress())
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render(fmt.Sprintf("  Andar %d Completo!", d.FloorNum)))
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

func (m Model) renderDungeonVictory(l layout, titleStyle lipgloss.Style) string {
	d := m.dungeonGame

	victoryTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true)
	rewardStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#73F59F"))

	var lines []string
	lines = append(lines, "")
	lines = append(lines, victoryTitle.Render("  === VITORIA! ==="))
	lines = append(lines, "")
	victorySprite := ItemSprite("victory")
	if victorySprite != "" {
		for _, sl := range strings.Split(victorySprite, "\n") {
			lines = append(lines, "      "+sl)
		}
	}
	lines = append(lines, "")
	lines = append(lines, victoryTitle.Render("  O Dragao Anciao foi derrotado!"))
	lines = append(lines, "")
	lines = append(lines, rewardStyle.Render(fmt.Sprintf("  + %d XP", d.TotalXP)))
	lines = append(lines, rewardStyle.Render(fmt.Sprintf("  + %d Ouro", d.TotalGold)))
	lines = append(lines, "")
	lines = append(lines, "  Aperte Enter para voltar.")

	return strings.Join(lines, "\n")
}

func (m Model) renderDungeonDefeat(l layout) string {
	d := m.dungeonGame

	defeatTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true)
	lossStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	var lines []string
	lines = append(lines, "")
	lines = append(lines, defeatTitle.Render("  === DERROTA ==="))
	lines = append(lines, "")
	defeatSprite := ItemSprite("defeat")
	if defeatSprite != "" {
		for _, sl := range strings.Split(defeatSprite, "\n") {
			lines = append(lines, "      "+sl)
		}
	}
	lines = append(lines, "")
	lines = append(lines, lossStyle.Render("  "+d.Message))
	if d.TotalXP > 0 {
		lines = append(lines, lossStyle.Render(fmt.Sprintf("  XP recuperado: %d", d.TotalXP/2)))
	}
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

// detectCombatEffect analisa os logs de combate e retorna o efeito visual adequado.
func detectCombatEffect(logs []string) string {
	for _, log := range logs {
		switch {
		case strings.Contains(log, "CRITICO"):
			return "critical"
		case strings.Contains(log, "atacou") || strings.Contains(log, "dano em"):
			return "attack"
		case strings.Contains(log, "defender"):
			return "defend"
		case strings.Contains(log, "Recuperou") || strings.Contains(log, "Cura"):
			return "heal"
		case strings.Contains(log, "Usou") && strings.Contains(log, "!"):
			return "skill"
		case strings.Contains(log, "fugiu"):
			return "flee"
		case strings.Contains(log, "Nao conseguiu fugir"):
			return "miss"
		case strings.Contains(log, "derrotado"):
			return "damage"
		}
	}
	return ""
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
