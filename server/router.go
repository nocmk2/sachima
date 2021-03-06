package server

import (
	"log"

	gormadapter "github.com/casbin/gorm-adapter"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nocmk2/sachima/server/component"
)

func router(r *gin.Engine, au *jwt.GinJWTMiddleware) {
	// defer component.DB.Close()
	adapter := gormadapter.NewAdapterByDB(component.DB)

	r.POST("/login", au.LoginHandler)
	r.GET("/test", testHandler)
	r.GET("/test2", test2Handler)
	r.GET("/test3/:arg", test3Handler)

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
		sachima.GET("/hello", Casbin("hello", "write", adapter), helloHandler)
		sachima.GET("/role", roleHandler)
		sachima.GET("/signup", signupHandler)
		sachima.GET("/featurelists", featurelistsHandler)
		sachima.GET("/rules", rulesHandler)
		sachima.GET("/rule/:name/:version", ruleHandler)
		sachima.GET("/featuredetail/:feature", featuredetailHandler)
		sachima.GET("/user/:ID", getUserByIDHandler)
		sachima.GET("/role/:ID", getRoleByIDHandler)
		sachima.GET("/object/:ID", getObjectByIDHandler)
		sachima.POST("/adduser", Casbin("config", "write", adapter), adduserHandler)
		sachima.GET("/getusers", Casbin("config", "write", adapter), getUsersHandler)
		sachima.GET("/getroles", Casbin("config", "write", adapter), getRolesHandler)
		sachima.GET("/getobjects", Casbin("config", "write", adapter), getObjectsHandler)
		sachima.GET("/getuserrole", Casbin("config", "write", adapter), getUserRoleHandler)
		sachima.GET("/getroleobject", Casbin("config", "write", adapter), getRoleObjectHandler)
	}
}
