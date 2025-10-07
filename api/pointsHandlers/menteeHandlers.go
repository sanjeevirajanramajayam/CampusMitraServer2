package pointshandlers

import (
	"bitresume/config"
	"bitresume/models"
	"github.com/gin-gonic/gin"
)

func HandleMentee(c *gin.Context) {
	var data models.Mentee
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mentorRollNo := data.MentorRollNo
	menteeRollNo := data.MenteeRollNo
	skillName := data.SkillName

	// 1. Check if mentor exists for the given skill
	var mentorExists bool
	err := config.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM mentee_skills 
			WHERE mentor_rollno = ? AND skill_name = ?
		)
	`, mentorRollNo, skillName).Scan(&mentorExists)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check mentor for the given skill"})
		return
	}

	// 2. If mentor exists for the skill, check if mentee exists under that mentor for the same skill
	if mentorExists {
		var menteeExists bool
		err = config.DB.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM mentee_skills 
				WHERE mentor_rollno = ? AND mentee_rollno = ? AND skill_name = ?
			)
		`, mentorRollNo, menteeRollNo, skillName).Scan(&menteeExists)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to check if mentee already exists under this mentor for the skill"})
			return
		}
		if menteeExists {
			c.JSON(200, gin.H{"message": "Mentee already exists under this mentor for the skill, no action taken"})
			return
		}
	}

	// 3. Insert record if mentee not found under this mentor for the skill
	stmt, err := config.DB.Prepare(`
		INSERT INTO mentee_skills (mentor_rollno, mentee_rollno, skill_name)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to prepare insert"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(mentorRollNo, menteeRollNo, skillName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert mentee"})
		return
	}

	c.JSON(201, gin.H{"message": "Mentee added successfully under mentor for the skill"})
}

func FetchMentorSkillStats(c *gin.Context) {
	mentorRollNo := c.Param("rollno")
	rows, err := config.DB.Query(`
		SELECT skill_name, COUNT(*) as mentee_count
		FROM mentee_skills
		WHERE mentor_rollno = ?
		GROUP BY skill_name
	`, mentorRollNo)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch skill stats"})
		return
	}
	defer rows.Close()

	type SkillStat struct {
		SkillName   string `json:"skill_name"`
		MenteeCount int    `json:"mentee_count"`
	}

	var stats []SkillStat
	for rows.Next() {
		var s SkillStat
		if err := rows.Scan(&s.SkillName, &s.MenteeCount); err != nil {
			c.JSON(500, gin.H{"error": "Error scanning stats"})
			return
		}
		stats = append(stats, s)
	}

	c.JSON(200, stats)
}

func FetchSkillWiseAvgMentees(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT  skill_name, ROUND(COUNT(*) * 1.0 / COUNT(DISTINCT mentor_rollno), 2) AS avg_mentees_per_mentor FROM mentee_skills GROUP BY skill_name `)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch skill-wise average mentee count"})
		return
	}
	defer rows.Close()

	type SkillAvg struct {
		SkillName           string  `json:"skill_name"`
		AvgMenteesPerMentor float64 `json:"avg_mentees_per_mentor"`
	}

	var results []SkillAvg

	for rows.Next() {
		var s SkillAvg
		if err := rows.Scan(&s.SkillName, &s.AvgMenteesPerMentor); err != nil {
			c.JSON(500, gin.H{"error": "Error scanning result"})
			return
		}
		results = append(results, s)
	}

	c.JSON(200, results)
}
