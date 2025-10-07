package naacreports

import (
    "bitresume/config"
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetAllCertificates(c *gin.Context) {
    year := c.Query("year")
    rollno := c.Query("rollno")
    
    query := `
    SELECT 
        ct.id,
        ct.rollno,
        l.user_name,
        l.department,
        ct.certificate_type,
        ct.status,
        ct.created_at
    FROM certificates_type ct
    JOIN login l ON ct.rollno = l.rollno
    WHERE ct.status = 'Verified'
    `
    
    params := []interface{}{}
    
    if year != "" {
        query += " AND YEAR(ct.created_at) = ?"
        params = append(params, year)
    }
    
    if rollno != "" {
        query += " AND ct.rollno = ?"
        params = append(params, rollno)
    }
    
    query += " ORDER BY ct.created_at DESC"
    
    rows, err := config.DB.Query(query, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()
    
    var certificates []map[string]interface{}
    for rows.Next() {
        var cert struct {
            ID              int    `json:"id"`
            RollNo          string `json:"rollno"`
            UserName        string `json:"user_name"`
            Department      string `json:"department"`
            CertificateType string `json:"certificate_type"`
            Status          string `json:"status"`
            CreatedAt       string `json:"created_at"`
        }
        
        err := rows.Scan(&cert.ID, &cert.RollNo, &cert.UserName, &cert.Department,
                        &cert.CertificateType, &cert.Status, &cert.CreatedAt)
        if err != nil {
            continue
        }
        
        certificates = append(certificates, map[string]interface{}{
            "id": cert.ID,
            "rollno": cert.RollNo,
            "user_name": cert.UserName,
            "department": cert.Department,
            "title": cert.CertificateType,
            "platform": "Event Platform",
            "certificate_type": cert.CertificateType,
            "award_level": "Institute",
            "competition_rank": "Participation",
            "team_or_individual": "Individual",
            "created_at": cert.CreatedAt,
        })
    }
    
    c.JSON(http.StatusOK, certificates)
}
