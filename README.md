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

## ✨ Funcionalidades Atuais (MVP)

- **Sistema de Necessidades Vitais**: Fome, sede, energia, felicidade e humor
- **Interface CLI Interativa**: Terminal UI responsivo e animado usando Bubble Tea
- **Estados Emocionais**: Seu TamaGO pode ficar feliz, irritado, deprimido ou furioso
- **Ciclo de Vida**: Sistema de degradação automática das estatísticas ao longo do tempo
- **Comandos Interativos**: Alimente, hidrate, faça carinho e interaja com seu pet
- **Animações**: Expressões faciais e emojis que mudam conforme o humor
- **Sistema de Morte**: Game Over se você não cuidar bem do seu TamaGO

---

## 🚀 Roadmap - Visão Futura

### Versão 1.0 - Sistema de Progressão
- [ ] Sistema de níveis e experiência (XP)
- [ ] Evolução do TamaGO com diferentes formas
- [ ] Estatísticas de atributos (Força, Defesa, Agilidade, etc.)

### Versão 2.0 - Sistema de Combate
- [ ] Batalhas entre TamaGOs (PvP)
- [ ] Batalhas contra NPCs/Monstros (PvE)
- [ ] Sistema de turnos tático
- [ ] Habilidades especiais e movimentos

### Versão 3.0 - Equipamentos e Customização
- [ ] Sistema de inventário
- [ ] Equipamentos (armas, armaduras, acessórios)
- [ ] Lojas para comprar itens
- [ ] Sistema de moedas/economia

### Versão 4.0 - Multiplayer e Rankings
- [ ] Ranking local de TamaGOs
- [ ] Ranking global online
- [ ] Sistema de amizades e trocas
- [ ] Torneios e eventos especiais

### Versão 5.0 - Mundo Expandido
- [ ] Missões e quests
- [ ] Exploração de dungeons
- [ ] Boss fights
- [ ] História e lore do mundo TamaGO

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
| `status` | - | Atualiza e mostra o status atual |
| `quit`, `exit`, `q` | - | Sai do jogo |

### Controles

- **Enter**: Envia o comando digitado
- **ESC / Ctrl+C**: Sai do jogo

### Dicas

- Mantenha todas as barras de status equilibradas
- Se a fome ou sede chegarem a zero, seu TamaGO pode morrer
- Um TamaGO feliz tem melhor desempenho
- Não o irrite muito, ou ele ficará furioso!
- As estatísticas diminuem automaticamente com o tempo, fique atento!

---

## 📁 Estrutura do Projeto

```
TamaGO/
├── tamago/
│   ├── main.go           # Ponto de entrada da aplicação
│   └── tamago            # Executável compilado
├── internal/
│   ├── model/
│   │   └── tama.go       # Modelo de dados do Tamagotchi
│   └── ui/
│       ├── tui.go        # Interface do terminal (Bubble Tea)
│       ├── view.go       # Renderização da interface
│       └── update.go     # Lógica de atualização e comandos
├── go.mod                # Dependências do projeto
├── go.sum                # Checksums das dependências
└── README.md             # Este arquivo
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
