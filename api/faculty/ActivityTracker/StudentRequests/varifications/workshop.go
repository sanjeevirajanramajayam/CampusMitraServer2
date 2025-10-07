package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetWorkshops() ([]facultymodel.Varification, error) {
	var projects []facultymodel.Varification

	getQuery := `
		SELECT 
    w.id,
    w.title,
    w.event_type,
    w.mode_of_delivery,
    w.event_nature,
    w.organised_by,
    w.location,
    w.start_date,
    w.end_date,
    w.participation_type,
    w.is_certificate,
    w.certificate,
    w.link,
    w.topics_covered,
    w.relevence,
    w.skills_gained,
    w.status,
    w.upload_type,
    l.user_name
FROM workshops AS w
INNER JOIN login AS l 
    ON l.rollno = w.rollno;
	
	`

	rows, err := config.DB.Query(getQuery)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data for proejcts for varifications")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r facultymodel.Varification
		err = rows.Scan(&r.Id, &r.Title, &r.Event_type, &r.ModeOfDelivary, &r.EventNature, &r.OrganisedBy, &r.Location, &r.StartDate, &r.EndDate, &r.Participation_type, &r.IsCertificate, &r.Certificate_pdf, &r.Link, &r.TopicsCovered, &r.Relevence, &r.SkillsGained, &r.Approval_status,&r.UploadType, &r.User_name)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			fmt.Println("Could not scan the data for projects varification")
			return nil, err
		}
		projects = append(projects, r)
	}
	return projects, nil
}