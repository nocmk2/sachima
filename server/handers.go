package server

import (
	"net/http"
	"time"

	"github.com/nocmk2/sachima/dur"

	"github.com/nocmk2/sachima/pass"
	"github.com/nocmk2/sachima/server/component"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/rule"
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

func adduserHandler(c *gin.Context) {
	db := component.DB
	var user User
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash := pass.HashAndSalt([]byte(user.Password))

	new := User{UserName: user.UserName, Password: hash, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}
	db.Create(&new)

	c.JSON(http.StatusOK, gin.H{"status": hash})
}

func getRolesHandler(c *gin.Context) {
	d := dur.ReadSQL("select * from roles", "sachima_local")
	c.JSON(http.StatusOK, d.AllRows())
}

func getUsersHandler(c *gin.Context) {
	d := dur.ReadSQL("select user_name as id,first_name as name from users", "sachima_local")
	c.JSON(http.StatusOK, d.AllRows())
}
