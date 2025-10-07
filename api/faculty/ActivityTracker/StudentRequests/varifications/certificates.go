package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetCertificates() ([]facultymodel.Varification, error) {
	var certificates []facultymodel.Varification

	query := `
		SELECT
    COALESCE(c.title, '') AS title,
    COALESCE(c.platform, '') AS platform,
    COALESCE(c.issue_date, '') AS issue_date,
    COALESCE(c.course_link, '') AS course_link,
    COALESCE(c.certificate_pdf, '') AS certificate_pdf,
    COALESCE(c.certiificate_id,'') AS certificate_id,
    '' AS event_name,
    '' AS event_code,
    '' AS participation_type,
    '' AS summary,
    '' AS winning_status,
    '' AS activity_type,
    '' AS duration,
    '' AS location,
    COALESCE(l.user_name, '') AS user_name,
    COALESCE(ct.certificate_type, '') AS certificate_type,
    COALESCE(ct.status, '') AS approval_status,
    COALESCE(ct.upload_type, '') AS upload_type
FROM certificate_onlinecourses AS c
INNER JOIN login AS l ON l.rollno = c.rollno
INNER JOIN certificates_type AS ct ON ct.id = c.certiificate_id

UNION ALL

SELECT
    '' AS title,
    '' AS platform,
    '' AS issue_date,
    '' AS course_link,
    COALESCE(c.certificate_pdf, '') AS certificate_pdf,
    COALESCE(c.certificate_id,'') AS certificate_id,
    COALESCE(c.event_name, '') AS event_name,
    COALESCE(c.event_code, '') AS event_code,
    COALESCE(c.participation_type, '') AS participation_type,
    COALESCE(c.summary, '') AS summary,
    COALESCE(c.did_you_win, '') AS winning_status,
    '' AS activity_type,
    '' AS duration,
    '' AS location,
    COALESCE(l.user_name, '') AS user_name,
    COALESCE(ct.certificate_type, '') AS certificate_type,
    COALESCE(ct.status, '') AS approval_status,
    COALESCE(ct.upload_type, '') AS upload_type
FROM certificates_events AS c
INNER JOIN login AS l ON l.rollno = c.rollno
INNER JOIN certificates_type AS ct ON ct.id = c.certificate_id

UNION ALL

SELECT
    '' AS title,
    '' AS platform,
    COALESCE(cv.issue_date, '') AS issue_date,
    '' AS course_link,
    COALESCE(cv.certificate_pdf, '') AS certificate_pdf,
    COALESCE(cv.certificate_id, '') AS certificate_id,
    '' AS event_name,
    '' AS event_code,
    '' AS participation_type,
    COALESCE(cv.summary, '') AS summary,
    '' AS winning_status,
    COALESCE(cv.activity_type, '') AS activity_type,
    COALESCE(cv.duration, '') AS duration,
    COALESCE(cv.location, '') AS location,
    COALESCE(l.user_name, '') AS user_name,
    COALESCE(ct.certificate_type, '') AS certificate_type,
    COALESCE(ct.status, '') AS approval_status,
    COALESCE(ct.upload_type, '') AS upload_type
FROM certificates_voluntree AS cv
INNER JOIN certificates_type AS ct ON ct.id = cv.certificate_id
INNER JOIN login AS l ON cv.rollno = l.rollno;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		fmt.Println("Could not execute the combined query for certificates")
		fmt.Println("Error:", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a facultymodel.Varification
		err = rows.Scan(
			&a.Title,
			&a.Platform,
			&a.Issue_date,
			&a.Course_link,
			&a.Certificate_pdf,
			&a.Id,
			&a.EventName,
			&a.EventCode,
			&a.Participation_type,
			&a.Summary,
			&a.WinningStatus,
			&a.Activity_type,
			&a.Duration,
			&a.Location,
			&a.User_name,
			&a.Certificate_type,
			&a.Approval_status,
            &a.UploadType,
		)
		if err != nil {
			fmt.Println("Error scanning the combined certificates query")
			fmt.Println("Error:", err.Error())
			return nil, err
		}
		certificates = append(certificates, a)
	}

	return certificates, nil
}