package studentdata

import (
	"bitresume/config"

	"github.com/gin-gonic/gin"
)

func HandleStudentData(c *gin.Context) {
	query := `SELECT
    l.rollno,
	l.year,
    l.user_name,
    ag.current_point,
    ag.current_rank,
    acg.cummulative_points
	FROM login l
	JOIN (
		SELECT
			ag1.rollno,
			ag1.current_point,
			ag1.current_rank
		FROM activity_graph ag1
		INNER JOIN (
			SELECT
				rollno,
				MAX(currdate) AS max_date
			FROM activity_graph
			GROUP BY rollno
		) latest_ag
			ON ag1.rollno = latest_ag.rollno
			AND ag1.currdate = latest_ag.max_date
	) ag
		ON l.rollno = ag.rollno
	JOIN (
		SELECT
			acg1.rollno,
			acg1.cummulative_points
		FROM achievement_graph acg1
		INNER JOIN (
			SELECT
				rollno,
				MAX(currdate) AS max_date
			FROM achievement_graph
			GROUP BY rollno
		) latest_acg
			ON acg1.rollno = latest_acg.rollno
			AND acg1.currdate = latest_acg.max_date
	) acg
		ON l.rollno = acg.rollno;`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()
	var results []struct {
		RollNo            string  `json:"rollno"`
		Year              string  `json:"year"`
		Name              string  `json:"user_name"`
		CurrentPoint      float32 `json:"current_point"`
		CurrentRank       string  `json:"current_rank"`
		CummulativePoints float32 `json:"cummulative_points"`
	}
	for rows.Next() {
		var result struct {
			RollNo            string  `json:"rollno"`
			Year              string  `json:"year"`
			Name              string  `json:"user_name"`
			CurrentPoint      float32 `json:"current_point"`
			CurrentRank       string  `json:"current_rank"`
			CummulativePoints float32 `json:"cummulative_points"`
		}
		if err := rows.Scan(&result.RollNo, &result.Year, &result.Name, &result.CurrentPoint, &result.CurrentRank, &result.CummulativePoints); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, result)
	}
	c.JSON(200, gin.H{"mentees": results})
}
func HandleMenteesData(c *gin.Context) {
	rollno := c.Param("rollno")
	query := `SELECT
    l.rollno,
	l.year,
    l.user_name,
    ag.current_point,
    ag.current_rank,
    acg.cummulative_points
	FROM login l
	JOIN (
		SELECT
			ag1.rollno,
			ag1.current_point,
			ag1.current_rank
		FROM activity_graph ag1
		INNER JOIN (
			SELECT
				rollno,
				MAX(currdate) AS max_date
			FROM activity_graph
			GROUP BY rollno
		) latest_ag
			ON ag1.rollno = latest_ag.rollno
			AND ag1.currdate = latest_ag.max_date
	) ag
		ON l.rollno = ag.rollno
	JOIN (
		SELECT
			acg1.rollno,
			acg1.cummulative_points
		FROM achievement_graph acg1
		INNER JOIN (
			SELECT
				rollno,
				MAX(currdate) AS max_date
			FROM achievement_graph
			GROUP BY rollno
		) latest_acg
			ON acg1.rollno = latest_acg.rollno
			AND acg1.currdate = latest_acg.max_date
	) acg
		ON l.rollno = acg.rollno
	WHERE l.mentor_id = ? order by current_point desc;
	`
	rows, err := config.DB.Query(query, rollno)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()
	var results []struct {
		RollNo            string  `json:"rollno"`
		Year              string  `json:"year"`
		Name              string  `json:"user_name"`
		CurrentPoint      float32 `json:"current_point"`
		CurrentRank       string  `json:"current_rank"`
		CummulativePoints float32 `json:"cummulative_points"`
	}
	for rows.Next() {
		var result struct {
			RollNo            string  `json:"rollno"`
			Year              string  `json:"year"`
			Name              string  `json:"user_name"`
			CurrentPoint      float32 `json:"current_point"`
			CurrentRank       string  `json:"current_rank"`
			CummulativePoints float32 `json:"cummulative_points"`
		}
		if err := rows.Scan(&result.RollNo, &result.Year, &result.Name, &result.CurrentPoint, &result.CurrentRank, &result.CummulativePoints); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			return
		}
		results = append(results, result)
	}
	c.JSON(200, gin.H{"mentees": results})
}
