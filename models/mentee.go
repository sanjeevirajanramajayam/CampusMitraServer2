package models       
type Mentee struct {
    MentorRollNo        string `json:"mentor_rollno"`
    MenteeRollNo        string `json:"mentee_rollno"`      
    SkillName           string `json:"skill_name"`       
}