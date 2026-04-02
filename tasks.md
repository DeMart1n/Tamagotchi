## Roadmap: Novas masmorras & Biomas

Uma visão organizada das tarefas e subtarefas para implementação das novas masmorras, biomas e sistemas relacionados.

---

### 1) Novas masmorras & biomas

- [x] Definir biomas e efeitos principais
	- [x] Florestal — +recuperação de fome; +encontros amistosos; -vulnerabilidade a fogo
	- [x] Gélido — -velocidade; +DEF gelo; chance de congelamento
	- [x] Vulcânico — +ATK inimigos; DOT (fogo over-time); -felicidade
	- [x] Abissal — +sorte inimiga; +chance de loot raro; -regeneração de HP do `Tama`

- [ ] Projetar layout por bioma (tipos de salas, tema visual, pontos de descanso/loja)

- [ ] Atualizar gerador de andares para suportar biomas (aleatoriedade/seed)

- [x] Criar pools de inimigos e encontros por bioma (incluindo inimigos raros exclusivos)

- [x] Adicionar assets ASCII por bioma (`internal/dungeon/ascii_art.go`)

- [ ] Mapear salas especiais: boss room, treasure, descanso/loja, eventos raros

---

### 2) Efeitos de bioma nos atributos

- [ ] Especificar tabela de modificadores por bioma (HP, ATK, DEF, VEL, LUCK, regen)

- [ ] Implementar aplicação dos modificadores em combate (`internal/dungeon/combat.go`)

- [ ] Aplicar modificadores fora de combate (ticks: fome/sono/felicidade enquanto no bioma)

- [ ] Balancear valores e adicionar testes unitários para validação

---

### 3) Habilidades especiais & classes

- [ ] Definir 4 classes iniciais e estatísticas base
	- [ ] `Guardiã` — Tank (HP ↑, DEF ↑)
	- [ ] `Guerreiro` — Dano (ATK ↑)
	- [ ] `Místico` — Suporte (buffs/heal)
	- [ ] `Rápida` — Crítico/Velocidade (VEL ↑, CRIT ↑)

- [ ] Projetar 3 habilidades por classe (2 ativas + 1 passiva)

- [ ] Especificar números base: efeitos %, duração, cooldowns, custos

- [ ] Implementar framework de skills e hooks em `combat.go` e `model/` (skill structs)

- [ ] UI/TUI: seleção de classe no onboarding e painel de habilidades visível

---

### 4) Sistema de missões & quests

- [ ] Definir tipos de missão: story (única), diárias, repeatable

- [ ] Criar 10 missões iniciais
	- [ ] 2 × Story
	- [ ] 3 × Diárias
	- [ ] 5 × Repeatable

- [ ] Criar esquema de dados persistente e serializar em `tamago_save.json`

- [ ] Implementar fluxo: gerar → aceitar → completar → reivindicar recompensa + notificações

- [ ] Integrar missões com dungeons, minigames e conquistas

---

### 5) Expandir loja & comércio com NPCs

- [ ] Adicionar categorias de loja: Consumíveis, Equipamento, Itens de qualidade, Itens únicos NPC

- [ ] Implementar restock dinâmico (rotação diária / por descanso)

- [ ] Sistema de preços baseado em raridade e oferta/demanda

- [ ] Interface TUI: compra / venda / confirmação (mostrar impacto em gold/inventário)

- [ ] Balanceamento da economia e testes de fluxo

---

### 6) Troca livre com NPC (3-por-1)

- [ ] Design da mecânica: NPC em descanso oferece 3 itens gratuitos; jogador escolhe 1

- [ ] Implementar fluxo de interação na área de descanso

- [ ] Balancear oferta: regras de raridade/nível e frequência de aparição

- [ ] Testes de UX para validar impacto no gameplay

---

### 7) Mais chefes, encontros & loot

- [ ] Criar 6 novos encontros temáticos

- [ ] Projetar 4 chefes (um por bioma) com mecânicas únicas e fases

- [ ] Definir tabela de drops e raridades, integrar com `AllEquipment`

- [ ] Implementar mecânicas de chefe (phases, enrage, triggers)

- [ ] Integrar drops ao `Inventory` e sistema de loot

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

- [ ] Escrever testes unitários para biomas, classes, quests e combat hooks

- [ ] Criar script de migração para `tamago_save.json` (novos campos)

- [ ] Atualizar `README.md` com comandos, lojas, classes e dungeons

- [ ] Planejar 5 playtests internos para feedback de balanceamento

- [ ] Revisar e polir mensagens/feedback de loot e conquistas
