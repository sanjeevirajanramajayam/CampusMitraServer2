package jobs
import (
	achievementgraph "bitresume/api/dashboard/achievement_graph"
	activitygraph "bitresume/api/dashboard/activity_graph"
	"bitresume/config"
	"log"
	"time"
)
func CallDailyTasksForAllDates() {
	currentDate := time.Now().Format("2006-01-02")
	DailyTask(currentDate) 
}
type Student struct {
	RollNo string
	Sem    int   
}
func GetStudentData() ([]Student, error) {
	var students []Student
	query := `SELECT rollno, sem FROM login WHERE role = 'student'`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.RollNo, &s.Sem); err != nil {
			log.Printf("Error scanning student row: %v", err)
			continue
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}
// DailyTask executes all daily tasks for a specific date and semester
func DailyTask(date string) {
	students, err := GetStudentData()
	if err != nil {
		log.Fatal("Error fetching student data:", err)
	}
	for _, s := range students {
		log.Printf("RollNo: %s, Sem: %d\n", s.RollNo, s.Sem)
		activitygraph.HandleInactivity(s.RollNo, date, s.Sem)
		achievementgraph.HandleInactivity(s.RollNo, date, s.Sem)
		activitygraph.HandleActivityGraphPoints(s.RollNo, s.Sem, date)
		achievementgraph.HandleAcheivemnetPoints(s.RollNo, date, s.Sem)
	}
	achievementgraph.HandleInstituteAvg(1, date)
}
