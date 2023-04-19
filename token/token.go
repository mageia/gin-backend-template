package token

import (
	"api-server/config"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractToken(c *gin.Context) string {
	token := c.Query("access_token")
	if token == "" {
		bearerToken := c.GetHeader("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			token = strings.Split(bearerToken, " ")[1]
		}
	}
	return token
}

func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.G.Auth.ApiSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId := claims["user_id"].(float64)
		return uint(userId), nil
	}

	return 0, nil
}
