package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	// API_SECRET is the secret key used to sign the JWT
	API_SECRET string = "mysecret"

	//TOKEN_HOUR_LIFESPAN is the second the token will be valid for
	TOKEN_HOUR_LIFESPAN int = 3600 * 24

	// LOG_LEVEL is the log level, valid values are "debug", "info", "warn", "error", "fatal", "panic"
	LOG_LEVEL = zerolog.InfoLevel
)

func init() {
	if os.Getenv("API_SECRET") != "" {
		API_SECRET = os.Getenv("API_SECRET")
	}

	if os.Getenv("TOKEN_HOUR_LIFESPAN") != "" {
		if t, e := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN")); e == nil {
			TOKEN_HOUR_LIFESPAN = t
		}
	}

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		LOG_LEVEL = zerolog.DebugLevel
	case "info":
		LOG_LEVEL = zerolog.InfoLevel
	case "warn":
		LOG_LEVEL = zerolog.WarnLevel
	case "error":
		LOG_LEVEL = zerolog.ErrorLevel
	case "fatal":
		LOG_LEVEL = zerolog.FatalLevel
	case "panic":
		LOG_LEVEL = zerolog.PanicLevel
	default:
		LOG_LEVEL = zerolog.InfoLevel
	}
}
