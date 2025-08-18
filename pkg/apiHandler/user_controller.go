package apiHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsq/music_bakcend_micorservice/pkg/model"
)

func Register(c *gin.Context) {

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// DAO: 调用 model 层
	user, err := model.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username may already exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "user": user.Username})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// 这里可以返回 JWT 或 session
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user.Username})
}
