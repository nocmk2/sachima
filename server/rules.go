package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/rule"
)

const jsonPath string = "data/rule.json"

func featureRaw() interface{} {
	r := rule.Rule{}
	r.ReadRuleFile(jsonPath)
	return r.Features()
}

func featurelists() []string {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	return rule1.FeatureList()
}

func featureName(name string) string {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	return rule1.FeatureName(name)
}

func featureBin(name string) map[string]int {
	rule1 := rule.Rule{}
	rule1.ReadRuleFile(jsonPath)
	// return rule1.FeatureBin(name)
	return map[string]int{"a": 100, "b": -99}
}

func featurelistsHandler(c *gin.Context) {
	f := featurelists()
	// f := featureRaw()
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"features": f,
	})
}

func featuresHandler(c *gin.Context) {
	// f := featurelists()
	f := featureRaw()
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"features": f,
	})
}

func featuredetailHandler(c *gin.Context) {
	fname := c.Param("feature")
	bin := featureBin(fname)
	cname := featureName(fname)
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)

	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"name":     cname,
		"bin":      bin,
	})
}
