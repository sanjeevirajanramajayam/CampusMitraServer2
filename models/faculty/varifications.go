package facultymodel

type 	Varification struct {
	Id                 int    `json:"id,omitempty"`
	UploadType         string `json:"upload_type,omitempty"`
	PaperTitle         string `json:"paper_title,omitempty"`
	ConferenceTitle    string `json:"conference_title,omitempty"`
	Location           string `json:"location,omitempty"`
	DateOfPresentation string `json:"date_of_presentation,omitempty"`
	Pdf                string `json:"pdf,omitempty"`
	Certificate        string `json:"certificate,omitempty"`
	Award              string `json:"award,omitempty"`
	Approval_status    string `json:"approval_status,omitempty"`
	User_name          string `json:"user_name,omitempty"`
	// Internship-specific fields
	CompanyName           string `json:"company_name,omitempty"`
	Roll                  string `json:"roll,omitempty"`
	Domain                string `json:"domain,omitempty"`
	InternshipType        string `json:"internship_type,omitempty"`
	IsStipend             int    `json:"is_stipend,omitempty"`
	ConsultedFacultyName  string `json:"consulted_faculty_name,omitempty"`
	IndustryMentorName    string `json:"industry_mentor_name,omitempty"`
	IndustryMentorContact string `json:"industry_mentor_contact,omitempty"`
	OfferLetter           string `json:"offer_letter,omitempty"`
	Report                string `json:"report,omitempty"`
	FacultyRemarks        string `json:"faculty_remarks,omitempty"`
	SkillGained           string `json:"skill_gained,omitempty"`
	Outcomes              string `json:"outcomes,omitempty"`

	// Patent-specific fields
	PatentTitle         string `json:"title,omitempty"`
	ApplicationNumber   string `json:"application_number,omitempty"`
	DateOfFiling        string `json:"date_of_filing,omitempty"`
	PatentDocs          string `json:"patent_docs,omitempty"`
	SupportingFiles     string `json:"supporting_files,omitempty"`
	LinkToPatentListing string `json:"link_to_patent_listing,omitempty"`
	UsecaseOfPatent     string `json:"usecase_of_patent,omitempty"`
	PatentStatus        string `json:"patent_status,omitempty"`

	// Other existing fields...
	EventName          string `json:"event_name,omitempty"`
	EventCode          string `json:"event_code,omitempty"`
	Participation_type string `json:"participation_type,omitempty"`
	Summary            string `json:"summary,omitempty"`
	WinningStatus      string `json:"winning_status,omitempty"`
	Title              string `json:"title,omitempty"`
	Platform           string `json:"platform,omitempty"`
	Issue_date         string `json:"issue_date,omitempty"`
	Course_link        string `json:"course_link,omitempty"`
	Certificate_pdf    string `json:"certificate_pdf,omitempty"`
	Certificate_type   string `json:"certificate_type,omitempty"`
	Certificate_id     int    `json:"certificate_id,omitempty"`
	Duration           string `json:"duration,omitempty"`
	Activity_type      string `json:"activity_type,omitempty"`
	Event_type         string `json:"event_type,omitempty"`
	ModeOfDelivary     string `json:"mode_of_delivary,omitempty"`
	EventNature        string `json:"event_nature,omitempty"`
	OrganisedBy        string `json:"organised_by,omitempty"`
	StartDate          string `json:"start_date,omitempty"`
	EndDate            string `json:"end_date,omitempty"`
	IsCertificate      int    `json:"is_certificate,omitempty"`
	Link               string `json:"link,omitempty"`
	TopicsCovered      string `json:"topic_covered,omitempty"`
	Relevence          string `json:"relevence,omitempty"`
	SkillsGained       string `json:"skills_gained,omitempty"`
	ProblemStatement   string `json:"problem_statement,omitempty"`
	Objective          string `json:"objective,omitempty"`

	// Fields from projects
	TitleIdea           string `json:"title_idea,omitempty"`
	StartTime           string `json:"start_time,omitempty"`
	EndTime             string `json:"end_time,omitempty"`
	IsTeamProject       int    `json:"is_team_project,omitempty"`
	ConsultedMentor     string `json:"consulted_mentor,omitempty"`
	ChangesFromIdea     string `json:"changes_from_idea,omitempty"`
	GithubLink          string `json:"github_link,omitempty"`
	ReportPdf           string `json:"report_pdf,omitempty"`
	DemoVideo           string `json:"demo_video,omitempty"`
	PresentedExternally string `json:"presented_externally,omitempty"`
	AwardsWon           string `json:"awards_won,omitempty"`
	Rollno              string `json:"rollno,omitempty"`
	MemberName          string `json:"member_name,omitempty"`
	Department          string `json:"department,omitempty"`
	TechNames           string `json:"tech_names,omitempty"`
}