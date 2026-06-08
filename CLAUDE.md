# CLAUDE.md — Tamagotchi CLI

## Visão Geral

Jogo de Tamagotchi em terminal escrito em **Go 1.25**, módulo `Pessoal`. TUI construída com a stack Charmbracelet: **Bubble Tea** (arquitetura MVU/Elm), **Lip Gloss** (estilo/layout) e **Bubbles** (componentes como text input).

Não há TypeScript, React nem nenhuma dependência de Node.js neste projeto.

---

## Estrutura de Diretórios

```
tamago/main.go              # Entry point — carrega save, inicia Bubble Tea
internal/
  model/                    # Domínio puro (zero I/O)
    tama.go                 # Struct Tama + AddXP + Reset + Class
    evolution.go            # Enum Stage + thresholds de XP
    skills.go               # Skill struct + AllSkills + GetActiveSkills (12 skills + passivas)
    classes.go              # ClassID + Class + AllClasses (Guardia/Guerreiro/Mistico/Rapida)
    achievements.go         # Achievement + DefaultAchievements + CheckAchievements
    events.go               # Event + EventResult + slice Events
    quests.go               # Quest + QuestProgress + lógica CRUD
  dungeon/                  # Subsistema RPG de dungeon
    stats.go                # CombatStats + DeriveCombatStats + Apply*
    combat.go               # Motor de combate por turnos
    dungeon.go              # Máquina de estados DungeonRun
    enemy.go                # EnemyTemplate + Enemy + SpawnEnemy + pools por andar
    equipment.go            # Equipment + Inventory + AllEquipment (12 itens)
    floor.go                # Floor + Room + GenerateFloor (biome-aware)
    items.go                # Item + ItemBag + ShopItems
    biome.go                # Biome enum + BiomeModifier + ModifierForBiome
    ascii_art.go            # Helpers de renderização ASCII
    dungeon_sprites.go      # Assets de sprite para dungeon
  ui/                       # Camada Bubble Tea
    tui.go                  # Model + Init/Update/View + handleCommand
    update.go               # Lógica de tick, handlers de ação (Feed, Water, Pet…)
    view.go                 # Cálculo de layout + renderização de barras de stat
    dungeon_view.go         # Renderização completa do dungeon TUI
    quest_view.go           # Painel de missões
    minigames.go            # GuessGame + ReactGame + enum GameMode
    sprites.go              # Arte sprite do Tama por estágio/estado
  persistence/
    save.go                 # SaveData + Save/Load/Delete + migrações de schema
```

---

## Arquitetura Central

### Bubble Tea (MVU)

O padrão é **Model → Update → View**, não OOP com estado mutável. O `Model` em `ui/tui.go` é passado por valor e retornado modificado em cada `Update`. Toda mutação de estado do jogo acontece em `Update` ou em funções chamadas por ele.

Três tickers paralelos inicializados em `Init()` via `tea.Batch`:
- `tickCmd()` — 500ms, incrementa frame de animação
- `lifeCycleTickCmd()` — 20s, decai stats vitais, atualiza quests, dispara eventos
- `autoSaveTickCmd()` — 60s, salva em disco

### GameMode (enum em `ui/minigames.go`)

```go
ModeNormal       // Tela principal
ModeGuess        // Mini-game adivinhação
ModeReact        // Mini-game reação
ModeDungeon      // Dungeon RPG
ModeClassSelect  // Seleção de classe no onboarding
```

O `View()` despacha o render com base em `gameMode` e flags booleanas como `showingQuests`.

---

## Domínio: `model.Tama`

Struct principal do jogo (`internal/model/tama.go`). Campos relevantes:

```go
// Vitais (0–100)
Hunger, Thirst, Sleepy, Happiness, Angry, Weight int

// Progressão
XP, Level int
Stage Stage   // Baby → Child → Teen → Adult → Elder

// Combate
MP, MPMax int
SkillsKnown []string   // IDs de habilidades (ex: ["pancada", "cura"])
Class ClassID          // (a implementar) — classe selecionada no onboarding

// Missões
QuestProgress []QuestProgress
```

`Reset(name string)` é chamado no Game Over para reiniciar o jogo do zero. Ao implementar classes, deve também limpar `Class` para forçar nova seleção.

---

## Dungeon: Motor de Combate

### Fluxo de um turno (`combat.go:ExecuteAction`)

