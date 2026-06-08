## Roadmap: Novas masmorras & Biomas

Uma visão organizada das tarefas e subtarefas para implementação das novas masmorras, biomas e sistemas relacionados.

---

### 1) Novas masmorras & biomas

- [x] Definir biomas e efeitos principais
	- [x] Florestal — +recuperação de fome; +encontros amistosos; -vulnerabilidade a fogo
	- [x] Gélido — -velocidade; +DEF gelo; chance de congelamento
	- [x] Vulcânico — +ATK inimigos; DOT (fogo over-time); -felicidade
	- [x] Abissal — +sorte inimiga; +chance de loot raro; -regeneração de HP do `Tama`

- [x] Projetar layout por bioma (tipos de salas, tema visual, pontos de descanso/loja)

- [x] Atualizar gerador de andares para suportar biomas (aleatoriedade/seed)

- [x] Criar pools de inimigos e encontros por bioma (incluindo inimigos raros exclusivos)

- [x] Adicionar assets ASCII por bioma (`internal/dungeon/ascii_art.go`)

- [x] Mapear salas especiais: boss room, treasure, descanso/loja, eventos raros

---

### 2) Efeitos de bioma nos atributos

- [x] Especificar tabela de modificadores por bioma (HP, ATK, DEF, VEL, LUCK, regen)

- [x] Implementar aplicação dos modificadores em combate (`internal/dungeon/combat.go`)

- [x] Aplicar modificadores fora de combate (ticks: fome/sono/felicidade enquanto no bioma)

- [x] Balancear valores e adicionar testes unitários para validação

---

### 3) Habilidades especiais & classes

- [x] Criar `internal/model/skills.go` com struct `Skill` estendida (IsPassive, Cooldown, BuffStat, BuffTurns)
- [x] Definir banco de habilidades: 4 genéricas + 12 específicas por classe
- [x] Implementar lógica de desbloqueio por classe (vs. estágio)
- [x] Atualizar struct `Tama` com campo `Class ClassID`
- [x] Calcular `MPMax` escalonado por nível/estágio e classe em `internal/dungeon/stats.go`
- [x] Incluir `ActionSkill` e `executeSkill` com cooldowns e buffs multi-turno em `internal/dungeon/combat.go`
- [x] Atualizar `internal/dungeon/dungeon.go` com filtro de skills ativas (remover passivas do menu)
- [x] Garantir compatibilidade com saves antigos (v3 migration: atribui classe Guerreiro)
- [x] Atualizar `internal/ui/dungeon_view.go` com exibição de cooldowns `[CD:X]`
- [x] Definir 4 classes com multiplicadores de stats
	- [x] `Guardiã` — Tank (HP×1.30, DEF×1.20)
	- [x] `Guerreiro` — Dano (ATK×1.30)
	- [x] `Místico` — Suporte (MP×1.40)
	- [x] `Rápida` — Velocidade/Crítico (VEL×1.30, LUCK×1.40)

- [x] Implementar 3 habilidades por classe (2 ativas + 1 passiva, 12 total)

- [x] Sistema de buffs multi-turno, cooldowns, dodge, MP regen per-turno

- [x] UI/TUI: `ModeClassSelect` no onboarding com navegação ↑↓ e seleção 1-4
- [x] Painel de habilidades (comando `hab`/`habilidades`) renderizando ativas e passivas
- [x] Schema v3 com migração automática de saves antigos

---

### 4) Sistema de missões & quests

- [x] Definir tipos de missão: story (única), diárias, repeatable

- [x] Criar 10 missões iniciais
	- [x] 2 × Story
	- [x] 3 × Diárias
	- [x] 5 × Repeatable

- [x] Criar esquema de dados persistente e serializar em `tamago_save.json`

- [x] Implementar fluxo: gerar → aceitar → completar → reivindicar recompensa + notificações

- [x] Integrar missões com dungeons, minigames e conquistas

---

### 5) Expandir loja & comércio com NPCs

- [x] Adicionar categorias de loja: Consumíveis, Equipamento, Itens de qualidade, Itens únicos NPC

- [x] Implementar restock dinâmico (rotação por andar / acesso ao entrar na loja)

- [x] Sistema de preços baseado em raridade e oferta/demanda

- [x] Interface TUI: compra / venda / confirmação (mostrar impacto em gold/inventário)

