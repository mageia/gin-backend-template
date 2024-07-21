package server

import (
	"api-server/config"
	"api-server/router"
	"api-server/rpc"
	"api-server/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunServer() {
	gin.SetMode(gin.ReleaseMode)
	url := fmt.Sprintf("http://%s:%d", config.G.Host, config.G.Port)
	log.Info().Msgf("Health check at %s/healthz", url)
	log.Info().Msgf("Swagger docs at %s/swagger/index.html", url)
	log.Info().Msg("Press Ctrl+C to quit")

	grpcServer := grpc.NewServer()
	rpc.RegisterHelloWroldServer(grpcServer, &service.HelloWorld{})
	reflection.Register(grpcServer)

	router := router.InitRouter()
	h2Handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 &&
			strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			router.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	s := &http.Server{Addr: fmt.Sprintf(":%d", config.G.Port), Handler: h2Handler}
	if e := s.ListenAndServe(); e != nil {
		log.Error().Err(e).Msg("ListenAndServe")
	}
}
