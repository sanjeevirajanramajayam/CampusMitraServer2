package models
type PsLevels struct {
	RollNo       string `json:"rollno"`
	SkillDomain  string `json:"skilldomain"`
	SkillName    string `json:"skillname"`
	SkillLevel   string `json:"skilllevel"`
	TotalLevels  int    `json:"totallevels"`
}
type Ps struct {
		RollNo      string `json:"rollno"`
		Points      float64    `json:"points"`
		SkillDomain string `json:"skilldomain"`
		SkillName 	string `json:"skillname"`
		SkillLevel 	string `json:"skilllevel"`
		Attempts  int `json:"attempts"`
		Sem int `json:"sem"`
		Currdate string `json:"currdate"`
}
type Points_Logs struct {
		RollNo      string `json:"rollno"`
		Source      string `json:"source"`
		Points      float64    `json:"points"`
		Description string `json:"description"`
		Sem int `json:"sem"`
		Currdate string `json:"currdate"`
}
type SemCount struct {
	Sem       int `json:"sem"`
	SemCount  int `json:"sem_count"`
}