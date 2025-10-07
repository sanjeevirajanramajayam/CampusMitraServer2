package Uploadsdelete

import (
	"bitresume/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteUploadRequest struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subType,omitempty"` // optional
}

func Uploadsdelete(c *gin.Context) {
	var req DeleteUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	switch req.Type {
	case "Paper Presentation":
		_, err := config.DB.Exec("DELETE FROM paperpresentation WHERE id = ?", req.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
			return
		}

	case "Patent":
		_, err := config.DB.Exec("DELETE FROM patents WHERE id = ?", req.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
			return
		}
	
	case "Internship":
		_, err := config.DB.Exec("DELETE FROM internships WHERE id = ?", req.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
			return
		}
	
	case "Project":
		_, err := config.DB.Exec("DELETE FROM projects WHERE id = ?", req.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
			return
		}
	
	case "Seminar / Workshop":
		_, err := config.DB.Exec("DELETE FROM workshops WHERE id = ?", req.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
			return
		}

	case "Certificate":
		switch req.SubType {
			case "online-course":
				_, err := config.DB.Exec("Delete from certificate_onlinecourses where certiificate_id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}

				_, err = config.DB.Exec("Delete from certificates_type where id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}
				
			case "hackathon":
				_, err := config.DB.Exec("Delete from certificates_events where certificate_id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}

				_, err = config.DB.Exec("Delete from certificates_type where id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}
				
			case "participation":
				_, err := config.DB.Exec("Delete from certificates_voluntree where certificate_id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}

				_, err = config.DB.Exec("Delete from certificates_type where id = ?", req.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete", "error": err.Error()})
					return
				}
				
			default:
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type"})
				return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
