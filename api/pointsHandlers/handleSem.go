package pointshandlers

import (
	"bitresume/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Existing HandleSem function to fetch semester data
func HandleSem(c *gin.Context) {
	query := `SELECT
    year,
		MAX(batch) AS batch,
		MAX(sem) AS sem
	FROM bitresume.login
	WHERE year IS NOT NULL
	AND batch IS NOT NULL
	AND sem IS NOT NULL
	GROUP BY year
	ORDER BY year;
	`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed", "details": err.Error()})
		return
	}
	defer rows.Close()

	type SemData struct {
		Year  string `json:"year"`
		Batch string `json:"batch"`
		Sem   int    `json:"sem"`
	}

	var results []SemData
	for rows.Next() {
		var sd SemData
		if err := rows.Scan(&sd.Year, &sd.Batch, &sd.Sem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse database result", "details": err.Error()})
			return
		}
		results = append(results, sd)
	}

	c.JSON(http.StatusOK, results)
}

// New HandleUpdateSem function to edit the semester for a batch
func HandleUpdateSem(c *gin.Context) {
	type UpdateRequest struct {
		Batch string `json:"batch"`
		Sem   int    `json:"sem"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	query := `UPDATE bitresume.login SET sem = ? WHERE batch = ?`
	result, err := config.DB.Exec(query, req.Sem, req.Batch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed", "details": err.Error()})
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected", "details": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No records found for the given batch to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Semester updated successfully for batch " + req.Batch, "rows_affected": rowsAffected})
}