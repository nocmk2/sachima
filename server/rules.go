package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/rule"
)

const jsonPath string = "data/rule.json"

func featurelists() []string {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	return rule1.FeatureList()
}

func featureName(name string) string {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	// return rule1.
	return "jijfiejf"
}

func featurelistsHandler(c *gin.Context) {
	f := featurelists()
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     f,
	})
}

func featuredetailHandler(c *gin.Context) {
	fname := c.Param("feature")
	cname := featureName(fname)
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     cname,
	})
}
