package model

import (
	"fmt"
	"strings"
	"time"
)

type QuestType string

const (
	QuestStory      QuestType = "story"
	QuestDaily      QuestType = "daily"
	QuestRepeatable QuestType = "repeatable"
)

type QuestCondType string

const (
	CondDungeonRuns     QuestCondType = "dungeon_runs"
	CondEnemiesDefeated QuestCondType = "enemies_defeated"
	CondBossesDefeated  QuestCondType = "bosses_defeated"
	CondFeeds           QuestCondType = "feeds"
	CondWaters          QuestCondType = "waters"
	CondPets            QuestCondType = "pets"
	CondExercises       QuestCondType = "exercises"
	CondEvolveStage     QuestCondType = "evolve_stage"
)

type QuestReward struct {
	XP        int
	Gold      int
	Hunger    int
	Thirst    int
	Happiness int
}

type Quest struct {
	ID          string
	Name        string
	Description string
	Type        QuestType
	CondType    QuestCondType
	CondTarget  int
	Reward      QuestReward
}

type QuestProgress struct {
	QuestID   string    `json:"quest_id"`
	Status    string    `json:"status"` // "active" | "completed" | "claimed"
	Baseline  int       `json:"baseline"`
	StartedAt time.Time `json:"started_at"`
}

var AllQuests = map[string]Quest{
	"primeira_aventura": {
		ID: "primeira_aventura", Name: "Primeira Aventura", Type: QuestStory,
		Description: "Complete sua primeira dungeon run",
		CondType:    CondDungeonRuns, CondTarget: 1,
		Reward: QuestReward{XP: 50, Gold: 20},
	},
	"evolucao_child": {
		ID: "evolucao_child", Name: "Crescendo!", Type: QuestStory,
		Description: "Evolua para o estagio Crianca",
		CondType:    CondEvolveStage, CondTarget: int(StageChild),
		Reward: QuestReward{XP: 100, Happiness: 20},
	},
	"daily_feed": {
		ID: "daily_feed", Name: "Alimentacao do Dia", Type: QuestDaily,
		Description: "Alimente seu pet 5 vezes hoje",
		CondType:    CondFeeds, CondTarget: 5,
		Reward: QuestReward{XP: 10, Hunger: 20},
	},
	"daily_water": {
		ID: "daily_water", Name: "Hidratacao do Dia", Type: QuestDaily,
		Description: "Regue seu pet 3 vezes hoje",
		CondType:    CondWaters, CondTarget: 3,
		Reward: QuestReward{XP: 10, Thirst: 20},
	},
	"daily_pet": {
		ID: "daily_pet", Name: "Carinho do Dia", Type: QuestDaily,
		Description: "Faca carinho 3 vezes hoje",
		CondType:    CondPets, CondTarget: 3,
		Reward: QuestReward{XP: 10, Happiness: 20},
	},
	"exterminador": {
		ID: "exterminador", Name: "Exterminador", Type: QuestRepeatable,
		Description: "Derrote 10 inimigos na dungeon",
		CondType:    CondEnemiesDefeated, CondTarget: 10,
		Reward: QuestReward{XP: 30, Gold: 15},
	},
	"explorador": {
		ID: "explorador", Name: "Explorador", Type: QuestRepeatable,
		Description: "Complete 3 dungeon runs",
		CondType:    CondDungeonRuns, CondTarget: 3,
		Reward: QuestReward{XP: 40, Gold: 25},
	},
	"carinhoso_rep": {
		ID: "carinhoso_rep", Name: "Carinhoso", Type: QuestRepeatable,
		Description: "Faca carinho 10 vezes",
		CondType:    CondPets, CondTarget: 10,
		Reward: QuestReward{XP: 20, Happiness: 15},
	},
	"alimentador": {
		ID: "alimentador", Name: "Alimentador", Type: QuestRepeatable,
		Description: "Alimente seu pet 15 vezes",
		CondType:    CondFeeds, CondTarget: 15,
		Reward: QuestReward{XP: 25, Hunger: 10},
	},
	"atleta": {
		ID: "atleta", Name: "Atleta", Type: QuestRepeatable,
		Description: "Faca 5 exercicios",
		CondType:    CondExercises, CondTarget: 5,
		Reward: QuestReward{XP: 30, Gold: 25},
	},
}

var questOrder = []string{
	"primeira_aventura", "evolucao_child",
	"daily_feed", "daily_water", "daily_pet",
	"exterminador", "explorador", "carinhoso_rep", "alimentador", "atleta",
}

func currentValue(tama *Tama, condType QuestCondType) int {
	switch condType {
	case CondDungeonRuns:
		return tama.TotalDungeonRuns
	case CondEnemiesDefeated:
		return tama.TotalEnemiesDefeated
	case CondBossesDefeated:
		return tama.TotalBossesDefeated
	case CondFeeds:
		return tama.TotalFeeds
	case CondWaters:
		return tama.TotalWaters
	case CondPets:
		return tama.TotalPets
	case CondExercises:
		return tama.TotalExercises
	case CondEvolveStage:
		return int(tama.Stage)
	default:
		return 0
	}
}

