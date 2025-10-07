package naacreports

import (
    "bitresume/config"
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetApprovedProjects(c *gin.Context) {
    rollno := c.Query("rollno")
    year := c.Query("year")
    department := c.Query("department")
    
    query := `
    SELECT 
        p.rollno,
        p.title_idea,
        p.summary,
        l.user_name,
        l.department,
        p.start_time,
        p.end_time,
        p.approval_status
    FROM projects p
    JOIN login l ON p.rollno = l.rollno
    WHERE p.approval_status = 'Approved'
    `
    
    params := []interface{}{}
    
    if rollno != "" {
        query += " AND p.rollno = ?"
        params = append(params, rollno)
    }
    
    if year != "" {
        query += " AND YEAR(p.created_at) = ?"
        params = append(params, year)
    }
    
    if department != "" {
        query += " AND l.department LIKE ?"
        params = append(params, "%"+department+"%")
    }
    
    query += " ORDER BY p.created_at DESC"
    
    rows, err := config.DB.Query(query, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()
    
    var projects []map[string]interface{}
    for rows.Next() {
        var p struct {
            RollNo        string `json:"rollno"`
            TitleIdea     string `json:"title_idea"`
            Summary       string `json:"summary"`
            UserName      string `json:"user_name"`
            Department    string `json:"department"`
            StartTime     string `json:"start_time"`
            EndTime       string `json:"end_time"`
            ApprovalStatus string `json:"approval_status"`
        }
        
        err := rows.Scan(&p.RollNo, &p.TitleIdea, &p.Summary, 
                        &p.UserName, &p.Department,
                        &p.StartTime, &p.EndTime, &p.ApprovalStatus)
        if err != nil {
            continue
        }
        
        projects = append(projects, map[string]interface{}{
            "rollno": p.RollNo,
            "title_idea": p.TitleIdea,
            "summary": p.Summary,
            "user_name": p.UserName,
            "department": p.Department,
            "program_name": "Bachelor of Technology",
            "program_code": "CSE001", // You can enhance this later
            "start_time": p.StartTime,
            "end_time": p.EndTime,
        })
    }
    
    c.JSON(http.StatusOK, projects)
}
