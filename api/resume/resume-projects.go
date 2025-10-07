package resume

import (
	"bitresume/config"
	"bitresume/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProjectsData returns each project (title, summary, GitHub link) with its tech stack list.
func GetProjectsData(c *gin.Context) {
	rollno := c.Param("rollno")
	// Query joins: projects + project_files + project_tech_stack
	query := `
		SELECT 
			p.id,
			p.title_idea,
			p.summary,
			COALESCE(f.github_link, '') AS github_link,
			COALESCE(t.tech_name, '')   AS tech_name
		FROM projects p
		LEFT JOIN project_files f      ON p.id = f.project_id
		LEFT JOIN project_tech_stack t ON p.id = t.project_id
		WHERE p.rollno = ? and p.approval_status='Approved'
		ORDER BY p.id;
	`

	rows, err := config.DB.Query(query, rollno)
	if err != nil {
		fmt.Println("Error querying joined tables:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch project data"})
		return
	}
	defer rows.Close()

	// Use a map to aggregate tech stacks for each project ID
	projectMap := make(map[int]*models.Project)

	for rows.Next() {
		var (
			id         int
			title      string
			summary    string
			githubLink string
			techName   string
		)
		if err := rows.Scan(&id, &title, &summary, &githubLink, &techName); err != nil {
			fmt.Println("Scan error:", err)
			continue
		}

		// If we haven't seen this project yet, create it
		if _, exists := projectMap[id]; !exists {
			projectMap[id] = &models.Project{
				Title:       title,
				Description: summary,
				Github:      githubLink,
				Stack:       []string{},
			}
		}

		// Append tech stack name if not empty
		if techName != "" {
			projectMap[id].Stack = append(projectMap[id].Stack, techName)
		}
	}

	// Convert map to slice
	var allProjects []models.Project
	for _, p := range projectMap {
		allProjects = append(allProjects, *p)
	}

	c.JSON(http.StatusOK, allProjects)
}
func GetAresOfExpertise(c *gin.Context){
	
}