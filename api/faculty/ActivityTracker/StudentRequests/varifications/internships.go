package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetInternship() ([]facultymodel.Varification, error) {
	var internships []facultymodel.Varification

	query := `
		SELECT 
			i.id,
			i.upload_type,
			i.company_name,
			i.roll,
			i.domain,
			i.internship_type,
			i.is_stipend,
			i.start_date,
			i.end_date,
			i.consulted_faculty_name,
			i.industry_mentor_name,
			i.industry_mentor_contact,
			i.offer_letter,
			i.report,
			i.skill_gained,
			i.outcomes,
			i.status,
			l.user_name
		FROM internships AS i
		INNER JOIN login AS l 
			ON i.rollno = l.rollno;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not fetch data for internships")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v facultymodel.Varification
		err := rows.Scan(
			&v.Id,
			&v.UploadType,
			&v.CompanyName,
			&v.Roll,
			&v.Domain,
			&v.InternshipType,
			&v.IsStipend,
			&v.StartDate,
			&v.EndDate,
			&v.ConsultedFacultyName,
			&v.IndustryMentorName,
			&v.IndustryMentorContact,
			&v.OfferLetter,
			&v.Report,
			&v.SkillGained,
			&v.Outcomes,
			&v.Approval_status,
			&v.User_name,
		)
		if err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		internships = append(internships, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return internships, nil
}