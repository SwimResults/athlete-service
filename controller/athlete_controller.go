package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sr-athlete/athlete-service/model"
	"sr-athlete/athlete-service/service"
)

func athleteController() {
	router.GET("/athlete", getAthletes)
	router.GET("/athlete/:id", getAthlete)
	router.DELETE("/athlete/:id", removeAthlete)
	router.POST("/athlete", addAthlete)
	router.PUT("/athlete", updateAthlete)
	// TODO: HEAD requests receive 404
}

func getAthletes(c *gin.Context) {
	athletes, err := service.GetAthletes()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athletes)
}

func getAthlete(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	athlete, err := service.GetAthleteById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athlete)
}

func removeAthlete(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveAthleteById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addAthlete(c *gin.Context) {
	var athlete model.Athlete
	if err := c.BindJSON(&athlete); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddAthlete(athlete)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func updateAthlete(c *gin.Context) {
	var athlete model.Athlete
	if err := c.BindJSON(&athlete); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateAthlete(athlete)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