1. Decrementa duração de buffs ativos e reverte os expirados (`tickActiveBuffs`)
2. Decrementa cooldowns de habilidades
3. Aplica efeitos de bioma por turno (DoT fogo, regen HP, MP regen de Místico)
4. Se `FrozenFor > 0`: pula ação do jogador
5. Executa ação escolhida: Atacar / Defender / Item / Habilidade / Fugir
6. Verifica morte do inimigo
7. Inimigo tenta esquivar (se passiva de Rápida ativa) → ataca se não esquivou
8. Verifica morte do jogador
9. Sincroniza HP/MP de volta para `OriginalStats`

### Fórmulas de dano

```
// Ataque básico do jogador
baseDmg = Ataque - Defesa_inimigo/2  (min 1)
variância = ±20%
crítico = Sorte% → ×1.5

// Dano de habilidade
dmg = Ataque * (Power/100) - Defesa_inimigo/2  (min 1)

// Cura de habilidade
heal = HPMax * Power/100
```

### Derivação de stats (`stats.go:DeriveCombatStats`)

1. **Stats base por estágio:**
   ```
   mult = stageMult(stage)   // Baby=0.6, Child=0.8, Teen=1.0, Adult=1.2, Elder=1.1
   HP   = (50 + level*10) * mult
   MP   = (20 + level*5)  * mult
   ATK  = (8  + level*3)  * mult
   DEF  = (3  + level*2)  * mult
   VEL  = 5 + level
   LUCK = 5 + level/2
   ```

2. **Bônus emocionais:** `Happiness > 70` → ATK+3, VEL+2 / `Hunger > 60` → HP+10

3. **Multiplicadores de classe:** Aplicados após base (ex: Guardia: HP×1.30, DEF×1.20, ATK×0.90)

4. **Passivas de stat em derivação:**
   - `resistencia` (Guardia): DEF ×1.15
   - `instinto` (Guerreiro): LUCK +15
   - `reflexos` (Rápida): LUCK +20

---

## Sistema de Habilidades (`model/skills.go`)

```go
type Skill struct {
    ID, Name, Description string
    Type      SkillType   // SkillDamage | SkillHeal | SkillBuff | SkillDebuff | SkillPassive
    Cost      int         // custo de MP
    Power     int         // multiplicador em %
    Element   string
    IsPassive bool        // Não aparece no menu; efeito pré-calculado em combate
    Cooldown  int         // Turnos de espera após uso
    BuffStat  string      // "atk" | "def" | "vel" | "dodge" — stat alvo do buff
    BuffTurns int         // Duração do buff em turnos
}
```

**Habilidades por classe:**
- **Guardia:** Muralha (DEF buff CD:2), Golpe de Escudo (120% ATK CD:1), Resistência (passiva DEF+15%)
- **Guerreiro:** Golpe Brutal (200% ATK CD:2), Fúria (ATK buff 3t CD:3), Instinto (passiva LUCK+15)
- **Místico:** Cura Profunda (35% heal CD:2), Barreira (DEF buff 3t CD:3), Fluxo Arcano (passiva MP regen +3/t)
- **Rápida:** Corte Veloz (110% ATK sem CD), Esquiva (dodge buff 2t CD:3), Reflexos (passiva LUCK+20)

As 4 habilidades genéricas (`pancada`, `cura`, `rugido`, `foco`) ainda existem como fallback.

---

## Persistência e Migrações

**Arquivo:** `tamago_save.json` (escrita atômica via tmp+rename)

```go
type SaveData struct {
    SchemaVersion int
    SavedAt       time.Time
    Tama          *model.Tama
    Inventory     json.RawMessage
}
```

`CurrentSchemaVersion = 3`. 

**Migrações automáticas:**
- **v0 → v1:** Envelope SaveData (flat → structured)
- **v1 → v2:** Sistema de quests (InitDefaultQuests)
- **v2 → v3:** Classes — saves antigos recebem `ClassGuerreiro` com suas skills base automaticamente

O sistema calcula **degradação offline**: ao carregar, computa ticks perdidos (tempo offline / 20s), capped em 50, aplicando decaimento de stats. Tama pode morrer offline.

---

## Sistema de Missões (`model/quests.go`)

10 missões (2 story, 3 daily, 5 repeatable). Tracking via baseline: `progresso = valorAtual - baseline`. Daily resets à meia-noite. Claim via comando `claim <id>`. Toggle de painel: comando `q` / `quests`. Ver `ui/quest_view.go` como modelo para novos painéis de toggle.

---

## Biomas (`dungeon/biome.go`)

