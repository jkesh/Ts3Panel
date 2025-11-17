package api

import (
	"Ts3Panel/database"
	"Ts3Panel/models"
	"Ts3Panel/utils"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hash, _ := utils.HashPassword(req.Password)
	user := models.User{Username: req.Username, Password: hash}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "User creation failed (username exists?)"})
		return
	}
	c.JSON(201, gin.H{"msg": "Created"})
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.Role)
	c.JSON(200, gin.H{"token": token})
}
