package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/athlete-service/dto"
	"github.com/swimresults/athlete-service/model"
	"github.com/swimresults/athlete-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func athleteController() {
	router.GET("/athlete", getAthletes)
	router.GET("/athlete/:id", getAthlete)
	router.GET("/athlete/name_year", getAthleteByNameAndYear)
	router.GET("/athlete/meet/:meet_id", getAthletesByMeeting)
	router.GET("/athlete/team/:team_id", getAthletesByTeam)
	router.GET("/athlete/team/:team_id/meet/:meet_id", getAthletesByTeamAndMeeting)

	router.DELETE("/athlete/:id", removeAthlete)
	router.POST("/athlete", addAthlete)
	router.POST("/athlete/import", importAthlete)
	router.POST("/athlete/participation", addParticipation)
	router.PUT("/athlete", updateAthlete)

	router.HEAD("/athlete", getAthletes)
	router.HEAD("/athlete/:id", getAthlete)
}

func getAthletes(c *gin.Context) {
	athletes, err := service.GetAthletes(extractPagingParams(c))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athletes)
}

func getAthletesByMeeting(c *gin.Context) {
	id := c.Param("meet_id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	athletes, err := service.GetAthletesByMeetingId(id, extractPagingParams(c))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athletes)
}

func getAthletesByTeam(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("team_id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given team_id was not of type ObjectID"})
		return
	}

	athletes, err := service.GetAthletesByTeamId(id, extractPagingParams(c))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athletes)
}

func getAthletesByTeamAndMeeting(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("team_id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given team_id was not of type ObjectID"})
		return
	}

	meeting := c.Param("meet_id")
	if meeting == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	athletes, err := service.GetAthletesByTeamAndMeeting(id, meeting, extractPagingParams(c))
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, athlete)
}

func getAthleteByNameAndYear(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given name was empty"})
		return
	}
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	athlete, err := service.GetAthleteByNameAndYear(name, year)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
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

	c.IndentedJSON(http.StatusCreated, r)
}

func addParticipation(c *gin.Context) {

	var data dto.AddParticipationRequestDto
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if data.MeetingId == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meeting is empty"})
		return
	}

	if data.AthleteId.IsZero() {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given athlete is empty"})
		return
	}

	r, err := service.AddParticipation(data.AthleteId, data.MeetingId)
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

func importAthlete(c *gin.Context) {
	var athlete *model.Athlete
	var request dto.ImportAthleteRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	athlete, r, err := service.ImportAthlete(request.Athlete, request.Meeting)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if r {
		c.IndentedJSON(http.StatusCreated, *athlete)
	} else {
		c.IndentedJSON(http.StatusOK, *athlete)
	}
}
