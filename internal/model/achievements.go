package model

type Achievement struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Unlocked    bool   `json:"unlocked"`
}

func DefaultAchievements() []Achievement {
	return []Achievement{
		{ID: "first_feed", Name: "Primeira Refeicao", Description: "Alimentou o Tama pela primeira vez", Icon: "[~]"},
		{ID: "first_pet", Name: "Primeiro Carinho", Description: "Fez carinho no Tama pela primeira vez", Icon: "<3>"},
		{ID: "first_exercise", Name: "Primeiro Exercicio", Description: "Exercitou o Tama pela primeira vez", Icon: "[!]"},
		{ID: "night_owl", Name: "Coruja Noturna", Description: "Sobreviveu 50 ticks", Icon: "[*]"},
		{ID: "centurion", Name: "Centuriao", Description: "Sobreviveu 100 ticks", Icon: "[C]"},
		{ID: "zen_master", Name: "Mestre Zen", Description: "Todos os stats acima de 70", Icon: "[Z]"},
		{ID: "marathon", Name: "Maratonista", Description: "Exercitou 20 vezes", Icon: "[M]"},
		{ID: "glutton", Name: "Glutao", Description: "Alimentou 30 vezes", Icon: "[G]"},
		{ID: "hydrated", Name: "Hidratado", Description: "Deu agua 30 vezes", Icon: "[o]"},
		{ID: "cuddler", Name: "Carinhoso", Description: "Fez carinho 25 vezes", Icon: "<3>"},
		{ID: "troublemaker", Name: "Encrenqueiro", Description: "Irritou o Tama 15 vezes", Icon: "[>]"},
		{ID: "evolved_child", Name: "Crescendo!", Description: "Evoluiu para Crianca", Icon: "[i]"},
		{ID: "evolved_teen", Name: "Adolescente", Description: "Evoluiu para Adolescente", Icon: "[I]"},
		{ID: "evolved_adult", Name: "Adulto", Description: "Evoluiu para Adulto", Icon: "[A]"},
		{ID: "evolved_elder", Name: "Anciao Sabio", Description: "Evoluiu para Anciao", Icon: "[S]"},
		{ID: "level_10", Name: "Nivel 10", Description: "Atingiu nivel 10", Icon: "[*]"},
		{ID: "first_dungeon", Name: "Aventureiro", Description: "Completou uma dungeon run", Icon: "[D]"},
		{ID: "dragon_slayer", Name: "Mata-Dragao", Description: "Derrotou o Dragao Anciao", Icon: "[X]"},
		{ID: "dungeon_master", Name: "Mestre da Masmorra", Description: "Completou 5 dungeon runs", Icon: "[#]"},
	}
}

func CheckAchievements(tama *Tama) []string {
	unlocked := []string{}

	checks := map[string]bool{
		"first_feed":     tama.TotalFeeds >= 1,
		"first_pet":      tama.TotalPets >= 1,
		"first_exercise": tama.TotalExercises >= 1,
		"night_owl":      tama.TotalTicks >= 50,
		"centurion":      tama.TotalTicks >= 100,
		"zen_master":     tama.Hunger > 70 && tama.Thirst > 70 && tama.Sleepy > 70 && tama.Happiness > 70 && tama.Angry < 30,
		"marathon":       tama.TotalExercises >= 20,
		"glutton":        tama.TotalFeeds >= 30,
		"hydrated":       tama.TotalWaters >= 30,
		"cuddler":        tama.TotalPets >= 25,
		"troublemaker":   tama.TotalAnnoys >= 15,
		"evolved_child":  tama.Stage >= StageChild,
		"evolved_teen":   tama.Stage >= StageTeen,
		"evolved_adult":  tama.Stage >= StageAdult,
		"evolved_elder":  tama.Stage >= StageElder,
		"level_10":       tama.Level >= 10,
		"first_dungeon":  tama.TotalDungeonRuns >= 1,
		"dragon_slayer":  tama.TotalBossesDefeated >= 1,
		"dungeon_master": tama.TotalDungeonRuns >= 5,
	}

	for i := range tama.Achievements {
		if tama.Achievements[i].Unlocked {
			continue
		}
		if checks[tama.Achievements[i].ID] {
			tama.Achievements[i].Unlocked = true
			unlocked = append(unlocked, tama.Achievements[i].Icon+" "+tama.Achievements[i].Name)
		}
	}

	return unlocked
}
