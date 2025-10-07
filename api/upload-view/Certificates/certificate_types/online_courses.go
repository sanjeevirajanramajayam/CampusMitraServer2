package certificatetypes

import (
	"bitresume/config"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)



func ReceiveDataOnlineCourse(c *gin.Context,id int) {
	rollno := c.PostForm("rollno")
	title := c.PostForm("title")
	platform := c.PostForm("platform")
	issue_date := c.PostForm("issue_date")
	start_date := c.PostForm("start_date")
	end_date := c.PostForm("end_date")
	course_link := c.PostForm("course_link")

	pdf, err := c.FormFile("certificate_pdf")
	if err != nil {
		c.JSON(400, gin.H{"error": "Certificate PDF file is missing or unreadable"})
		return
	}

	const folderPath = "uploads/certificates/online_courses"
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create directory for certificate storage"})
		return
	}

	savePathPdf := filepath.Join(folderPath, pdf.Filename)
	if err := c.SaveUploadedFile(pdf, savePathPdf); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save uploaded certificate PDF"})
		return
	}

	query := `
		INSERT INTO certificate_onlinecourses (
			certiificate_id,
			rollno,
			title,
			platform,
			issue_date,
			start_date,
			end_date,
			certificate_pdf,
			course_link,
			created_at
		) VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, CURRENT_DATE)
	`

	_, err = config.DB.Exec(query,id, rollno, title, platform, issue_date, start_date, end_date, savePathPdf, course_link)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save online course certificate to the database"})
		return
	}

	c.JSON(200, gin.H{"message": "Online course certificate uploaded successfully"})
}