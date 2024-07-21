package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var G Config

type Config struct {
	Host string `mapstructure:"listen_host"`
	Port int    `mapstructure:"listen_port"`

	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	Auth struct {
		ApiSecret   string `mapstructure:"api_secret"`
		TokenExpire int    `mapstructure:"token_expire"`
	} `mapstructure:"auth"`

	DB struct {
		Type string `mapstructure:"type"`
		URL  string `mapstructure:"url"`
	} `mapstructure:"db"`
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if e := viper.ReadInConfig(); e != nil {
		log.Fatal().Err(e).Msg(fmt.Sprintf("Config file not found: %s", viper.ConfigFileUsed()))
	}

	if e := viper.Unmarshal(&G); e != nil {
		log.Fatal().Err(e).Msg("Error reading config file")
	}

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		FormatCaller: func(i interface{}) string {
			_, f := path.Split(i.(string))
			return "[" + fmt.Sprintf("%-20s", f) + "]"
		},
	}).With().Caller().Timestamp().Stack().Logger()

	LOG_LEVEL := zerolog.InfoLevel
	switch G.Log.Level {
	case "debug":
		LOG_LEVEL = zerolog.DebugLevel
	case "warn":
		LOG_LEVEL = zerolog.WarnLevel
	case "error":
		LOG_LEVEL = zerolog.ErrorLevel
	case "fatal":
		LOG_LEVEL = zerolog.FatalLevel
	case "panic":
		LOG_LEVEL = zerolog.PanicLevel
	}
	zerolog.SetGlobalLevel(LOG_LEVEL)
}
