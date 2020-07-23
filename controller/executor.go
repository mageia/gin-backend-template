package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
)

func ReadS3(s3Addr, objectName string) {
	u, err := url.Parse(s3Addr)
	if err != nil {
		log.Error().Err(err).Msg("Parse s3 address")
		return
	}

	fmt.Println(u.Scheme, u.Host, u.User)
}

func Executor(c *gin.Context) {
	ReadS3(os.Getenv("S3_ADDR"), "mageia/1.csv")
	c.JSON(200, "OK")
}
