package server

import (
	"github.com/gin-gonic/gin"
)

func testHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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
	fname := c.Param("feature")
	c.JSON(200, gin.H{
		"text": fname,
	})
}
