package internship

import (
	"bitresume/config"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

const maxPdfSize = 5 * 1024 * 1024 // 5MB

func ReceiveInternshipData(c *gin.Context) {
	rollno := c.PostForm("rollno")
	company_name := c.PostForm("company_name")
	roll := c.PostForm("roll")
	domain := c.PostForm("domain")
	internship_type := c.PostForm("internship_type")
	is_stipend_str := c.PostForm("is_stipend")
	start_date := c.PostForm("start_date")
	end_date := c.PostForm("end_date")
	consulted_faculty_name := c.PostForm("consulted_faculty_name")
	industry_mentor_name := c.PostForm("industry_mentor_name")
	industry_mentor_contact := c.PostForm("industry_mentor_contact")
	skill_gained := c.PostForm("skill_gained")
	outcomes := c.PostForm("outcomes")

	// Convert stipend string to int (0 or 1)
	is_stipend, err := strconv.Atoi(is_stipend_str)
	if err != nil || (is_stipend != 0 && is_stipend != 1) {
		c.JSON(400, gin.H{"error": "Invalid value for is_stipend, expected 0 or 1"})
		return
	}

	// Handle offer letter
	offer_letter, err := c.FormFile("offer_letter")
	if err != nil {
		c.JSON(400, gin.H{"error": "Missing or invalid offer_letter", "details": err.Error()})
		return
	}

	if err := os.MkdirAll("uploads/internships/offer_letter", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "Could not create directory for offer_letter", "details": err.Error()})
		return
	}

	savePathOfferLetter := filepath.Join("uploads/internships/offer_letter", offer_letter.Filename)
	if err := c.SaveUploadedFile(offer_letter, savePathOfferLetter); err != nil {
		c.JSON(500, gin.H{"error": "Could not save offer_letter", "details": err.Error()})
		return
	}

	// Handle report (optional)
	var savePathReport string
	report, err := c.FormFile("report")
	if err == nil && report != nil {
		if report.Size > maxPdfSize {
			c.JSON(400, gin.H{"error": "report file size exceeds 5MB"})
			return
		}

		if err := os.MkdirAll("uploads/internships/report", os.ModePerm); err != nil {
			c.JSON(500, gin.H{"error": "Could not create directory for report", "details": err.Error()})
			return
		}

		savePathReport = filepath.Join("uploads/internships/report", report.Filename)
		if err := c.SaveUploadedFile(report, savePathReport); err != nil {
			c.JSON(500, gin.H{"error": "Could not save report", "details": err.Error()})
			return
		}
	}

	uploadType := "internship"

	// Insert into database
	query := `
		INSERT INTO internships (
			upload_type,
			rollno,
			company_name,
			roll,
			domain,
			internship_type,
			is_stipend,
			start_date,
			end_date,
			consulted_faculty_name,
			industry_mentor_name,
			industry_mentor_contact,
			offer_letter,
			report,
			faculty_remarks,
			skill_gained,
			outcomes,
			submitted_on
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_DATE)
	`

	_, err = config.DB.Exec(
		query,
		uploadType,
		rollno,
		company_name,
		roll,
		domain,
		internship_type,
		is_stipend,
		start_date,
		end_date,
		consulted_faculty_name,
		industry_mentor_name,
		industry_mentor_contact,
		savePathOfferLetter,
		savePathReport,
		"", // faculty_remarks default
		skill_gained,
		outcomes,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Could not insert into database", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Internship successfully submitted"})
}