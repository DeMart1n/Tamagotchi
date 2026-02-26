package model

type EventType int

const (
	EventPositive EventType = iota
	EventNegative
	EventInteractive
)

type Event struct {
	Name        string
	Description string
	Type        EventType
	// Efeitos diretos (aplicados automaticamente)
	HungerDelta    int
	ThirstDelta    int
	HappinessDelta int
	AngryDelta     int
	SleepyDelta    int
	XPDelta        int
	// Para eventos interativos
	Choice1Label  string
	Choice2Label  string
	Choice1Result EventResult
	Choice2Result EventResult
}

type EventResult struct {
	Message        string
	HungerDelta    int
	ThirstDelta    int
	HappinessDelta int
	AngryDelta     int
	XPDelta        int
}

var Events = []Event{
	// Positivos
	{
		Name:        "Comida encontrada!",
		Description: "Achou comida no chão!",
		Type:        EventPositive,
		HungerDelta: 10,
		XPDelta:     3,
	},
	{
		Name:           "Borboleta!",
		Description:    "Uma borboleta pousou no nariz!",
		Type:           EventPositive,
		HappinessDelta: 20,
		XPDelta:        3,
	},
	{
		Name:        "Chuva refrescante",
		Description: "Uma chuva leve caiu! Refrescante!",
		Type:        EventPositive,
		ThirstDelta: 15,
		XPDelta:     2,
	},
	{
		Name:           "Brisa agradável",
		Description:    "Uma brisa agradável passou!",
		Type:           EventPositive,
		HappinessDelta: 10,
		AngryDelta:     -10,
		XPDelta:        2,
	},
	// Negativos
	{
		Name:           "Tempestade!",
		Description:    "Uma tempestade assustou o Tama!",
		Type:           EventNegative,
		HappinessDelta: -15,
		AngryDelta:     5,
	},
	{
		Name:           "Tropeçou!",
		Description:    "O Tama tropeçou e caiu!",
		Type:           EventNegative,
		HappinessDelta: -5,
		AngryDelta:     10,
	},
	{
		Name:        "Barulho alto!",
		Description: "Um barulho alto assustou o Tama!",
		Type:        EventNegative,
		AngryDelta:  15,
		SleepyDelta: -10,
	},
	{
		Name:           "Calor intenso",
		Description:    "O sol está muito forte!",
		Type:           EventNegative,
		ThirstDelta:    -10,
		HappinessDelta: -5,
	},
	// Interativos
	{
		Name:         "Um gato apareceu!",
		Description:  "Um gato apareceu! O que fazer?",
		Type:         EventInteractive,
		Choice1Label: "Fazer amizade",
		Choice2Label: "Ignorar",
		Choice1Result: EventResult{
			Message:        "Fizeram amizade! O gato ronronou.",
			HappinessDelta: 20,
			XPDelta:        10,
		},
		Choice2Result: EventResult{
			Message:        "O gato foi embora...",
			HappinessDelta: -5,
		},
	},
	{
		Name:         "Fruta na árvore!",
		Description:  "Tem uma fruta no alto da árvore! O que fazer?",
		Type:         EventInteractive,
		Choice1Label: "Escalar",
		Choice2Label: "Deixar pra lá",
		Choice1Result: EventResult{
			Message:     "Conseguiu a fruta! Deliciosa!",
			HungerDelta: 20,
			XPDelta:     8,
		},
		Choice2Result: EventResult{
			Message: "A fruta caiu sozinha depois... Que pena.",
		},
	},
	{
		Name:         "Outro Tama apareceu!",
		Description:  "Outro Tamagotchi apareceu! O que fazer?",
		Type:         EventInteractive,
		Choice1Label: "Brincar junto",
		Choice2Label: "Brigar",
		Choice1Result: EventResult{
			Message:        "Brincaram juntos! Foi muito divertido!",
			HappinessDelta: 25,
			AngryDelta:     -10,
			XPDelta:        12,
		},
		Choice2Result: EventResult{
			Message:        "A briga foi feia...",
			AngryDelta:     20,
			HappinessDelta: -10,
		},
	},
}
