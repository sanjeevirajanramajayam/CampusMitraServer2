
package activitymaster

import (
	"bitresume/config"
	activitymastermodels "bitresume/models/Activitymaster"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSurveys(c *gin.Context) {
	rollno := c.Param("rollno")
	var surveys []activitymastermodels.SurveyDetail
	fmt.Println("rollno:",rollno)

	QUERY := `
		SELECT 
			sd.activity_id,
			sd.publishing_department,
			sd.description,
			sd.start_date,
			sd.end_date,
			sd.link_or_location,
			sd.target_year,
			sd.target_department,
			sd.all_students
		FROM survey_details sd
		JOIN login l 
			ON l.rollno = ?
		WHERE 
			(
				sd.all_students = 1
				OR (
					(sd.target_year = l.year OR sd.target_year = 'All Years')
					AND (sd.target_department = l.department OR sd.target_department = 'All Departments')
				)
			)
			AND l.role = 'student';
	`

	rows, err := config.DB.Query(QUERY, rollno)
	if err != nil {
		fmt.Println("The ERROR is:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the data from the database"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r activitymastermodels.SurveyDetail
		err = rows.Scan(
			&r.ActivityID,
			&r.PublishingDept,
			&r.Description,
			&r.StartDate,
			&r.EndDate,
			&r.LinkOrLocation,
			&r.TargetYear,
			&r.TargetDepartment,
			&r.AllStudents,
		)
		if err != nil {
			fmt.Println("The ERROR while scanning is:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error scanning the data"})
			return
		}
		surveys = append(surveys, r)
	}


	// if surveys == nil {
	// 	surveys = []activitymastermodels.SurveyDetail{}
	// }

	c.JSON(http.StatusOK, surveys)
}