package server

import (
	"executor/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RunServer() {
	ConfigServer := &http.Server{Addr: ":9999", Handler: router.InitRouter()}
	gin.SetMode(gin.ReleaseMode)
	log.Info().Msg(fmt.Sprintf("Server start at http://0.0.0.0:%d", 9999))
	if err := ConfigServer.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("ListenAndServe")
	}
}