4 biomas: Florestal, Glacial, Vulcânico, Abissal. Cada um aplica `BiomeModifier` com ~18 campos cobrindo: bônus/penalidades de HP/ATK/DEF/VEL/LUCK em %, DoT de fogo, freeze chance, multiplicadores de taxa de vitais fora de combate.

---

## Convenções do Projeto

- **Go idiomático**: nomes claros, funções com responsabilidade única, sem `any`
- **Sem comentários redundantes** — o código deve ser autoexplicativo; comentar apenas o *porquê*
- **Constantes nomeadas** para valores mágicos (ex: `MaxHunger = 100`, não `100` solto)
- **Retorno antecipado** para reduzir nesting
- **Sem tratamento de erros silencioso** — erros devem ser propagados ou logados
- **Sem over-engineering** — YAGNI; não criar abstrações antes da segunda ocorrência

---

## Status da Implementação (TASKS.md, Seção 3–5)

**✅ Sistema de Classes e Habilidades — CONCLUÍDO**

Implementado:
1. ✅ `internal/model/classes.go` — 4 classes com multiplicadores de stats
2. ✅ Estender `Skill` com `IsPassive`, `Cooldown`, `BuffStat`, `BuffTurns`
3. ✅ `Tama.Class ClassID` — novo campo + migração v3
4. ✅ `DeriveCombatStats` — aplica multiplicadores de classe e passivas de stat
5. ✅ `Combat` — buffs multi-turno (`ActiveBuff`), cooldowns, dodge, MP regen
6. ✅ TUI — `ModeClassSelect` no onboarding, painel `habilidades` (comando `hab`)
7. ✅ `ui/skills_view.go` — painel renderizando ativas e passivas
8. ✅ Schema v3 — migração automática de saves antigos

**✅ Sistema de Loja Expandido — CONCLUÍDO**

Implementado:
1. ✅ `ItemCategory` (4 categorias) + `Rarity` em `Item` (`internal/dungeon/items.go`)
2. ✅ `Equipment.Price` com tabela por raridade (`internal/dungeon/equipment.go`)
3. ✅ `ShopEntry` tipo unificado (item ou equipment)
4. ✅ `DungeonRun` novos campos: `ShopTab`, `ShopMode`, `ShopStock`, `ShopConfirm`, `ShopPending`
5. ✅ `generateShopStock(floorNum)` — geração dinâmica: Comum (andares 1–3), Raro (4–6), Lendário (7+)
6. ✅ 9 itens consumíveis (Consumível + Qualidade + NPC Único com stock=1)
7. ✅ 12 equipamentos com preços (Comum 30–50g, Incomum 80–120g, Raro 200–300g, Lendário 600–800g)
8. ✅ Navegação: cursor ↑↓/kj, tabs ←→/hl, modo buy/sell (v), confirmação (s/n)
9. ✅ TUI nova: abas, modo visual, items com raridade e descrição, confirmação inline

**✅ Mais chefes, encontros & loot — CONCLUÍDO**

Implementado:
1. ✅ 6 encontros temáticos por bioma (`BiomeEnemyPools` com ranges de floor):
   - Florestal: Druida das Raizes (fl. 1-2), Lobo Sombrio (fl. 3-4)
   - Gélido: Espirito de Gelo (fl. 2-3), Urso Glacial (fl. 3-4)
   - Vulcânico: Elemental de Magma (fl. 2-3)
   - Abissal: Sombra Profunda (fl. 3-4)
2. ✅ `RandomEnemyForFloor` corrigida para usar bioma
3. ✅ 4 chefes únicos por bioma (Guardiao da Floresta, Lich do Gelo Eterno, Lorde das Chamas, Devorador do Abismo)
4. ✅ `BossForBiome(biome, playerLevel)` substituindo `BossForFloor5`
5. ✅ Mecânica de fases: 3 fases + enrage (triggers em 75%, 50%, 25% HP)
6. ✅ `Combat` com `BossPhase int` e `BossEnraged bool`
7. ✅ `checkBossPhaseTransition()` e `applyBossPhaseEffects()` em `combat.go`
8. ✅ Efeitos únicos por chefe (buffs permanentes, self-heal, múltiplos debuffs)
9. ✅ 4 equipamentos lendários exclusivos de chefe em `AllEquipment`
10. ✅ `bossLootForBiome()` em `floor.go`
11. ✅ `handleCombatResult` com drop garantido para chefes
12. ✅ Victory message dinâmica com nome real do chefe
13. ✅ Testes: `boss_test.go` com validação de pools, bosses e drops (todos passando)
