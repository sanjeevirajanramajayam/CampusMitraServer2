package naacreports

import (
    "bitresume/config"
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetExtensionActivities(c *gin.Context) {
    year := c.Query("year")
    department := c.Query("department")
    
    // Get aggregated data from student uploads (approved only)
    query := `
    SELECT 
        seu.activity_name,
        seu.organizing_unit,
        'Government Initiative' as scheme_name,
        YEAR(seu.participation_date) as activity_year,
        COUNT(DISTINCT seu.student_rollno) as student_participants,
        SUM(seu.hours_contributed) as total_hours
    FROM student_extension_uploads seu
    WHERE seu.approval_status = 'Approved'
    `
    
    params := []interface{}{}
    
    if year != "" {
        query += " AND YEAR(seu.participation_date) = ?"
        params = append(params, year)
    }
    
    if department != "" {
        query += " AND seu.organizing_unit LIKE ?"
        params = append(params, "%"+department+"%")
    }
    
    query += " GROUP BY seu.activity_name, seu.organizing_unit, YEAR(seu.participation_date)"
    query += " ORDER BY activity_year DESC, student_participants DESC"
    
    rows, err := config.DB.Query(query, params...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()
    
    var activities []map[string]interface{}
    for rows.Next() {
        var activity struct {
            ActivityName        string `json:"activity_name"`
            OrganizingUnit      string `json:"organizing_unit"`
            SchemeName          string `json:"scheme_name"`
            ActivityYear        int    `json:"activity_year"`
            StudentParticipants int    `json:"student_participants"`
            TotalHours          int    `json:"total_hours"`
        }
        
        err := rows.Scan(&activity.ActivityName, &activity.OrganizingUnit, 
                        &activity.SchemeName, &activity.ActivityYear,
                        &activity.StudentParticipants, &activity.TotalHours)
        if err != nil {
            continue
        }
        
        activities = append(activities, map[string]interface{}{
            "activity_name": activity.ActivityName,
            "organizing_unit": activity.OrganizingUnit,
            "scheme_name": activity.SchemeName,
            "activity_year": activity.ActivityYear,
            "student_participants": activity.StudentParticipants,
            "total_hours": activity.TotalHours,
        })
    }
    
    c.JSON(http.StatusOK, activities)
}
