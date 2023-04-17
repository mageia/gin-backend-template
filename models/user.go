package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Avatar   string
}

func (u *User) SaveUser() (*User, error) {
  if e := DB.Create(&u).Error; e != nil {
    return nil, e
  }

  return u, nil
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
