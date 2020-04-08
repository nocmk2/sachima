package server

import (
	"log"
	"net/http"
	"time"

	"github.com/nocmk2/sachima/pass"
	"github.com/nocmk2/sachima/server/component"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/rule"
	"golang.org/x/crypto/bcrypt"
)

func signupHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

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

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

func testHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	time.Sleep(7 * time.Second)
	c.JSON(200, gin.H{
		"text": "Test Sachima",
	})
}

func test2Handler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(200, gin.H{
		"text": "Test2 Sachima",
	})
}

func test3Handler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fname := c.Param("arg")
	c.JSON(200, gin.H{
		"text": fname,
	})
}

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

//NewUser info
// type NewUser struct {
// 	UserName string `form:"user" json:"username" binding:"required"`
// 	Password string `form:"password" json:"password" binding:"required"`
// 	Email    string `json:"email"`
// }

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func adduserHandler(c *gin.Context) {
	var user User
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := pass.HashAndSalt([]byte(user.Password))
	log.Println(hash)
	log.Println("save to database")
	log.Println(component.DB)

	c.JSON(http.StatusOK, gin.H{"status": hash})

	// c.JSON(200, gin.H{
	// 	"status":  "posted",
	// 	"message": message,
	// 	"nick":    nick,
	// })
}
