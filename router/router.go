package router

import (
	"api-server/controller"
	"api-server/middlewares"
	"api-server/models"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InitRouter() *gin.Engine {
	e, err := models.GetEnforcer()
	if err != nil {
		log.Error().Err(err).Msg("GetEnforcer")
		return nil
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(func(c *gin.Context) {
		if len(c.Errors) != 0 {
			c.AbortWithStatusJSON(400, gin.H{"message": c.Errors.String()})
			return
		}
		c.Next()
	})

	r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
	auth.POST("/logout", middlewares.JwtAuthMiddleware(e), controller.Logout)

	admin := api.Group("/admin")
	admin.Use(middlewares.JwtAuthMiddleware(e))
	admin.GET("/user", controller.RetrieveCurrentUser)
	admin.PUT("/user", controller.UpdateCurrentUser)

	admin.POST("/policies", controller.PostPolicy)
	admin.GET("/policies", controller.GetPolicy)
	admin.PUT("/policies", controller.PutPolicy)
	admin.DELETE("/policies", controller.DelPolicy)

	return r
}
