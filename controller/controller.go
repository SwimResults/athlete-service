package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func Run() {

	exampleController()

	port := os.Getenv("SR_EXAMPLE_PORT")

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}
