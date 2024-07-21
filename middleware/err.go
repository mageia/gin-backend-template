package middleware

import "github.com/gin-gonic/gin"

type baseErrorResponse struct {
	Error string `json:"error"`
}

func ErrProcessor(c *gin.Context) {
	c.Next()

	if len(c.Errors) != 0 {
		c.AbortWithStatusJSON(400, baseErrorResponse{c.Errors.String()})
		return
	}
}
