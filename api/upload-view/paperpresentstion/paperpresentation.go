package paperpresentstion

import (
	"bitresume/config"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ReceivePaperPresentationData(c *gin.Context){
	rollno := c.PostForm("rollno")
	paper_title := c.PostForm("paper_title")
	conference_title := c.PostForm("conference_title")
	location := c.PostForm("location")
	date_of_presentation := c.PostForm("date_of_presentation")
	award := c.PostForm("award")

	pdf, err := c.FormFile("pdf")
	if err != nil {
		c.JSON(500, "Pdf not receive")
		return
	}

	savePathPdf := filepath.Join("uploads/paperpresentation/presentation_pdf", pdf.Filename)

	if err := os.MkdirAll("uploads/paperpresentation/presentation_pdf", os.ModePerm); err != nil {
		c.JSON(500 , "could not save the pdf")
		return
	}
	if err := c.SaveUploadedFile(pdf,savePathPdf); err != nil{
		c.JSON(500, "Could not save file")
		return
	}

	certificate, err := c.FormFile("certificate")
	if err != nil {
		c.JSON(500, "certificate not receive")
		return
	}
	savePathCertificate := filepath.Join("uploads/paperpresentation/presentation_certificate", pdf.Filename)

	if err := os.MkdirAll("uploads/paperpresentation/presentation_certificate", os.ModePerm); err != nil {
		c.JSON(500 , "could not save the pdf")
		return
	}
	if err := c.SaveUploadedFile(certificate,savePathCertificate); err != nil{
		c.JSON(500, "Could not save file")
		return
	}
	uploadType := "paperpresentation"

	query := `
		insert into paperpresentation 
		(
			upload_type,
			rollno,
			paper_title,
			conference_title,
			location,
			date_of_presentation,
			pdf,
			certificate,
			award,
			approval_status,
			submitted_on
		)
		values (?,?,?,?,?,?,?,?,?,?,CURRENT_DATE)
	`

	_,err =config.DB.Exec(query,uploadType,rollno,paper_title,conference_title,location,date_of_presentation,savePathPdf,savePathCertificate,award,"Pending")

	if err != nil {
		c.JSON(500, "could not upload to db")
		return
	}

	c.JSON(200,"success")

}