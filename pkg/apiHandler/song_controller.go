package apiHandler

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsq/music_bakcend_micorservice/database"
	"github.com/gsq/music_bakcend_micorservice/pkg/model"
)

var song model.Song

func Root(c *gin.Context) {
	htmlPath := filepath.Join("static", "pages", "index.html")
	c.File(htmlPath)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "live",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "music-service",
		"version":   "1.0.0",
	})
}

func GetSongs(c *gin.Context) {
	songs := model.GetSongs()
	c.JSON(http.StatusOK, songs)
}

func GetSongbyName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	song, err := model.GetSongByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	// 返回 JSON
	c.JSON(http.StatusOK, song)
}

func SearchSongs(c *gin.Context) {
	keyword := c.Query("q") // URL 中 ?q=关键词

	songs, err := model.SearchSongByKeyword(keyword)
	//fmt.Println(songs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// RecordPlayHistory 记录播放历史
func RecordPlayHistory(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	songIDStr := c.Param("songId")
	if songIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "songId is required"})
		return
	}

	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid songId"})
		return
	}

	err = model.CreatePlayHistory(userID.(uint), uint(songID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除Redis缓存
	cacheKey := "play_history:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	database.Rdb.Del(database.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "播放历史记录成功"})
}

// DeletePlayHistory 删除用户的播放历史记录
func DeletePlayHistory(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	songIDStr := c.Param("songId")
	if songIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "songId is required"})
		return
	}

	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid songId"})
		return
	}

	err = model.DeletePlayHistory(userID.(uint), uint(songID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除Redis缓存
	cacheKey := "play_history:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	database.Rdb.Del(database.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "播放历史删除成功"})
}

// GetUserPlayHistory 获取用户播放历史（使用Redis缓存）
func GetUserPlayHistory(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 从Redis获取缓存
	cacheKey := "play_history:" + strconv.FormatUint(uint64(userID.(uint)), 10)
	cacheData, err := database.Rdb.Get(database.Ctx, cacheKey).Result()
	
	if err == nil {
		// 缓存命中
		var histories []model.PlayHistory
		if err := json.Unmarshal([]byte(cacheData), &histories); err == nil {
			c.JSON(http.StatusOK, histories)
			return
		}
	}

	// 缓存未命中或解析失败，从数据库获取
	// 获取查询参数中的限制数量
	limitStr := c.Query("limit")
	limit := 20 // 默认限制20条
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 从数据库获取播放历史
	histories, err := model.GetUserPlayHistory(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将数据存入Redis缓存，设置过期时间为1小时
	if data, err := json.Marshal(histories); err == nil {
		database.Rdb.Set(database.Ctx, cacheKey, data, time.Hour)
	}

	// 返回数据
	c.JSON(http.StatusOK, histories)
}
