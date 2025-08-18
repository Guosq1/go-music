package apiGateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gsq/music_bakcend_micorservice/pkg/apihandler"
)

func Setup() *gin.Engine {

	r := gin.Default()
	r.Static("/pages", "./static/pages")
	r.Static("/resources", "./static/resources")
	registerRoutes(r)

	return r
}

func registerRoutes(r *gin.Engine) {
	r.GET("/", apiHandler.Root)

	r.GET("/health", apiHandler.HealthCheck)

	api := r.Group("/api")
	{
		api.GET("/songs", apiHandler.GetSongs)
		api.GET("/songs/:name", apiHandler.GetSongbyName)
		api.POST("/login", apiHandler.Login)
		api.POST("/register", apiHandler.Register)
		api.GET("/songs/search", apiHandler.SearchSongs)
	}

}
