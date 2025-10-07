package facultymodel
import "time"

type ActivityApproval struct {
	EventName        string     `json:"event_name"`
	EventCode       string     `json:"event_code"`
	EventType        string     `json:"event_type"`
	EventStartDate   *time.Time `json:"event_start_date"`
	Domain           *string    `json:"domain"`
	ProblemStatement *string    `json:"problem_statement"`
	ApplicantName    *string    `json:"applicant_name"`
	Rollno           *string    `json:"rollno"`
	Verified         *string    `json:"verified"`
	SubmittedDate    *time.Time `json:"submitted_date"`
}