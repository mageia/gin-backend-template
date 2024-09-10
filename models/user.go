package models

import (
	"api-server/config"
	"html"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Role     string `gorm:"not null;default:'user'"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"uniqueIndex;not null"`
	Avatar   string
}

func (u *User) BeforeSave(*gorm.DB) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPass)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func (u *User) SafeResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
		"avatar":   u.Avatar,
		"email":    u.Email,
	}
}

func LoginCheck(username, password string) (string, error) {
	var u User
	if DB.Model(User{}).Where("username = ?", username).Or("email = ?", username).First(&u).Error != nil {
		return "", gorm.ErrRecordNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", err
	}

	token, err := GenerateToken(&u)
	if err != nil {
		return "", err
	}

	return token, nil
}

type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(u *User) (string, error) {
	claims := Claims{
		UserId:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.G.Auth.TokenExpire)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.G.Auth.ApiSecret))
}
