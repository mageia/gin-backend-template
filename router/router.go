package router

import (
	"executor/controller"
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

	router.POST("/executor", controller.Executor)
	router.POST("/importer", controller.Importer)

	return router
}
