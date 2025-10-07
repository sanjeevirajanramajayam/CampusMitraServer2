package activitygraph

import (
	"bitresume/config"
	"bitresume/models"
	"fmt"

	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FetchActivityGraphData returns full activity graph records for   student
func FetchActivityGraphData(c *gin.Context) {
	var records []models.ActGph
	var r models.ActGph
	rollno := c.Param("rollno")
	rows, err := config.DB.Query("SELECT rollno, current_point, current_rank, sem, currdate FROM activity_graph WHERE rollno = ?", rollno)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&r.RollNo, &r.Current_point, &r.Current_rank, &r.Sem, &r.Currdate)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}
	c.JSON(http.StatusAccepted, records)
}

// FetchDataRank fetches the most recent rank for a student
func FetchDataRank(rollno string) (models.ActGph, error) {
	var r models.ActGph
	stmt, err := config.DB.Prepare("SELECT current_rank FROM activity_graph WHERE rollno = ? ORDER BY currdate DESC LIMIT 1")
	if err != nil {
		log.Printf("Prepare failed: %v", err)
		return r, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(rollno)
	err = row.Scan(&r.Current_rank)
	if err != nil {
		log.Printf("Failed to scan rank for rollno %s: %v", rollno, err)
		return r, err
	}

	return r, nil
}

// FetchLastPoints fetches the most recent point total for a student

func FetchLastPoints(rollno string) (models.ActGph, error) {
	var r models.ActGph
	stmt, err := config.DB.Prepare("SELECT current_point FROM activity_graph WHERE rollno = ? ORDER BY currdate DESC LIMIT 1")
	if err != nil {
		log.Printf("Prepare failed: %v", err)
		return r, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(rollno)
	err = row.Scan(&r.Current_point)
	if err != nil {
		log.Printf("Failed to scan last points for rollno %s: %v", rollno, err)
		return r, err
	}
	// fmt.Print("Last points: ", r.Current_point)
	return r, nil
}

// HandleActivityGraphPoints calculates and updates current points and rank for a student called in cron job
func HandleActivityGraphPoints(rollno string, sem int, currdate string) {
	stmt, err := config.DB.Prepare("SELECT SUM(points) FROM points_logs WHERE rollno = ? AND currdate = ?")
	if err != nil {
		log.Printf("Failed to prepare points sum query: %v", err)
		return
	}
	defer stmt.Close()
	var r models.ActGph
	row := stmt.QueryRow(rollno, currdate)
	err = row.Scan(&r.Current_point)
	if err != nil {
		log.Printf("Failed to scan SUM of points for %s on %s: %v", rollno, currdate, err)
		return
	}
	prevPoints, err := FetchLastPoints(rollno)
	if err != nil {
		log.Printf("Failed to fetch previous points for %s: %v", rollno, err)
		return
	}
	var newpoints float64
	var rank string
	newpoints = float64(prevPoints.Current_point) + float64(r.Current_point)

	if newpoints > 100{
		newpoints = 100 // Need to add bonus table
	}
	if newpoints < 70 {
		newpoints = 70 //Need to add continuos inactivity table 
	}

	switch {
	case newpoints >= 90:
		rank = "TITANIUM"
	case newpoints >= 80:
		rank = "GOLD"
	default:
		rank = "SILVER"
	}

	insertStmt, reqerr := config.DB.Prepare("INSERT INTO activity_graph(rollno, current_point, current_rank, sem, currdate) VALUES (?, ?, ?, ?, ?)")
	if reqerr != nil {
		log.Printf("Failed to prepare insert: %v", reqerr)
		return
	}
	defer insertStmt.Close()

	_, execErr := insertStmt.Exec(rollno, newpoints, rank, sem, currdate)
	if execErr != nil {
		log.Printf("Failed to insert activity graph data: %v", execErr)
	}
}
// Handle inactivity for a student called in ----------------------cron job
func HandleInactivity(rollno string, currDate string, sem int) error {
    var count int
    query := `SELECT COUNT(*) FROM points_logs WHERE rollno = ? AND DATE(currdate) = DATE(?)`
    err := config.DB.QueryRow(query, rollno, currDate).Scan(&count)
    if err != nil {
        return fmt.Errorf("query error: %w", err)
    }              
	fmt.Print("Count: ", count)
    if count == 0 {
        rank, rankerr := FetchDataRank(rollno)
        if rankerr != nil {
            return fmt.Errorf("rank fetch error: %w", rankerr)
        }
        var newpoints float64
        if rank.Current_rank == "TITANIUM" {
            newpoints = -0.4
        } else if rank.Current_rank == "GOLD" {
            newpoints = -0.3
        } else {
            newpoints = -0.2
        }
        insertQuery := `INSERT INTO points_logs(rollno, source, points, description, sem, currdate) 
                        VALUES (?, 'inactivity', ?, 'inactivity', ?, ?)`
        _, err = config.DB.Exec(insertQuery, rollno, newpoints, sem, currDate)
        if err != nil {
            return fmt.Errorf("insert penalty error: %w", err)
        }
        fmt.Println("Penalty point added for", rollno)
    } else {
        fmt.Println("Point already added today for", rollno)
    }

    return nil
}