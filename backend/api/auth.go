package api

import (
	"Ts3Panel/database"
	"Ts3Panel/models"
	"Ts3Panel/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		jsonError(c, http.StatusBadRequest, "username is required")
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "password hash failed")
		return
	}
	user := models.User{Username: req.Username, Password: hash, Role: "user"}

	var userCount int64
	if err := database.DB.Model(&models.User{}).Count(&userCount).Error; err == nil && userCount == 0 {
		user.Role = "admin"
	}

	if err := database.DB.Create(&user).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "User creation failed (username exists?)")
		return
	}
	jsonMessage(c, http.StatusCreated, "Created")
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Username = strings.TrimSpace(req.Username)

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		jsonError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		jsonError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, "Token generation failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
