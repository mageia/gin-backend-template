package auth_jwt

import (
	"api-server/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.TOKEN_HOUR_LIFESPAN)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.API_SECRET))
}
