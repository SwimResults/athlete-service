package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sr-athlete/athlete-service/model"
	"sr-athlete/athlete-service/service"
)

func teamController() {
	router.GET("/team", getTeams)
	router.GET("/team/:id", getTeam)
	router.POST("/team", addTeam)

	router.HEAD("/team", getTeams)
	router.HEAD("/team/:id", getTeam)
}

func getTeams(c *gin.Context) {
	teams, err := service.GetTeams()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, teams)
}

func getTeam(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	team, err := service.GetTeamById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, team)
}

func addTeam(c *gin.Context) {
	var team model.Team
	if err := c.BindJSON(&team); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddTeam(team)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
