package pointshandlers

import (
	achievementgraph "bitresume/api/dashboard/achievement_graph"
	activitygraph "bitresume/api/dashboard/activity_graph"
	"bitresume/config"
	"bitresume/models"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
) 
//Main function for points if points is come by his activity
func HandlePointlogs(rollno ,source string ,points int ,desc string,sem int,currdate string) error  { //This is for all other than ps 
	var newpoints float64
	rank, rankerr := activitygraph.FetchDataRank(rollno)
	if rankerr != nil {
		return	rankerr
	}
	if source == "PS" {
		// HandlePs(rollno)
		if rank.Current_rank == "TITANIUM" {
			if points > 0 {
				newpoints = float64(points) * 0.5 / 300.0
			} else if points == 0 {
				newpoints = 0
			} else {
				newpoints = -1

			}
		} else if rank.Current_rank == "GOLD" {
			if points > 0 {
				newpoints = float64(points) * 1 / 300.0
			} else if points == 0 {
				newpoints = 0
			} else {
				newpoints = -0.5
			}
		} else {
			if points > 0 {
				newpoints = float64(points) * 2 / 300.0
			} else if points == 0 {
				newpoints = 0
			} else {
				newpoints = -0.5
			}
		}
	}
	newpoints = math.Round(newpoints*100) / 100
	stmp, reqerr := config.DB.Prepare("INSERT INTO points_logs(rollno,source,points,description,sem,currdate) values (?,?,?,?,?,?)")
	if reqerr != nil {
		return reqerr
	}
	_, execErr := stmp.Exec(rollno, source, newpoints, desc, sem, currdate)
	if execErr != nil {
		return execErr
	}
	if points > 0 {
		achievementgraph.HandlePointlogs2(rollno, newpoints, sem, currdate) //to calculate the achievement points
	}	
	return nil
	}
