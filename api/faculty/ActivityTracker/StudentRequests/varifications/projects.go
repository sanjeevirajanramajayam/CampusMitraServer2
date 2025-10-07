package studentrequests

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
)

func GetProjects() ([]facultymodel.Varification, error) {
	var projects []facultymodel.Varification
	query := `
	SELECT
		p.id,
		p.upload_type,
		p.title_idea,
		p.summary,
		p.problem_statement,
		p.objective,
		p.start_time,
		p.end_time,
		p.is_team_project,
		p.consulted_mentor,
		p.approval_status,
		COALESCE(pf.github_link, '') AS github_link,
		COALESCE(pf.report_pdf, '') AS report_pdf,
		COALESCE(pf.demo_video, '') AS demo_video,
		COALESCE(pp.presented_externally, '') AS presented_externally,
		COALESCE(pp.awards_won, '') AS awards_won,
		COALESCE(GROUP_CONCAT(DISTINCT ptm.rollno ORDER BY ptm.rollno SEPARATOR ', '), '') AS rollnos,
		COALESCE(GROUP_CONCAT(DISTINCT ptm.member_name ORDER BY ptm.member_name SEPARATOR ', '), '') AS member_names,
		COALESCE(GROUP_CONCAT(DISTINCT ptm.department ORDER BY ptm.department SEPARATOR ', '), '') AS departments,
		COALESCE(GROUP_CONCAT(DISTINCT pts.tech_name ORDER BY pts.tech_name SEPARATOR ', '), '') AS tech_names,
		l.user_name
	FROM projects AS p
	LEFT JOIN login AS l ON p.rollno = l.rollno
	LEFT JOIN project_evaluation AS pe ON pe.project_id = p.id
	LEFT JOIN project_files AS pf ON pf.project_id = p.id
	LEFT JOIN project_presentations AS pp ON pp.project_id = p.id
	LEFT JOIN project_team_members AS ptm ON ptm.project_id = p.id
	LEFT JOIN project_tech_stack AS pts ON pts.project_id = p.id
	GROUP BY
		p.id,
		p.upload_type,
		p.title_idea,
		p.summary,
		p.problem_statement,
		p.objective,
		p.start_time,
		p.end_time,
		p.is_team_project,
		p.consulted_mentor,
		p.approval_status,
		pf.github_link,
		pf.report_pdf,
		pf.demo_video,
		pp.presented_externally,
		pp.awards_won,
		l.user_name;
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		fmt.Println("Error querying projects:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r facultymodel.Varification
		err = rows.Scan(
			&r.Id,
			&r.UploadType,
			&r.TitleIdea,
			&r.Summary,
			&r.ProblemStatement,
			&r.Objective,
			&r.StartTime,
			&r.EndTime,
			&r.IsTeamProject,
			&r.ConsultedMentor,
			&r.Approval_status,
			&r.GithubLink,
			&r.ReportPdf,
			&r.DemoVideo,
			&r.PresentedExternally,
			&r.AwardsWon,
			&r.Rollno,
			&r.MemberName,
			&r.Department,
			&r.TechNames,
			&r.User_name,
		)
		if err != nil {
			fmt.Println("Error scanning project row:", err)
			return nil, err
		}
		projects = append(projects, r)
	}
	return projects, nil
}
