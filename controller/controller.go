package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
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

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}
