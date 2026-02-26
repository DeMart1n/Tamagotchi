# Dungeon RPG — Documentacao do Sistema

## Como Jogar

1. No jogo, atinja **Level 3** (ou use `setlvl 3` para testar)
2. Digite `dungeon` ou `d` para entrar na masmorra
3. No menu da masmorra:
   - `1` — Entrar na Masmorra
   - `2` — Ver Inventario
   - `3` — Voltar ao modo normal
4. `ESC` sai da masmorra a qualquer momento

### Comando Secreto
- `setlvl <numero>` — Seta o level do Tama diretamente (ex: `setlvl 10`)

---

## Estrutura de Arquivos

```
internal/dungeon/
  stats.go         Stats de combate derivados do Tama
  equipment.go     Equipamentos, inventario, raridades
  items.go         Itens consumiveis (pocoes)
  enemy.go         Templates de inimigos, pools por andar
  floor.go         Geracao de andares e salas
  combat.go        Motor de combate por turnos
  dungeon.go       Maquina de estados principal
  ascii_art.go     ASCII art dos inimigos

internal/model/
  tama.go          +Inventory, +TotalDungeonRuns, +TotalBossesDefeated
  achievements.go  +3 conquistas de dungeon

internal/ui/
  minigames.go     +ModeDungeon
  view.go          +dungeonWidth/Height no layout
  dungeon_view.go  Renderizacao completa da dungeon
  tui.go           Integracao: routing, input, serializacao

tamago/main.go     Migracao de achievements para saves antigos
```

---

## Stats de Combate

Derivados do Level/Stage do Tama + bonus de equipamento:

| Stat       | Formula Base                  |
|------------|-------------------------------|
| HP Max     | (50 + Level*10) * stageMult   |
| Ataque     | (8 + Level*3) * stageMult     |
| Defesa     | (3 + Level*2) * stageMult     |
| Velocidade | 5 + Level                     |
| Sorte      | 5 + Level/2                   |

**Multiplicadores de estagio:**
- Baby: 0.6x
- Crianca: 0.8x
- Adolescente: 1.0x
- Adulto: 1.2x
- Anciao: 1.1x

**Bonus condicionais:**
- Felicidade > 70: +3 ATK, +2 VEL
- Fome > 60: +10 HP

---

## Equipamentos (12 total)

3 slots: Arma, Armadura, Acessorio
4 raridades com cores: Comum (cinza), Incomum (verde), Raro (azul), Lendario (dourado)

### Armas
| Nome             | Raridade  | Bonus               |
|------------------|-----------|----------------------|
| Espada de Ferro  | Comum     | +5 ATK               |
| Espada Afiada    | Incomum   | +10 ATK, +2 VEL      |
| Machado de Guerra| Raro      | +18 ATK, +3 DEF      |
| Katana Lendaria  | Lendario  | +25 ATK, +5 VEL, +5 LUCK |

### Armaduras
| Nome              | Raridade  | Bonus                    |
|-------------------|-----------|--------------------------|
| Armadura de Couro | Comum     | +4 DEF, +5 HP            |
| Armadura de Ferro | Incomum   | +8 DEF, +10 HP           |
| Armadura de Mithril| Raro     | +14 DEF, +20 HP, +2 VEL  |
| Armadura de Dragao| Lendario  | +22 DEF, +35 HP, +3 VEL  |

### Acessorios
| Nome              | Raridade  | Bonus                         |
|-------------------|-----------|-------------------------------|
| Anel da Sorte     | Comum     | +8 LUCK                       |
| Amuleto de Vigor  | Incomum   | +15 HP, +3 DEF                |
| Bracelete da Furia| Raro      | +10 ATK, +5 VEL               |
| Coroa do Rei      | Lendario  | +8 ATK, +8 DEF, +20 HP, +10 LUCK |

**Drops sao aleatorios** — filtrados por raridade maxima do andar:
- Andares 1-2: ate Incomum
- Andar 3: ate Raro
- Andares 4-5: todas as raridades

---

## Itens Consumiveis

