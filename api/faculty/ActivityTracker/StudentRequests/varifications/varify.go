package studentrequests

import (
	"bitresume/config"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PostVarification(c *gin.Context) {
	feedback := c.PostForm("feedback")
	id := c.PostForm("id")
	rejected,_ := strconv.ParseBool(c.PostForm("rejected"))
	upload_type := c.PostForm("upload_type")
	verified,_ := strconv.ParseBool(c.PostForm("verified"))

	if upload_type == "patents"{
		fmt.Println("Varified: ", verified)
		fmt.Print("Rejected: ", rejected)

		var status string
		if verified == true {
			status = "Approved"
		} else if rejected == true {
			status = "Rejected"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE patents
			SET patent_status = ?, faculty_remarks = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, feedback, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Patent status updated successfully",
		})
	}else if upload_type == "certificate"{
		var status string
		if verified == true {
			status = "Verified"
		} else if rejected == true {
			status = "Rejected"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE certificates_type
			SET status = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Patent status updated successfully",
		})
	}else if upload_type == "paperpresentation"{
		var status string
		if verified == true {
			status = "Approved"
		} else if rejected == true {
			status = "Not Approved"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE paperpresentation
			SET approval_status = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Patent status updated successfully",
		})
	}else if upload_type == "workshop"{
		var status string
		if verified == true {
			status = "Approved"
		} else if rejected == true {
			status = "Rejected"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE workshops
			SET status = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Patent status updated successfully",
		})
	}else if upload_type == "internship"{
		var status string
		if verified == true {
			status = "Approved"
		} else if rejected == true {
			status = "Rejected"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE internships
			SET status = ?, faculty_remarks = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, feedback, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Patent status updated successfully",
		})
	}else if upload_type == "project"{
		tier := c.PostForm("tier")
		feedback := c.PostForm("feedback")
		projectId := c.PostForm("id")

		var status string
		if verified == true {
			status = "Approved"
		} else if rejected == true {
			status = "Rejected"
		} else {
			status = "Pending"
		}

		query := `
			UPDATE projects
			SET approval_status = ?, complexity = ?
			WHERE id = ?
		`

		_, err := config.DB.Exec(query, status, tier, id)
		if err != nil {
			fmt.Print("error:",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update patent status",
			})
			return
		}

		query = `
			insert into project_evaluation(
				project_id,
				faculty_remarks,
				upload_date
			) values (?,?,?)
		`

		_,err = config.DB.Exec(query, projectId,feedback,time.Now().Format("2006-01-02"))
		if err != nil {
			fmt.Println("Error: ",err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project status"})
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Project status updated successfully",
		})
	}
}