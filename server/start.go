package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/server/component"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// User demo
type User struct {
	UserName  string `form:"user" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string
	LastName  string
	Email     string `json:"email"`
}

// Start auth server
func Start() {
	defer component.DB.Close()
	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// return origin == "https://github.com"
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	if port == "" {
		port = "8000"
	}

	router(r, Jwt())

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
