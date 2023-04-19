package server

import (
	"api-server/config"
	"api-server/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RunServer() {
	gin.SetMode(gin.ReleaseMode)

	log.Info().Msg(fmt.Sprintf("Server start at http://0.0.0.0:%d", config.G.Port))
	s := &http.Server{Handler: router.InitRouter(), Addr: fmt.Sprintf(":%d", config.G.Port)}
	if err := s.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("ListenAndServe")
	}
}
