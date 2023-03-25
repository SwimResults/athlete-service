package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func Run() {

	athleteController()
	teamController()

	port := os.Getenv("SR_ATHLETE_PORT")

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}
