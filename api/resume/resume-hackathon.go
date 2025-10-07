package resume

import (
	"bitresume/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHackathonData(c *gin.Context) {
	rollno := c.Param("rollno")
	query := `
        SELECT e.image_url, e.event_name, ce.did_you_win
        FROM certificates_events ce
        JOIN events e ON ce.event_code = e.event_code
        WHERE rollno = ?;
    `

	rows, err := config.DB.Query(query, rollno)
	if err != nil {
		fmt.Println("error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not execute query on hackathon table"})
		return
	}
	defer rows.Close()
	// define slice
	var results []struct {
		ImgUrl    string `json:"img_url"`
		EventName string `json:"event_name"`
		DidYouWin string `json:"did_you_win"`
	}

	for rows.Next() {
		var result struct {
			ImgUrl    string `json:"img_url"`
			EventName string `json:"event_name"`
			DidYouWin string `json:"did_you_win"`
		}

		var didYouWinInt string
		if err := rows.Scan(&result.ImgUrl, &result.EventName, &didYouWinInt); err != nil {
			c.JSON(500, gin.H{"error": "Failed to scan row", "details": err.Error()})
			fmt.Print(err)
			return
		}

		// convert int → string if needed
		if didYouWinInt == "Winner"{
			result.DidYouWin = "Yes"
		} else {
			result.DidYouWin = "No"
		}

		results = append(results, result)
	}

	// check if no rows found
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No hackathon data found"})
		return
	}

	c.JSON(http.StatusOK, results)
}
