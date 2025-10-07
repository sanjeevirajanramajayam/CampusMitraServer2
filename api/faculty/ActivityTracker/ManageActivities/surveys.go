package manageactivities

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveSurveyData(c *gin.Context, activity_id int) {
	// FIX: Use c.PostForm for all fields
	publishingDepartment := c.PostForm("publishingDepartment")
	description := c.PostForm("description")
	linkorlocation := c.PostForm("linkorlocation")
	start_date := c.PostForm("start_date")
	end_date := c.PostForm("end_date")
	year_type := c.PostForm("year_type")
	target_dept := c.PostForm("target_dept")
	all_students := c.PostForm("all_students")

	allStudents := 0
	if all_students == "1" {
		allStudents = 1
	}

	// The column names in your query match the frontend FormData keys perfectly.
	query := `
		INSERT INTO survey_details (
			activity_id,
			publishing_department,
			description,
			start_date,
			end_date,
			link_or_location,
			target_year,
			target_department,
			all_students
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query,activity_id, publishingDepartment, description, start_date, end_date, linkorlocation, year_type, target_dept, allStudents)
	if err != nil {
		fmt.Println("Database insert error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert survey data into the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey data successfully inserted into the database"})
}

func GetSurveyData()([]facultymodel.Activity, error){
	var survey []facultymodel.Activity

	rows, err := config.DB.Query(`SELECT 
    	al.activity_title,
    	al.activity_type,
    	s.description,
    	s.start_date,
    	s.end_date,
    	s.link_or_location,
    	s.target_year,
    	s.all_students
	FROM 
    	survey_details AS s
	INNER JOIN 
    	activity_list AS al ON al.id = s.activity_id;`)


	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var r facultymodel.Activity
		err = rows.Scan(&r.ActivityTitle,&r.ActivityType,&r.Description,&r.StartDate,&r.EndDate,&r.LinkOrLocation,&r.TargetYear,&r.AllStudents)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return nil , err
		} 
		survey = append(survey, r)
	}

	return survey, nil
}