package router

import (
	"api-server/controller"
	"api-server/middlewares"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(func(c *gin.Context) {
		if len(c.Errors) != 0 {
			c.AbortWithStatusJSON(400, gin.H{"message": c.Errors.String()})
			return
		}

		c.Next()
	})

	r.GET("/healthz", controller.Health)

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
	auth.POST("/logout", controller.Logout)

	admin := api.Group("/admin")
	admin.Use(middlewares.JwtAuthMiddleware())

	admin.GET("/user", controller.RetrieveCurrentUser)
	admin.PUT("/user", controller.UpdateCurrentUser)

	return r
}
