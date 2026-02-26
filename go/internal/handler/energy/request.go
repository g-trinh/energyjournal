package energy

type SaveEnergyLevelsRequest struct {
	Date               string `json:"date"`
	Physical           int    `json:"physical"`
	Mental             int    `json:"mental"`
	Emotional          int    `json:"emotional"`
	SleepQuality       int    `json:"sleepQuality,omitempty"`
	StressLevel        int    `json:"stressLevel,omitempty"`
	PhysicalActivity   string `json:"physicalActivity,omitempty"`
	Nutrition          string `json:"nutrition,omitempty"`
	SocialInteractions string `json:"socialInteractions,omitempty"`
	TimeOutdoors       string `json:"timeOutdoors,omitempty"`
	Notes              string `json:"notes,omitempty"`
}
