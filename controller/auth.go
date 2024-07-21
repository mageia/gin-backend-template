package controller

import (
	"api-server/models"
	"api-server/schema"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register
// @Description User Register
// @Tags Auth
// @Router /auth/register [post]
// @Param schema.ReqRegister body schema.ReqRegister true "User Register"
// @Success 200 {object} schema.ResRegister
func Register(c *gin.Context) {
	var input schema.ReqRegister
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	u := models.User{Username: input.Username, Password: input.Password}

	if e := models.DB.Create(&u).Error; e != nil {
		c.Error(fmt.Errorf("Registration failed: %s", e.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

// @Summary Login
// @Description User Login
// @Tags Auth
// @Router /auth/login [post]
// @Param schema.ReqLogin body schema.ReqLogin true "User Login"
// @Success 200 {object} schema.ResLogin
func Login(c *gin.Context) {
	var input schema.ReqLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(errors.New("Check Login Input failed"))
		return
	}

	u := models.User{Username: input.Username, Password: input.Password}
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.Error(errors.New("Invalid username or password"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// @Summary Logout
// @Description User Logout
// @Tags Auth
// @Router /auth/logout [post]
// @Success 200 {object} schema.ResLogout
func Logout(c *gin.Context) { c.JSON(200, gin.H{"message": "logout success"}) }
