package controller

import (
	"api-server/models"
	"api-server/token"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RetrieveCurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var u models.User
	if e := models.DB.First(&u, userId).Error; e != nil {
		c.Error(e)
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
		c.Error(err)
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if input.Avatar == nil {
		c.JSON(200, gin.H{"message": "Nothing to do"})
		c.Error(err)
		return
	}

	var u models.User
	if e := models.DB.First(&u, userId).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		c.Error(err)
		return
	}

	if e := models.DB.Model(&u).Update("avatar", *input.Avatar).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": e.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

type Policy struct {
	Sub string `json:"sub" binding:"required"`
	Obj string `json:"obj" binding:"required"`
	Act string `json:"act" binding:"required"`
}

func PostPolicy(c *gin.Context) {
	var policy Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.Error(errors.New("Check policy input failed"))
		return
	}

	e, _ := models.GetEnforcer()
	if ok, err := e.AddPolicy(policy.Sub, policy.Obj, policy.Act); !ok || err != nil {
		c.Error(errors.New("Add policy failed"))
		return
	}

	_ = e.LoadPolicy()
	c.JSON(http.StatusOK, gin.H{"policy": policy})
}

func GetPolicy(c *gin.Context) {
	e, _ := models.GetEnforcer()
	policies, err := e.GetPolicy()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"policies": policies})
}

func PutPolicy(c *gin.Context) {
	var policy struct {
		Old Policy `json:"old" binding:"required"`
		New Policy `json:"new" binding:"required"`
	}
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.Error(errors.New("Check policy input failed"))
		return
	}
	e, _ := models.GetEnforcer()
	if ok, err := e.UpdatePolicy([]string{
		policy.Old.Sub,
		policy.Old.Obj,
		policy.Old.Act,
	}, []string{
		policy.New.Sub,
		policy.New.Obj,
		policy.New.Act,
	}); !ok || err != nil {
		c.Error(errors.New("Update policy failed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"policy": policy.New})
}

func DelPolicy(c *gin.Context) {
	var policy Policy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.Error(errors.New("Check policy input failed"))
		return
	}

	e, _ := models.GetEnforcer()
	if ok, err := e.RemovePolicy(policy.Sub, policy.Obj, policy.Act); !ok || err != nil {
		c.Error(errors.New("Remove policy failed"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
