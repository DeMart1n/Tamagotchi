# Projeto Tamago - Resumo Arquitetural (Sênior)

Este documento registra o estado técnico e as decisões de design do projeto Tamago para referência futura do Gemini CLI.

## 🏗️ Arquitetura Core
O projeto utiliza o padrão **MVU (Model-View-Update)** através do framework `Bubble Tea` (Charmbracelet). A estrutura é modular e fortemente desacoplada:

- **Framework TUI:** `github.com/charmbracelet/bubbletea` para orquestração e `lipgloss` para estilização.
- **Fluxo de Dados:** Unidirecional (Events -> Update -> Model -> View).

## 📂 Estrutura de Pacotes (`internal/`)
- **`model/`**: Contém as estruturas de dados puras (`Tama`), lógica de evolução, XP e o sistema de conquistas. É agnóstico à interface.
- **`ui/`**: Implementa a interface Bubble Tea. Gerencia múltiplos estados (Modo Normal, Mini-games, Transições).
- **`persistence/`**: Abstração de I/O. Implementa **escrita atômica** (escreve em `.tmp` e renomeia) para evitar corrupção e suporta **migração de schema** para compatibilidade com saves antigos.
- **`dungeon/`**: Um sub-sistema RPG completo com sua própria máquina de estados, sistema de combate e inventário, integrado via composição ao modelo do `Tama`.

## 🛠️ Padrões e Decisões Técnicas
- **Execução Nativa:** O projeto deve ser executado nativamente via Go para garantir suporte total a TTY e interatividade do terminal (TUI).
- **Persistência:** O estado é mantido em `tamago_save.json` com versionamento de schema.
- **UX:** Implementa lógica de auto-relaunch no macOS para garantir que o jogo sempre rode em uma janela de terminal apropriada.

## 📌 Status Atual
- Base de código limpa e modularizada.
- Integração Sim-Life + RPG (Dungeon) funcional.
- Persistência robusta com cálculos de degradação offline implementados.

---
*Nota: Decisões sobre containerização (Docker) foram descartadas em favor da portabilidade nativa do binário Go e melhor experiência TUI.*
