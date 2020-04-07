package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Capitalize(name string) string {
	return strings.ToUpper(name)
}

func main1() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	r.GET("/hello/:name/*title", func(c *gin.Context) {
		name := c.Param("name")
		title := c.Param("title")
		log.Println(title)
		message := title + " " + Capitalize(name)
		c.String(http.StatusOK, message)
	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message+" "+c.FullPath())
	})

	// Querystring parameters
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	// Multipart/Urlencoded Form
	r.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	//Another example: query + post form
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		c.String(http.StatusOK, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
		log.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})

	//Map as querystring or postform parameters
	r.POST("/postmap", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		fmt.Printf("ids: %v; names: %v", ids, names)
	})

	// upload file
	/*
	 	curl -X POST http://localhost:8080/upload \
	   -F "file=@/Users/appleboy/test.zip" \
	   -H "Content-Type: multipart/form-data"
	*/
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	r.Run()
}
