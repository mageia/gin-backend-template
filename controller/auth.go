package controller

import (
	"api-server/models"
	"api-server/schema"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 用户注册
// @Description 处理用户注册请求，创建新用户账户
// @Tags Auth
// @Param user body schema.ReqRegister true "用户注册信息"
// @Success 200 {object} schema.ResRegister "注册成功"
// @Router /auth/register [post]
func AuthRegister(c *gin.Context) {
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

// @Summary 用户登录
// @Description 处理用户登录请求，验证凭据并返回访问令牌
// @Tags Auth
// @Param credentials body schema.ReqLogin true "用户登录凭据"
// @Success 200 {object} schema.ResLogin "登录成功，返回访问令牌"
// @Router /auth/login [post]
func AuthLogin(c *gin.Context) {
	var input schema.ReqLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(errors.New("Check Login Input failed"))
		return
	}

	token, err := models.LoginCheck(input.Username, input.Password)
	if err != nil {
		c.Error(errors.New("Invalid username or password"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// @Summary 用户登出
// @Description 处理用户登出请求，使当前会话失效
// @Tags Auth
// @Router /auth/logout [post]
// @Security ApiKeyAuth
func AuthLogout(c *gin.Context) {
	//TODO: 实现吊销访问令牌的逻辑
	c.JSON(200, gin.H{"message": "logout success"})
}
