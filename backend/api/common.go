package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func jsonError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func jsonMessage(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"msg": msg})
}

func parsePathInt(c *gin.Context, key string) (int, bool) {
	v, err := strconv.Atoi(c.Param(key))
	if err != nil {
		jsonError(c, http.StatusBadRequest, "Invalid "+key)
		return 0, false
	}
	return v, true
}

func parseQueryIntWithDefault(c *gin.Context, key string, defaultValue int) (int, bool) {
	raw := c.DefaultQuery(key, strconv.Itoa(defaultValue))
	v, err := strconv.Atoi(raw)
	if err != nil {
		jsonError(c, http.StatusBadRequest, "Invalid "+key)
		return 0, false
	}
	return v, true
}

func ensureAdmin(c *gin.Context) bool {
	if c.GetString("role") != "admin" {
		jsonError(c, http.StatusForbidden, "Permission denied: Admins only")
		return false
	}
	return true
}
