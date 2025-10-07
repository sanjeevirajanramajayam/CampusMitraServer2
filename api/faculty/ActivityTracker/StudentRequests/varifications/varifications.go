package studentrequests

import (
	facultymodel "bitresume/models/faculty"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVerifications(c *gin.Context) {
	allVarifications := make([]facultymodel.Varification, 0)
	certificates, err := GetCertificates()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the certificate folder")
		return
	}
	allVarifications = append(allVarifications, certificates...)
	workshops, err := GetWorkshops()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the workshops folder")
		return
	}
	allVarifications = append(allVarifications, workshops...)

	projects, err := GetProjects()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the projects folder")
		return
	}
		allVarifications = append(allVarifications, projects...)
	paperpresentation, err := GetPaperpresentation()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the paperpresentation folder")
		return
	}
	// if err == nil {
		allVarifications = append(allVarifications, paperpresentation...)
	// }

	internships, err := GetInternship()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the internship folder")
		return
	}
	allVarifications = append(allVarifications, internships...)

	patents, err := GetPatents()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		fmt.Println("Could not get the data from the patent folder")
		return
	}
	allVarifications = append(allVarifications, patents...)

	c.JSON(http.StatusAccepted, allVarifications)
}