package model

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
)

type Tama struct {
	Name      string
	Hunger    int
	Thirst    int
	Sleepy    int
	Happiness int
	Angry     int
	Sleeping  bool
	Dead      bool
	Depressed bool
	PissedOf  bool
}
