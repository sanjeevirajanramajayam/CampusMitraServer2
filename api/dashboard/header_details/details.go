package headerdetails

import (
	"bitresume/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchDataRank(c *gin.Context) {
	rollno := c.Param("rollno")
	type Response struct {
		CurrentRank        string `json:"current_rank"`
		TotalPoints        float64 `json:"total_points"`
		PositivePointsCount int `json:"positive_points_count"`
		PenaltyCount        int `json:"penalty_count"`
	}
	var res Response
	// 1. Fetch latest rank
	err := config.DB.QueryRow("SELECT current_rank FROM activity_graph WHERE rollno = ? ORDER BY currdate DESC LIMIT 1", rollno).
		Scan(&res.CurrentRank)
	if err != nil {
		log.Printf("Error fetching rank: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching rank"})
		return
	}
	// 2. Fetch total achievement points
	err = config.DB.QueryRow("SELECT cummulative_points FROM achievement_graph WHERE rollno = ? ORDER BY currdate DESC LIMIT 1", rollno).
		Scan(&res.TotalPoints)
	if err != nil {
		log.Printf("Error fetching total points: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching total points"})
		return
	}
	// 3. Count of positive points
	err = config.DB.QueryRow("SELECT COUNT(*) FROM points_logs WHERE rollno = ? AND points >= 0", rollno).
		Scan(&res.PositivePointsCount)
	if err != nil {
		log.Printf("Error fetching positive points count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching positive points count"})
		return
	}
	// 4. Count of penalties
	err = config.DB.QueryRow("SELECT COUNT(*) FROM points_logs WHERE rollno = ? AND points < 0", rollno).
		Scan(&res.PenaltyCount)
	if err != nil {
		log.Printf("Error fetching penalty count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching penalty count"})
		return
	}

	// Final JSON response
	c.JSON(http.StatusOK, res)
}

