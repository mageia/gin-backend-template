package router

import (
	"api-server/controller"
	"api-server/docs"
	"api-server/middleware"
	"api-server/models"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/rs/zerolog/log"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func InitRouter() *gin.Engine {
	e, err := models.GetEnforcer()
	if err != nil {
		log.Error().Err(err).Msg("GetEnforcer")
		return nil
	}

	r := gin.Default()
	ginprometheus.NewPrometheus("gin").Use(r)
	r.Use(cors.Default())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestId())
	r.Use(middleware.ErrorHandler())

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "Api Server Template"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Api Server Template"
	r.GET("/swagger/*any", func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
		ginSwagger.WrapHandler(swaggerfiles.Handler)(c)
	})

	r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
	auth.POST("/logout", middleware.JwtAuth(e), controller.Logout)

	admin := api.Group("/admin")
	admin.Use(middleware.JwtAuth(e))
	admin.GET("/user", controller.RetrieveCurrentUser)
	admin.PUT("/user", controller.UpdateCurrentUser)

	admin.POST("/policies", controller.PostPolicy)
	admin.GET("/policies", controller.GetPolicy)
	admin.PUT("/policies", controller.PutPolicy)
	admin.DELETE("/policies", controller.DelPolicy)

	return r
}
