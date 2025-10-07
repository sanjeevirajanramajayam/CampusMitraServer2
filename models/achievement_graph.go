package models

type Institute_avg struct{
	Cummulative_points float64 `json:"cummulative_points"`
	Points float64 `json:"points"`
	Sem int `json:"sem"`
	Currdate string `json:"currdate"`
}

type Point_logs2 struct{
	Rollno string `json:"rollno"`
	Points float64 `json:"cummulative_points"`
	Sem int `json:"sem"`
	Currdate string `json:"currdate"`
}

type Achievementgraph struct{
	Rollno string `json:"rollno"`
	Cummulative_points float64 `json:"cummulative_points"`
	Points_earned float64 `json:"points_earned"`
	Sem int `json:"sem"`
	Currdate string `json:"currdate"`
 }