package apiHandler

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
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