- [x] Balanceamento da economia e testes de fluxo

---

### 6) Troca livre com NPC (3-por-1)

- [x] Design da mecânica: NPC em descanso oferece 3 itens gratuitos; jogador escolhe 1

- [x] Implementar fluxo de interação na área de descanso

- [x] Balancear oferta: regras de raridade/nível e frequência de aparição

- [x] Testes de UX para validar impacto no gameplay

---

### 7) Mais chefes, encontros & loot

- [x] Criar 6 novos encontros temáticos

- [x] Projetar 4 chefes (um por bioma) com mecânicas únicas e fases

- [x] Definir tabela de drops e raridades, integrar com `AllEquipment`

- [x] Implementar mecânicas de chefe (phases, enrage, triggers)

- [x] Integrar drops ao `Inventory` e sistema de loot

---

### 8) PvE balanceamento & escalonamento

- [ ] Definir fórmula de scaling por nível do `Tama` e floor

- [ ] Simular runs e ajustar XP / gold / drop rates

- [ ] Adicionar configurações de dificuldade: Fácil / Normal / Difícil

- [ ] Definir KPIs de balanceamento (winrate, tempo médio, gold/level)

---

### 9) Sistema de amigos & trocas entre jogadores

- [ ] Especificar friend list (adicionar / remover / limites)

- [ ] Implementar fluxo de presentes e trocas (offline-friendly)

- [ ] Proteções anti-abuso: cooldowns e limites diários

- [ ] Protótipo offline: troca por token local

---

### 10) PvP batalhas (online — bônus)

- [ ] Definir regras PvP (turnos, itens permitidos, condições de vitória)

- [ ] Implementar modo local (pass-and-play) para testes sem rede

- [ ] Planejar matchmaker mínimo e sincronização (se houver backend)

- [ ] Definir recompensas/placares PvP (XP/elo/medalhas)

---

### 11) Ranking global & leaderboards (bônus)

- [ ] Especificar métricas (melhor run, MMR, vitórias PvP, tempo)

- [ ] Definir protocolo mínimo para submissão de scores

- [ ] Implementar protótipo local (top-10 salvo localmente)

---

### 12) Matchmaking / infra online (bônus)

- [ ] Escolher stack (REST + WebSocket recomendado)

- [ ] Definir endpoints essenciais: `/auth`, `/submit-score`, `/leaderboard`, `/friends`

- [ ] Criar protótipo de server em Go (esqueleto) e scripts básicos de deploy (opcional)

---

### 13) Expansão de conquistas & sincronização

- [ ] Listar novas conquistas relacionadas às features adicionadas

- [ ] Implementar migrador de saves para adicionar/validar novas conquistas em `tamago_save.json`

- [ ] Notificações & UI para exibir unlocks e histórico de conquistas

---

### 14) UI melhorias & polimento

- [ ] Atualizar HUD para mostrar bioma atual, classe selecionada e missão ativa

- [ ] Criar nova arte ASCII/estados para chefes e biomas (`internal/dungeon/ascii_art.go`)

- [ ] Onboarding para explicar classes, biomas e funcionamento da loja

- [ ] Melhorar feedback de combate (dano, buffs, debuffs, loot)

---

### 15) API / backend scaffolding (bônus)

- [ ] Criar esqueleto de API para autenticação, leaderboard e friends

- [ ] Priorizar endpoints: `/auth`, `/submit-score`, `/leaderboard`, `/friends`

- [ ] Documentar contrato e exemplos no `README.md`

---

### 16) Testes, documentação e migração

- [x] Criar esquema de dados persistente com `SaveData` envelope e `schema_version`
- [x] Implementar migração automática de saves antigos (formato flat → envelope)
- [x] Escrita atômica via `.tmp` + `os.Rename` para evitar corrupção por crash
- [x] Expandir `save_test.go` com testes de roundtrip, degradação offline, migração legacy e atomicidade
- [x] Escrever testes unitários para quests (`internal/model/quests_test.go`)
- [ ] Escrever testes unitários para biomas, classes e combat hooks

- [x] Criar script de migração para `tamago_save.json` (novos campos: schema v1→v2 para sistema de quests)

- [ ] Atualizar `README.md` com comandos, lojas, classes e dungeons

- [ ] Planejar 5 playtests internos para feedback de balanceamento

- [ ] Revisar e polir mensagens/feedback de loot e conquistas
