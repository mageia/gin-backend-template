package middleware

import "github.com/gin-gonic/gin"

type baseErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) != 0 {
			c.AbortWithStatusJSON(400, baseErrorResponse{400, c.Errors.Last().Error()})
			return
		}
	}
}
