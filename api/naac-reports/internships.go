package naacreports

import (
    "bitresume/config"
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetAllInternships(c *gin.Context) {
    year := c.Query("year")
    rollno := c.Query("rollno")
    
    query := `
    SELECT 
        i.rollno,
        i.company_name,
        i.domain,
        i.start_date,
        i.end_date,
        l.user_name,
        l.department
    FROM internships i
    JOIN login l ON i.rollno = l.rollno
    WHERE i.status = 'Approved'
    `
    
    params := []interface{}{}
    
    if year != "" {
        query += " AND YEAR(i.start_date) = ?"
        params = append(params, year)
    }
    
    if rollno != "" {
        query += " AND i.rollno = ?"
        params = append(params, rollno)
    }
    
    query += " ORDER BY i.start_date DESC"
    
    rows, err := config.DB.Query(query, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()
    
    var internships []map[string]interface{}
    for rows.Next() {
        var internship struct {
            RollNo      string `json:"rollno"`
            CompanyName string `json:"company_name"`
            Domain      string `json:"domain"`
            StartDate   string `json:"start_date"`
            EndDate     string `json:"end_date"`
            UserName    string `json:"user_name"`
            Department  string `json:"department"`
        }
        
        err := rows.Scan(&internship.RollNo, &internship.CompanyName, &internship.Domain,
                        &internship.StartDate, &internship.EndDate, &internship.UserName, &internship.Department)
        if err != nil {
            continue
        }
        
        internships = append(internships, map[string]interface{}{
            "rollno": internship.RollNo,
            "company_name": internship.CompanyName,
            "domain": internship.Domain,
            "start_date": internship.StartDate,
            "end_date": internship.EndDate,
            "user_name": internship.UserName,
            "department": internship.Department,
            "program_name": "Bachelor of Technology",
            "program_code": "CSE001",
        })
    }
    
    c.JSON(http.StatusOK, internships)
}
