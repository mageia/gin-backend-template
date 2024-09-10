package controller

import (
	"api-server/models"
	"api-server/token"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags Admin/User
// @Router /admin/user [get]
// @Success 200 {object} models.User "用户信息"
func RetrieveCurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var u models.User
	if err := models.DB.First(&u, userId).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u.SafeResponse()})
}

// UpdateUserInput 更新用户输入结构
type UpdateUserInput struct {
	Avatar *string `json:"avatar"`
}

// @Summary 更新当前用户信息
// @Description 更新当前登录用户的信息
// @Tags Admin/User
// @Router /admin/user [put]
// @Param user body UpdateUserInput true "用户更新信息"
// @Success 200 {object} models.User "更新后的用户信息"
func UpdateCurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if input.Avatar == nil {
		c.JSON(http.StatusOK, gin.H{"message": "无需更新"})
		return
	}

	var u models.User
	if err := models.DB.First(&u, userId).Error; err != nil {
		c.Error(err)
		return
	}

	if err := models.DB.Model(&u).Update("avatar", *input.Avatar).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u.SafeResponse()})
}

// Policy 策略结构
type Policy struct {
	Sub string `json:"sub" binding:"required"`
	Obj string `json:"obj" binding:"required"`
	Act string `json:"act" binding:"required"`
}

// @Summary 添加新策略
// @Description 添加新的访问控制策略
// @Tags Admin/Policy
// @Router /admin/policies [post]
// @Param policy body Policy true "新策略信息"
func PostPolicy(c *gin.Context) {
	var policy Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.Error(errors.New("检查策略输入失败"))
		return
	}

	e, _ := models.GetEnforcer()
	if ok, err := e.AddPolicy(policy.Sub, policy.Obj, policy.Act); !ok || err != nil {
		c.Error(errors.New("添加策略失败"))
		return
	}

	_ = e.LoadPolicy()
	c.JSON(http.StatusOK, gin.H{"policy": policy})
}

// @Summary 获取所有策略
// @Description 获取系统中所有的访问控制策略
// @Tags Admin/Policy
// @Router /admin/policies [get]
// @Success 200 {array} []string "策略列表"
// @Failure 500 {object} string "获取策略失败"
func GetPolicies(c *gin.Context) {
	e, _ := models.GetEnforcer()
	policies, err := e.GetPolicy()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"policies": policies})
}

type ReqUpdatePolicy struct {
	Old Policy `json:"old" binding:"required"`
	New Policy `json:"new" binding:"required"`
}

// @Summary 更新策略
// @Description 更新现有的访问控制策略
// @Tags Admin/Policy
// @Router /admin/policies [put]
// @Param policyUpdate body ReqUpdatePolicy true "策略更新信息"
// @Success 200 {object} Policy "更新后的策略"
func UpdatePolicy(c *gin.Context) {
	var policyUpdate ReqUpdatePolicy
	if err := c.ShouldBindJSON(&policyUpdate); err != nil {
		c.Error(errors.New("检查策略输入失败"))
		return
	}

	e, _ := models.GetEnforcer()
	if ok, err := e.UpdatePolicy(
		[]string{policyUpdate.Old.Sub, policyUpdate.Old.Obj, policyUpdate.Old.Act},
		[]string{policyUpdate.New.Sub, policyUpdate.New.Obj, policyUpdate.New.Act},
	); !ok || err != nil {
		c.Error(errors.New("更新策略失败"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"policy": policyUpdate.New})
}

// @Summary 删除策略
// @Description 删除指定的访问控制策略
// @Tags Admin/Policy
// @Router /admin/policies [delete]
// @Param policy body Policy true "要删除的策略信息"
// @Success 204 "删除成功"
func DeletePolicy(c *gin.Context) {
	var policy Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.Error(errors.New("检查策略输入失败"))
		return
	}

	e, _ := models.GetEnforcer()
	if ok, err := e.RemovePolicy(policy.Sub, policy.Obj, policy.Act); !ok || err != nil {
		c.Error(errors.New("删除策略失败"))
		return
	}

	c.Status(http.StatusNoContent)
}
