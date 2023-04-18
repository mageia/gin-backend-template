package controller

import (
	"api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	u := models.User{
		Username: input.Username,
		Password: input.Password,
	}

	if e := models.DB.Create(&u).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("ShouldBindJSON")
		c.JSON(400, gin.H{"message": "Check Login Input failed"})
		return
	}

	u := models.User{Username: input.Username, Password: input.Password}
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		log.Error().Err(err).Msg("LoginCheck")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func Logout(c *gin.Context) {}
