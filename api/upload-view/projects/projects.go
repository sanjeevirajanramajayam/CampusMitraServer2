package projects

import (
	"bitresume/config"
	projectModal "bitresume/models/Project"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func PostProjects(c *gin.Context) {
	// Read basic form data
	submitter_roll_no := c.PostForm("submitter_roll_no")
	title_idea := c.PostForm("title_idea")
	project_abstract := c.PostForm("project_abstract")
	problem_statement := c.PostForm("problem_statement")
	objective := c.PostForm("objective")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	is_team_project := c.PostForm("is_team_project")
	consulted_mentor := c.PostForm("consulted_mentor")
	github_link := c.PostForm("github_link")
	presented_externally := c.PostForm("presented_externally")
	awards_won := c.PostForm("awards_won")

	// Handle file uploads
	demo_video, err := c.FormFile("demo_video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Demo video upload failed", "details": err.Error()})
		return
	}
	saveVIDEOPath := filepath.Join("uploads/projects/demovideos", demo_video.Filename)
	if err := c.SaveUploadedFile(demo_video, saveVIDEOPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save demo video", "details": err.Error()})
		return
	}

	report_pdf, err := c.FormFile("report_pdf")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report PDF upload failed", "details": err.Error()})
		return
	}
	savePDFPath := filepath.Join("uploads/projects/report_PDF", report_pdf.Filename)
	if err := c.SaveUploadedFile(report_pdf, savePDFPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save report PDF", "details": err.Error()})
		return
	}

	// Parse boolean flags
	IsTeamProject := 0
	if is_team_project == "true" {
		IsTeamProject = 1
	}
	ConsultedMentor := 0
	if consulted_mentor == "true" {
		ConsultedMentor = 1
	}
	PresentedExternally := 0
	if presented_externally == "true"{
		PresentedExternally = 1
	}
	// Parse team members JSON (sent as string in form-data)
	var teamMembers []projectModal.TeamMember
	teamMembersJSON := c.PostForm("team_members")
	if teamMembersJSON != "" {
		if err := json.Unmarshal([]byte(teamMembersJSON), &teamMembers); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team_members format", "details": err.Error()})
			return
		}
	}

	// Insert into projects table
	projectsUpload := `
		insert into projects(
			upload_type,
			rollno,
			title_idea,
			summary,
			problem_statement,
			objective,
			start_time,
			end_time,
			is_team_project,
			consulted_mentor,
			approval_status,
			complexity
		) values (?,?,?,?,?,?,?,?,?,?,?,?)
	`

	res, err := config.DB.Exec(
		projectsUpload,
		"Project",
		submitter_roll_no,
		title_idea,
		project_abstract,
		problem_statement,
		objective,
		start_time,
		end_time,
		IsTeamProject,
		ConsultedMentor,
		"Pending",
		"T1",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert project", "details": err.Error()})
		return
	}

	projectID, _ := res.LastInsertId()

	// Insert team members + tech stack
	// ✅ Handle team vs individual
	if IsTeamProject == 1 {
		for _, member := range teamMembers {
			// Skip empty placeholder rows
			if member.Name == "" && member.RollNumber == "" {
				continue
			}

			// Insert into project_team_members
			teamInsert := `
			insert into project_team_members (project_id, rollno, member_name, department)
			values (?,?,?,?)
		`
			_, err := config.DB.Exec(teamInsert, projectID, member.RollNumber, member.Name, member.Department)
			if err != nil {
				fmt.Println("Error inserting team member:", err.Error())
				continue
			}

			// Insert tech stack for this member
			for _, tech := range member.TechStack {
				techInsert := `
				insert into project_tech_stack (project_id, tech_name)
				values (?,?)
			`
				_, err := config.DB.Exec(techInsert, projectID, tech)
				if err != nil {
					fmt.Println("Error inserting tech stack:", err.Error())
					continue
				}
			}
		}
	} else {
		// ✅ Individual project → only insert tech stack
		if len(teamMembers) > 0 {
			for _, tech := range teamMembers[0].TechStack {
				techInsert := `
				insert into project_tech_stack (project_id, tech_name)
				values (?,?)
			`
				_, err := config.DB.Exec(techInsert, projectID, tech)
				if err != nil {
					fmt.Println("Error inserting tech stack:", err.Error())
					continue
				}
			}
		}
	}

	PresentationUpload := `
		 insert into project_presentations(
			project_id,
    		presented_externally,
    		awards_won
		 )values(?,?,?)
	`

	_ , err = config.DB.Exec(PresentationUpload,projectID,PresentedExternally,awards_won)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	ProjectFilesUpload := `
		insert into project_files(
			project_id,
    		github_link,
    		report_pdf,
   			demo_video
		) values (?,?,?,?)
	`

	_ , err = config.DB.Exec(ProjectFilesUpload,projectID,github_link,savePDFPath,saveVIDEOPath)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully uploaded"})
}
