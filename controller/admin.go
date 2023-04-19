package controller

import (
	"api-server/models"
	"api-server/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetrieveCurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var u models.User
	if e := models.DB.First(&u, userId).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

type UpdateUserInput struct {
	Avatar *string `json:"avatar"`
}

func UpdateCurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if input.Avatar == nil {
		c.JSON(200, gin.H{"message": "Nothing to do"})
		return
	}

	var u models.User
	if e := models.DB.First(&u, userId).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	if e := models.DB.Model(&u).Update("avatar", *input.Avatar).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}
