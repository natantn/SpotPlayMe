package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/natantn/SpotPlayMe/integrations/spotify"
)

func GetToken(c *gin.Context) {
	token := integrations.AccessToken{}
	token.GenerateToken()
	c.JSON(200, token)
}
