package model

import (
	"encoding/json"
	"time"
)

const (
	MaxHunger = 100
	MinHunger = 0

	MaxThirst = 100
	MinThirst = 0

	MaxSleepy = 100
	MinSleepy = 0

	MaxHappiness = 100
	MinHappiness = 0

	MaxAngry = 100
	MinAngry = 0

	MaxWeight = 100
	MinWeight = 0

	Underweight     = 20
	NormalWeightMin = 30
	NormalWeightMax = 70
	Overweight      = 80
)

type Tama struct {
	Name        string    `json:"name"`
	Hunger      int       `json:"hunger"`
	Thirst      int       `json:"thirst"`
	Sleepy      int       `json:"sleepy"`
	Happiness   int       `json:"happiness"`
	Angry       int       `json:"angry"`
	Weight      int       `json:"weight"`
	Sleeping    bool      `json:"sleeping"`
	Dead        bool      `json:"dead"`
	Depressed   bool      `json:"depressed"`
	PissedOf    bool      `json:"pissed_of"`
	Overweight  bool      `json:"overweight"`
	Underweight bool      `json:"underweight"`
	LastSaved   time.Time `json:"last_saved"`

	// XP e Evolução
	XP    int   `json:"xp"`
	Level int   `json:"level"`
	Stage Stage `json:"stage"`

	// Contadores de ação (para conquistas)
	TotalFeeds     int `json:"total_feeds"`
	TotalWaters    int `json:"total_waters"`
	TotalPets      int `json:"total_pets"`
	TotalSleeps    int `json:"total_sleeps"`
	TotalExercises int `json:"total_exercises"`
	TotalAnnoys    int `json:"total_annoys"`
	TotalTicks     int `json:"total_ticks"`

	// Conquistas
	Achievements []Achievement `json:"achievements"`

	// Dungeon RPG
	Inventory           json.RawMessage `json:"inventory,omitempty"`
	TotalDungeonRuns    int             `json:"total_dungeon_runs"`
	TotalBossesDefeated int             `json:"total_bosses_defeated"`
	CurrentBiome        string          `json:"current_biome"`

	// Habilidades
	MP          int      `json:"mp"`
	MPMax       int      `json:"mp_max"`
	SkillsKnown []string `json:"skills_known"`

	// Missões
	QuestProgress        []QuestProgress `json:"quest_progress,omitempty"`
	TotalQuestsDone      int             `json:"total_quests_done"`
	LastDailyReset       time.Time       `json:"last_daily_reset"`
	TotalEnemiesDefeated int             `json:"total_enemies_defeated"`
}

// AddXP adiciona XP e faz level up se necessário. Retorna true se houve level up.
func (t *Tama) AddXP(amount int) bool {
	t.XP += amount
	leveled := false
	for t.XP >= XPForNextLevel(t.Level) {
		t.XP -= XPForNextLevel(t.Level)
		t.Level++
		leveled = true
	}
	t.Stage = StageForLevel(t.Level)
	return leveled
}
