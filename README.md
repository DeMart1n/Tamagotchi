# TamaGO 🎮

<div align="center">

**Um Tamagotchi RPG desenvolvido em Go com interface CLI interativa**

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-orange?style=for-the-badge)](https://github.com)

</div>

---

## 📖 Sobre o Projeto

**TamaGO** é um jogo de Tamagotchi RPG desenvolvido em Go que combina a nostalgia dos pets virtuais clássicos com mecânicas de RPG modernas. Atualmente em fase de **MVP**, o projeto oferece uma experiência interativa via linha de comando onde você cuida do seu TamaGO, mantendo-o alimentado, hidratado, feliz e descansado.

O objetivo é evoluir o TamaGO para um RPG completo com sistema de níveis, batalhas, equipamentos e rankings globais!

---

## ✨ Funcionalidades

### 🐾 Sistema de Pet Virtual
- **Necessidades Vitais**: Fome, sede, energia, felicidade e humor
- **Interface CLI Interativa**: Terminal UI responsivo e animado usando Bubble Tea
- **Estados Emocionais**: Seu TamaGO pode ficar feliz, irritado, deprimido ou furioso
- **Ciclo de Vida**: Sistema de degradação automática das estatísticas ao longo do tempo
- **Evolução**: Baby → Criança → Adolescente → Adulto → Ancião
- **Sistema de Morte**: Game Over se você não cuidar bem do seu TamaGO
- **Achievements**: Conquistas desbloqueáveis por marcos no jogo

### ⚔️ Sistema de Dungeon RPG
- **Masmorra de 5 andares** com combate por turnos, salas de descanso e tesouro
- **Combate tático**: Atacar, Defender, usar Itens ou Fugir
- **12 equipamentos** em 3 slots (Arma, Armadura, Acessório) com 4 raridades
- **9 tipos de inimigos** que escalam com o level do jogador
- **Boss Fight**: Dragão Ancião no andar 5
- **Inventário persistente**: Equipamentos e ouro salvos entre sessões
- **Stats derivados**: HP, ATK, DEF, VEL e Sorte baseados no level/estágio do Tama
- **Itens consumíveis**: Poções e elixirs comprados em salas de descanso
- **3 conquistas exclusivas** de dungeon (Aventureiro, Mata-Dragão, Mestre da Masmorra)

> 📄 Documentação completa do sistema de dungeon: [DUNGEON.md](DUNGEON.md)

---

## 🚀 Roadmap

### ✅ Implementado
- [x] Sistema de necessidades vitais e estados emocionais
- [x] Interface CLI com Bubble Tea
- [x] Sistema de níveis e experiência (XP)
- [x] Evolução do TamaGO com diferentes estágios
- [x] Stats de combate (HP, ATK, DEF, VEL, Sorte)
- [x] Dungeon RPG com 5 andares e boss fight
- [x] Sistema de combate por turnos
- [x] Equipamentos com raridades (Comum → Lendário)
- [x] Sistema de inventário e economia (ouro)
- [x] Itens consumíveis (poções e elixirs)
- [x] Achievements / Conquistas
- [x] Persistência de dados (save/load)

### 🔮 Futuro
- [ ] Batalhas PvP entre TamaGOs
- [ ] Ranking global online
- [ ] Novas dungeons e biomas
- [ ] Habilidades especiais e classes
- [ ] Missões e quests
- [ ] Sistema de amizades e trocas

---

## 🛠️ Tecnologias Utilizadas

- **[Go 1.25](https://go.dev/)** - Linguagem de programação principal
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Framework TUI para interfaces de terminal
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Estilização e layout para terminal
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - Componentes TUI reutilizáveis

---

## 📦 Instalação

### Pré-requisitos

- Go 1.25 ou superior instalado no sistema
- Terminal com suporte a cores e Unicode

### Passos

1. **Clone o repositório**
```bash
git clone https://github.com/seu-usuario/TamaGO.git
cd TamaGO
```

2. **Instale as dependências**
```bash
go mod download
```

3. **Compile o projeto**
```bash
cd tamago
go build -o tamago
```

4. **Execute o jogo**
```bash
./tamago
```

Ou compile e execute diretamente:
```bash
go run main.go
```

---

## 🎮 Como Jogar

### Comandos Disponíveis

| Comando | Atalho | Descrição |
|---------|--------|-----------|
| `feed` | `f` | Alimenta o TamaGO (reduz fome) |
| `water` | `w` | Dá água ao TamaGO (reduz sede) |
| `pet` | `p` | Faz carinho no TamaGO (aumenta felicidade) |
| `sleep` | `s` | Coloca o TamaGO para dormir (recupera energia) |
| `annoy` | `a` | Irrita o TamaGO (aumenta raiva) |
| `dungeon` | `d` | Entra na masmorra (requer Level 3+) |
| `status` | - | Atualiza e mostra o status atual |
| `quit`, `exit`, `q` | - | Sai do jogo |

### Comandos Secretos

| Comando | Descrição |
|---------|-----------|
| `setlvl <n>` | Seta o level do Tama diretamente (ex: `setlvl 10`) |

### Controles

- **Enter**: Envia o comando digitado
- **ESC**: Sai do jogo / Sai da dungeon
- **Ctrl+C**: Sai do jogo

### Dicas

- Mantenha todas as barras de status equilibradas
- Se a fome ou sede chegarem a zero, seu TamaGO pode morrer
- Um TamaGO feliz tem melhor desempenho (bônus de stats na dungeon!)
- Não o irrite muito, ou ele ficará furioso!
- As estatísticas diminuem automaticamente com o tempo, fique atento!
- Alcance Level 3 para desbloquear a Dungeon RPG

---

## 📁 Estrutura do Projeto

```
TamaGO/
├── tamago/
│   └── main.go              # Ponto de entrada e migração de saves
├── internal/
│   ├── model/
│   │   ├── tama.go          # Modelo de dados do Tamagotchi
│   │   ├── evolution.go     # Sistema de evolução por estágios
│   │   ├── achievements.go  # Conquistas desbloqueáveis
│   │   └── events.go        # Sistema de eventos
│   ├── ui/
│   │   ├── tui.go           # Interface do terminal (Bubble Tea)
│   │   ├── view.go          # Renderização da interface
│   │   ├── update.go        # Lógica de atualização e comandos
│   │   ├── minigames.go     # Modos de jogo (dungeon, etc.)
│   │   └── dungeon_view.go  # Renderização da dungeon
│   ├── dungeon/
│   │   ├── dungeon.go       # Máquina de estados principal
│   │   ├── combat.go        # Motor de combate por turnos
│   │   ├── enemy.go         # Inimigos e pools por andar
│   │   ├── floor.go         # Geração de andares e salas
│   │   ├── equipment.go     # Equipamentos e inventário
│   │   ├── items.go         # Itens consumíveis
│   │   ├── stats.go         # Stats de combate derivados
│   │   └── ascii_art.go     # ASCII art dos inimigos
│   └── persistence/
│       └── save.go          # Sistema de save/load
├── DUNGEON.md               # Documentação do sistema de dungeon
├── go.mod                   # Dependências do projeto
├── go.sum                   # Checksums das dependências
└── README.md                # Este arquivo
```

---

## 🎨 Preview

```
╭──────────────────────╮
│   🎮 TAMAGOTCHI CLI   │
╰──────────────────────╯

╭──────────────╮          ╭────────────────────────────╮
│              │          │       STATUS VITAIS        │
│   ✨ 😄 ✨   │         │                            │
│              │          │ Fome       ▇▇▇▇▇▇░░░   70% │
│    TamaGo     │         │ Sede       ▇▇▇▇▇▇▇░░   80% │
│   Radiante!  │          │ Energia    ▇▇▇▇░░░░░   45% │
│              │          │ Felicidade ▇▇▇▇▇▇▇▇░   90% │
╰──────────────╯          │ Calma      ▇▇▇▇▇▇▇▇▇   95% │
                          ╰────────────────────────────╯

📢 ✨ Olá! Cuide bem do TamaGo!
───────────────────────────────────────────
> feed_

ESC/Ctrl+C: Sair • Comandos: (f)eed, (w)ater, (p)et, (s)leep
```

---

## 🤝 Contribuindo

Contribuições são sempre bem-vindas! Se você tem ideias para melhorar o TamaGO:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/MinhaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona MinhaFeature'`)
4. Push para a branch (`git push origin feature/MinhaFeature`)
5. Abra um Pull Request

---

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 👥 Autores

Desenvolvido com ❤️ por [Cauã De Martin](https://github.com/DeMart1n) e [Luiz Felipe Arcanjo](https://github.com/luiz0ar)

---

## 🌟 Agradecimentos

- [Charm](https://charm.sh/) pela incrível biblioteca Bubble Tea
- Comunidade Go pela linguagem fantástica
- Tamagotchi original pela inspiração nostálgica

---

<div align="center">

**Se você gostou do projeto, deixe uma ⭐!**

[Reportar Bug](https://github.com/seu-usuario/TamaGO/issues) · [Solicitar Feature](https://github.com/seu-usuario/TamaGO/issues) · [Documentação](https://github.com/seu-usuario/TamaGO/wiki)

</div>
