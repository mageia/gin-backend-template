package main

import (
	"api-server/config"
	"api-server/server"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	h       bool
	v       bool
	V       bool
	Version = "v0.1.0"
)

func main() {
	flag.Parse()

	if h {
		flag.PrintDefaults()
		return
	}

	if v {
		fmt.Println(Version)
		return
	}

	if V {
		fmt.Println(Version)
		return
	}

	server.RunServer()
}

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&V, "V", false, "show version and configure options then exit")

	zerolog.SetGlobalLevel(config.LOG_LEVEL)
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatCaller: func(i interface{}) string {
			_, f := path.Split(i.(string))
			return "[" + fmt.Sprintf("%-20s", f) + "]"
		},
	}).With().Caller().Timestamp().Stack().Logger()
}
