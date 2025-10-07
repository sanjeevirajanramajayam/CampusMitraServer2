// models/models.go
package projectModal

// Project represents the main project details in the 'Projects' table.
type TechStack struct {
	TechName string `json:"tech_name"`
}

type TeamMember struct {
	Name       string   `json:"name"`
	RollNumber string   `json:"rollNumber"`
	Department string   `json:"department"`
	TechStack  []string `json:"techStack"`
}