Comprados na loja de descanso, usados em combate:

| Nome            | Tipo      | Efeito          | Preco |
|-----------------|-----------|-----------------|-------|
| Pocao Pequena   | Cura HP   | +25 HP          | 10    |
| Pocao Grande    | Cura HP   | +50 HP          | 25    |
| Elixir de Forca | Buff ATK  | +8 ATK (1 turno)| 20    |

---

## Inimigos (por andar)

Todos os inimigos **escalam com o level do jogador** (+8% por level acima de 3).

| Andar | Inimigos              | HP Base  | ATK Base | DEF | XP Drop |
|-------|-----------------------|----------|----------|-----|---------|
| 1     | Slime                 | 20-25    | 6-8      | 2   | 15      |
| 1     | Rato Gigante          | 18-22    | 7-9      | 1   | 18      |
| 2     | Esqueleto             | 30-38    | 10-12    | 5   | 30      |
| 2     | Morcego Vampiro       | 25-32    | 11-13    | 3   | 28      |
| 3     | Goblin Guerreiro      | 45-55    | 14-16    | 8   | 45      |
| 3     | Aranha Venenosa       | 40-50    | 15-17    | 6   | 42      |
| 4     | Cavaleiro Negro       | 60-75    | 22-26    | 14  | 70      |
| 4     | Mago Sombrio          | 50-65    | 25-28    | 8   | 75      |
| 5     | **Dragao Anciao (BOSS)** | 200   | 35       | 18  | 200     |

**Escalamento por level:** Exemplo com Level 10 (escala 1.56x):
- Slime: HP ~35, ATK ~11
- Dragao Anciao: HP ~312, ATK ~55

---

## Estrutura dos Andares

### Andares 1-4
```
[Combate] → [Combate] → [Descanso ou Tesouro*] → [Combate]
```
*40% chance de tesouro, 60% descanso

### Andar 5
```
[Combate] → [Descanso] → [Combate] → [BOSS]
```

---

## Sistema de Combate

Turnos: jogador age → verifica morte inimigo → inimigo age → verifica morte jogador.

### Acoes do Jogador
1. **Atacar** — dano = ATK - DEF/2, variancia +-20%, critico = Sorte%  (1.5x dano)
2. **Defender** — dano recebido reduzido em 50% neste turno
3. **Item** — usar pocao/elixir do inventario
4. **Fugir** — 40% base + 2% por vantagem de velocidade (impossivel contra Boss)

### IA do Inimigo
- Sempre ataca
- 30% chance de defender quando HP < 25%

---

## Maquina de Estados

```
MenuPrincipal → [1] StartDungeon → Combate/Descanso/Tesouro
              → [2] Inventario → voltar
              → [3] PhaseDone (volta ao normal)

Combate → vitoria → CombatResult → proximo sala
        → derrota → PhaseDerrota → PhaseDone
        → fuga   → proximo sala

FimAndar → proximo andar
Vitoria (boss derrotado) → PhaseDone
Derrota → PhaseDone (XP parcial: metade do acumulado)
```

---

## Conquistas de Dungeon

| ID             | Nome              | Condicao                | Icone |
|----------------|-------------------|-------------------------|-------|
| first_dungeon  | Aventureiro       | Completou 1 run         | 🗡️   |
| dragon_slayer  | Mata-Dragao       | Derrotou o boss         | 🐉   |
| dungeon_master | Mestre da Masmorra| Completou 5 runs        | 🏰   |

---

## Persistencia

- **Inventario** (equipamentos, ouro, mochila) salva como `json.RawMessage` no Tama
- Auto-save a cada 60s, save ao sair
- Saves antigos recebem as novas conquistas automaticamente na inicializacao
- Itens consumiveis (pocoes) **nao persistem** entre runs — sao comprados durante a run

---

## Requisitos

- **Level minimo: 3** para entrar na dungeon
- O jogo roda apenas em terminal real (TTY) — nao funciona em sandbox

## Como Rodar

```bash
cd ~/Developer/Tamagotchi
go run tamago/main.go
```
