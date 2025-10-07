package workshops

import (
	"bitresume/config"
	"os"
	"path/filepath"
	// "strconv" // Still needed if you want to support "1"/"0" for boolean, though direct string check is better

	"github.com/gin-gonic/gin"
)

const maxSize = 5 * 1024 * 1024 // 5MB

func ReceiveWorkshopData(c *gin.Context) {
	// Corrected to use snake_case for c.PostForm to match frontend payload
	rollno := c.PostForm("rollno")
	eventTitle := c.PostForm("event_title")
	eventType := c.PostForm("event_type")
	deliveryMode := c.PostForm("delivery_mode")
	eventNature := c.PostForm("event_nature")
	organizedBy := c.PostForm("organized_by")
	location := c.PostForm("location") // 'location' was already correct as per debug log
	eventStartDate := c.PostForm("event_start_date")
	eventEndDate := c.PostForm("event_end_date")
	participationRole := c.PostForm("participation_role")
	isCertificateProvidedStr := c.PostForm("is_certificate_provided")
	eventLink := c.PostForm("event_link")
	topicsCovered := c.PostForm("topics_covered")
	skillsGained := c.PostForm("skills_gained")
	relevance := c.PostForm("relevance") // 'relevance' was already correct

	// Convert isCertificateProvided string ("true", "false", "1", "0") to int (TINYINT compatible)
	isCert := 0 // Default to 0 (false)
	if isCertificateProvidedStr == "true" || isCertificateProvidedStr == "1" {
		isCert = 1
	}

	// Handle optional certificate upload
	var savePathCertificate string
	certificateFile, err := c.FormFile("certificate") // Frontend sends field named "certificate"
	if err == nil && certificateFile != nil {
		if certificateFile.Size > maxSize {
			c.JSON(400, gin.H{"error": "Certificate file size exceeds 5MB limit"})
			return
		}
		uploadDir := "uploads/workshops"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create directory for certificate", "details": err.Error()})
			return
		}
		savePathCertificate = filepath.Join(uploadDir, certificateFile.Filename)
		if err := c.SaveUploadedFile(certificateFile, savePathCertificate); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save certificate file", "details": err.Error()})
			return
		}
	} else if err != nil && err.Error() != "http: no such file" {
		// Handle other errors from FormFile, if not just "no such file"
		c.JSON(500, gin.H{"error": "Error processing certificate file", "details": err.Error()})
		return
	}

	// Use nil for certificate path if it's empty, to insert SQL NULL
	var certificateValueToSave interface{}
	if savePathCertificate == "" {
		certificateValueToSave = nil
	} else {
		certificateValueToSave = savePathCertificate
	}

	upload_type := "workshop"

	// Prepare SQL insert
	// Omit `submitted_on` from the insert list to use the database's DEFAULT CURRENT_TIMESTAMP
	// Ensure `relevence` matches the column name in your DB schema (it has 'e' in schema image).
	query := `
	INSERT INTO workshops (
		upload_type,
		rollno,
		title,
		event_type,
		mode_of_delivery,
		event_nature,
		organised_by,
		location,
		start_date,
		end_date,
		participation_type,
		is_certificate,
		certificate,
		link,
		topics_covered,
		relevence, 
		skills_gained
		-- submitted_on is omitted to use DB default
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
` // 16 placeholders

	// Execute insert
	_, err = config.DB.Exec(
		query,
		upload_type,
		rollno,        
		eventTitle,            
		eventType,             
		deliveryMode,          
		eventNature,           
		organizedBy,           
		location,              
		eventStartDate,        
		eventEndDate,          
		participationRole,     
		isCert,
		certificateValueToSave,
		eventLink,             
		topicsCovered,         
		relevance,             
		skillsGained,          
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Database insertion failed", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Workshop data successfully submitted"})
}