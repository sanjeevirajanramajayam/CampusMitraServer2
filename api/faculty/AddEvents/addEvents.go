package addevents

import (
	"bitresume/config"
	"encoding/json"
	"strconv"

	// "bitresume/models"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AddEvents handles the submission of new event data via multipart/form-data
func AddEvents(c *gin.Context) {
	// Parse basic form fields
	eventName := c.PostForm("event_name")
	eventType := c.PostForm("type")
	deadline := c.PostForm("deadline")
	location := c.PostForm("location")
	applyLink := c.PostForm("apply_link")
	domains := c.PostForm("domains")
	description := c.PostForm("description")
	rules := c.PostForm("rules")
	constraints := c.PostForm("constraints")
	minTeamSize := c.PostForm("min_team_size")
	maxTeamSize := c.PostForm("max_team_size")
	noOfRounds := c.PostForm("no_of_rounds")
	onlineRounds := c.PostForm("online_rounds")
	offlineRounds := c.PostForm("offline_rounds")
	finalPrice1 := c.PostForm("final_prize1")
	finalPrice2 := c.PostForm("final_prize2")
	finalPrice3 := c.PostForm("final_prize3")
	roundData := c.PostForm("roundsData")

	var rounds []facultymodel.Round
	if err := json.Unmarshal([]byte(roundData), &rounds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("===============================")
	for _, r := range rounds {
		fmt.Println("Round No:", r.RoundNumber)
		fmt.Println("Description:", r.Description) // log new field
		fmt.Println("Start date:", r.StartDate)
		fmt.Println("End date:", r.EndDate)
		fmt.Println("Reward Points:", r.Rewardpoints.Year1, r.Rewardpoints.Year2, r.Rewardpoints.Year3, r.Rewardpoints.Year4)
	}

	var imageURL string
	file, err := c.FormFile("image")
	if err == nil {
		imageURL = "uploads/events/" + file.Filename
		if err := c.SaveUploadedFile(file, imageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	var id int
	err = config.DB.QueryRow(`SELECT COALESCE(MAX(id), 0) FROM events`).Scan(&id)
	if err != nil {
		log.Println("Error fetching max id:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch max id"})
		return
	}

	eventCode := fmt.Sprintf("%02dBIT%d", time.Now().Year()%100, id+1)

	// Insert rounds (NOTE: column name is intentionally "decription" per your migration)
	for _, r := range rounds {
		_, err := config.DB.Exec(`
			INSERT INTO event_rounds_dates (
				event_code, round_number, description, start_date, end_date, year1_rp, year2_rp, year3_rp, year4_rp
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			eventCode, r.RoundNumber, r.Description, r.StartDate, r.EndDate,
			r.Rewardpoints.Year1, r.Rewardpoints.Year2, r.Rewardpoints.Year3, r.Rewardpoints.Year4,
		)
		if err != nil {
			log.Println("Error inserting rounds:", err)
		}
	}

	stmt, err := config.DB.Prepare(`		
		INSERT INTO events (
			event_name, event_code, type, deadline, min_team_size, max_team_size,
			no_of_rounds, online_rounds, offline_rounds, location, apply_link,
			domains, description, rules, constraints, final_prize1, final_prize2, final_prize3, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Println("Error preparing insert:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		eventName, eventCode, eventType, deadline, minTeamSize, maxTeamSize,
		noOfRounds, onlineRounds, offlineRounds, location, applyLink,
		domains, description, rules, constraints, finalPrice1, finalPrice2, finalPrice3, imageURL,
	)
	if err != nil {
		log.Println("Error executing insert:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Event inserted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Event added successfully"})
}
func FetchEvents(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	query := `
		SELECT e.id, e.event_name, e.event_code, e.type, e.deadline, e.min_team_size, e.max_team_size,
			e.no_of_rounds, e.online_rounds, e.offline_rounds, e.location, e.apply_link,
			e.domains, e.image_url, e.description AS event_description, e.rules, e.constraints,
			e.final_prize1, e.final_prize2, e.final_prize3,
			r.round_number, r.description AS round_description, r.start_date, r.end_date, 
			r.year1_rp, r.year2_rp, r.year3_rp, r.year4_rp
		FROM events AS e
		LEFT JOIN event_rounds_dates AS r 
			ON e.event_code = r.event_code order by end_date desc
		LIMIT ? OFFSET ?;
`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error executing query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	// Use map to group events by EventCode
	eventMap := make(map[string]*facultymodel.Event)

	for rows.Next() {
		var event facultymodel.Event
		var round facultymodel.Rounds

		err := rows.Scan(
			&event.ID, &event.EventName, &event.EventCode, &event.Type, &event.Deadline,
			&event.MinTeamSize, &event.MaxTeamSize, &event.NoOfRounds, &event.OnlineRounds,
			&event.OfflineRounds, &event.Location, &event.ApplyLink, &event.Domains,
			&event.ImageURL, &event.Description, &event.Rules, &event.Constraints,
			&event.FinalPrize1, &event.FinalPrize2, &event.FinalPrize3,
			&round.RoundNumber, &round.Description, &round.StartDate, &round.EndDate,
			&round.Year1RP, &round.Year2RP, &round.Year3RP, &round.Year4RP,
		)
		if err != nil {
			log.Println("Error scanning event:", err)
			continue
		}
		if existingEvent, found := eventMap[event.EventCode]; found {
			if round.RoundNumber != 0 {
				existingEvent.Rounds = append(existingEvent.Rounds, round)
			}
		} else {
			if round.RoundNumber != 0 {
				event.Rounds = []facultymodel.Rounds{round}
			}
			eventMap[event.EventCode] = &event
		}
	}

	events := make([]facultymodel.Event, 0, len(eventMap))
	for _, e := range eventMap {
		events = append(events, *e)
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
func CheckApplied(c *gin.Context) {
	rollno := c.Query("rollno")
	eventCode := c.Query("event_code")
	query := `
        SELECT COUNT(*)
        FROM (
            SELECT rollno, event_code
            FROM register_events
            UNION
            SELECT member_rollno AS rollno, event_code
            FROM requested_events
        ) AS combined
        WHERE rollno = ? AND event_code = ?;
    `
	var count int
	err := config.DB.QueryRow(query, rollno, eventCode).Scan(&count)
	if err != nil {
		log.Println("Error checking applied status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check applied status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applied": count > 0,
	})
}
func DeleteEvent(c *gin.Context){
	id:= c.Param("id")
	query:=`Delete from events where id=?`
	_, err := config.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting event:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
}