package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/athlete-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func certificateController() {
	router.GET("/certificate", getCertificates)

	router.GET("/certificate/amount", getCertificatesAmount)
	router.GET("/certificate/meet/:meet_id/amount", getCertificatesAmountByMeeting)

	router.GET("/certificate/:id", getCertificate)
	router.GET("/certificate/athlete/:athlete_id", getCertificatesByAthlete)
	router.GET("/certificate/athlete/:athlete_id/meet/:meet_id", getCertificatesByAthleteAndMeeting)

	router.DELETE("/certificate/:id", removeCertificate)
	router.POST("/certificate", addCertificate)
	router.POST("/certificate/import", importCertificate)
	router.PUT("/certificate", updateCertificate)
}

func getCertificates(c *gin.Context) {
	certificates, err := service.GetCertificates()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, certificates)
}

func getCertificatesAmount(c *gin.Context) {
	starts, err := service.GetCertificatesAmount()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, starts)
}

func getCertificatesAmountByMeeting(c *gin.Context) {
	meeting := c.Param("meet_id")
	if meeting == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	starts, err := service.GetCertificatesAmountByMeeting(meeting)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, starts)
}

func getCertificatesByAthlete(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("athlete_id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given athlete_id was not of type ObjectID"})
		return
	}

	certificates, err := service.GetCertificatesByAthleteId(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, certificates)
}

func getCertificatesByAthleteAndMeeting(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("athlete_id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given athlete_id was not of type ObjectID"})
		return
	}

	meeting := c.Param("meet_id")
	if meeting == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	certificates, err := service.GetCertificatesByAthleteIdAndMeeting(id, meeting)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, certificates)
}

func getCertificate(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	certificate, err := service.GetCertificateById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, certificate)
}

func removeCertificate(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveCertificateById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addCertificate(c *gin.Context) {
	var certificate model.Certificate
	if err := c.BindJSON(&certificate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddCertificate(certificate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, r)
}

func updateCertificate(c *gin.Context) {
	var certificate model.Certificate
	if err := c.BindJSON(&certificate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateCertificate(certificate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func importCertificate(c *gin.Context) {
	var certificate *model.Certificate
	var request dto.ImportCertificateRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	certificate, err := service.ImportCertificate(request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, certificate)

}