// GetCurrentProgress retorna o progresso atual de uma quest (valor atual - baseline).
func GetCurrentProgress(tama *Tama, qp QuestProgress) int {
	q, ok := AllQuests[qp.QuestID]
	if !ok {
		return 0
	}
	v := currentValue(tama, q.CondType)
	prog := v - qp.Baseline
	if prog < 0 {
		prog = 0
	}
	return prog
}

// GetQuestOrder retorna os IDs das quests na ordem de exibição.
func GetQuestOrder() []string {
	return questOrder
}

// InitDefaultQuests inicializa QuestProgress com todas as quests como "active".
func InitDefaultQuests(tama *Tama) {
	existing := map[string]bool{}
	for _, qp := range tama.QuestProgress {
		existing[qp.QuestID] = true
	}
	now := time.Now()
	for _, id := range questOrder {
		if existing[id] {
			continue
		}
		q := AllQuests[id]
		tama.QuestProgress = append(tama.QuestProgress, QuestProgress{
			QuestID:   id,
			Status:    "active",
			Baseline:  currentValue(tama, q.CondType),
			StartedAt: now,
		})
	}
	if tama.LastDailyReset.IsZero() {
		tama.LastDailyReset = now
	}
}

// ResetDailyQuests reseta quests diárias se passou da meia-noite desde o último reset.
func ResetDailyQuests(tama *Tama) {
	now := time.Now()
	if tama.LastDailyReset.IsZero() {
		tama.LastDailyReset = now
		return
	}
	last := tama.LastDailyReset
	if now.Year() == last.Year() && now.Month() == last.Month() && now.Day() == last.Day() {
		return
	}
	for i := range tama.QuestProgress {
		q, ok := AllQuests[tama.QuestProgress[i].QuestID]
		if !ok || q.Type != QuestDaily {
			continue
		}
		tama.QuestProgress[i].Status = "active"
		tama.QuestProgress[i].Baseline = currentValue(tama, q.CondType)
		tama.QuestProgress[i].StartedAt = now
	}
	tama.LastDailyReset = now
}

// UpdateQuestProgress verifica progresso de todas as quests ativas.
// Retorna nomes das quests recém-completadas.
func UpdateQuestProgress(tama *Tama) []string {
	completed := []string{}
	for i := range tama.QuestProgress {
		qp := &tama.QuestProgress[i]
		if qp.Status != "active" {
			continue
		}
		q, ok := AllQuests[qp.QuestID]
		if !ok {
			continue
		}
		if GetCurrentProgress(tama, *qp) >= q.CondTarget {
			qp.Status = "completed"
			completed = append(completed, q.Name)
		}
	}
	return completed
}

// ClaimReward reivindica a recompensa de uma quest completada.
// addGold é um callback para adicionar gold ao inventário da dungeon (evita dependência circular).
func ClaimReward(tama *Tama, questID string, addGold func(int)) (string, bool) {
	for i := range tama.QuestProgress {
		qp := &tama.QuestProgress[i]
		if qp.QuestID != questID || qp.Status != "completed" {
			continue
		}
		q, ok := AllQuests[questID]
		if !ok {
			return "[!] Missao nao encontrada.", false
		}

		if q.Reward.XP > 0 {
			tama.AddXP(q.Reward.XP)
		}
		if q.Reward.Hunger > 0 {
			tama.Hunger = min(tama.Hunger+q.Reward.Hunger, MaxHunger)
		}
		if q.Reward.Thirst > 0 {
			tama.Thirst = min(tama.Thirst+q.Reward.Thirst, MaxThirst)
		}
		if q.Reward.Happiness > 0 {
			tama.Happiness = min(tama.Happiness+q.Reward.Happiness, MaxHappiness)
		}
		if q.Reward.Gold > 0 && addGold != nil {
			addGold(q.Reward.Gold)
		}
		tama.TotalQuestsDone++

		if q.Type == QuestRepeatable {
			qp.Status = "active"
			qp.Baseline = currentValue(tama, q.CondType)
			qp.StartedAt = time.Now()
		} else {
			qp.Status = "claimed"
		}

		return fmt.Sprintf("[MISSAO] %s concluida! %s", q.Name, buildRewardMsg(q.Reward)), true
	}
	return "[!] Missao nao disponivel para reivindicar.", false
}

func buildRewardMsg(r QuestReward) string {
	parts := []string{}
	if r.XP > 0 {
		parts = append(parts, fmt.Sprintf("+%d XP", r.XP))
	}
	if r.Gold > 0 {
		parts = append(parts, fmt.Sprintf("+%d Ouro", r.Gold))
	}
	if r.Hunger > 0 {
		parts = append(parts, fmt.Sprintf("+%d Fome", r.Hunger))
	}
	if r.Thirst > 0 {
		parts = append(parts, fmt.Sprintf("+%d Sede", r.Thirst))
	}
	if r.Happiness > 0 {
		parts = append(parts, fmt.Sprintf("+%d Felicidade", r.Happiness))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}
