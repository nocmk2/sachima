package server

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	// "github.com/nocmk2/sachima/rule"
)

// func Features(){
// rule1 := Rule{rulePath: jsonPath}
// 	rule1.lazyInit()
// 	log.Println(rule1.srcsql)
// 	d := dur.ReadSQL(rule1.srcsql, "hawaii")
// 	log.Println(d.Rows())
// 	// d.Add(scores, "dxscore")
// 	rule1.cal(d)
// }

func rulesHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"rule":     "a",
	})
}
