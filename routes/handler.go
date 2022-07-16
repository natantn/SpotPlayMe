package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/natantn/SpotPlayMe/controllers"
)

func HandleRequests() {
	r := gin.Default()

	r.GET("/token", controllers.GetToken)

	r.Run()
}
