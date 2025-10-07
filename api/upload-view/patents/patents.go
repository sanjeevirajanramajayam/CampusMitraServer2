package patents

import (
	"bitresume/config"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

const maxPdfSize = 5 * 1024 * 1024 // 5 MB

func ReceivePatentsData(c *gin.Context) {
	rollno := c.PostForm("rollno")
	title := c.PostForm("title")
	application_number := c.PostForm("application_number")
	date_of_filing := c.PostForm("date_of_filing")
	link_to_patent_listing := c.PostForm("link_to_patent_listing")
	summary := c.PostForm("summary")
	usecase_of_patent := c.PostForm("usecase_of_patent")

	// Create directories if they don't exist
	if err := os.MkdirAll("uploads/patents/patent_docs", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "Could not create directory for patent_docs", "details": err.Error()})
		return
	}
	if err := os.MkdirAll("uploads/patents/supporting_files", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "Could not create directory for supporting_files", "details": err.Error()})
		return
	}

	// Handle patent_docs file
	patent_docs, err := c.FormFile("patent_docs")
	if err != nil {
		c.JSON(400, gin.H{"error": "patent_docs file is missing or invalid", "details": err.Error()})
		return
	}
	if patent_docs.Size > maxPdfSize {
		c.JSON(400, gin.H{"error": "patent_docs file size exceeds 5MB"})
		return
	}
	savePathPatentDocs := filepath.Join("uploads/patents/patent_docs", patent_docs.Filename)
	if err := c.SaveUploadedFile(patent_docs, savePathPatentDocs); err != nil {
		c.JSON(500, gin.H{"error": "Could not save patent_docs file", "details": err.Error()})
		return
	}

	supporting_files, _ := c.FormFile("supporting_files")
	var savePathSupportingFiles string
	if supporting_files != nil {
		if supporting_files.Size > maxPdfSize {
			c.JSON(400, gin.H{"error": "supporting_files file size exceeds 5MB"})
			return
		}
		savePathSupportingFiles = filepath.Join("uploads/patents/supporting_files", supporting_files.Filename)
		if err := c.SaveUploadedFile(supporting_files, savePathSupportingFiles); err != nil {
			c.JSON(500, gin.H{"error": "Could not save supporting_files", "details": err.Error()})
			return
		}
	}

	uploadType := "patents"


	// Insert into database
	query := `
	INSERT INTO patents
	(upload_type,rollno, title, application_number, date_of_filing, patent_docs, supporting_files, link_to_patent_listing, summary, usecase_of_patent, faculty_remarks, patent_status, submission_date)
	VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_DATE)`

	_, err = config.DB.Exec(query,uploadType,
		rollno,
		title,
		application_number,
		date_of_filing,
		savePathPatentDocs,
		savePathSupportingFiles,
		link_to_patent_listing,
		summary,
		usecase_of_patent,
		"",
		"Pending",
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload to database", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Patent successfully uploaded"})
}