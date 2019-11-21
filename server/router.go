package server

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func router(r *gin.Engine, au *jwt.GinJWTMiddleware) {
	r.POST("/login", au.LoginHandler)
	r.GET("/test", testHandler)
	r.GET("/test2", test2Handler)

	r.NoRoute(au.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", au.RefreshHandler)
	auth.Use(au.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/role", roleHandler)
		auth.GET("/signup", signupHandler)
		auth.GET("/rules", rulesHandler)
	}
}
