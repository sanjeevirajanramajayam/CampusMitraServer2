package routes

import (
	activitymaster "bitresume/api/ActivityMaster"
	auth "bitresume/api/auth"
	achievementgraph "bitresume/api/dashboard/achievement_graph"
	activitygraph "bitresume/api/dashboard/activity_graph"
	headerdetails "bitresume/api/dashboard/header_details"
	manageactivities "bitresume/api/faculty/ActivityTracker/ManageActivities"
	studentrequests "bitresume/api/faculty/ActivityTracker/StudentRequests/varifications"
	addevents "bitresume/api/faculty/AddEvents"
	studentdata "bitresume/api/faculty/StudentData"
	dashBoardfaculty "bitresume/api/faculty/dashboardfaculty"
	naacreports "bitresume/api/naac-reports"
	pointshandlers "bitresume/api/pointsHandlers"
	registerevents "bitresume/api/registerEvents"
	"bitresume/api/resume"
	certificates "bitresume/api/upload-view/Certificates"
	Uploadsdelete "bitresume/api/upload-view/delete"
	"bitresume/api/upload-view/internship"
	"bitresume/api/upload-view/paperpresentstion"
	"bitresume/api/upload-view/patents"
	"bitresume/api/upload-view/projects"
	dashboard "bitresume/api/upload-view/upload_view_dashboard"
	"bitresume/api/upload-view/workshops"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Auth routes
	r.GET("/api/auth/google/login", auth.GoogleLogin)
	r.GET("/api/auth/google/callback", auth.GoogleCallback)
	r.GET("/api/auth/me", auth.Me)
	r.GET("/api/auth/logout", auth.Logout)

	// Student routes (previously protected)
	r.POST("/api/points_logs/ps/attempts", pointshandlers.HandlePs)
	r.POST("/api/points_logs/ps/levels", pointshandlers.HandlePsLevelStatus)
	r.POST("/api/ps/mentor_mentee/", pointshandlers.HandleMentee)
	r.POST("/api/mentee/add", pointshandlers.HandleMentee)
	r.POST("/api/projects", projects.PostProjects)
	r.POST("/api/patents", patents.ReceivePatentsData)
	r.POST("/api/internships", internship.ReceiveInternshipData)
	r.POST("/api/workshops", workshops.ReceiveWorkshopData)
	r.GET("/api/fetch/header_details/:rollno", headerdetails.FetchDataRank)
	r.POST("/api/paper-presentations", paperpresentstion.ReceivePaperPresentationData)
	r.POST("/api/certificates/online-course", certificates.ReceiveCertificateData)
	r.POST("/api/certificates/events", certificates.ReceiveCertificateData)
	r.POST("/api/certificates/participation", certificates.ReceiveCertificateData)
	r.POST("/api/addregisterevents", registerevents.HandleRegisterEvents)
	r.GET("/api/events/registered/:rollno", registerevents.GetRegisteredEvents)
	r.GET("/api/events/requested_events/:rollno", registerevents.GetRequestedEvents)
	r.GET("/api/events/registered_events/:rollno", registerevents.GetRegisteredEvents)
	r.PUT("/api/events/registered_events/approve_reject", registerevents.HandleRequestEventsApproveReject)
	r.GET("/api/checkapplied", addevents.CheckApplied)
	r.GET("/api/resume/getprojects/:rollno", resume.GetProjectsData)
	r.GET("/api/resume/getcertificates/:rollno", resume.GetCertificatesData)
	r.GET("/api/activitymaster/getsurveydata/:rollno", activitymaster.GetSurveys)
	r.GET("/api/activitymaster/getsessiondata/:rollno", activitymaster.GetSessionsByRollNo)
	r.GET("/api/uploadview/getuploaddetails/:rollno", dashboard.UploadViewDashboard)
	r.GET("/api/resume/gethackathondata/:rollno", resume.GetHackathonData)
	r.GET("/api/resume/getinternshipdata/:rollno", resume.GetInternshipData)
	r.DELETE("/api/uploadview/deleteupload", Uploadsdelete.Uploadsdelete)
	r.PUT("/api/header/updateprofile", headerdetails.UpdateProfile)
	r.GET("/api/header/getprofile", headerdetails.GetProfileDetails)

	// Faculty routes (previously protected)
	r.GET("/api/manageactivities", manageactivities.GetActivityData)
	r.GET("/api/manageactivities/approvels/:rollno", manageactivities.HandleActivityApprovals)
	r.PUT("/api/manageactivities/approvels_reject", manageactivities.HandleApproveReject)
	r.GET("/api/dashboard/leardeardborad/:rollno", dashBoardfaculty.Leaderboard)
	r.GET("/api/dashboard/prioritylearners/:rollno", dashBoardfaculty.HandlePriorityLearners)
	r.GET("/api/studentrequests/varifications", studentrequests.GetVerifications)
	r.POST("/api/manageactivities/createActivity", manageactivities.ReceiveActivityData)
	r.GET("/api/manageactivities/receiveActivities", manageactivities.GetActivityData)
	r.GET("/api/manageactivities/progressgrpah/:rollno", manageactivities.HandleProgressGraph)
	r.POST("/api/studentrequests/varifications", studentrequests.PostVarification)

	// Admin routes (previously protected)
	r.DELETE("/api/deleteevents/:id", addevents.DeleteEvent)
	r.POST("/api/addevents/create", addevents.AddEvents)
	r.GET("/api/events/fetchregisteredteams/:eventcode", registerevents.HandleRegisteredTeams)
	r.GET("/api/studentdata/fetchstudentdata", studentdata.HandleStudentData)

	// Public routes (already public)
	r.GET("/api/activitymaster/fetch", addevents.FetchEvents)
	r.GET("/api/studentdata/fetchmentees/:rollno", studentdata.HandleMenteesData)
	r.GET("/api/activity_graph/fetchData/:rollno", activitygraph.FetchActivityGraphData)
	r.GET("/api/achievement_graph/fetchData/:rollno", achievementgraph.HandleFetchAchievementGraph)
	r.GET("/api/achievement_graph/institute_avg/fetchData", achievementgraph.HandleFetchInstituteAvg)
	r.GET("/api/ps/attempts/:rollno", pointshandlers.HandleFetchPsAttempts)
	r.GET("/api/ps/levels_status/:rollno", pointshandlers.HandleFetchPsLevels)
	r.GET("/api/mentor/details/:rollno", pointshandlers.FetchMentorSkillStats)
	r.GET("/api/mentor/institute_avg/fetchData", pointshandlers.FetchSkillWiseAvgMentees)
	r.GET("/api/sem_wise_totaldays", pointshandlers.HandleSemDays)
	r.GET("/api/handlesem", pointshandlers.HandleSem)
	r.PUT("/api/updatesem", pointshandlers.HandleUpdateSem)

	// NAAC Report routes - ADD THESE LINES
	r.GET("/api/projects/approved", naacreports.GetApprovedProjects)
	r.GET("/api/certificates/all", naacreports.GetAllCertificates)
	r.GET("/api/internships/all", naacreports.GetAllInternships)
	r.GET("/api/activities/extension", naacreports.GetExtensionActivities) // NEW

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "sih-backend-api",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
}
