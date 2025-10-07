package certificatetypes

import (
	"bitresume/config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ReceiveVoluntreeData(c *gin.Context,id int){
	rollno := c.PostForm("rollno")
	certificate_pdf, err := c.FormFile("certificate_pdf")
	issue_date := c.PostForm("issue_date")
	summary := c.PostForm("summary")
	activity_type := c.PostForm("activity_type")
	duration := c.PostForm("duration")
	location := c.PostForm("location")

	if err != nil {
		c.JSON(500 , "could not get the pdf")
		return
	}

	savePathPdf := filepath.Join("uploads/certificates/participation",certificate_pdf.Filename)

	if err := os.MkdirAll("uploads/certificates/participation",os.ModePerm); err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(500, "could not find the file directory")
		return
	}

	if err := c.SaveUploadedFile(certificate_pdf,savePathPdf); err!= nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(500, "Could not save the file")
		return
	}

	query := `
		insert into certificates_voluntree
		(
			certificate_id,
			rollno,
			activity_type,
			duration,
			issue_date,
			certificate_pdf,
			summary,
			location,
			faculty_name,
			faculty_id,
			faculty_reamrks,
			submission_date
		) values (?,?,?,?,?,?,?,?,?,?,?, current_date)
	`

	_,err = config.DB.Exec(query,id,rollno,activity_type,duration,issue_date,savePathPdf,summary,location,"","","")

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(500, "Could not insert into db")
		return
	}

	c.JSON(200, "Success")
}