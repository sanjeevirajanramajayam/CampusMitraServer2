package uploadview

type UploadView struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Complexity  *string `json:"complexity,omitempty"`
	Status      string  `json:"status"`
	UploadedOn  string  `json:"uploaded_on"`
	RollNo      string  `json:"rollno"`
	Subtype     *string  `json:"SUb-type,omitempty"`
}
