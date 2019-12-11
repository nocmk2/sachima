package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func roleHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":        claims[identityKey],
		"userName":      user.(*User).UserName,
		"text":          "Hello World.",
		"role":          "view",
		"lastlogintime": "1111",
	})
}
