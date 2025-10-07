
package manageactivities

import (
	"bitresume/config"
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveActivityData(c *gin.Context) {
	// FIX: Use c.PostForm to read data from the request body
	activity_type := c.PostForm("activity_type")
	activity_title := c.PostForm("activity_title")
	if activity_type == "Sessions"{
		activity_type = "Session"
	}
	
	fmt.Println("============================")
	fmt.Println("Activity_type", activity_type)

	query := `
		INSERT INTO activity_list (activity_type, activity_title, created_at)
		VALUES (?, ?, NOW())`

	res, err := config.DB.Exec(query, activity_type, activity_title)

	if err != nil {
		fmt.Println("Error inserting into activity_list:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not insert into the activity_list database"})
		return
	}

	lastid, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	activity_id := int(lastid)

	// This condition will now work correctly
	if activity_type == "Survey" {
		ReceiveSurveyData(c, activity_id)
		return // Ensure we stop here after handling the survey
	}else if activity_type == "Workshop"{
		ReceiveWorkshopData(c, activity_id)
		return
	}else if activity_type == "Meeting"{
		ReceiveMeetingData(c,activity_id)
		return
	}else if activity_type == "Session"{
		ReceiveSessionData(c,activity_id)
		return
	}

	// If no type matched, send a response.
	c.JSON(http.StatusOK, gin.H{"message": "Activity list updated, no specific details handled."})
}

func GetActivityData(c *gin.Context){
	  // Create an empty slice to hold all activities
    allActivities := make([]facultymodel.Activity, 0)

    // Fetch workshops
    workshops, err := GetWrokshopData()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch workshop details"})
        return
    }
    if workshops != nil {
        allActivities = append(allActivities, workshops...)
    }

    // Fetch surveys (assuming you have a GetSurveyData function)
    surveys, err := GetSurveyData() 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch survey details"})
        return
    }
    if surveys != nil {
        allActivities = append(allActivities, surveys...)
    }

	// Fetch meetings
	meetings, err := GetMeetingData()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Could not fetch meeting details"})
		return
	}
	if meetings != nil{
		allActivities = append(allActivities, meetings...)
	}

	sessions, err := GetsessionData()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message" : "Could not fetch Session details"})
		return
	}
	if sessions != nil {
		allActivities = append(allActivities, sessions...)
	}

    c.JSON(http.StatusOK, allActivities)
}   
func HandleProgressGraph(c *gin.Context) {

}