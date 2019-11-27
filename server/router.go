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
	r.GET("/feature/:feature", test3Handler)

	r.NoRoute(au.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	sachima := r.Group("/sachima")
	// Refresh time can be longer than token timeout
	sachima.GET("/refresh_token", au.RefreshHandler)
	sachima.Use(au.MiddlewareFunc())
	{
		sachima.GET("/hello", helloHandler)
		sachima.GET("/role", roleHandler)
		sachima.GET("/signup", signupHandler)
		sachima.GET("/featurelists", featurelistsHandler)
		sachima.GET("/featuredetail/:feature", featuredetailHandler)
	}
}
