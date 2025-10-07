package achievementgraph

import (
	"bitresume/config"
	"bitresume/models"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandlePointlogs2(rollno string, point float64 , sem int, currdate string) {
	stmt, err := config.DB.Prepare("INSERT INTO point_logs2 (rollno, points, sem, currdate) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return
	}
	defer stmt.Close()
	_, execErr := stmt.Exec(rollno, point, sem, currdate)
	if execErr != nil {
		log.Printf("Failed to execute insert: %v", execErr)
	}
}
func FetchLastPoints(rollno string) (models.Achievementgraph, error) {
	var r models.Achievementgraph
	stmt, err := config.DB.Prepare("SELECT cummulative_points FROM achievement_graph WHERE rollno = ? ORDER BY currdate DESC LIMIT 1")
	if err != nil {
		log.Printf("Prepare failed: %v", err)
		return r, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(rollno)
	err = row.Scan(&r.Cummulative_points)
	if err != nil {
		log.Printf("Failed to scan last points for rollno %s: %v", rollno, err)
		return r, err
	}
	fmt.Print("Last points fetched: ", r.Cummulative_points)
	return r, nil
}
func HandleInactivity(rollno string, currdate string, sem int) error {
	// Step 1: Check if today's point already added
	var count int
	checkQuery := `SELECT COUNT(*) FROM point_logs2 WHERE rollno = ? AND currdate = ?`
	err := config.DB.QueryRow(checkQuery, rollno, currdate).Scan(&count)
	fmt.Print("Count: ", count)
	if err != nil {
		return fmt.Errorf("check query error: %w", err)
	}
	if count > 0 {
		fmt.Println("Point already added today for", rollno)
		return nil
	}
	pointsQuery := `
		SELECT points 
		FROM point_logs2 
		WHERE rollno = ? AND DATE(currdate) < ? 
		ORDER BY DATE(currdate) DESC 
		LIMIT 5
	`
	rows, err := config.DB.Query(pointsQuery, rollno, currdate)
	if err != nil {
		return fmt.Errorf("fetch last 5 points error: %w", err)
	}
	defer rows.Close()

	var total float64
	var entries int

	for rows.Next() {
		var point float64
		if err := rows.Scan(&point); err != nil {
			return fmt.Errorf("scan point error: %w", err)
		}
		total += point
		entries++
	}

	var avg float64
	if entries > 0 {
		avg = ((total / float64(entries))*0.2) //in this i need to modify the avg points (i should take some percent of avg points)
	} else {
		avg = 0.0
	}

	insertQuery := `INSERT INTO point_logs2 (rollno, points, sem, currdate) VALUES (?, ?, ?, ?)`
	_, err = config.DB.Exec(insertQuery, rollno, avg, sem, currdate)
	if err != nil {
		return fmt.Errorf("insert point error: %w", err)
	}
	fmt.Printf("Inserted avg point %.2f for rollno %s on %s\n", avg, rollno, currdate)
	return nil
}
 
// HandleAchievementPoints calculates and updates current points and rank for a student called in ----------cron job
func HandleAcheivemnetPoints(rollno string, currdate string, sem int) {
	fmt.Print("Rollno: ", rollno)
	fmt.Print("Currdate: ", currdate)
	fmt.Print("Sem: ", sem)
	stmt, err := config.DB.Prepare("SELECT SUM(points) FROM point_logs2 WHERE rollno = ? AND currdate = ?")
	if err != nil {
		log.Printf("Failed to prepare points sum query: %v", err)
		return
	}
	defer stmt.Close()
	var r models.Achievementgraph
	row := stmt.QueryRow(rollno, currdate)
	err = row.Scan(&r.Points_earned)
	if err != nil {
		log.Printf("Failed to scan SUM of points for %s on %s: %v", rollno, currdate, err)
		return
	}
	prevPoints, err := FetchLastPoints(rollno)
	if err != nil {
		log.Printf("Failed to fetch previous points for %s: %v", rollno, err)
		return
	}
	newpoints := prevPoints.Cummulative_points + r.Points_earned
	insertStmt, reqerr := config.DB.Prepare("INSERT INTO achievement_graph(rollno, cummulative_points, points_earned, sem, currdate) VALUES (?, ?, ?, ?, ?)")
	if reqerr != nil {
		log.Printf("Failed to prepare insert: %v", reqerr)
		return
	}
	defer insertStmt.Close()
	_, execErr := insertStmt.Exec(rollno, newpoints, r.Points_earned, sem, currdate)
	if execErr != nil {
		log.Printf("Failed to insert activity graph data: %v", execErr)
	}
	HandleInstituteAvg(sem, currdate)
}
// calculate the average points for the institute called in -----------------------cron job
func FetchLastInstituteAvg() (float64, error) {
	var lastCumulative float64
	stmt, err := config.DB.Prepare("SELECT cummulative_points FROM institute_avg ORDER BY id DESC LIMIT 1")
	if err != nil {
		log.Printf("Prepare failed: %v", err)
		return 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	err = row.Scan(&lastCumulative)
	if err != nil {
		log.Printf("Failed to scan last cumulative_points: %v", err)
		return 0, err
	}
	// fmt.Println("Last cumulative_points:", lastCumulative)
	return lastCumulative, nil
}

func HandleInstituteAvg(sem int, currdate string) { 
	lastCumulative, err := FetchLastInstituteAvg()
	if err != nil {
		log.Printf("Failed to fetch last cumulative_points: %v", err)
		return
	}
	stmt, err := config.DB.Prepare("SELECT AVG(points_earned) FROM achievement_graph WHERE currdate = ?")
	if err != nil {
		log.Printf("Failed to prepare points sum query: %v", err)
		return
	}
	defer stmt.Close()
	var r models.Institute_avg
	row := stmt.QueryRow(currdate)
	err = row.Scan(&r.Points)
	CurrentCumulative := lastCumulative + r.Points
	// r.Points i need to divide by the number of students
	if err != nil {
		log.Printf("Failed to scan SUM of points for %s: %v", currdate, err)
		return
	}
	insertStmt, reqerr := config.DB.Prepare("INSERT INTO institute_avg(cummulative_points,points,sem, currdate) VALUES (?,?, ?, ?)")
	if reqerr != nil {
		log.Printf("Failed to prepare insert: %v", reqerr)
		return
	}
	defer insertStmt.Close()
	_, execErr := insertStmt.Exec(CurrentCumulative,r.Points, sem, currdate)
	if execErr != nil {
		log.Printf("Failed to insert activity graph data: %v", execErr)
	}
}

func HandleFetchAchievementGraph(c *gin.Context) {
	var records []models.Achievementgraph
	var r models.Achievementgraph
	rollno := c.Param("rollno")
	rows, err := config.DB.Query("SELECT cummulative_points, points_earned, sem, currdate FROM achievement_graph WHERE rollno = ?", rollno)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&r.Cummulative_points, &r.Points_earned, &r.Sem, &r.Currdate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}

	c.JSON(http.StatusAccepted, records)
}

func HandleFetchInstituteAvg(c *gin.Context) {
	var records1 []models.Institute_avg

	rows1, err1 := config.DB.Query("SELECT cummulative_points, points, sem, currdate FROM institute_avg")
	if err1 != nil {
		c.JSON(500, gin.H{"error": err1.Error()})
		return
	}
	defer rows1.Close()

	for rows1.Next() {
		var r1 models.Institute_avg
		err := rows1.Scan(&r1.Cummulative_points, &r1.Points, &r1.Sem, &r1.Currdate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records1 = append(records1, r1)
	}

	c.JSON(http.StatusAccepted, records1)
}

