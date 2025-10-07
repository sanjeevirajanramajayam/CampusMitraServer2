package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetPatents() ([]facultymodel.Varification, error) {
	var patents []facultymodel.Varification

	query := `
		SELECT 
			p.id,
			p.upload_type,
			p.title,
			p.application_number,
			p.date_of_filing,
			p.patent_docs,
			p.supporting_files,
			p.link_to_patent_listing,
			p.summary,
			p.usecase_of_patent,
			p.patent_status,
			l.user_name
		FROM patents AS p
		INNER JOIN login AS l
			ON p.rollno = l.rollno;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not fetch data for patents")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v facultymodel.Varification
		err := rows.Scan(
			&v.Id,
			&v.UploadType,
			&v.Title,
			&v.ApplicationNumber,
			&v.DateOfFiling,
			&v.PatentDocs,
			&v.SupportingFiles,
			&v.LinkToPatentListing,
			&v.Summary,
			&v.UsecaseOfPatent,
			&v.PatentStatus,
			&v.User_name,
		)
		if err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		patents = append(patents, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patents, nil
}