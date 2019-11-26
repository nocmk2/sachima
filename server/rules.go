package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/rule"
)

const jsonPath string = "data/rule.json"

func features() []string {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	return rule1.FeatureList()
}

func rulesHandler(c *gin.Context) {
	f := features()
	// log.Println(string(f))
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     f,
	})
}
