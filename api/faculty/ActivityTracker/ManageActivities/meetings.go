
package manageactivities

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReceiveMeetingData now correctly handles the form data for a "Meeting" activity.
func ReceiveMeetingData(c *gin.Context, activity_id int) {
	// --- Common fields for all activities ---
	publishingDepartment := c.PostForm("publishingDepartment")
	description := c.PostForm("description")
	linkorlocation := c.PostForm("linkorlocation")
	host := c.PostForm("host")
	
	// --- Fields specific to timed events (Meeting, Sessions) ---
	dateOfMeeting := c.PostForm("date_of_meeting") // e.g., "2025-08-20"
	startTime := c.PostForm("start_time")         // e.g., "10:00"
	endTime := c.PostForm("end_time")             // e.g., "12:00"	

	// --- Audience fields ---
	year_type := c.PostForm("year_type")
	target_dept := c.PostForm("target_dept")
	allStudentsStr := c.PostForm("all_students")

	// Convert "true"/"false" string from FormData to a boolean or integer for the DB
	allStudents, err := strconv.ParseBool(allStudentsStr)
	if err != nil {
		// Default to false if parsing fails
		allStudents = false
	}

	// --- Logging for Debugging ---
	// fmt.Println("--- Received Meeting Data ---")
	// fmt.Println("Activity ID:", activity_id)
	// fmt.Println("Publishing Department:", publishingDepartment)
	// fmt.Println("Host:", host)
	// fmt.Println("Description:", description)
	// fmt.Println("Date of Meeting:", dateOfMeeting)
	// fmt.Println("Start Time:", startTime)
	// fmt.Println("End Time:", endTime)
	// fmt.Println("Location/Link:", linkorlocation)
	// fmt.Println("Target Year:", year_type)
	// fmt.Println("Target Department:", target_dept)
	// fmt.Println("All Students:", allStudents)
	// fmt.Println("---------------------------")


	// Your SQL query should match the fields you're actually using.
	// Note: Renamed 'start_time' and 'end_time' to match frontend and be more descriptive.
	query := `
		INSERT INTO meeting_details(
			activity_id,
			publishing_department,
			host,
			description,
			start_time,          -- Storing time as string e.g., "10:00"
			end_time,            -- Storing time as string e.g., "12:00"
			date_of_meeting,     -- Storing date as string e.g., "2025-08-20"
			link_or_location,
			target_year,
			target_department,
			all_students
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`


	// Execute the query with the correct variables
	_, dbErr := config.DB.Exec(query,
		activity_id,
		publishingDepartment,
		host,
		description,
		startTime, // Use the time-only variable
		endTime,   // Use the time-only variable
		dateOfMeeting,
		linkorlocation,
		year_type,
		target_dept,
		allStudents, // Use the parsed boolean/int
	)

	if dbErr != nil {
		fmt.Println("Database Error:", dbErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not insert meeting details into the database.", "error": dbErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Meeting activity created successfully."})
}

func GetMeetingData()([]facultymodel.Activity,error){
	var meeting []facultymodel.Activity

	rows, err := config.DB.Query(`SELECT 
    	al.activity_title,
    	al.activity_type,
        s.host,
    	s.description,
    	s.start_time,
    	s.end_time,
        s.date_of_meeting,
    	s.link_or_location,
    	s.target_year,
    	s.all_students
	FROM 
    	meeting_details AS s
	INNER JOIN 
    	activity_list AS al ON al.id = s.activity_id;`)

	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var m facultymodel.Activity
		err = rows.Scan(&m.ActivityTitle,&m.ActivityType,&m.Host,&m.Description,&m.StartDate,&m.EndDate,&m.DateofMeeting,&m.LinkOrLocation,&m.TargetYear,&m.AllStudents)
		if err != nil{
			fmt.Println("Error: ", err.Error())
			return nil, err
		}
		meeting = append(meeting, m)
	}

	return meeting, nil
}
