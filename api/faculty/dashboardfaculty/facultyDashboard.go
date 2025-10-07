package dashBoardfaculty

import (
	"bitresume/config"

	"github.com/gin-gonic/gin"
)

func Leaderboard(c *gin.Context){
	query:=`SELECT
    ag.rollno,
	l.user_name,
    l.department,
    ag.current_point,
    ag.current_rank,
    COALESCE(p.project_count, 0) AS project_count
	FROM
		activity_graph ag
	INNER JOIN (
		SELECT
			rollno,
			MAX(currdate) AS max_date
		FROM
			activity_graph
		GROUP BY
			rollno
	) AS latest 
		ON ag.rollno = latest.rollno 
		AND ag.currdate = latest.max_date
	JOIN login l 
		ON ag.rollno = l.rollno 
		AND l.mentor_id = ?
	LEFT JOIN (
		SELECT 
			rollno, 
			COUNT(*) AS project_count
		FROM 
			projects       
		GROUP BY 
			rollno
	) AS p 
		ON ag.rollno = p.rollno;
	`
	rows, err := config.DB.Query(query, c.Param("rollno"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()
	var results []struct {
		RollNo        string `json:"rollno"`
		Name        string `json:"user_name"`
		Department        string `json:"department"`
		CurrentPoint  float32    `json:"current_point"`
		CurrentRank   string    `json:"current_rank"`
		ProjectCount  int    `json:"project_count"`
	}
	for rows.Next() {
		var result struct {
			RollNo       string `json:"rollno"`
			Name         string `json:"user_name"`
			Department   string `json:"department"`
			CurrentPoint float32    `json:"current_point"`
			CurrentRank  string    `json:"current_rank"`
			ProjectCount int    `json:"project_count"`
		}
		if err := rows.Scan(&result.RollNo,&result.Name,&result.Department, &result.CurrentPoint, &result.CurrentRank, &result.ProjectCount); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, result)
	}
	c.JSON(200, gin.H{"leaderboard": results})
}
func HandlePriorityLearners(c *gin.Context){
	query:=`SELECT
    ag.rollno,
    l.user_name,
    l.department,
    ag.current_point,
    ag.current_rank,
    COALESCE(p.project_count, 0) AS project_count
	FROM
		activity_graph ag
	INNER JOIN (
		SELECT
			rollno,
			MAX(currdate) AS max_date
		FROM
			activity_graph
		GROUP BY
			rollno
	) AS latest 
		ON ag.rollno = latest.rollno 
		AND ag.currdate = latest.max_date
	JOIN login l 
		ON ag.rollno = l.rollno 
		AND l.mentor_id = ?
	LEFT JOIN (
		SELECT 
			rollno, 
			COUNT(*) AS project_count
		FROM 
			projects       
		GROUP BY 
			rollno
	) AS p 
		ON ag.rollno = p.rollno where ag.current_rank='Silver' order by current_point`
	rows, err := config.DB.Query(query, c.Param("rollno"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()
	var results []struct {
		RollNo        string  `json:"rollno"`
		Name          string  `json:"user_name"`
		Department    string  `json:"department"`
		CurrentPoint  float32 `json:"current_point"`
		CurrentRank   string  `json:"current_rank"`
		ProjectCount  int     `json:"project_count"`
	}
	for rows.Next() {
		var result struct {
			RollNo       string  `json:"rollno"`
			Name         string  `json:"user_name"`
			Department   string  `json:"department"`
			CurrentPoint float32 `json:"current_point"`
			CurrentRank  string  `json:"current_rank"`
			ProjectCount int     `json:"project_count"`
		}
		if err := rows.Scan(&result.RollNo, &result.Name, &result.Department, &result.CurrentPoint, &result.CurrentRank, &result.ProjectCount); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, result)
		c.JSON(200, gin.H{"prioritylearners": results})
	}
}