package router

import (
	"api-server/controller"
	"api-server/middlewares"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) != 0 {
			c.AbortWithStatusJSON(400, gin.H{"code": "failed", "message": c.Errors.String()})
		}
	})

	router.GET("/healthz", controller.Health)

	api := router.Group("/api/v1")
	api.POST("/login", controller.Login)
	api.POST("/register", controller.Register)
	api.POST("/logout", controller.Logout)

	admin := router.Group("/api/v1/admin")
	admin.Use(middlewares.JwtAuthMiddleware())
	admin.GET("/user", controller.CurrentUser)

	return router
}
