package models

import (
	"api-server/token"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"uniqueIndex;not null"`
	Avatar   string
}

func (u* User) BeforeSave(*gorm.DB) error {
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
  if err != nil {
    return err
  }

  u.Password = string(hashedPass)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

  return nil
}

func LoginCheck(username, password string)  (string ,error){
  var u User
  if DB.Model(User{}).Where("username = ?", username).Or("email = ?", username).First(&u).Error != nil {
    return "", gorm.ErrRecordNotFound
  }

  if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
    return "", err
  }

  token, err := token.GenerateToken(u.ID)
  if err != nil {
    return "", err
  }

  return token, nil
}
