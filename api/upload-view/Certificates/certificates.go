package certificates

import (
	certificatetypes "bitresume/api/upload-view/Certificates/certificate_types"
	"bitresume/config"
	"github.com/gin-gonic/gin"
)

func ReceiveCertificateData(c *gin.Context) {
	certificateType := c.PostForm("certificate_type")
	rollno := c.PostForm("rollno")
	uploadType := "certificate"
	query := `
		INSERT INTO certificates_type (
			upload_type,
			rollno,
			certificate_type
		) VALUES (?, ?, ?)
	`

	res , err := config.DB.Exec(query,uploadType, rollno, certificateType)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save certificate type to the database"})
		return
	}
	lastID, _ := res.LastInsertId()
	certificate_id := int(lastID)

	if certificateType == "online-course" {
		certificatetypes.ReceiveDataOnlineCourse(c,certificate_id)
		return
	}else if certificateType == "hackathon"{
		certificatetypes.ReceiveEventsData(c,certificate_id)
		return
	}else if certificateType == "participation"{
		certificatetypes.ReceiveVoluntreeData(c,certificate_id)
		return
	}

	c.JSON(400, gin.H{"error": "Unsupported certificate type"})
}