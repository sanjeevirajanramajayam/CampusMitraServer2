package activitymaster

import (
	"bitresume/config"
	activitymastermodels "bitresume/models/Activitymaster"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSessionsByRollNo(c *gin.Context) {
	rollno := c.Param("rollno")

	// -------------------- Fetch Sessions --------------------
	sessionQuery := `
		SELECT 
			s.activity_id,
			s.publishing_department,
			s.host,
			s.description,
			s.start_time,
			s.end_time,
			s.date_of_session,
			s.link_or_location
		FROM session_details s
		JOIN session_students ss
			ON s.activity_id = ss.activity_id
		WHERE ss.student_rollno = ?;
	`

	sessionRows, err := config.DB.Query(sessionQuery, rollno)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query sessions"})
		return
	}
	defer sessionRows.Close()

	var sessions []activitymastermodels.SessionDetails
	for sessionRows.Next() {
		var session activitymastermodels.SessionDetails
		if err := sessionRows.Scan(
			&session.ActivityID,
			&session.PublishingDept,
			&session.Host,
			&session.Description,
			&session.StartTime,
			&session.EndTime,
			&session.DateOfSession,
			&session.LinkOrLocation,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan session"})
			return
		}
		sessions = append(sessions, session)
	}
	if err := sessionRows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating sessions"})
		return
	}

	// -------------------- Fetch Eligible Meetings --------------------
	meetingQuery := `
		SELECT 
			md.activity_id,
			md.publishing_department,
			md.host,
			md.description,
			md.start_time,
			md.end_time,
			md.date_of_meeting,
			md.link_or_location,
			md.target_year,
			md.target_department,
			md.all_students
		FROM meeting_details md
		JOIN login l 
			ON (
				md.all_students = 1
				OR (md.target_year = l.year AND md.target_department = l.department)
				OR (md.target_year = l.year AND md.target_department = 'All Departments')
				OR (md.target_year = 'All Years' AND md.target_department = l.department)
			)
		WHERE l.rollno = ?;
	`

	meetingRows, err := config.DB.Query(meetingQuery, rollno)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query meetings"})
		return
	}
	defer meetingRows.Close()

	var meetings []activitymastermodels.EligibleMeeting
	for meetingRows.Next() {
		var meeting activitymastermodels.EligibleMeeting
		if err := meetingRows.Scan(
			&meeting.ActivityID,
			&meeting.PublishingDept,
			&meeting.Host,
			&meeting.Description,
			&meeting.StartTime,
			&meeting.EndTime,
			&meeting.DateOfMeeting,
			&meeting.LinkOrLocation,
			&meeting.TargetYear,
			&meeting.TargetDepartment,
			&meeting.AllStudents,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan meeting"})
			return
		}
		meetings = append(meetings, meeting)
	}
	if err := meetingRows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating meetings"})
		return
	}

	// -------------------- Final JSON Response --------------------
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"meetings": meetings,
	})
}