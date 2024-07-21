package middleware

import (
	"api-server/config"
	"api-server/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuth(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("access_token")
		if tokenString == "" {
			bearerToken := c.GetHeader("Authorization")
			if len(strings.Split(bearerToken, " ")) == 2 {
				tokenString = strings.Split(bearerToken, " ")[1]
			}
		}

		claims := models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.G.Auth.ApiSecret), nil
		})

		if err != nil {
			var message string
			if err == jwt.ErrSignatureInvalid {
				message = "Invalid token signature"
			} else {
				message = "Invalid token"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": message})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		c.Set("user", &models.User{
			ID:       claims.UserId,
			Username: claims.Username,
			Role:     claims.Role,
		})

		if ok, err := e.Enforce(claims.Role, c.Request.URL.Path, c.Request.Method); !ok || err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			return
		}

		c.Next()
	}
}
