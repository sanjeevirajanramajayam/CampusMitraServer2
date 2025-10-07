
package resume

import (
	"bitresume/config"
	"bitresume/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCertificatesData(c *gin.Context) {
	rollno := c.Param("rollno")
	var certificates []models.Certificates
	rows, err := config.DB.Query("select event_name from certificates_events where rollno = ?", rollno)
	if err != nil {
		fmt.Print("error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the data from certificates_events table"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var p models.Certificates
		err = rows.Scan(&p.Title)
		if err != nil {
			fmt.Print("Error:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not scan the data from certificates_events table"})
			return
		}
		certificates = append(certificates, p)
	}
	rows1, err := config.DB.Query("select title from certificate_onlinecourses where rollno = ?", rollno)
	if err != nil {
		fmt.Print("error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch the data from certificate_onlinecourses table"})
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		var p models.Certificates
		err = rows1.Scan(&p.Title)
		if err != nil {
			fmt.Print("Error:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not scan the data from certificate_onlinecourses table"})
			return
		}
		certificates = append(certificates, p)
	}

	c.JSON(http.StatusAccepted, certificates)
}
