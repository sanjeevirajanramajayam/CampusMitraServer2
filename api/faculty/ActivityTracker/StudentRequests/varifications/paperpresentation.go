package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetPaperpresentation() ([]facultymodel.Varification, error) {
	var paperPresentation []facultymodel.Varification

	query := `
		SELECT 
			p.id,
			p.upload_type,
			p.paper_title,
			p.conference_title,
			p.location,
			p.date_of_presentation,
			p.pdf,
			p.certificate,
			p.award,
			p.approval_status,
			l.user_name
		FROM paperpresentation AS p
		INNER JOIN login AS l 
			ON p.rollno = l.rollno;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not fetch the data for paperpresentation for verifications")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v facultymodel.Varification
		err := rows.Scan(
			&v.Id,
			&v.UploadType,
			&v.PaperTitle,
			&v.ConferenceTitle,
			&v.Location,
			&v.DateOfPresentation,
			&v.Pdf,
			&v.Certificate,
			&v.Award,
			&v.Approval_status,
			&v.User_name,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err.Error())
			return nil, err
		}
		paperPresentation = append(paperPresentation, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return paperPresentation, nil
}