
// /api/resume/internship.go

package resume

import (
	"bitresume/config"
	"bitresume/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInternshipData(c *gin.Context) {
	rollno := c.Param("rollno")
	var internships []models.Internship

	rows, err := config.DB.Query("SELECT company_name, domain, start_date, end_date, is_stipend, roll FROM internships WHERE rollno = ?", rollno)
	if err != nil {
		fmt.Print("error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not execute query on internships table"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		// 1. Declare variables to scan into. Dates are simple strings.
		var companyName, domain, role, startDate, endDate string
		var isStipend int

		// 2. Scan all columns. Dates are now read as plain strings.
		err = rows.Scan(&companyName, &domain, &startDate, &endDate, &isStipend, &role)
		if err != nil {
			fmt.Print("Error scanning internship row:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not scan the data from internships table"})
			return
		}

		// 3. Create the struct. The Duration is now a simple combination of the two strings.
		internship := models.Internship{
			CompanyName: companyName,
			Domain:      domain,
			Roll:        role,
			Duration:    startDate + " to " + endDate, // e.g., "2025-05-22 to 2025-05-30"
			IsPaid:      (isStipend == 1),
		}

		internships = append(internships, internship)
	}

	if err = rows.Err(); err != nil {
		fmt.Print("Error during row iteration:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error processing internship data rows"})
		return
	}

	// 4. Send the successful response.
	c.JSON(http.StatusOK, internships)
}
