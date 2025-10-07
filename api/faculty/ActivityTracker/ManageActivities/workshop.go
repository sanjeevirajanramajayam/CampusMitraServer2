package manageactivities

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveWorkshopData(c *gin.Context, activity_id int) {
	publishingDepartment := c.PostForm("publishingDepartment")
	description := c.PostForm("description")
	linkorlocation := c.PostForm("linkorlocation")
	start_date := c.PostForm("start_date")
	end_date := c.PostForm("end_date")
	host := c.PostForm("host")
	year_type := c.PostForm("year_type")
	target_dept := c.PostForm("target_dept")
	all_students := c.PostForm("all_students")

	query := `
		insert into workshop_details(
			activity_id,
			publishing_department,
			host,
			description,
			start_date,
			end_date,
			link_or_location,
			target_year,
			target_department,
			all_students
		) values (?,?,?,?,?,?,?,?,?,?)
	`

	_, err := config.DB.Exec(query, activity_id, publishingDepartment, host, description, start_date, end_date, linkorlocation, year_type, target_dept, all_students)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not insert into wrokshop_details"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully inserted into database"})
}

func GetWrokshopData()([]facultymodel.Activity, error){
	var workshops []facultymodel.Activity

	rows, err := config.DB.Query(`SELECT 
    			al.activity_title,
    			al.activity_type,
        		s.host,
    			s.description,
    			s.start_date,
    			s.end_date,
    			s.link_or_location,
    			s.target_year,
    			s.all_students
				FROM 
    			workshop_details AS s
				INNER JOIN 
    			activity_list AS al ON al.id = s.activity_id;`)

	if err != nil{
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var r facultymodel.Activity
		err = rows.Scan(&r.ActivityTitle,&r.ActivityType,&r.Host,&r.Description,&r.StartDate,&r.EndDate,&r.LinkOrLocation,&r.TargetYear,&r.AllStudents)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return nil, err
		}

		workshops = append(workshops, r)
	}

	return workshops, nil
}