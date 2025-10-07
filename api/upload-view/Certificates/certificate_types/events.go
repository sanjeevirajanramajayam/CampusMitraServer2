package certificatetypes

import (
	"bitresume/config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ReceiveEventsData(c *gin.Context,id int) {
	rollno := c.PostForm("rollno")
	eventName := c.PostForm("event_name")
	eventCode := c.PostForm("event_code")
	issue_date := c.PostForm("issue_date")
	summary := c.PostForm("summary")
	participationType := c.PostForm("participation_type")
	didYouWin := c.PostForm("did_you_win")

	// Retrieve the uploaded certificate file
	certificateFile, err := c.FormFile("certificate_pdf")
	if err != nil {
		c.JSON(400, gin.H{"error": "Certificate PDF is required"})
		return
	}

	// Define save path
	saveDir := "uploads/certificates/events"
	savePath := filepath.Join(saveDir, certificateFile.Filename)

	// Ensure the folder exists
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create directory for saving certificate"})
		return
	}

	// Save the uploaded file
	if err := c.SaveUploadedFile(certificateFile, savePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save certificate file"})
		return
	}

	// Insert into the database
	query := `
		INSERT INTO certificates_events (
			certificate_id,
			rollno,
			event_name,
			event_code,
			issue_date,
			participation_type,
			certificate_pdf,
			summary,
			did_you_win,
			faculty_name,
			faculty_id,
			faculty_remarks,
			submission_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_DATE)
	`

	fmt.Println("Certificate_id: ",id)

	_, err = config.DB.Exec(query,
		id,
		rollno,
		eventName,
		eventCode,
		issue_date,
		participationType,
		savePath,
		summary,
		didYouWin,
		"",     // faculty_name (to be filled later)
		"",     // faculty_id
		"",     // faculty_remarks
	)

	if err != nil {
		fmt.Println("Error: ", err.Error())
		c.JSON(500, gin.H{"error": "Failed to insert certificate data into database"})
		return
	}

	c.JSON(200, gin.H{"message": "Certificate uploaded successfully"})
}