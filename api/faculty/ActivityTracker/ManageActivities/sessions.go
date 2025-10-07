package manageactivities

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ReceiveSessionData handles storing session details and specific student roll numbers
func ReceiveSessionData(c *gin.Context, activity_id int) {
	// Step 1: Collect all form values from the request
	publishingDepartment := c.PostForm("publishingDepartment")
	description := c.PostForm("description")
	linkorlocation := c.PostForm("linkorlocation")
	date_of_meeting := c.PostForm("date_of_meeting")
	start_time := c.PostForm("start_time")
	end_time := c.PostForm("end_time")
	session_with := c.PostForm("session_with")
	specific_rollno := c.PostForm("specific_rollno")

	// Step 2: Convert comma-separated roll numbers into a slice
	SpecificRollnos := strings.Split(specific_rollno, ",")

	// Step 3: Insert session details into the session_details table
	sessionQuery := `
		INSERT INTO session_details (
			activity_id,
			publishing_department,
			host,
			description,
			start_time,
			end_time,
			date_of_session,
			link_or_location
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := config.DB.Exec(sessionQuery, activity_id, publishingDepartment, session_with, description, start_time, end_time, date_of_meeting, linkorlocation)
	if err != nil {
		fmt.Println("Error inserting into session_details:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session details"})
		return
	}

	// Step 4: Insert each student roll number into session_students table
	rollQuery := `
		INSERT INTO session_students (
			activity_id,
			student_rollno
		) VALUES (?, ?);
	`

	for _, rollno := range SpecificRollnos {
		// Trim whitespace in case roll numbers have spaces after commas
		rollno = strings.TrimSpace(rollno)

		if rollno == "" {
			continue // Skip empty entries
		}

		_, err := config.DB.Exec(rollQuery, activity_id, rollno)
		if err != nil {
			fmt.Printf("Error inserting roll number '%s' into session_students: %s\n", rollno, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save student roll numbers"})
			return
		}
	}

	// Step 5: Send success response back to frontend
	c.JSON(http.StatusOK, gin.H{"message": "Session data successfully saved"})
}


func GetsessionData() ([]facultymodel.Activity, error) {
	// Slice to store all fetched session records
	var sessions []facultymodel.Activity

	// Step 1: Execute SQL query to get session details joined with activity_list
	rows, err := config.DB.Query(`
		SELECT 
    al.activity_title,
    al.activity_type,
    s.host,
    s.description,
    s.start_time,
    s.end_time,
    s.date_of_session,
    s.link_or_location,
    ss.student_rollno
FROM 
    session_details AS s
INNER JOIN 
    activity_list AS al ON al.id = s.activity_id
INNER JOIN 
    session_students AS ss ON ss.activity_id = s.activity_id;
	`)
	if err != nil {
		// Log the error and return
		fmt.Println(" Error fetching session data from database:", err)
		return nil, fmt.Errorf("failed to fetch session data: %v", err)
	}
	defer rows.Close() // Ensure the DB rows are closed properly after use

	// Step 2: Loop through the results and scan each row into the Activity struct
	for rows.Next() {
		var r facultymodel.Activity

		err = rows.Scan(
			&r.ActivityTitle,
			&r.ActivityType,
			&r.Host,
			&r.Description,
			&r.StartDate,
			&r.EndDate,
			&r.DateofMeeting,
			&r.LinkOrLocation,
			&r.SpecificRollno,
		)
		if err != nil {
			fmt.Println(" Error scanning session row:", err)
			return nil, fmt.Errorf("failed to scan session data: %v", err)
		}

		sessions = append(sessions, r) // Add to result slice
	}

	// Step 3: Return the complete slice of session data
	return sessions, nil
} 