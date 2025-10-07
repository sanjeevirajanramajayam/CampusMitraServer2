package activitymastermodels

type SurveyDetail struct {
    ActivityID         int    `json:"activity_id"`
    PublishingDept     string `json:"publishing_department"`
    Description        string `json:"description"`
    StartDate          string `json:"start_date"`
    EndDate            string `json:"end_date"`
    LinkOrLocation     string `json:"link_or_location"`
    TargetYear         string `json:"target_year"`
    TargetDepartment   string `json:"target_department"`
    AllStudents        int    `json:"all_students"`
}

type SessionDetails struct {
	ActivityID          int    `json:"activity_id"`
	PublishingDept      string `json:"publishing_department"`
	Host                string `json:"host"`
	Description         string `json:"description"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	DateOfSession       string `json:"date_of_session"`
	LinkOrLocation      string `json:"link_or_location"`
}

type EligibleMeeting struct {
    ActivityID          int    `json:"activity_id"`
    PublishingDept      string `json:"publishing_department"`
    Host                string `json:"host"`
    Description         string `json:"description"`
    StartTime           string `json:"start_time"`
    EndTime             string `json:"end_time"`
    DateOfMeeting       string `json:"date_of_meeting"`
    LinkOrLocation      string `json:"link_or_location"`
    TargetYear          string `json:"target_year"`
    TargetDepartment    string `json:"target_department"`
    AllStudents         bool   `json:"all_students"`
}