func HandlePs(c *gin.Context) { //if attempted itself
	var data models.Ps
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	rollno := data.RollNo
	points := data.Points //rewards points for that level
	domain := data.SkillDomain 
	skillname := data.SkillName
	skilllevel := data.SkillLevel
	desc := skillname + " " + skilllevel
	attempts := data.Attempts
	currdate := data.Currdate
	sem := data.Sem
	source := "PS"
	var newpoints float64
	// Fetch current rank
	rank, rankerr := activitygraph.FetchDataRank(rollno)
	if rankerr != nil { 
		c.JSON(500, gin.H{"error": rankerr.Error()})
		return
	}

	// Calculate new points based on rank
	switch rank.Current_rank {
	case "TITANIUM":
		if points > 0 {
			newpoints = float64(points) * 0.5 / 300.0
		} else if points == 0 {
			newpoints = 0
		} else {
			newpoints = -1
		}
	case "GOLD":
		if points > 0 {
			newpoints = float64(points) * 1 / 300.0
		} else if points == 0 {
			newpoints = 0
		} else {
			newpoints = -0.5
		}
	default:
		if points > 0 {   //Silver
			newpoints = float64(points) * 2 / 300.0
		} else if points == 0 {   //Fail in that level
			newpoints = 0
		} else {
			newpoints = -0.5 //Attempted but not went
		}
	}
	newpoints = math.Round(newpoints*100) / 100
	// Insert into points_logs
	stmp, reqerr := config.DB.Prepare(`INSERT INTO points_logs(rollno, source, points, description, sem, currdate) VALUES (?, ?, ?, ?, ?, ?)`)
	if reqerr != nil {
		log.Println("Error preparing points_logs insert:", reqerr)
		c.JSON(500, gin.H{"error": reqerr.Error()})
		return
	}
	_, execErr := stmp.Exec(rollno, source, newpoints, desc, sem, currdate)
	if execErr != nil {
		log.Println("Error executing points_logs insert:", execErr)
		c.JSON(500, gin.H{"error": execErr.Error()})
		return
	}
	if points > 0 {
		achievementgraph.HandlePointlogs2(rollno, newpoints, sem, currdate)
	}
	// Check if the skill-level already exists in ps_status
	var exists bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM ps_status WHERE rollno = ? AND skill_name = ? AND skill_level = ?)`
	err := config.DB.QueryRow(checkQuery, rollno, skillname, skilllevel).Scan(&exists)
	if err != nil {
		log.Println("Error checking existing skill:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		// Update attempts if record exists
		updateStmt, err := config.DB.Prepare(`UPDATE ps_status SET attempts = ? WHERE rollno = ? AND skill_name = ? AND skill_level = ?`)
		if err != nil {
			log.Println("Error preparing update:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = updateStmt.Exec(attempts, rollno, skillname, skilllevel)
		if err != nil {
			log.Println("Error executing update:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		log.Println("Updated existing skill level record.")
	} else {
		// Insert new record if not exists
		insertStmt, err := config.DB.Prepare(`INSERT INTO ps_status (rollno, skill_domain, skill_name, skill_level, attempts) VALUES (?, ?, ?, ?, ?)`)
		if err != nil {
			log.Println("Error preparing insert:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		_, err = insertStmt.Exec(rollno, domain, skillname, skilllevel, attempts)
		if err != nil {
			log.Println("Error executing insert:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		log.Println("Inserted new skill level record.")
	}

	// Return success message
	log.Println("Successfully handled PS log for", rollno, "Skill:", skillname, "Level:", skilllevel)
	c.JSON(200, gin.H{
		"message": "PS data logged successfully",
		"rollno":  rollno,
		"skill":   skillname,
		"level":   skilllevel,
		"points":  newpoints,
		"operation": func() string {
			if exists {
				return "updated"
			}
			return "inserted"
		}(),
	})
}

func HandlePsLevelStatus(c *gin.Context) { //If completed only
	var data models.PsLevels
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Check if record exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM ps_level_status 
		WHERE rollno = ? AND skilldomain = ? AND skillname = ?
	)`
	err := config.DB.QueryRow(query, data.RollNo, data.SkillDomain, data.SkillName).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check record"})
		return
	}
	if exists {
		// 2. Update if exists
		updateQuery := `
			UPDATE ps_level_status 
			SET levels_completed = ?, total_levels = ?, updated_at = NOW() 
			WHERE rollno = ? AND skilldomain = ? AND skillname = ?
		`
		_, err := config.DB.Exec(updateQuery, data.SkillLevel, data.TotalLevels, data.RollNo, data.SkillDomain, data.SkillName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Skill updated successfully"})
	} else {
		// 3. Insert if not exists
		insertQuery := `
			INSERT INTO ps_level_status (rollno, skilldomain, skillname, levels_completed, total_levels, created_at)
			VALUES (?, ?, ?, ?, ?, NOW())
		`
		_, err := config.DB.Exec(insertQuery, data.RollNo, data.SkillDomain, data.SkillName, data.SkillLevel, data.TotalLevels)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new skill"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "New skill added"})
	}
}
func HandleFetchPsAttempts(c *gin.Context){
	var records []models.Ps
	// var r models.Ps
	rollno := c.Param("rollno")
	rows, err := config.DB.Query("SELECT  skill_domain, skill_name, skill_level, attempts FROM ps_status WHERE rollno = ?", rollno)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r models.Ps
		err := rows.Scan(&r.SkillDomain, &r.SkillName, &r.SkillLevel, &r.Attempts)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}
	c.JSON(http.StatusAccepted, records)
}                                       
func HandleFetchPsLevels(c *gin.Context) {
	var records []models.PsLevels
	// var r models.PsLevels
	rollno := c.Param("rollno")
	rows, err := config.DB.Query("SELECT skilldomain, skillname, levels_completed, total_levels FROM ps_level_status WHERE rollno = ?", rollno)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r models.PsLevels
		err := rows.Scan(&r.SkillDomain, &r.SkillName, &r.SkillLevel, &r.TotalLevels)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		records = append(records, r)
	}
	c.JSON(http.StatusAccepted, records)
}
func HandleSemDays(c *gin.Context) {
	var results []models.SemCount
	rows, err := config.DB.Query("SELECT sem, COUNT(*) as sem_count FROM activity_graph GROUP BY sem ORDER BY sem")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sc models.SemCount
		if err := rows.Scan(&sc.Sem, &sc.SemCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse database result", "details": err.Error()})
			return
		}
		results = append(results, sc)
	}

	c.JSON(http.StatusOK, results)
}