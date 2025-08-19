package apiHandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsq/music_bakcend_micorservice/myredis"
	"github.com/gsq/music_bakcend_micorservice/pkg/model"
	"github.com/gsq/music_bakcend_micorservice/utils"
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
		Remember bool   `json:"remember"`
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

	var expireTime time.Duration
	if req.Remember {
		expireTime = time.Hour * 24 * 30
	} else {
		expireTime = time.Minute * 30
	}

	tokenString, err := utils.GenerateJWT(user.ID, expireTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	err = myredis.Rdb.Set(myredis.Ctx, tokenString, user.ID, expireTime).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存token失败"})
		return
	}

	maxAge := 1800
	if req.Remember {
		maxAge = 30 * 24 * 3600
	}

	c.SetCookie("jwt_token", tokenString, maxAge, "/", "", false, true)
	// 这里可以返回 JWT 或 session
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful", "user": user.Username,
		"remember": req.Remember,
	})
}
