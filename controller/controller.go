package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swimresults/athlete-service/service"
	"net/http"
	"os"
	"strconv"
)

var router = gin.Default()

func Run() {

	port := os.Getenv("SR_ATHLETE_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set SR_ATHLETE_PORT.")
		return
	}

	athleteController()
	teamController()

	router.GET("/actuator", actuator)

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}

func extractPagingParams(c *gin.Context) service.Paging {
	limit := 1000
	offset := 0
	query := ""
	limit, _ = strconv.Atoi(c.Query("limit"))
	offset, _ = strconv.Atoi(c.Query("offset"))
	query = c.Query("query")
	return service.Paging{Limit: limit, Offset: offset, Query: query}
}

func actuator(c *gin.Context) {

	state := "OPERATIONAL"

	if !service.PingDatabase() {
		state = "DATABASE_DISCONNECTED"
	}
	c.String(http.StatusOK, state)
}
