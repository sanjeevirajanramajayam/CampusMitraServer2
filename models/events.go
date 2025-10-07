package models

type RoundData struct {
	StartDate    string            `json:"start_date"`
	EndDate      string            `json:"end_date"`
	RewardPoints map[string]string `json:"reward_points"`
}

type Event struct {
	EventName     string
	Type          string
	Deadline      string
	MinTeamSize   int
	MaxTeamSize   int
	NoOfRounds    int
	OnlineRounds  string
	OfflineRounds string
	Location      string
	ApplyLink     string
	Domains       string
	Description   string
	Rules         string
	Constraints   string
	ImageURL      string
	FinalPrizes   map[string]string
	Rounds        []RoundData
}
