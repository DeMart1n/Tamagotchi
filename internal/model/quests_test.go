package model

import (
	"testing"
	"time"
)

func TestInitDefaultQuests(t *testing.T) {
	tama := &Tama{
		Name:  "TestPet",
		Level: 1,
		Stage: StageBaby,
	}

	InitDefaultQuests(tama)

	if len(tama.QuestProgress) != len(questOrder) {
		t.Errorf("Expected %d quests, got %d", len(questOrder), len(tama.QuestProgress))
	}

	for _, qp := range tama.QuestProgress {
		if qp.Status != "active" {
			t.Errorf("Quest %s should start as active, got %s", qp.QuestID, qp.Status)
		}
		if qp.Baseline != currentValue(tama, AllQuests[qp.QuestID].CondType) {
			t.Errorf("Quest %s baseline not set correctly", qp.QuestID)
		}
	}
}

func TestUpdateQuestProgress_DailyCompletion(t *testing.T) {
	tama := &Tama{
		Name:       "TestPet",
		Level:      1,
		Stage:      StageBaby,
		TotalFeeds: 0,
	}

	InitDefaultQuests(tama)

	// Simular 5 alimentações
	tama.TotalFeeds = 5

	completed := UpdateQuestProgress(tama)

	found := false
	for _, name := range completed {
		if name == "Alimentacao do Dia" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 'Alimentacao do Dia' to be completed after 5 feeds")
	}

	// Verificar status da quest
	for _, qp := range tama.QuestProgress {
		if qp.QuestID == "daily_feed" && qp.Status != "completed" {
			t.Errorf("Expected daily_feed to be completed, got %s", qp.Status)
		}
	}
}

func TestResetDailyQuests_AfterMidnight(t *testing.T) {
	tama := &Tama{
		Name:  "TestPet",
		Level: 1,
		Stage: StageBaby,
	}

	InitDefaultQuests(tama)

	// Simular que foi reseted ontem
	yesterday := time.Now().AddDate(0, 0, -1)
	tama.LastDailyReset = yesterday

	// Marcar uma quest diária como completada
	for i := range tama.QuestProgress {
		if tama.QuestProgress[i].QuestID == "daily_feed" {
			tama.QuestProgress[i].Status = "completed"
			break
		}
	}

	ResetDailyQuests(tama)

	// Verificar se a quest foi resetada
	for _, qp := range tama.QuestProgress {
		if qp.QuestID == "daily_feed" {
			if qp.Status != "active" {
				t.Errorf("Expected daily_feed to reset to active, got %s", qp.Status)
			}
			break
		}
	}
}

func TestClaimReward_AppliesStats(t *testing.T) {
	tama := &Tama{
		Name:      "TestPet",
		Level:     1,
		Stage:     StageBaby,
		XP:        0,
		Hunger:    50,
		Happiness: 50,
	}

	InitDefaultQuests(tama)

	// Marcar uma quest como completa
	for i := range tama.QuestProgress {
		if tama.QuestProgress[i].QuestID == "daily_feed" {
			tama.QuestProgress[i].Status = "completed"
			break
		}
	}

	initialXP := tama.XP
	initialHunger := tama.Hunger

	msg, ok := ClaimReward(tama, "daily_feed", func(g int) {
		// Gold would be added to inventory here
	})

	if !ok {
		t.Errorf("Expected claim to succeed, got error: %s", msg)
	}

	if tama.XP <= initialXP {
		t.Errorf("Expected XP to increase, but got same or less")
	}

	if tama.Hunger <= initialHunger {
		t.Errorf("Expected Hunger to increase, but got same or less")
	}

	if tama.TotalQuestsDone != 1 {
		t.Errorf("Expected TotalQuestsDone=1, got %d", tama.TotalQuestsDone)
	}
}

func TestRepeatableQuestReset(t *testing.T) {
	tama := &Tama{
		Name:       "TestPet",
		Level:      5,
		Stage:      StageChild,
		TotalPets:  0,
	}

	InitDefaultQuests(tama)

	// Marcar quest repeatável como completa
	for i := range tama.QuestProgress {
		if tama.QuestProgress[i].QuestID == "carinhoso_rep" {
			tama.QuestProgress[i].Status = "completed"
			break
		}
	}

	// Reivindicar recompensa
	_, ok := ClaimReward(tama, "carinhoso_rep", nil)
	if !ok {
		t.Errorf("Expected claim to succeed")
	}

	// Verificar se resetou para active
	for _, qp := range tama.QuestProgress {
		if qp.QuestID == "carinhoso_rep" {
			if qp.Status != "active" {
				t.Errorf("Expected repeatable quest to reset to active, got %s", qp.Status)
			}
			// Baseline deve ser reset ao valor atual
			currentVal := currentValue(tama, AllQuests[qp.QuestID].CondType)
			if qp.Baseline != currentVal {
				t.Errorf("Expected baseline to be reset to current value %d, got %d", currentVal, qp.Baseline)
			}
			break
		}
	}
}

func TestGetCurrentProgress(t *testing.T) {
	tama := &Tama{
		Name:       "TestPet",
		Level:      1,
		TotalFeeds: 10,
	}

	InitDefaultQuests(tama)

	// Find daily_feed quest
	var dailyFeedQP *QuestProgress
	for i := range tama.QuestProgress {
		if tama.QuestProgress[i].QuestID == "daily_feed" {
			dailyFeedQP = &tama.QuestProgress[i]
			break
		}
	}

	if dailyFeedQP == nil {
		t.Fatalf("daily_feed quest not found")
	}

	// Set baseline to 0 (starting point)
	dailyFeedQP.Baseline = 0

	// Progress should be current value (10) - baseline (0) = 10
	prog := GetCurrentProgress(tama, *dailyFeedQP)
	if prog != 10 {
		t.Errorf("Expected progress=10, got %d", prog)
	}
}
