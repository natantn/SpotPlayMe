package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/natantn/SpotPlayMe/controllers"
)

func HandleRequests() {
	r := gin.Default()

	r.GET("/spotify/token", controllers.GetToken)
	r.GET("/spotify/sync", controllers.SyncPlaylists)

	r.Run()
}
