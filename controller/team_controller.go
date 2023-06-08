package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/athlete-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func teamController() {
	router.GET("/team", getTeams)
	router.GET("/team/:id", getTeam)
	router.GET("/team/meet/:meet_id", getTeamsByMeeting)
	router.POST("/team", addTeam)
	router.POST("/team/import", importTeam)

	router.HEAD("/team", getTeams)
	router.HEAD("/team/:id", getTeam)
}

func getTeams(c *gin.Context) {
	teams, err := service.GetTeams(extractPagingParams(c))
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

func getTeamsByMeeting(c *gin.Context) {
	id := c.Param("meet_id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	teams, err := service.GetTeamsByMeeting(id, extractPagingParams(c))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, teams)
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

	c.IndentedJSON(http.StatusCreated, r)
}

func importTeam(c *gin.Context) {
	var team model.Team
	if err := c.BindJSON(&team); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	team, r, err := service.ImportTeam(team)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if r {
		c.IndentedJSON(http.StatusCreated, team)
	} else {
		c.IndentedJSON(http.StatusOK, team)
	}
}
