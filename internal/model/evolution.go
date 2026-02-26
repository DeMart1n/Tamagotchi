package model

type Stage int

const (
	StageBaby Stage = iota
	StageChild
	StageTeen
	StageAdult
	StageElder
)

func (s Stage) String() string {
	switch s {
	case StageBaby:
		return "Baby"
	case StageChild:
		return "Criança"
	case StageTeen:
		return "Adolescente"
	case StageAdult:
		return "Adulto"
	case StageElder:
		return "Ancião"
	default:
		return "???"
	}
}

func (s Stage) Avatar(frame int) string {
	switch s {
	case StageBaby:
		frames := []string{"(°◡°)", "(°o°)", "(°◡°)"}
		return frames[frame%len(frames)]
	case StageChild:
		frames := []string{"(・_・)", "(・o・)", "(・_・)"}
		return frames[frame%len(frames)]
	case StageTeen:
		frames := []string{"(¬‿¬)", "(¬_¬)", "(¬‿¬)"}
		return frames[frame%len(frames)]
	case StageAdult:
		frames := []string{"(◕‿◕)", "(◕_◕)", "(◕‿◕)"}
		return frames[frame%len(frames)]
	case StageElder:
		frames := []string{"(ᵔ‿ᵔ)", "(ᵔ_ᵔ)", "(ᵔ‿ᵔ)"}
		return frames[frame%len(frames)]
	default:
		return "(・_・)"
	}
}

// DecayRate retorna o multiplicador de decay por estágio
// Baby decai menos, Elder decai mais
func (s Stage) DecayRate() int {
	switch s {
	case StageBaby:
		return 1
	case StageChild:
		return 1
	case StageTeen:
		return 1
	case StageAdult:
		return 2
	case StageElder:
		return 2
	default:
		return 1
	}
}

// Thresholds de XP para level up
var LevelThresholds = []int{
	100, 250, 500, 800, 1200, 1700, 2300, 3000, 3800, 4700,
	5700, 6800, 8000, 9300, 10700, 12200, 13800, 15500, 17300, 19200,
	21200, 23300, 25500, 27800, 30200,
}

func XPForNextLevel(level int) int {
	if level < len(LevelThresholds) {
		return LevelThresholds[level]
	}
	return LevelThresholds[len(LevelThresholds)-1] + (level-len(LevelThresholds)+1)*3000
}

func StageForLevel(level int) Stage {
	switch {
	case level >= 25:
		return StageElder
	case level >= 15:
		return StageAdult
	case level >= 8:
		return StageTeen
	case level >= 3:
		return StageChild
	default:
		return StageBaby
	}
}
