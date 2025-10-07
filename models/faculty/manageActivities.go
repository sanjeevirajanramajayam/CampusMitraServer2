package facultymodel

type ManageActivities struct {
	Activity_title  string `json:"activity_title"`
	Activity_type   string `json:"activity_type"`
	Description     string `json:"description"`
	Start_date      string `json:"start_date"`
	End_date        string `json:"end_date"`
	Linkorlocation  string `json:"linkorlocation"`
	All_students    string `json:"all_students"`
	Specific_rollno string `json:"specific_rollno"`
	Year_type       string `json:"year_type"`
	Target_dept     string `json:"target_dept"`
}
type Event struct {
    ID            int64          `json:"id"`
    EventName     string         `json:"event_name"`
    EventCode     string         `json:"event_code"`
    Type          string         `json:"type"`
    Deadline      string         `json:"deadline"`
    MinTeamSize   int            `json:"min_team_size"`
    MaxTeamSize   int            `json:"max_team_size"`
    NoOfRounds    int            `json:"no_of_rounds"`
    OnlineRounds  int  `json:"online_rounds"`
    OfflineRounds int `json:"offline_rounds"`
    Location       string `json:"location"`
    ApplyLink      string `json:"apply_link"`
    Domains        string `json:"domains"`
    ImageURL       string `json:"image_url"`
    Description    string `json:"description"`
    Rules          string `json:"rules"`
    Constraints    string `json:"constraints"`
    FinalPrize1    string `json:"final_prize1"`
    FinalPrize2    string `json:"final_prize2"`
    FinalPrize3    string `json:"final_prize3"`

    // Add rounds slice
    Rounds []Rounds `json:"rounds"`
}

type Rounds struct {
    RoundNumber int    `json:"round_number"`
	Description string `json:"description"`
    StartDate   string `json:"start_date"`
    EndDate     string `json:"end_date"`
    Year1RP     string `json:"year1_rp"`
    Year2RP     string `json:"year2_rp"`
    Year3RP     string `json:"year3_rp"`
    Year4RP     string `json:"year4_rp"`
}


type Round struct {
	RoundNumber  int    `json:"round_no" `
	Description  string  `json:"description"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Rewardpoints Rewardpoints `json:"reward_points"`
}

type Rewardpoints struct{
	Year1        string	`json:"year1"`
	Year2        string	`json:"year2"`
	Year3        string	`json:"year3"`
	Year4        string	`json:"year4"`
}
type EventRoundDates struct {
	EventCode   string `json:"event_code"`
	RoundNumber int    `json:"round_number"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Year1RP     string `json:"year1_rp"`
	Year2RP     string `json:"year2_rp"`
	Year3RP     string `json:"year3_rp"`
	Year4RP     string `json:"year4_rp"`
}
type Activity struct {
	ActivityTitle  string `json:"activity_title"`
	ActivityType   string `json:"activity_type"`
	Description    string `json:"description"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	DateofMeeting  string `json:"date_of_meeting,omitempty"`
	LinkOrLocation string `json:"link_or_location"`
	TargetYear     string `json:"TargetYear,omitempty"`
	AllStudents    int    `json:"all_students,omitempty"`
	Host           string `json:"host,omitempty"`
	SpecificRollno string `json:"specific_rollno,omitempty"`
}
type RegisteredStudent struct {
	RollNo      string `json:"rollno"`
	EventCode   string `json:"event_code"`
}