package energy

type SaveEnergyLevelsRequest struct {
	Date      string `json:"date"`
	Physical  int    `json:"physical"`
	Mental    int    `json:"mental"`
	Emotional int    `json:"emotional"`
}
